package service

import (
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/coinbase_service"
	"github.com/autonomousdotai/handshake-exchange/integration/solr_service"
	"github.com/go-errors/errors"
	"github.com/shopspring/decimal"
	"time"
)

type OfferService struct {
	dao      *dao.OfferDao
	userDao  *dao.UserDao
	transDao *dao.TransactionDao
}

func (s OfferService) GetOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	price, _ := decimal.NewFromString(offer.Price)
	percentage, _ := decimal.NewFromString(offer.Percentage)

	price, fiatPrice, fiatAmount, err := s.GetQuote(offer.Type, offer.Amount, offer.Currency, offer.FiatCurrency)
	if offer.Type == bean.OFFER_TYPE_SELL && price.Equal(common.Zero) {
		if ce.SetError(api_error.GetDataFailed, err) {
			return
		}
		markup := fiatAmount.Mul(percentage)
		fiatAmount = fiatAmount.Add(markup)
	}
	offer.Price = fiatPrice.Round(2).String()
	offer.FiatAmount = fiatAmount.Round(2).String()

	return
}

func (s OfferService) CreateOffer(userId string, offerBody bean.Offer) (offer bean.Offer, ce SimpleContextError) {
	// Offer type
	if offerBody.Type != bean.OFFER_TYPE_BUY && offerBody.Type != bean.OFFER_TYPE_SELL {
		ce.SetStatusKey(api_error.UnsupportedOfferType)
		return
	}

	if offerBody.Type == bean.OFFER_TYPE_BUY {
		// Need to set address to receive crypto
		if offerBody.UserAddress == "" {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		offerBody.Status = bean.OFFER_STATUS_ACTIVE
	} else {
		if offerBody.RefundAddress == "" {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		offerBody.Status = bean.OFFER_STATUS_CREATED
	}

	currencyInst := bean.CurrencyMapping[offerBody.Currency]
	if currencyInst.Code == "" {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}

	// Minimum amount
	amount, errFmt := decimal.NewFromString(offerBody.Amount)
	if ce.SetError(api_error.InvalidRequestBody, errFmt) {
		return
	}
	if currencyInst.Code == bean.ETH.Code {
		if amount.LessThan(decimal.NewFromFloat(0.01)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if currencyInst.Code == bean.BTC.Code {
		if amount.LessThan(decimal.NewFromFloat(0.1)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}

	if offerBody.Percentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(offer.Percentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		offerBody.Percentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		offerBody.Percentage = "0"
	}

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)
	offerBody.UID = userId

	if profile.ActiveOffers[currencyInst.Code] {
		ce.SetStatusKey(api_error.TooManyOffer)
		return
	}
	profile.ActiveOffers[currencyInst.Code] = true

	transCountTO := s.transDao.GetTransactionCount(offerBody.UID, offerBody.Currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, transCountTO) {
		return
	}
	transCount := transCountTO.Object.(bean.TransactionCount)

	// Only BTC need to generate address to transfer in
	if offerBody.Currency == bean.BTC.Code {
		addressResponse, err := coinbase_service.GenerateAddress(currencyInst.Code)
		if err != nil {
			ce.SetError(api_error.ExternalApiFailed, err)
			return
		}
		offerBody.SystemAddress = addressResponse.Data.Address
	}

	offerBody.TransactionCount = transCount
	offer, err := s.dao.AddOffer(offerBody, profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offer.CreatedAt = time.Now()
	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) ActiveOffer(address string, amountStr string) (offer bean.Offer, ce SimpleContextError) {
	addressMapTO := s.dao.GetOfferAddress(address)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, addressMapTO) {
		return
	}
	addressMap := addressMapTO.Object.(bean.OfferAddressMap)

	offerTO := s.dao.GetOffer(addressMap.Offer)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	if offer.Status != bean.OFFER_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}
	inputAmount, _ := decimal.NewFromString(amountStr)
	offerAmount, _ := decimal.NewFromString(offer.Amount)
	sub := offerAmount.Sub(inputAmount)
	if sub.Equal(common.Zero) {
		// Good
		offer.Status = bean.OFFER_STATUS_ACTIVE
		err := s.dao.UpdateOfferActive(offer)
		if ce.SetError(api_error.UpdateDataFailed, err) {
			return
		}

		solr_service.UpdateObject(bean.NewSolrFromOffer(offer))
	} else {
		// TODO Process to refund?
	}

	return
}

func (s OfferService) CloseOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	if offer.Status != bean.OFFER_STATUS_ACTIVE && offer.Status != bean.OFFER_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}
	offerProfile := s.getOfferProfile(offer, profile, &ce)
	offerProfile.ActiveOffers[offer.Currency] = false

	if ce.HasError() {
		return
	}

	offer.Status = bean.OFFER_STATUS_CLOSED
	err := s.dao.UpdateOfferClose(offer, offerProfile)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) ShakeOffer(userId string, offerId string, body bean.OfferShakeRequest) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if profile.UserId == offer.UID {
		ce.SetStatusKey(api_error.OfferPayMyself)
		return
	}

	if offer.Status != bean.OFFER_STATUS_ACTIVE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	offer.ToUID = userId
	if offer.Type == bean.OFFER_TYPE_SELL {
		// Only BTC needs to check
		if body.Address == "" && offer.Currency == bean.BTC.Code {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		_, numberErr := decimal.NewFromString(body.FiatAmount)
		if ce.SetError(api_error.InvalidNumber, numberErr) {
			return
		}
		offer.FiatAmount = body.FiatAmount
		offer.UserAddress = body.Address
		offer.Status = bean.OFFER_STATUS_SHAKE
	} else {
		// Only BTC needs to check
		if body.Address == "" && offer.Currency == bean.BTC.Code {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		offer.RefundAddress = body.Address
		offer.Status = bean.OFFER_STATUS_SHAKING
	}

	err := s.dao.UpdateOffer(offer, offer.GetUpdateOfferShaking())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) RejectShakeOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	offerProfile := s.getOfferProfile(offer, profile, &ce)
	offerProfile.ActiveOffers[offer.Currency] = false

	if ce.HasError() {
		return
	}

	if profile.UserId != offer.UID && profile.UserId != offer.ToUID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if offer.Status != bean.OFFER_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offer.Status = bean.OFFER_STATUS_REJECTED
	transCount := s.getFailedTransCount(offer)
	err := s.dao.UpdateOfferReject(offer, offerProfile, transCount)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) CompleteShakeOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if offer.Type == bean.OFFER_TYPE_SELL {
		if offer.UID != userId {
			ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
		}
	} else {
		if offer.ToUID != userId {
			ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
		}
	}

	offerProfile := s.getOfferProfile(offer, profile, &ce)
	offerProfile.ActiveOffers[offer.Currency] = false

	if ce.HasError() {
		return
	}

	if profile.UserId != offer.UID && profile.UserId != offer.ToUID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if offer.Status != bean.OFFER_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	offer.Status = bean.OFFER_STATUS_COMPLETED
	transCount := s.getSuccessTransCount(offer)
	err := s.dao.UpdateOfferCompleted(offer, offerProfile, transCount)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) WithdrawOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if profile.UserId != offer.UID && profile.UserId != offer.ToUID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if offer.Status != bean.OFFER_STATUS_CLOSED && offer.Status != bean.OFFER_STATUS_REJECTED && offer.Status != bean.OFFER_STATUS_COMPLETED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	if offer.Type == bean.OFFER_TYPE_SELL {
		if offer.Status == bean.OFFER_STATUS_REJECTED || offer.Status == bean.OFFER_STATUS_CLOSED {
			if offer.UID != userId {
				ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
				return
			}
		} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
			if offer.ToUID != userId {
				ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
				return
			}
		}
	} else {
		if offer.Status == bean.OFFER_STATUS_REJECTED {
			if offer.ToUID != userId {
				ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
				return
			}
		} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
			if offer.UID != userId {
				ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
				return
			}
		}
	}

	// Only BTC can transfer
	//var externalId string
	if offer.Currency == bean.BTC.Code {
		if offer.Status == bean.OFFER_STATUS_COMPLETED {
			if offer.UserAddress != "" {
				//Transfer
				description := fmt.Sprintf("Transfer to userId %s offerId %s status %s", userId, offer.Id, offer.Status)
				coinbaseResponse, err := coinbase_service.SendTransaction(offer.UserAddress, offer.Amount, offer.Currency, description, offer.Id)
				if ce.SetError(api_error.ExternalApiFailed, err) {
					return
				}
				offer.Provider = bean.OFFER_PROVIDER_COINBASE
				offer.ProviderData = coinbaseResponse
				//externalId = coinbaseResponse.Id
			} else {
				ce.SetStatusKey(api_error.InvalidRequestBody)
				return
			}
		} else if offer.Status == bean.OFFER_STATUS_REJECTED || offer.Status == bean.OFFER_STATUS_CLOSED {
			if offer.RefundAddress != "" {
				//Refund
				description := fmt.Sprintf("Refund to userId %s offerId %s status %s", userId, offer.Id, offer.Status)
				coinbaseResponse, err := coinbase_service.SendTransaction(offer.RefundAddress, offer.Amount, offer.Currency, description, offer.Id)
				if ce.SetError(api_error.ExternalApiFailed, err) {
					return
				}
				offer.Provider = bean.OFFER_PROVIDER_COINBASE
				offer.ProviderData = coinbaseResponse
			} else {
				ce.SetStatusKey(api_error.InvalidRequestBody)
				return
			}
		}

	}

	offer.Status = bean.OFFER_STATUS_WITHDRAW
	err := s.dao.UpdateOfferWithdraw(offer)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) getSuccessTransCount(offer bean.Offer) bean.TransactionCount {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, offer.Currency)
	var transCount bean.TransactionCount
	if !transCountTO.HasError() {
		transCount = transCountTO.Object.(bean.TransactionCount)
	}
	transCount.Currency = offer.Currency
	transCount.Success += 1

	return transCount
}

