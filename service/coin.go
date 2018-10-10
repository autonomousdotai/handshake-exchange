package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type CoinService struct {
	dao     *dao.CoinDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CoinService) GetCoinQuote(amountStr string, currency string, fiatLocalCurrency string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	amount := common.StringToDecimal(amountStr)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_COINBASE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	price := amount.Mul(cryptoPrice)

	configTO := s.miscDao.GetSystemConfigFromCache(bean.COIN_ORDER_LIMIT)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
		return
	}
	systemConfig := configTO.Object.(bean.SystemConfig)
	// There is no free start on
	limit := common.StringToDecimal(systemConfig.Value)

	codFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_COD)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, codFeeTO) {
		return
	}
	codFee := codFeeTO.Object.(bean.SystemFee)
	codFeePercentage := decimal.NewFromFloat(codFee.Value).Round(3)

	bankFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_BANK)
	if price.GreaterThan(limit) {
		bankFeeTO = dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_BANK_HIGH)
	}
	if ce.FeedDaoTransfer(api_error.GetDataFailed, bankFeeTO) {
		return
	}
	bankFee := bankFeeTO.Object.(bean.SystemFee)
	bankFeePercentage := decimal.NewFromFloat(bankFee.Value).Round(3)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}

	codValue, codFeeValue := dao.AddFeePercentage(price, codFeePercentage)
	bankValue, bankFeeValue := dao.AddFeePercentage(price, bankFeePercentage)

	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)

	codValueLocal := codValue.Mul(rateNumber)
	codFeeLocal := codFeeValue.Mul(rateNumber)
	bankValueLocal := bankValue.Mul(rateNumber)
	bankFeeLocal := bankFeeValue.Mul(rateNumber)

	coinQuote.FiatCurrency = bean.USD.Code
	coinQuote.FiatLocalCurrency = fiatLocalCurrency
	coinQuote.FeePercentage = bankFeePercentage.String()
	coinQuote.RawFiatAmount = price.String()
	coinQuote.Price = cryptoPrice.String()

	coinQuote.FiatAmount = bankValue.RoundBank(2).String()
	coinQuote.Fee = bankFeeValue.RoundBank(2).String()

	coinQuote.FiatLocalAmount = bankValueLocal.RoundBank(2).String()
	coinQuote.FeeLocal = bankFeeLocal.RoundBank(2).String()

	if !price.GreaterThan(limit) {
		coinQuote.FeePercentageCOD = codFeePercentage.String()
		coinQuote.FiatAmountCOD = codValue.RoundBank(2).String()
		coinQuote.FeeCOD = codFeeValue.RoundBank(2).String()

		coinQuote.FiatLocalAmountCOD = codValueLocal.RoundBank(2).String()
		coinQuote.FeeLocalCOD = codFeeLocal.RoundBank(2).String()
	}

	coinQuote.Limit = limit.String()

	coinPoolTO := s.dao.GetCoinPool(currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinPoolTO) {
		return
	}
	coinPool := coinPoolTO.Object.(bean.CoinPool)
	usage := common.StringToDecimal(coinPool.Usage)
	usageLimit := common.StringToDecimal(coinPool.Limit)
	if usageLimit.LessThan(usage.Add(amount)) {
		ce.SetStatusKey(api_error.CreditOutOfStock)
		return
	}

	return
}

func (s CoinService) ListCoinCenter(country string) (coinCenters []bean.CoinCenter, ce SimpleContextError) {
	coinCenterTO := s.dao.ListCoinCenter(country)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinCenterTO) {
		return
	}

	coinCenters = make([]bean.CoinCenter, 0)
	for _, item := range coinCenterTO.Objects {
		coinCenter := item.(bean.CoinCenter)
		coinCenters = append(coinCenters, coinCenter)
	}

	return
}

func (s CoinService) AddOrder(userId string, orderBody bean.CoinOrder) (order bean.CoinOrder, ce SimpleContextError) {

	orderTest, testOfferCE := s.GetCoinQuote(orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency)
	if ce.FeedContextError(api_error.GetDataFailed, testOfferCE) {
		return
	}

	orderBody.UID = userId
	needToCheckAmount := false
	if orderTest.FiatAmount != orderBody.FiatAmount {
		needToCheckAmount = true
	}
	if orderBody.Type == bean.COIN_ORDER_TYPE_COD && orderTest.FiatLocalAmountCOD != orderBody.FiatLocalAmount {
		needToCheckAmount = true
	}
	if needToCheckAmount {
		notOk := true
		testFiatAmount := common.StringToDecimal(orderTest.FiatAmount)
		if orderBody.Type == bean.COIN_ORDER_TYPE_COD {
			testFiatAmount = common.StringToDecimal(orderTest.FiatLocalAmountCOD)
		}
		inputFiatAmount := common.StringToDecimal(orderBody.FiatAmount)
		if orderBody.Type == bean.COIN_ORDER_TYPE_COD {
			inputFiatAmount = common.StringToDecimal(orderBody.FiatLocalAmount)
		}
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

	if orderBody.Currency == bean.ETH.Code {
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

	coinPoolTO := s.dao.GetCoinPool(order.Currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinPoolTO) {
		return
	}
	coinPool := coinPoolTO.Object.(bean.CoinPool)
	usage := common.StringToDecimal(coinPool.Usage)
	limit := common.StringToDecimal(coinPool.Limit)
	if limit.LessThan(usage.Add(amount)) {
		ce.SetStatusKey(api_error.CreditOutOfStock)
		return
	}

	setupCoinOrder(&orderBody, orderTest)
	err := s.dao.AddCoinOrder(&orderBody)
	order = orderBody
	if strings.Contains(err.Error(), "out of stock") {
		if ce.SetError(api_error.CreditOutOfStock, err) {
			return
		}
	} else if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	order.CreatedAt = time.Now().UTC()
	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) UpdateOrderReceipt(orderId string, coinOrder bean.CoinOrderUpdateInput) (order bean.CoinOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinOrder)
	order.ReceiptURL = coinOrder.ReceiptURL
	order.Status = bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING

	err := s.dao.UpdateCoinStoreReceipt(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func setupCoinOrder(order *bean.CoinOrder, coinQuote bean.CoinQuote) {
	order.Price = coinQuote.Price
	order.RawFiatAmount = coinQuote.RawFiatAmount
	order.Status = bean.COIN_ORDER_STATUS_PENDING
	order.Fee = coinQuote.Fee
	order.FeePercentage = coinQuote.FeePercentage
	order.Price = coinQuote.Price

	if order.Type == bean.COIN_ORDER_TYPE_COD {
		order.Fee = coinQuote.FeeCOD
		order.FeePercentage = coinQuote.FeePercentageCOD
	}

	order.Duration = int64(30 * 60) // 30 minutes
}
