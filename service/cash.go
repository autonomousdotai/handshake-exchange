package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type CashService struct {
	dao     *dao.CashDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CashService) GetCashStore(userId string) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}

	if cashTO.Found {
		cash = cashTO.Object.(bean.CashStore)
	} else {
		ce.NotFound = true
	}

	return
}

func (s CashService) AddCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}

	var err error
	if cashTO.Found {
		ce.SetStatusKey(api_error.CashStoreExists)
	} else {
		body.UID = userId
		cash = body

		err = s.dao.AddCashStore(&cash)
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
			return
		}
		ce.NotFound = false
		solr_service.UpdateObject(bean.NewSolrFromCashStore(cash))
	}

	return
}

func (s CashService) UpdateCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}
	cash = cashTO.Object.(bean.CashStore)
	cash.Status = body.Status
	cash.Address = body.Address
	cash.BusinessType = body.BusinessType
	cash.Center = body.Center
	cash.Information = body.Information
	cash.Latitude = body.Latitude
	cash.Longitude = body.Longitude
	cash.Name = body.Name
	cash.Phone = body.Phone

	err := s.dao.UpdateCashStore(&cash)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}
	solr_service.UpdateObject(bean.NewSolrFromCashStore(cash))

	return
}

func (s CashService) GetProposeCashOrder(amountStr string, currency string, fiatCurrency string) (offer bean.CashOrder, ce SimpleContextError) {
	cryptoRateTO := s.miscDao.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_COINBASE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)

	systemFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_CASH_BUY_CRYPTO)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemFeeTO) {
		return
	}
	systemFee := systemFeeTO.Object.(bean.SystemFee)

	storeFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_CASH_STORE_SELL_CRYPTO)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemFeeTO) {
		return
	}
	storeFee := storeFeeTO.Object.(bean.SystemFee)

	amount, _ := decimal.NewFromString(amountStr)
	percentage, err := CreditServiceInst.GetCreditPoolPercentageByCache(currency, amount)
	if err != nil {
		if strings.Contains(err.Error(), "not enough") {
			ce.SetStatusKey(api_error.CreditOutOfStock)
			return
		}
		ce.SetError(api_error.GetDataFailed, err)
		return
	}

	price := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	totalWOFee := amount.Mul(price)

	feePercentage := decimal.NewFromFloat(systemFee.Value).Round(10)
	storeFeePercentage := decimal.NewFromFloat(storeFee.Value).Round(10)
	externalFeePercentage := decimal.NewFromFloat(float64(percentage)).Div(decimal.NewFromFloat(100))

	total, internalFee := dao.AddFeePercentage(totalWOFee, feePercentage)
	_, cashFee := dao.AddFeePercentage(totalWOFee, storeFeePercentage)
	_, externalFee := dao.AddFeePercentage(totalWOFee, externalFeePercentage)

	total = total.Add(cashFee)
	total = total.Add(externalFee)

	offer.FiatAmount = total.Round(2).String()
	offer.FiatCurrency = bean.USD.Code

	offer.Amount = amountStr
	offer.Currency = currency
	offer.Price = price.String()
	offer.Fee = internalFee.Round(2).String()
	offer.FeePercentage = feePercentage.String()
	offer.StoreFee = cashFee.Round(2).String()
	offer.StoreFeePercentage = storeFeePercentage.String()
	offer.ExternalFee = externalFee.Round(2).String()
	offer.ExternalFeePercentage = externalFeePercentage.String()

	if fiatCurrency != "" {
		rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatCurrency)
		if !rateTO.HasError() {
			rate := rateTO.Object.(bean.CurrencyRate)
			rateNumber := decimal.NewFromFloat(rate.Rate)
			tmpAmount := total.Mul(rateNumber)
			offer.FiatLocalAmount = tmpAmount.Round(2).String()
			offer.FiatLocalCurrency = fiatCurrency
			offer.LocalStoreFee = cashFee.Mul(rateNumber).Round(2).String()
		}
	}

	return
}

