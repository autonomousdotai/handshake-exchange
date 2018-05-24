package service

import (
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/gdax_service"
	"github.com/autonomousdotai/handshake-exchange/integration/stripe_service"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go"
	"time"
	"strconv"
	"os"
)

type CreditCardService struct {
	dao      *dao.CreditCardDao
	miscDao  *dao.MiscDao
	userDao  *dao.UserDao
	transDao *dao.TransactionDao
}

func (s CreditCardService) GetProposeInstantOffer(amountStr string, currency string) (offer bean.InstantOffer, ce SimpleContextError) {
	cryptoRateTO := s.miscDao.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_COINBASE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	systemFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_INSTANT_BUY_CRYPTO)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemFeeTO) {
		return
	}
	systemFee := systemFeeTO.Object.(bean.SystemFee)

	price := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	amount, _ := decimal.NewFromString(amountStr)
	totalWOFee := amount.Mul(price)
	feePercentage := decimal.NewFromFloat(systemFee.Value).Round(10)
	total, fee := dao.AddFeePercentage(totalWOFee, feePercentage)

	offer.FiatAmount = total.Round(2).String()
	offer.FiatCurrency = bean.USD.Code
	offer.Amount = amountStr
	offer.Currency = currency
	offer.Price = price.String()
	offer.Fee = fee.Round(2).String()
	offer.FeePercentage = feePercentage.String()

	return
}