func (s OfferService) getFailedTransCount(offer bean.Offer) bean.TransactionCount {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, offer.Currency)
	var transCount bean.TransactionCount
	if !transCountTO.HasError() {
		transCount = transCountTO.Object.(bean.TransactionCount)
	}
	transCount.Currency = offer.Currency
	transCount.Failed += 1

	return transCount
}

//func (s OfferService) EndOffers() (ce SimpleContextError) {
//	transferMaps, err := s.dao.ListTransferMaps()
//	if ce.SetError(api_error.GetDataFailed, err) {
//		return
//	}
//
//	for _, transferMap := range transferMaps {
//		resp, err := coinbase_service.GetTransaction(transferMap.ExternalId, transferMap.Currency)
//		if err == nil {
//			if resp.Status == "completed" {
//				offer := s.endOffer(transferMap.Offer, resp, &ce)
//				if !ce.HasError() {
//					solr_service.UpdateObject(bean.NewSolrFromOffer(offer))
//				}
//			} else {
//				s.dao.UpdateTickTransferMap(transferMap)
//			}
//		}
//	}
//
//	return
//}
//
//func (s OfferService) endOffer(offerId string, providerData interface{}, ce *SimpleContextError) (offer bean.Offer) {
//	offerTO := s.dao.GetOffer(offerId)
//	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
//		return
//	}
//	offer = offerTO.Object.(bean.Offer)
//	if offer.Status != bean.OFFER_STATUS_COMPLETING {
//		ce.SetStatusKey(api_error.OfferStatusInvalid)
//	}
//
//	offer.ProviderData = providerData
//	offer.Status = bean.OFFER_STATUS_COMPLETED
//	err := s.dao.UpdateOfferCompleted(offer)
//	if ce.SetError(api_error.UpdateDataFailed, err) {
//		return
//	}
//
//	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))
//
//	return
//}