func (s CashService) GetProposeCashAmount(currency string, fiatCurrency string) (offer bean.CashOrder, ce SimpleContextError) {
	_, amountStr, err := CreditServiceInst.GetCreditPoolLeastAmountByCache(currency)
	if err != nil {
		if strings.Contains(err.Error(), "not enough") {
			ce.SetStatusKey(api_error.CreditOutOfStock)
			return
		}
		ce.SetError(api_error.GetDataFailed, err)
		return
	}

	offer, ce = s.GetProposeCashOrder(amountStr, currency, fiatCurrency)
	offer.StoreFee = ""
	offer.Price = ""
	leastAmount := common.StringToDecimal(offer.Amount)
	fiatAmount := common.StringToDecimal(offer.FiatAmount)

	offer.FiatAmount = fiatAmount.Div(leastAmount).RoundBank(2).String()
	if fiatCurrency != "" {
		fiatLocalAmount := common.StringToDecimal(offer.FiatLocalAmount)
		offer.FiatLocalAmount = fiatLocalAmount.Div(leastAmount).RoundBank(2).String()
	}
	offer.Amount = "1"

	return
}

func (s CashService) AddOrder(userId string, orderBody bean.CashOrder) (order bean.CashOrder, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)
	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}
	cash := cashTO.Object.(bean.CashStore)

	orderTest, testOfferCE := s.GetProposeCashOrder(orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency)
	if ce.FeedContextError(api_error.GetDataFailed, testOfferCE) {
		return
	}

	orderBody.UID = userId
	if orderTest.FiatAmount != orderBody.FiatAmount {
		notOk := true
		testFiatAmount := common.StringToDecimal(orderTest.FiatAmount)
		inputFiatAmount := common.StringToDecimal(orderBody.FiatAmount)
		if inputFiatAmount.GreaterThanOrEqual(testFiatAmount) {
			notOk = false
		} else {
			delta := testFiatAmount.Sub(inputFiatAmount)
			deltaPercentage := delta.Div(testFiatAmount)
			if deltaPercentage.LessThanOrEqual(decimal.NewFromFloat(0.01)) {
				notOk = false
			}
		}

		if notOk {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}

	if orderTest.Currency != orderBody.Currency {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if orderTest.Currency == bean.ETH.Code {
		if !common.CheckETHAddress(orderBody.Address) {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	} else {
		if common.CheckETHAddress(orderBody.Address) {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}

	if orderBody.Currency != bean.ETH.Code && orderBody.Currency != bean.BTC.Code && orderBody.Currency != bean.BCH.Code {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}

	// Minimum amount
	amount, _ := decimal.NewFromString(orderBody.Amount)
	if orderBody.Currency == bean.ETH.Code {
		if amount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if orderBody.Currency == bean.BTC.Code {
		if amount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if orderBody.Currency == bean.BCH.Code {
		if amount.LessThan(bean.MIN_BCH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}

	// Make buy order
	isSuccess := false
	var creditTrans *bean.CreditTransaction

	creditTrans = &bean.CreditTransaction{
		ToUID:      userId,
		Amount:     orderBody.Amount,
		Currency:   orderBody.Currency,
		Percentage: common.StringToDecimal(orderTest.ExternalFeePercentage).Mul(common.StringToDecimal("100")).String(), // Convert to 3%
	}
	transCE := CreditServiceInst.AddCreditTransaction(creditTrans)
	if ce.SetError(api_error.CreditOutOfStock, transCE.CheckError()) {
		isSuccess = false
	} else {
		if creditTrans.Id != "" {
			isSuccess = true
		} else {
			isSuccess = false
		}
	}

	if isSuccess {
		orderBody.UserInfo = map[string]string{
			"name":  cash.Name,
			"phone": cash.Phone,
		}
		setupCashOrder(&orderBody, orderTest, *creditTrans)
		err := s.dao.AddCashOrder(&orderBody)
		order = orderBody
		if ce.SetError(api_error.AddDataFailed, err) {
			return
		}

		s.dao.UpdateNotificationCashOrder(order)
		order.CreatedAt = time.Now().UTC()
		solr_service.UpdateObject(bean.NewSolrFromCashOrder(order, cash))
	}

	return
}

func (s CashService) FinishOrder(refCode string, amount string, fiatCurrency string) (order bean.CashOrder, overSpent string, ce SimpleContextError) {
	cashOrderRefCodeTO := s.dao.GetCashOrderRefCode(refCode)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderRefCodeTO) {
		return
	}
	orderRefCode := cashOrderRefCodeTO.Object.(bean.CashOrderRefCode)

	cashOrderTO := s.dao.GetCashOrderByPath(orderRefCode.OrderRef)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order = cashOrderTO.Object.(bean.CashOrder)
	cashTO := s.dao.GetCashStore(order.UID)
	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}
	cash := cashTO.Object.(bean.CashStore)

	storePaymentTO := s.dao.GetCashStorePayment(order.Id)
	if storePaymentTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}
	inputAmount := common.StringToDecimal(amount)
	totalInputAmount := inputAmount
	var storePayment bean.CashStorePayment
	if storePaymentTO.Found {
		storePayment = storePaymentTO.Object.(bean.CashStorePayment)
		paymentAmount := common.StringToDecimal(storePayment.FiatAmount)
		totalInputAmount = totalInputAmount.Add(paymentAmount)
	} else {
		storePayment = bean.CashStorePayment{
			Order:        order.Id,
			FiatAmount:   amount,
			FiatCurrency: fiatCurrency,
		}
	}

	if fiatCurrency == order.FiatCurrency {
		storeFeeAmount := common.StringToDecimal(order.StoreFee)
		orderAmount := common.StringToDecimal(order.FiatAmount)
		checkAmount := orderAmount.Sub(storeFeeAmount)
		if totalInputAmount.LessThan(checkAmount) {
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_UNDER
		} else if totalInputAmount.GreaterThan(checkAmount) {
			overSpent = totalInputAmount.Sub(checkAmount).String()
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_OVER
			storePayment.OverSpent = overSpent
		} else {
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_MATCHED
		}
	} else if fiatCurrency == order.FiatLocalCurrency {
		storeFeeAmount := common.StringToDecimal(order.LocalStoreFee)
		orderAmount := common.StringToDecimal(order.FiatLocalAmount)
		checkAmount := orderAmount.Sub(storeFeeAmount)
		if totalInputAmount.LessThan(checkAmount) {
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_UNDER
		} else if totalInputAmount.GreaterThan(checkAmount) {
			overSpent = totalInputAmount.Sub(checkAmount).String()
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_OVER
			storePayment.OverSpent = overSpent
		} else {
			storePayment.Status = bean.CASH_STORE_PAYMENT_STATUS_MATCHED
		}
	}

	if storePaymentTO.Found {
		s.dao.UpdateCashStorePayment(&storePayment, inputAmount)
	} else {
		s.dao.AddCashStorePayment(&storePayment)
	}
	if storePayment.Status == "" || storePayment.Status == bean.CASH_STORE_PAYMENT_STATUS_UNDER {
		ce.SetStatusKey(api_error.InvalidAmount)
		return
	}

	if order.Status != bean.CASH_ORDER_STATUS_TRANSFERRING {
		ce.SetStatusKey(api_error.CashOrderStatusInvalid)
		return
	}

	orderRef := dao.GetCashOrderItemPath(order.Id)
	fiatAmount := common.StringToDecimal(order.RawFiatAmount)
	revenue := common.StringToDecimal(order.ExternalFee).Add(fiatAmount)
	fee := common.Zero

	ccCE := CreditServiceInst.FinishCreditTransaction(order.Currency, order.ProviderData.(string), orderRef, revenue, fee)
	if ccCE.HasError() {
		if ce.SetError(api_error.ExternalApiFailed, ccCE.CheckError()) {
			return
		}
	}

	if order.Currency == bean.ETH.Code {
		txHash, outNonce, outAddress, onChainErr := ReleaseContractFund(*s.miscDao, order.Address, order.Amount, order.Id, 1, "ETH_LOW_ADMIN_KEYS")

		order.ProviderWithdrawData = txHash
		if onChainErr != nil {
			order.ProviderWithdrawData = onChainErr.Error()
		}
		order.ProviderWithdrawDataExtra = map[string]interface{}{
			"nonce":   fmt.Sprintf("%d", outNonce),
			"address": outAddress,
		}
	} else {
		coinbaseTx, errWithdraw := coinbase_service.SendTransaction(order.Address, order.Amount, order.Currency,
			fmt.Sprintf("Withdraw tx = %s", order.Id), order.Id)
		if errWithdraw == nil {
			order.ProviderWithdrawData = coinbaseTx.Id
		} else {
			order.ProviderWithdrawData = errWithdraw.Error()
		}
	}

	order.Status = bean.CASH_ORDER_STATUS_SUCCESS
	err := s.dao.FinishCashOrder(&order, &cash)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}
	solr_service.UpdateObject(bean.NewSolrFromCashOrder(order, cash))

	return
}

func (s CashService) UpdateOrderReceipt(orderId string, cashOrder bean.CashOrderUpdateInput) (order bean.CashOrder, ce SimpleContextError) {
	cashOrderTO := s.dao.GetCashOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order = cashOrderTO.Object.(bean.CashOrder)
	order.ReceiptURL = cashOrder.ReceiptURL
	order.Status = bean.CASH_ORDER_STATUS_TRANSFERRING

	err := s.dao.UpdateCashStoreReceipt(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	cashTO := s.dao.GetCashStore(order.UID)
	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}
	cash := cashTO.Object.(bean.CashStore)
	solr_service.UpdateObject(bean.NewSolrFromCashOrder(order, cash))

	return
}

func (s CashService) RejectOrder(orderId string) (order bean.CashOrder, ce SimpleContextError) {
	cashOrderTO := s.dao.GetCashOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order = cashOrderTO.Object.(bean.CashOrder)
	cashTO := s.dao.GetCashStore(order.UID)
	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}
	cash := cashTO.Object.(bean.CashStore)

	storePaymentTO := s.dao.GetCashStorePayment(orderId)
	if storePaymentTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}

	if order.Status != bean.CASH_ORDER_STATUS_PROCESSING {
		ce.SetStatusKey(api_error.CashOrderStatusInvalid)
		return
	}

	orderRef := dao.GetCashOrderItemPath(order.Id)
	ccCE := CreditServiceInst.RevertCreditTransaction(order.Currency, order.ProviderData.(string), orderRef)
	if ccCE.HasError() {
		if ce.SetError(api_error.ExternalApiFailed, ccCE.CheckError()) {
			return
		}
	}

	order.Status = bean.CASH_ORDER_STATUS_CANCELLED
	err := s.dao.FinishCashOrder(&order, &cash)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}
	solr_service.UpdateObject(bean.NewSolrFromCashOrder(order, cash))

	return
}