func (s CreditCardService) PayInstantOffer(userId string, offerBody bean.InstantOffer) (offer bean.InstantOffer, ce SimpleContextError) {
	offerTest, testOfferCE := s.GetProposeInstantOffer(offerBody.Amount, offerBody.Currency)
	if ce.FeedContextError(api_error.GetDataFailed, testOfferCE) {
		return
	}
	if offerTest.FiatAmount != offerBody.FiatAmount {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerTest.Currency != offerBody.Currency {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)
	// Your CC got problem
	if profile.CreditCardStatus != bean.CREDIT_CARD_STATUS_OK {
		ce.SetStatusKey(api_error.InvalidCC)
		return
	}

	var err error
	var paymentMethodData bean.CreditCardInfo
	b, _ := json.Marshal(&offerBody.PaymentMethodData)
	json.Unmarshal(b, &paymentMethodData)

	saveCard := false
	var token string
	if paymentMethodData.Token == "" {
		token, err = stripe_service.CreateToken(paymentMethodData.CCNum, paymentMethodData.ExpirationDate, paymentMethodData.CVV)
		if ce.SetError(api_error.ExternalApiFailed, err) {
			return
		}
		saveCard = true
	}
	if paymentMethodData.Token == "true" {
		ccLimitCE := UserServiceInst.CheckCCLimit(offerBody.UID, offerBody.FiatAmount)
		if ccLimitCE.HasError() {
			ce.SetError(api_error.CCOverLimit, ccLimitCE.Error)
			return
		}
		paymentMethodData.Token = profile.CreditCard.Token
	}

	fiatAmount, _ := decimal.NewFromString(offerBody.FiatAmount)
	stripeCharge, err := stripe_service.Charge(token, paymentMethodData.Token, fiatAmount, fmt.Sprintf("Buy %s %s", offerBody.Amount, offerBody.Currency))
	if ce.SetError(api_error.ExternalApiFailed, err) {
		return
	}
	if stripeCharge.Status == "failed" {
		ce.SetStatusKey(api_error.ExternalApiFailed)
		return
	}

	// Make buy order
	isSuccess := false

	var ccTran bean.CCTransaction
	var gdaxResponse bean.GdaxOrderResponse

	setupCCTransaction(&ccTran, offerBody, stripeCharge)
	ccTran, err = s.dao.AddCCTransaction(ccTran)
	if ce.SetError(api_error.AddDataFailed, err) {
	} else {
		isSuccess = true
		gdaxResponse, err = gdax_service.PlaceOrder(offerBody.Amount, offerBody.Currency, offerTest.Price)
		if ce.SetError(api_error.ExternalApiFailed, err) {
			isSuccess = false
		} else {
			isSuccess = true
		}
	}

	if !isSuccess {
		// If failed, do refund
		stripeRefund, err := stripe_service.Refund(ccTran.ExternalId)
		if ce.SetError(api_error.ExternalApiFailed, err) {
			return
		}
		if stripeRefund.Status == "failed" {
			ce.SetStatusKey(api_error.ExternalApiFailed)
			return
		}
		ccTran.Status = bean.CC_TRANSACTION_STATUS_REFUNDED
		s.dao.UpdateCCTransactionStatus(ccTran)
	} else {
		setupInstantOffer(&offerBody, offerTest, gdaxResponse)
		offerBody.PaymentMethod = bean.INSTANT_OFFER_PAYMENT_METHOD_CC
		offerBody.PaymentMethodRef = dao.GetCCTransactionItemPath(offerBody.UID, ccTran.Id)

		transaction := bean.NewTransactionFromInstantOffer(offerBody)
		offer, err = s.dao.AddInstantOffer(offerBody, transaction, gdaxResponse.Id)
		if ce.SetError(api_error.AddDataFailed, err) {
			return
		}
		ccTran.DataRef = dao.GetInstantOfferItemPath(offer.UID, offer.Id)
		s.dao.UpdateCCTransaction(ccTran)
	}

	if isSuccess {
		if saveCard {
			token, _ = s.saveCreditCard(userId, paymentMethodData)
		} else {
			token = paymentMethodData.Token
		}
		// Update CC Track amount
		s.userDao.UpdateUserCCLimitAmount(userId, token, fiatAmount)
	}

	paymentMethodData.CCNum = ""
	paymentMethodData.CVV = ""
	paymentMethodData.Token = ""
	offer.PaymentMethodData = paymentMethodData

	return
}

func (s CreditCardService) FinishInstantOffers() (finishedInstantOffers []bean.InstantOffer, ce SimpleContextError) {
	pendingOffers, err := s.dao.ListPendingInstantOffer()
	if ce.SetError(api_error.GetDataFailed, err) {
		return
	} else {
		for _, pendingOffer := range pendingOffers {
			gdaxResponse, err := gdax_service.GetOrder(pendingOffer.ProviderId)
			isDone := false
			if err == nil {
				if gdaxResponse.Status == "done" {
					offer := s.finishInstantOffer(&pendingOffer, gdaxResponse, &ce)
					if ce.CheckError() != nil {
						// return
					} else {
						isDone = true
						finishedInstantOffers = append(finishedInstantOffers, offer)
					}
				}
			}

			if !isDone {
				// Over duration
				if time.Now().UTC().Sub(pendingOffer.CreatedAt).Seconds() > float64(pendingOffer.Duration) {
					s.cancelInstantOffer(&pendingOffer, &ce)
					if ce.CheckError() != nil {
						// return
					}
				}
			}
		}
	}

	return
}

func (s CreditCardService) saveCreditCard(userId string, paymentMethodData bean.CreditCardInfo) (string, error) {
	ccNum := paymentMethodData.CCNum[len(paymentMethodData.CCNum)-4:]
	profileTO := s.userDao.GetProfile(userId)
	profile := profileTO.Object.(bean.Profile)
	// Need to create another token to save customer
	token, err := stripe_service.CreateToken(paymentMethodData.CCNum, paymentMethodData.ExpirationDate, paymentMethodData.CVV)
	if err == nil {
		token, _ = stripe_service.CreateCustomer(profile.UserId, token)
		ccUserLimit, err := UserServiceInst.GetUserCCLimitFirstLevel()
		if err == nil {
			err = s.userDao.UpdateProfileCreditCard(userId, bean.UserCreditCard{
				CCNumber:       ccNum,
				ExpirationDate: paymentMethodData.ExpirationDate,
				Token:          token,
			}, ccUserLimit)
		}
	}

	return token, err
}

func (s CreditCardService) finishInstantOffer(pendingOffer *bean.PendingInstantOffer, gdaxResponse bean.GdaxOrderResponse, ce *SimpleContextError) (offer bean.InstantOffer) {
	offerTO := s.dao.GetInstantOffer(pendingOffer.UID, pendingOffer.InstantOffer)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}

	offer = offerTO.Object.(bean.InstantOffer)
	offer.ProviderData = gdaxResponse
	offer.Status = bean.INSTANT_OFFER_STATUS_SUCCESS

	ccTranTO := s.dao.GetCCTransactionByPath(offer.PaymentMethodRef)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	ccTran := ccTranTO.Object.(bean.CCTransaction)
	_, err := stripe_service.Capture(ccTran.ExternalId)
	if ce.SetError(api_error.ExternalApiFailed, err) {
		return
	}
	ccTran.Status = bean.CC_TRANSACTION_STATUS_CAPTURED

	s.dao.UpdateCCTransactionStatus(ccTran)

	transTO := s.transDao.GetTransactionByPath(offer.TransactionRef)
	var trans bean.Transaction
	if transTO.HasError() {
		// Just one to make sure we don't lost anything
		trans = bean.NewTransactionFromInstantOffer(offer)
	} else {
		trans = transTO.Object.(bean.Transaction)
	}
	trans.Status = bean.TRANSACTION_STATUS_SUCCESS

	gdaxWithdrawResponse, errWithdraw := gdax_service.WithdrawCrypto(offer.Amount, offer.Currency, offer.Address)
	if errWithdraw == nil {
		offer.ProviderWithdrawData = gdaxWithdrawResponse
	}

	_, err = s.dao.UpdateInstantOffer(offer, trans)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	return
}

