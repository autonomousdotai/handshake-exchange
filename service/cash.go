package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/shopspring/decimal"
	"strings"
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

	err := s.dao.UpdateCashStore(&cash)
	if ce.SetError(api_error.AddDataFailed, err) {
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
		}
	}

	return
}

func (s CashService) AddOrder(userId string, orderBody bean.CashOrder) (order bean.CashOrder, ce SimpleContextError) {
	orderTest, testOfferCE := s.GetProposeCashOrder(orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency)
	if ce.FeedContextError(api_error.GetDataFailed, testOfferCE) {
		return
	}

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
		setupStoreOrder(&orderBody, orderTest, *creditTrans)
		// order, err = s.dao.AddStoreOrder(orderBody)
	}

	return
}

func setupStoreOrder(order *bean.CashOrder, orderTest bean.CashOrder, creditTrans bean.CreditTransaction) {
	fiatAmount, _ := decimal.NewFromString(order.FiatAmount)
	fee, _ := decimal.NewFromString(orderTest.ExternalFee)

	order.RawFiatAmount = fiatAmount.Sub(fee).String()
	order.Status = bean.INSTANT_OFFER_STATUS_PROCESSING
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
	order.FiatLocalAmount = orderTest.FiatAmount

	// duration, _ := strconv.Atoi(os.Getenv("CC_LIMIT_DURATION"))
	order.Duration = int64(24 * 3600) // 24 hours
}