func (s CashService) ListCashCenter(country string) (cashCenters []bean.CashCenter, ce SimpleContextError) {
	cashCenterTO := s.dao.ListCashCenter(country)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashCenterTO) {
		return
	}

	cashCenters = make([]bean.CashCenter, 0)
	for _, item := range cashCenterTO.Objects {
		cashCenter := item.(bean.CashCenter)
		cashCenters = append(cashCenters, cashCenter)
	}

	return
}

func (s CashService) SyncCashStoreToSolr(id string) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}
	cash = cashTO.Object.(bean.CashStore)

	solr_service.UpdateObject(bean.NewSolrFromCashStore(cash))

	return
}

func (s CashService) SyncCashOrderToSolr(id string) (cashOrder bean.CashOrder, ce SimpleContextError) {
	cashOrderTO := s.dao.GetCashOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	cashOrder = cashOrderTO.Object.(bean.CashOrder)
	cashTO := s.dao.GetCashStore(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}
	cash := cashTO.Object.(bean.CashStore)

	solr_service.UpdateObject(bean.NewSolrFromCashOrder(cashOrder, cash))

	return
}

func setupCashOrder(order *bean.CashOrder, orderTest bean.CashOrder, creditTrans bean.CreditTransaction) {
	fiatAmount := common.StringToDecimal(order.FiatAmount)
	fee := common.StringToDecimal(orderTest.Fee)
	externalFee := common.StringToDecimal(orderTest.ExternalFee)
	storeFee := common.StringToDecimal(orderTest.StoreFee)

	order.RawFiatAmount = fiatAmount.Sub(fee).Sub(storeFee).Sub(externalFee).String()
	order.Status = bean.CASH_ORDER_STATUS_PROCESSING
	order.Type = bean.INSTANT_OFFER_TYPE_BUY
	order.Provider = bean.INSTANT_OFFER_PROVIDER_CREDIT
	order.ProviderData = creditTrans.Id
	order.Fee = orderTest.Fee
	order.FeePercentage = orderTest.FeePercentage
	order.StoreFee = orderTest.StoreFee
	order.StoreFeePercentage = orderTest.StoreFeePercentage
	order.ExternalFee = orderTest.ExternalFee
	order.ExternalFeePercentage = orderTest.ExternalFeePercentage
	order.Price = orderTest.Price
	order.FiatLocalCurrency = orderTest.FiatLocalCurrency
	order.FiatLocalAmount = orderTest.FiatLocalAmount
	order.LocalStoreFee = orderTest.LocalStoreFee

	// duration, _ := strconv.Atoi(os.Getenv("CC_LIMIT_DURATION"))
	order.Duration = int64(24 * 3600) // 24 hours
}