func (s OfferService) GetQuote(quoteType string, amountStr string, currency string, fiatCurrency string) (price decimal.Decimal, fiatPrice decimal.Decimal,
	fiatAmount decimal.Decimal, err error) {
	amount, numberErr := decimal.NewFromString(amountStr)
	to := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatCurrency)
	if numberErr != nil {
		err = numberErr
	}
	rate := to.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)
	tmpAmount := amount.Mul(rateNumber)

	if quoteType == "buy" {
		resp, errResp := coinbase_service.GetBuyPrice(currency)
		err = errResp
		if err != nil {
			return
		}
		price, _ = decimal.NewFromString(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else if quoteType == "sell" {
		resp, errResp := coinbase_service.GetSellPrice(currency)
		err = errResp
		if err != nil {
			return
		}
		price, _ := decimal.NewFromString(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else {
		err = errors.New(api_error.InvalidQueryParam)
	}

	return
}

func (s OfferService) SyncToSolr(offerId string) (offer bean.Offer, ce SimpleContextError) {
	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	solr_service.UpdateObject(bean.NewSolrFromOffer(offer))

	return
}

func (s OfferService) getOfferProfile(offer bean.Offer, profile bean.Profile, ce *SimpleContextError) (offerProfile bean.Profile) {
	if profile.UserId == offer.UID {
		offerProfile = profile
	} else {
		offerProfileTO := s.userDao.GetProfile(offer.UID)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, offerProfileTO) {
			return
		}
		offerProfile = offerProfileTO.Object.(bean.Profile)
	}

	return
}