func (s CreditCardService) cancelInstantOffer(pendingOffer *bean.PendingInstantOffer, ce *SimpleContextError) {
	offerTO := s.dao.GetInstantOffer(pendingOffer.UID, pendingOffer.InstantOffer)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}

	offer := offerTO.Object.(bean.InstantOffer)
	offer.Status = bean.INSTANT_OFFER_STATUS_CANCELLED

	ccTranTO := s.dao.GetCCTransactionByPath(offer.PaymentMethodRef)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	ccTran := ccTranTO.Object.(bean.CCTransaction)
	_, err := stripe_service.Refund(ccTran.ExternalId)
	if ce.SetError(api_error.ExternalApiFailed, err) {
		// return
	} else {
		ccTran.Status = bean.CC_TRANSACTION_STATUS_REFUNDED
		s.dao.UpdateCCTransactionStatus(ccTran)
	}

	transTO := s.transDao.GetTransactionByPath(offer.TransactionRef)
	var trans bean.Transaction
	if transTO.HasError() {
		// Just one to make sure we don't lost anything
		trans = bean.NewTransactionFromInstantOffer(offer)
	} else {
		trans = transTO.Object.(bean.Transaction)
	}
	trans.Status = bean.TRANSACTION_STATUS_FAILED

	_, err = s.dao.UpdateInstantOffer(offer, trans)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	// Decrease amount track
	fiatAmount, _ := decimal.NewFromString(offer.FiatAmount)
	profileTO := s.userDao.GetProfile(offer.UID)
	if !profileTO.HasError() {
		profile := profileTO.Object.(bean.Profile)
		if profile.CreditCard.Token != "" {
			s.userDao.UpdateUserCCLimitAmount(offer.UID, profile.CreditCard.Token, fiatAmount.Mul(common.NegativeOne))
		}
	}

}

func setupInstantOffer(offer *bean.InstantOffer, offerTest bean.InstantOffer, gdaxResponse bean.GdaxOrderResponse) {
	fiatAmount, _ := decimal.NewFromString(offer.FiatAmount)
	fee, _ := decimal.NewFromString(offerTest.Fee)

	offer.RawFiatAmount = fiatAmount.Sub(fee).String()
	offer.Status = bean.INSTANT_OFFER_STATUS_PROCESSING
	offer.Type = bean.INSTANT_OFFER_TYPE_BUY
	offer.Provider = bean.INSTANT_OFFER_PROVIDER_GDAX
	offer.ProviderData = gdaxResponse
	offer.Fee = offerTest.Fee
	offer.FeePercentage = offerTest.FeePercentage
	offer.Price = offerTest.Price
	duration, _ := strconv.Atoi(os.Getenv("CC_LIMIT_DURATION"))
	offer.Duration = int64(duration)
}

func setupCCTransaction(ccTran *bean.CCTransaction, offerBody bean.InstantOffer, stripeCharge *stripe.Charge) {
	ccTran.Status = bean.CC_TRANSACTION_STATUS_PURCHASED
	ccTran.Provider = bean.CC_PROVIDER_STRIPE
	if stripeCharge.Tx != nil {
		ccTran.ProviderData = stripeCharge.Tx.ID
	}
	ccTran.Currency = bean.USD.Code
	ccTran.Amount = offerBody.FiatAmount
	ccTran.UID = offerBody.UID
	ccTran.Type = bean.CC_TRANSACTION_TYPE
	ccTran.ExternalId = stripeCharge.ID
}
