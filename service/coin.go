package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/bitpay_service"
	"github.com/ninjadotorg/handshake-exchange/integration/bitstamp_service"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/crypto_service"
	"github.com/ninjadotorg/handshake-exchange/integration/slack_integration"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/integration/twilio_service"
	"github.com/ninjadotorg/handshake-exchange/service/email"
	"github.com/shopspring/decimal"
	"os"
	"strings"
	"time"
)

type CoinService struct {
	dao     *dao.CoinDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CoinService) GetCoinQuote(userId string, amountStr string, currency string, fiatLocalCurrency string, level string, check string, userCheck string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	amount := common.StringToDecimal(amountStr)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_BITSTAMP)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	price := amount.Mul(cryptoPrice)

	fiatValidLocalCurrency := fiatLocalCurrency
	if fiatValidLocalCurrency != bean.USD.Code && fiatValidLocalCurrency != bean.HKD.Code && fiatValidLocalCurrency != bean.VND.Code {
		fiatValidLocalCurrency = bean.USD.Code
	}

	userLimitObj, limitCE := s.GetUserLimit(userId, fiatValidLocalCurrency, level)
	if limitCE.HasError() {
		ce.SetError(api_error.InvalidRequestBody, limitCE.Error)
		return
	}
	userLimit := common.StringToDecimal(userLimitObj.Limit)
	userUsage := common.StringToDecimal(userLimitObj.Usage)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}
	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)
	localPrice := price.Mul(rateNumber)

	if userCheck == "" && userLimit.LessThan(userUsage.Add(localPrice)) {
		ce.SetStatusKey(api_error.CoinOverLimit)
		return
	}

	codFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_COD)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, codFeeTO) {
		return
	}
	codFee := codFeeTO.Object.(bean.SystemFee)
	codFeePercentage := decimal.NewFromFloat(codFee.Value).Round(3)

	bankFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_BANK)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, bankFeeTO) {
		return
	}
	bankFee := bankFeeTO.Object.(bean.SystemFee)
	bankFeePercentage := decimal.NewFromFloat(bankFee.Value).Round(3)

	codValue, codFeeValue := dao.AddFeePercentage(price, codFeePercentage)
	bankValue, bankFeeValue := dao.AddFeePercentage(price, bankFeePercentage)

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

	userLimitCOD := common.StringToDecimal(userLimitObj.LimitCOD)
	if !price.GreaterThan(userLimitCOD) {
		coinQuote.FeePercentageCOD = codFeePercentage.String()
		coinQuote.FiatAmountCOD = codValue.RoundBank(2).String()
		coinQuote.FeeCOD = codFeeValue.RoundBank(2).String()

		coinQuote.FiatLocalAmountCOD = codValueLocal.RoundBank(2).String()
		coinQuote.FeeLocalCOD = codFeeLocal.RoundBank(2).String()
	}

	coinQuote.Limit = userLimitCOD.String()

	if check == "" {
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
	}

	return
}

func (s CoinService) GetCoinQuoteReverse(userId string, fiatLocalAmountStr string, currency string, fiatLocalCurrency string,
	orderType string, level string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	fiatLocalAmount := common.StringToDecimal(fiatLocalAmountStr)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}
	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)

	fiatAmount := fiatLocalAmount.Div(rateNumber)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_BITSTAMP)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)

	fiatValidLocalCurrency := fiatLocalCurrency
	if fiatValidLocalCurrency != bean.USD.Code && fiatValidLocalCurrency != bean.HKD.Code && fiatValidLocalCurrency != bean.VND.Code {
		fiatValidLocalCurrency = bean.USD.Code
	}

	userLimitObj, limitCE := s.GetUserLimit(userId, fiatValidLocalCurrency, level)
	if limitCE.HasError() {
		ce.SetError(api_error.InvalidRequestBody, limitCE.Error)
	}
	userLimit := common.StringToDecimal(userLimitObj.Limit)
	userUsage := common.StringToDecimal(userLimitObj.Usage)

	if userLimit.LessThan(userUsage.Add(fiatLocalAmount)) {
		ce.SetStatusKey(api_error.CoinOverLimit)
		return
	}

	userLimitCOD := common.StringToDecimal(userLimitObj.LimitCOD)
	coinQuote.Currency = currency
	if orderType == bean.COIN_ORDER_TYPE_COD {
		coinQuote.Type = bean.COIN_ORDER_TYPE_COD
	} else {
		coinQuote.Type = bean.COIN_ORDER_TYPE_BANK
	}

	codFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_COD)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, codFeeTO) {
		return
	}
	codFee := codFeeTO.Object.(bean.SystemFee)
	codFeePercentage := decimal.NewFromFloat(codFee.Value).Round(3)

	bankFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_BANK)
	if fiatLocalAmount.GreaterThan(userLimitCOD) {
		coinQuote.Type = bean.COIN_ORDER_TYPE_BANK
	}
	if ce.FeedDaoTransfer(api_error.GetDataFailed, bankFeeTO) {
		return
	}
	bankFee := bankFeeTO.Object.(bean.SystemFee)
	bankFeePercentage := decimal.NewFromFloat(bankFee.Value).Round(3)

	codValue, _ := dao.RemoveFeePercentage(fiatAmount, codFeePercentage)
	bankValue, _ := dao.RemoveFeePercentage(fiatAmount, bankFeePercentage)

	amount := common.Zero
	if coinQuote.Type == bean.COIN_ORDER_TYPE_BANK {
		amount = bankValue.Div(cryptoPrice).Round(6)
		coinQuote.Amount = amount.String()
	} else {
		amount = codValue.Div(cryptoPrice).Round(6)
		coinQuote.Amount = amount.String()
	}

	coinQuote.FiatAmount = fiatAmount.RoundBank(2).String()
	coinQuote.FiatCurrency = bean.USD.Code
	coinQuote.Limit = userLimitCOD.String()

	coinQuote.FiatLocalCurrency = fiatLocalCurrency
	if coinQuote.Type == bean.COIN_ORDER_TYPE_COD {
		coinQuote.FiatLocalAmountCOD = fiatLocalAmountStr
	} else {
		coinQuote.FiatLocalAmount = fiatLocalAmountStr
	}

	if check == "" {
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
	}

	return
}

func (s CoinService) GetCoinSellingQuote(userId string, amountStr string, currency string, fiatLocalCurrency string, level string, check string, userCheck string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	amount := common.StringToDecimal(amountStr)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_BITSTAMP)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	price := amount.Mul(cryptoPrice)

	fiatValidLocalCurrency := fiatLocalCurrency
	if fiatValidLocalCurrency != bean.USD.Code && fiatValidLocalCurrency != bean.HKD.Code && fiatValidLocalCurrency != bean.VND.Code {
		fiatValidLocalCurrency = bean.USD.Code
	}

	userLimitObj, limitCE := s.GetSellingUserLimit(userId, fiatValidLocalCurrency, level)
	if limitCE.HasError() {
		ce.SetError(api_error.InvalidRequestBody, limitCE.Error)
		return
	}
	userLimit := common.StringToDecimal(userLimitObj.Limit)
	userUsage := common.StringToDecimal(userLimitObj.Usage)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}
	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)
	localPrice := price.Mul(rateNumber)

	if userCheck == "" && userLimit.LessThan(userUsage.Add(localPrice)) {
		ce.SetStatusKey(api_error.CoinOverLimit)
		return
	}

	bankFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_SELLING_ORDER_BANK)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, bankFeeTO) {
		return
	}
	bankFee := bankFeeTO.Object.(bean.SystemFee)
	bankFeePercentage := decimal.NewFromFloat(bankFee.Value).Round(3)

	bankValue, bankFeeValue := dao.AddFeePercentage(price, bankFeePercentage)

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

	userLimitCOD := common.StringToDecimal(userLimitObj.LimitCOD)
	coinQuote.Limit = userLimitCOD.String()

	if check == "" {
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
	}

	return
}

func (s CoinService) GetCoinSellingQuoteReverse(userId string, fiatLocalAmountStr string, currency string, fiatLocalCurrency string,
	orderType string, level string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	fiatLocalAmount := common.StringToDecimal(fiatLocalAmountStr)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}
	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)

	fiatAmount := fiatLocalAmount.Div(rateNumber)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_BITSTAMP)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)

	fiatValidLocalCurrency := fiatLocalCurrency
	if fiatValidLocalCurrency != bean.USD.Code && fiatValidLocalCurrency != bean.HKD.Code && fiatValidLocalCurrency != bean.VND.Code {
		fiatValidLocalCurrency = bean.USD.Code
	}

	userLimitObj, limitCE := s.GetSellingUserLimit(userId, fiatValidLocalCurrency, level)
	if limitCE.HasError() {
		ce.SetError(api_error.InvalidRequestBody, limitCE.Error)
	}
	userLimit := common.StringToDecimal(userLimitObj.Limit)
	userUsage := common.StringToDecimal(userLimitObj.Usage)

	if userLimit.LessThan(userUsage.Add(fiatLocalAmount)) {
		ce.SetStatusKey(api_error.CoinOverLimit)
		return
	}

	userLimitCOD := common.StringToDecimal(userLimitObj.LimitCOD)
	coinQuote.Currency = currency
	coinQuote.Type = bean.COIN_ORDER_TYPE_BANK

	bankFeeTO := dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_SELLING_ORDER_BANK)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, bankFeeTO) {
		return
	}
	bankFee := bankFeeTO.Object.(bean.SystemFee)
	bankFeePercentage := decimal.NewFromFloat(bankFee.Value).Round(3)

	bankValue, _ := dao.RemoveFeePercentage(fiatAmount, bankFeePercentage)

	amount := common.Zero
	if coinQuote.Type == bean.COIN_ORDER_TYPE_BANK {
		amount = bankValue.Div(cryptoPrice).Round(6)
		coinQuote.Amount = amount.String()
	}

	coinQuote.FiatAmount = fiatAmount.RoundBank(2).String()
	coinQuote.FiatCurrency = bean.USD.Code
	coinQuote.Limit = userLimitCOD.String()

	coinQuote.FiatLocalCurrency = fiatLocalCurrency
	if coinQuote.Type == bean.COIN_ORDER_TYPE_COD {
		coinQuote.FiatLocalAmountCOD = fiatLocalAmountStr
	} else {
		coinQuote.FiatLocalAmount = fiatLocalAmountStr
	}

	if check == "" {
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

func (s CoinService) ListCoinBank(country string) (coinBanks []bean.CoinBank, ce SimpleContextError) {
	coinBankTO := s.dao.ListCoinBank(country)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinBankTO) {
		return
	}

	coinBanks = make([]bean.CoinBank, 0)
	for _, item := range coinBankTO.Objects {
		coinBank := item.(bean.CoinBank)
		coinBanks = append(coinBanks, coinBank)
	}

	return
}

func (s CoinService) AddOrder(userId string, orderBody bean.CoinOrder) (order bean.CoinOrder, ce SimpleContextError) {
	orderTest, testOfferCE := s.GetCoinQuote(userId, orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency, orderBody.Level, "1", "")
	if ce.FeedContextError(testOfferCE.StatusKey, testOfferCE) {
		return
	}
	orderBody.UID = userId
	needToCheckAmount := false
	limit := common.StringToDecimal(orderTest.Limit)
	fiatAmount := common.StringToDecimal(orderBody.FiatAmount)
	if orderTest.FiatAmount != orderBody.FiatAmount {
		needToCheckAmount = true
	}
	if fiatAmount.GreaterThan(limit) {
		if orderTest.FiatAmount != orderBody.FiatAmount {
			needToCheckAmount = true
		}
	} else {
		if orderBody.Type == bean.COIN_ORDER_TYPE_COD && orderTest.FiatLocalAmountCOD != orderBody.FiatLocalAmount {
			needToCheckAmount = true
		} else if orderTest.FiatLocalAmount != orderBody.FiatLocalAmount {
			needToCheckAmount = true
		}
	}

	if needToCheckAmount {
		notOk := true
		testFiatAmount := common.StringToDecimal(orderTest.FiatAmount)
		inputFiatAmount := common.StringToDecimal(orderBody.FiatAmount)

		if fiatAmount.GreaterThan(limit) {
		} else {
			if orderBody.Type == bean.COIN_ORDER_TYPE_COD {
				testFiatAmount = common.StringToDecimal(orderTest.FiatLocalAmountCOD)
			} else {
				testFiatAmount = common.StringToDecimal(orderTest.FiatLocalAmount)
			}

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

	coinPoolTO := s.dao.GetCoinPool(orderBody.Currency)
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

	s.setupCoinOrder(&orderBody, orderTest)
	err := s.dao.AddCoinOrder(&orderBody)
	order = orderBody
	if err != nil {
		if strings.Contains(err.Error(), "out of stock") {
			if ce.SetError(api_error.CreditOutOfStock, err) {
				return
			}
		} else if ce.SetError(api_error.AddDataFailed, err) {
			return
		}
	}

	order.CreatedAt = time.Now().UTC()
	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))
	if order.Type == bean.COIN_ORDER_TYPE_COD {
		s.NotifyNewCoinOrder(order)
	}

	return
}

func (s CoinService) AddSellingOrder(userId string, orderBody bean.CoinSellingOrder) (order bean.CoinSellingOrder, ce SimpleContextError) {
	orderTest, testOfferCE := s.GetCoinSellingQuote(userId, orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency, orderBody.Level, "1", "")
	if ce.FeedContextError(testOfferCE.StatusKey, testOfferCE) {
		return
	}
	orderBody.UID = userId
	needToCheckAmount := false
	limit := common.StringToDecimal(orderTest.Limit)
	fiatAmount := common.StringToDecimal(orderBody.FiatAmount)
	if orderTest.FiatAmount != orderBody.FiatAmount {
		needToCheckAmount = true
	}
	if fiatAmount.GreaterThan(limit) {
		if orderTest.FiatAmount != orderBody.FiatAmount {
			needToCheckAmount = true
		}
	} else {
		if orderTest.FiatLocalAmount != orderBody.FiatLocalAmount {
			needToCheckAmount = true
		}
	}

	if needToCheckAmount {
		notOk := true
		testFiatAmount := common.StringToDecimal(orderTest.FiatAmount)
		inputFiatAmount := common.StringToDecimal(orderBody.FiatAmount)

		if fiatAmount.GreaterThan(limit) {
		} else {
			testFiatAmount = common.StringToDecimal(orderTest.FiatLocalAmount)
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

	// Is it a good address?
	addressTO := s.dao.GetCoinGenerateAddress(orderBody.Currency, orderBody.Address)
	if addressTO.HasError() || !addressTO.Found {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	coinPoolTO := s.dao.GetCoinSellingPool(orderBody.Currency)
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

	s.setupCoinSellingOrder(&orderBody, orderTest)
	err := s.dao.AddCoinSellingOrder(&orderBody)
	order = orderBody
	if err != nil {
		if strings.Contains(err.Error(), "out of stock") {
			if ce.SetError(api_error.CreditOutOfStock, err) {
				return
			}
		} else if ce.SetError(api_error.AddDataFailed, err) {
			return
		}
	}

	order.CreatedAt = time.Now().UTC()
	s.dao.UpdateNotificationCoinSellingOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))

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

	err := s.dao.UpdateCoinOrderReceipt(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))
	fmt.Println(order.Type)
	if order.Type == bean.COIN_ORDER_TYPE_BANK {
		s.NotifyNewCoinOrder(order)
	}

	return
}

func (s CoinService) UpdateOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinOrder)
	if order.Type == bean.COIN_ORDER_TYPE_COD && order.Status != bean.COIN_ORDER_STATUS_PENDING {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}
	if order.Type == bean.COIN_ORDER_TYPE_BANK && order.Status != bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}
	order.Status = bean.COIN_ORDER_STATUS_PROCESSING

	err := s.dao.UpdateCoinOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) UpdateSellingOrder(orderId string) (order bean.CoinSellingOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinSellingOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinSellingOrder)
	if order.Status != bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}
	order.Status = bean.COIN_ORDER_STATUS_PROCESSING

	err := s.dao.UpdateCoinSellingOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinSellingOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))

	return
}

func (s CoinService) CloseSellingOrder(orderId string) (order bean.CoinSellingOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinSellingOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	if order.Status != bean.COIN_ORDER_STATUS_PROCESSING {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	order = coinOrderTO.Object.(bean.CoinSellingOrder)
	order.Status = bean.COIN_ORDER_STATUS_SUCCESS

	err := s.dao.UpdateCoinSellingOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinSellingOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))

	return
}

func (s CoinService) CancelOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	order = s.cancelCoinOrder(orderId, bean.COIN_ORDER_STATUS_CANCELLED, &ce)
	return
}

func (s CoinService) CancelSellingOrder(orderId string) (order bean.CoinSellingOrder, ce SimpleContextError) {
	order = s.cancelCoinSellingOrder(orderId, bean.COIN_ORDER_STATUS_CANCELLED, &ce)
	return
}

func (s CoinService) RejectOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	order = s.cancelCoinOrder(orderId, bean.COIN_ORDER_STATUS_REJECTED, &ce)
	return
}

func (s CoinService) RemoveExpiredOrder() (ce SimpleContextError) {
	coinRefCodeTO := s.dao.ListCoinOrderRefCode()
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinRefCodeTO) {
		return
	}
	coinRefCodes := coinRefCodeTO.Objects

	for _, item := range coinRefCodes {
		coinRefCode := item.(bean.CoinOrderRefCode)
		if coinRefCode.CreatedAt.Add(time.Minute * time.Duration(coinRefCode.Duration)).Before(time.Now().UTC()) {
			s.expireOrder(coinRefCode.Order)
		}
	}

	return
}

func (s CoinService) SellingRemoveExpiredOrder() (ce SimpleContextError) {
	coinRefCodeTO := s.dao.ListCoinSellingOrderRefCode()
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinRefCodeTO) {
		return
	}
	coinRefCodes := coinRefCodeTO.Objects

	for _, item := range coinRefCodes {
		coinRefCode := item.(bean.CoinOrderRefCode)
		if coinRefCode.CreatedAt.Add(time.Minute * time.Duration(coinRefCode.Duration)).Before(time.Now().UTC()) {
			s.expireSellingOrder(coinRefCode.Order)
		}
	}

	return
}

func (s CoinService) FinishOrder(id string, amount string, fiatCurrency string) (order bean.CoinOrder, overSpent string, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinOrder)
	amount = order.FiatAmount
	fiatCurrency = order.FiatCurrency

	storePaymentTO := s.dao.GetCoinPayment(order.Id)
	if storePaymentTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, storePaymentTO)
		return
	}
	inputAmount := common.StringToDecimal(amount)
	totalInputAmount := inputAmount
	var coinPayment bean.CoinPayment
	if storePaymentTO.Found {
		coinPayment = storePaymentTO.Object.(bean.CoinPayment)
		paymentAmount := common.StringToDecimal(coinPayment.FiatAmount)
		totalInputAmount = totalInputAmount.Add(paymentAmount)
	} else {
		coinPayment = bean.CoinPayment{
			Order:        order.Id,
			FiatAmount:   amount,
			FiatCurrency: fiatCurrency,
		}
	}

	orderAmount := common.Zero
	if fiatCurrency == order.FiatCurrency {
		orderAmount = common.StringToDecimal(order.FiatAmount)
	} else if fiatCurrency == order.FiatLocalCurrency {
		orderAmount = common.StringToDecimal(order.FiatLocalAmount)
	}

	if totalInputAmount.LessThan(orderAmount) {
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_UNDER
	} else if totalInputAmount.GreaterThan(orderAmount) {
		overSpent = totalInputAmount.Sub(orderAmount).String()
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_OVER
		coinPayment.OverSpent = overSpent
	} else {
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_MATCHED
	}

	if storePaymentTO.Found {
		s.dao.UpdateCoinPayment(&coinPayment, inputAmount)
	} else {
		s.dao.AddCoinPayment(&coinPayment)
	}
	if coinPayment.Status == "" || coinPayment.Status == bean.COIN_PAYMENT_STATUS_UNDER {
		ce.SetStatusKey(api_error.InvalidAmount)
		return
	}

	if order.Status != bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING && order.Status != bean.COIN_ORDER_STATUS_PROCESSING && order.Status != bean.COIN_ORDER_STATUS_TRANSFER_FAILED {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	var externalId string
	var provider string
	if os.Getenv("ENVIRONMENT") == "dev" {
		if order.Currency == bean.ETH.Code {
			txHash, onChainErr := crypto_service.SendTransaction(order.Address, order.Amount, order.Currency)
			provider = bean.ETH_WALLET_NETWORK
			order.ProviderWithdrawData = txHash
			order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
			if onChainErr != nil {
				order.ProviderWithdrawData = onChainErr.Error()
				order.Status = bean.COIN_ORDER_STATUS_TRANSFER_FAILED
			}
		} else {
			order.ProviderWithdrawData = "xxx"
			order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
		}
	} else if os.Getenv("ENVIRONMENT") == "production" {
		// coinbaseTx, errWithdraw := coinbase_service.SendTransaction(order.Address, order.Amount, order.Currency,
		// fmt.Sprintf("Withdraw tx = %s", order.Id), order.Id)
		bitstampTx, errWithdraw := bitstamp_service.SendTransaction(order.Address, order.Amount, order.Currency,
			fmt.Sprintf("Withdraw tx = %s", order.Id), order.Id)

		if errWithdraw == nil {
			provider = bean.BTC_WALLET_BITSTAMP
			if bitstampTx.Id == 0 {
				order.ProviderWithdrawData = "Out of coin"
				order.Status = bean.COIN_ORDER_STATUS_TRANSFER_FAILED
			} else {
				externalId = fmt.Sprintf("%d", bitstampTx.Id)
				order.ProviderWithdrawData = externalId
				order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
			}
		} else {
			order.ProviderWithdrawData = errWithdraw.Error()
			order.Status = bean.COIN_ORDER_STATUS_TRANSFER_FAILED
		}
	}

	err := s.dao.FinishCoinOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	if order.Status == bean.COIN_ORDER_STATUS_TRANSFERRING {
		s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
			Provider:         provider,
			ProviderResponse: order.ProviderWithdrawData,
			DataType:         bean.OFFER_ADDRESS_MAP_COIN_ORDER,
			DataRef:          dao.GetCoinOrderItemPath(order.Id),
			UID:              order.UID,
			Description:      "",
			Amount:           order.Amount,
			Currency:         order.Currency,
			ExternalId:       externalId,
			TxHash:           order.ProviderWithdrawData.(string),
		})
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) FinishSellingOrder(id string, amount string, currency string, txHash string) (order bean.CoinSellingOrder, overSpent string, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinSellingOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinSellingOrder)
	orderAmount := common.StringToDecimal(order.Amount)
	orderCurrency := order.Currency

	storePaymentTO := s.dao.GetCoinSellingPayment(order.Id)
	if storePaymentTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, storePaymentTO)
		return
	}
	inputAmount := common.StringToDecimal(amount)
	totalInputAmount := inputAmount
	var coinPayment bean.CoinSellingPayment
	if storePaymentTO.Found {
		coinPayment = storePaymentTO.Object.(bean.CoinSellingPayment)
		paymentAmount := common.StringToDecimal(coinPayment.Amount)
		totalInputAmount = totalInputAmount.Add(paymentAmount)
	} else {
		coinPayment = bean.CoinSellingPayment{
			Order:    order.Id,
			Amount:   amount,
			Currency: currency,
		}
	}

	if currency != orderCurrency {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if totalInputAmount.LessThan(orderAmount) {
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_UNDER
	} else if totalInputAmount.GreaterThan(orderAmount) {
		overSpent = totalInputAmount.Sub(orderAmount).String()
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_OVER
		coinPayment.OverSpent = overSpent
	} else {
		coinPayment.Status = bean.COIN_PAYMENT_STATUS_MATCHED
	}

	if storePaymentTO.Found {
		s.dao.UpdateCoinSellingPayment(&coinPayment, inputAmount)
	} else {
		s.dao.AddCoinSellingPayment(&coinPayment)
	}
	if coinPayment.Status == "" {
		ce.SetStatusKey(api_error.InvalidAmount)
		return
	}

	if order.Status != bean.COIN_ORDER_STATUS_PENDING &&
		order.Status != bean.COIN_ORDER_STATUS_TRANSFERRING &&
		order.Status != bean.COIN_ORDER_STATUS_TRANSFER_FAILED {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	order.Status = bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING
	order.TxHash = txHash
	if coinPayment.Status == bean.COIN_PAYMENT_STATUS_UNDER {
		order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
		s.dao.UpdateCoinSellingOrder(&order)
	} else {
		err := s.dao.FinishCoinSellingOrder(&order)
		if ce.SetError(api_error.AddDataFailed, err) {
			return
		}
		s.NotifyNewCoinSellingOrder(order)
	}

	s.dao.UpdateNotificationCoinSellingOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))

	return
}

func (s CoinService) FinishCoinOrderPendingTransfer(ref string, txHash string) (order bean.CoinOrder, ce SimpleContextError) {
	cashOrderTO := s.dao.GetCoinOrderByPath(ref)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order = cashOrderTO.Object.(bean.CoinOrder)
	order.Status = bean.COIN_ORDER_STATUS_SUCCESS
	order.TxHash = txHash

	err := s.dao.UpdateCoinOrderReceipt(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) AddCoinReview(direction string, review bean.CoinReview) (ce SimpleContextError) {
	if direction == "buy" {
		cashOrderTO := s.dao.GetCoinOrder(review.Order)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
			return
		}
		order := cashOrderTO.Object.(bean.CoinOrder)
		if order.Reviewed {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		order.Reviewed = true

		s.dao.UpdateCoinOrderReview(&order)
		solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))
	} else if direction == "sell" {
		cashOrderTO := s.dao.GetCoinSellingOrder(review.Order)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
			return
		}
		order := cashOrderTO.Object.(bean.CoinSellingOrder)
		if order.Reviewed {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		order.Reviewed = true

		s.dao.UpdateCoinSellingOrderReview(&order)
		solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))
	}

	coinReviewCountTO := s.dao.GetCoinReviewCount(direction)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinReviewCountTO) {
		return
	}
	coinReviewCount := coinReviewCountTO.Object.(bean.CoinReviewCount)
	coinReviewCount.Count += 1

	err := s.dao.AddCoinReview(direction, &review, &coinReviewCount)
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}

	return
}

func (s CoinService) GetUserLimit(uid string, currency string, level string) (limit bean.CoinUserLimit, ce SimpleContextError) {
	limitTO := s.dao.GetCoinUserLimit(uid)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, limitTO) {
		return
	}

	configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "1", currency))
	if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
		return
	}
	systemConfig := configTO.Object.(bean.SystemConfig)
	limit.LimitCOD = systemConfig.Value

	if limitTO.Found {
		limit = limitTO.Object.(bean.CoinUserLimit)
		limit.LimitCOD = systemConfig.Value
		if limit.Currency != currency && !common.StringToDecimal(limit.Usage).Equal(common.Zero) {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		if limit.Level != level || limit.Currency != currency {
			configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "2", currency))
			if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
				return
			}
			systemConfig := configTO.Object.(bean.SystemConfig)

			limit.Currency = currency
			limit.Level = level
			limit.Limit = systemConfig.Value
			s.dao.UpdateCoinUserLimitLevel(&limit)
		}
		// do nothing
	} else {
		configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "2", currency))
		if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
			return
		}
		systemConfig := configTO.Object.(bean.SystemConfig)

		limit.Level = level
		limit.UID = uid
		limit.Limit = systemConfig.Value
		limit.Usage = common.Zero.String()
		limit.Currency = currency
		s.dao.AddCoinUserLimit(&limit)
	}

	return
}

func (s CoinService) GetSellingUserLimit(uid string, currency string, level string) (limit bean.CoinUserLimit, ce SimpleContextError) {
	limitTO := s.dao.GetCoinSellingUserLimit(uid)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, limitTO) {
		return
	}

	configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "1", currency))
	if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
		return
	}
	systemConfig := configTO.Object.(bean.SystemConfig)
	limit.LimitCOD = systemConfig.Value

	if limitTO.Found {
		limit = limitTO.Object.(bean.CoinUserLimit)
		limit.LimitCOD = systemConfig.Value
		if limit.Currency != currency && !common.StringToDecimal(limit.Usage).Equal(common.Zero) {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		if limit.Level != level || limit.Currency != currency {
			configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "2", currency))
			if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
				return
			}
			systemConfig := configTO.Object.(bean.SystemConfig)

			limit.Currency = currency
			limit.Level = level
			limit.Limit = systemConfig.Value
			s.dao.UpdateCoinSellingUserLimitLevel(&limit)
		}
		// do nothing
	} else {
		configTO := s.miscDao.GetSystemConfigFromCache(fmt.Sprintf("%s_%s_%s", bean.COIN_ORDER_LIMIT, "2", currency))
		if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
			return
		}
		systemConfig := configTO.Object.(bean.SystemConfig)

		limit.Level = level
		limit.UID = uid
		limit.Limit = systemConfig.Value
		limit.Usage = common.Zero.String()
		limit.Currency = currency
		s.dao.AddCoinSellingUserLimit(&limit)
	}

	return
}

func (s CoinService) ResetCoinUserLimit() (ce SimpleContextError) {
	userLimitTO := s.dao.ListCoinUserLimit()
	if ce.FeedDaoTransfer(api_error.GetDataFailed, userLimitTO) {
		return
	}
	for _, item := range userLimitTO.Objects {
		userLimit := item.(bean.CoinUserLimit)
		s.dao.ResetCoinUserLimit(userLimit.UID)
	}

	return
}

func (s CoinService) ResetCoinSellingUserLimit() (ce SimpleContextError) {
	userLimitTO := s.dao.ListCoinSellingUserLimit()
	if ce.FeedDaoTransfer(api_error.GetDataFailed, userLimitTO) {
		return
	}
	for _, item := range userLimitTO.Objects {
		userLimit := item.(bean.CoinUserLimit)
		s.dao.ResetCoinSellingUserLimit(userLimit.UID)
	}

	return
}

func (s CoinService) OrderCallNotification() (hasOrder bool, ce SimpleContextError) {
	hasOrder = false
	coinOrderTO := dao.CoinDaoInst.ListCoinOrders("pending", "cod", "", 10, nil)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	if len(coinOrderTO.Objects) > 0 {
		hasOrder = true
	}

	if !hasOrder {
		coinOrderTO := dao.CoinDaoInst.ListCoinOrders("fiat_transferring", "bank", "", 10, nil)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
			return
		}
		if len(coinOrderTO.Objects) > 0 {
			hasOrder = true
		}
	}

	if !hasOrder {
		coinOrderTO := dao.CoinDaoInst.ListCoinSellingOrders("fiat_transferring", "", 10, nil)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
			return
		}
		if len(coinOrderTO.Objects) > 0 {
			hasOrder = true
		}
	}

	if hasOrder {
		s.NotifyPendingOrder()
	}

	return
}

func (s CoinService) SyncCoinOrderToSolr(id string) (coinOrder bean.CoinOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	coinOrder = coinOrderTO.Object.(bean.CoinOrder)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(coinOrder))

	return
}

func (s CoinService) SyncCoinSellingOrderToSolr(id string) (coinOrder bean.CoinSellingOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinSellingOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	coinOrder = coinOrderTO.Object.(bean.CoinSellingOrder)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(coinOrder))

	return
}

func (s CoinService) GenerateAddress(currency string) (address string, ce SimpleContextError) {
	addressResponse, errCoinBase := coinbase_service.GenerateAddress(currency)
	if errCoinBase != nil {
		ce.SetError(api_error.ExternalApiFailed, errCoinBase)
		return
	}
	address = addressResponse.Data.Address
	err := s.dao.AddCoinGenerateAddress(&bean.CoinGeneratedAddress{
		Currency: currency,
		Address:  address,
	})
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	return
}

func (s CoinService) TrackAddressDeposit() (ce SimpleContextError) {
	for _, code := range []string{bean.BTC.Code, bean.ETH.Code, bean.BCH.Code} {
		addressTrackingTO := s.dao.ListCoinAddressTracking(code)
		if !addressTrackingTO.HasError() {
			for _, item := range addressTrackingTO.Objects {
				addressTracking := item.(bean.CoinAddressTracking)
				var provider string

				coinOrderTO := s.dao.GetCoinSellingOrder(addressTracking.Order)
				if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
					return
				}
				order := coinOrderTO.Object.(bean.CoinSellingOrder)

				amount := common.Zero
				isValid := false
				if addressTracking.Currency == bean.ETH.Code {
					provider = bean.CRYPTO_WALLET_NETWORK
				} else if addressTracking.Currency == bean.BTC.Code {
					provider = bean.CRYPTO_WALLET_NETWORK

					tmpTxHash, balance, txCount, bitPayErr := bitpay_service.GetBTCAddress(addressTracking.Address)
					if bitPayErr == nil {
						if txCount == 1 {
							order.TxHash = tmpTxHash
							amount = decimal.NewFromFloat(balance)
							isValid = true
						}
					}
				} else if addressTracking.Currency == bean.BCH.Code {
					provider = bean.CRYPTO_WALLET_NETWORK

					tmpTxHash, balance, txCount, bitPayErr := bitpay_service.GetBTCAddress(addressTracking.Address)
					if bitPayErr == nil {
						if txCount == 1 {
							order.TxHash = tmpTxHash
							amount = decimal.NewFromFloat(balance)
							isValid = true
						}
					}
				}

				if isValid {
					_, logErr := s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
						Provider:         provider,
						ProviderResponse: "",
						DataType:         bean.OFFER_ADDRESS_MAP_COIN_SELLING_ORDER,
						DataRef:          addressTracking.Order,
						UID:              order.UID,
						Description:      "",
						Amount:           amount.String(),
						Currency:         addressTracking.Currency,
						ExternalId:       "",
						TxHash:           order.TxHash,
					})
					order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
					fmt.Println(logErr)

					err := s.dao.UpdateCoinSellingOrder(&order)
					if ce.SetError(api_error.AddDataFailed, err) {
						return
					}

					s.dao.UpdateNotificationCoinSellingOrder(order)
					solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))
				}
			}
		}
	}

	return
}

func (s CoinService) expireOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	order = s.cancelCoinOrder(orderId, bean.COIN_ORDER_STATUS_EXPIRED, &ce)
	return
}

func (s CoinService) expireSellingOrder(orderId string) (order bean.CoinSellingOrder, ce SimpleContextError) {
	order = s.cancelCoinSellingOrder(orderId, bean.COIN_ORDER_STATUS_EXPIRED, &ce)
	return
}

func (s CoinService) setupCoinOrder(order *bean.CoinOrder, coinQuote bean.CoinQuote) {
	order.Price = coinQuote.Price
	order.RawFiatAmount = coinQuote.RawFiatAmount
	order.Status = bean.COIN_ORDER_STATUS_PENDING
	order.Fee = coinQuote.Fee
	order.FeePercentage = coinQuote.FeePercentage
	order.Price = coinQuote.Price

	if order.Type == bean.COIN_ORDER_TYPE_COD {
		order.Fee = coinQuote.FeeCOD
		order.FeePercentage = coinQuote.FeePercentageCOD
	} else {
		order.Duration = int64(30) // 30 minutes
	}
}

func (s CoinService) setupCoinSellingOrder(order *bean.CoinSellingOrder, coinQuote bean.CoinQuote) {
	order.Price = coinQuote.Price
	order.RawFiatAmount = coinQuote.RawFiatAmount
	order.Status = bean.COIN_ORDER_STATUS_PENDING
	order.Fee = coinQuote.Fee
	order.FeePercentage = coinQuote.FeePercentage
	order.Price = coinQuote.Price

	if order.Currency == bean.ETH.Code {
		order.Duration = int64(30) // 30 minutes
	} else {
		order.Duration = int64(90) // 90 minutes
	}
}

func (s CoinService) cancelCoinOrder(orderId string, status string, ce *SimpleContextError) (order bean.CoinOrder) {
	coinOrderTO := s.dao.GetCoinOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinOrder)

	if order.Status != bean.COIN_ORDER_STATUS_PENDING && (status == bean.COIN_ORDER_STATUS_CANCELLED || status == bean.COIN_ORDER_STATUS_EXPIRED) {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}
	if order.Status != bean.COIN_ORDER_STATUS_PROCESSING && status == bean.COIN_ORDER_STATUS_REJECTED {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	order.Status = status
	err := s.dao.CancelCoinOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) cancelCoinSellingOrder(orderId string, status string, ce *SimpleContextError) (order bean.CoinSellingOrder) {
	coinOrderTO := s.dao.GetCoinSellingOrder(orderId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	order = coinOrderTO.Object.(bean.CoinSellingOrder)

	if order.Status != bean.COIN_ORDER_STATUS_PENDING && (status == bean.COIN_ORDER_STATUS_CANCELLED || status == bean.COIN_ORDER_STATUS_EXPIRED) {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}
	if order.Status != bean.COIN_ORDER_STATUS_PROCESSING && status == bean.COIN_ORDER_STATUS_REJECTED {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	order.Status = status
	err := s.dao.CancelCoinSellingOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	if order.Currency != bean.ETH.Code {
		s.dao.RemoveCoinAddressTracking(order.Currency, order.Address)
	}

	s.dao.UpdateNotificationCoinSellingOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinSellingOrder(order))

	return
}

func (s CoinService) NotifyNewCoinOrder(order bean.CoinOrder) error {
	host := os.Getenv("FRONTEND_HOST")
	content := fmt.Sprintf("[%s] [ORDER] You have new BUY order, please check ref code: %s", strings.ToUpper(order.Type), order.RefCode)
	// _, err := twilio_service.SendSMS(os.Getenv("COIN_ORDER_TO_NUMBER"), smsBody)
	if os.Getenv("ENVIRONMENT") == "dev" {
		content = "TEST -- " + content
	}

	body := `Hi,

There is new BUY order please check the following link:
%s/internal-admin-dashboard?tab=%s&refCode=%s
`

	orderType := "Bank"
	if order.Type == bean.COIN_ORDER_TYPE_COD {
		orderType = "Cod"
	}
	body = fmt.Sprintf(body, host, fmt.Sprintf("buyCoin%s", orderType), order.RefCode)

	slack_integration.SendSlack(content)
	err := email.SendEmail("System", "dojo@ninja.org", "Admin", os.Getenv("COIN_ORDER_TO_EMAIL"), content, body)

	return err
}

func (s CoinService) NotifyNewCoinSellingOrder(order bean.CoinSellingOrder) error {
	host := os.Getenv("FRONTEND_HOST")
	content := fmt.Sprintf("[%s] [ORDER] You have new SELL order, please check ref code: %s", strings.ToUpper(order.Type), order.RefCode)
	// _, err := twilio_service.SendSMS(os.Getenv("COIN_ORDER_TO_NUMBER"), smsBody)
	if os.Getenv("ENVIRONMENT") == "dev" {
		content = "TEST -- " + content
	}

	body := `Hi,

There is new SELL order please check the following link:
%s/internal-admin-dashboard?tab=%s&refCode=%s
`

	body = fmt.Sprintf(body, host, "sellCoinBank", order.RefCode)

	slack_integration.SendSlack(content)
	err := email.SendEmail("System", "dojo@ninja.org", "Admin", os.Getenv("COIN_ORDER_TO_EMAIL"), content, body)

	return err
}

func (s CoinService) NotifyPendingOrder() error {
	content := fmt.Sprintf("You have pending BUY/SELL orders, please process ASAP")
	if os.Getenv("ENVIRONMENT") == "dev" {
		content = "TEST -- " + content
	}

	slack_integration.SendSlack(content)
	_, err := twilio_service.SendSMS(os.Getenv("COIN_ORDER_TO_NUMBER"), content)
	err = email.SendEmail("System", "dojo@ninja.org", "Admin", os.Getenv("COIN_ORDER_TO_EMAIL"), content, "")

	return err
}

func (s CoinService) InitBank() {
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "VP Bank - NH Viet Nam Thinh Vuong",
		Bank:    "VP Bank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "ACB - NH A Chau",
		Bank:    "ACB",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Agribank - NH Nong Nghiep va Phat Trien Nong Thon Viet Nam",
		Bank:    "Agribank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "ANZ Viet Nam",
		Bank:    "ANZ Viet Nam",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "BIDV - NH Dau Tu va Phat Trien Viet Nam",
		Bank:    "BIDV",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "CitiBank Viet Nam",
		Bank:    "CitiBank Viet Nam",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "DongA Bank - NH Dong A",
		Bank:    "DongA Bank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Eximbank - NH Xuat Nhap Khau",
		Bank:    "Eximbank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "HDBank - NH Phat Trien TP HCM",
		Bank:    "HDBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "HSBC Viet Nam",
		Bank:    "HSBC Viet Nam",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "MaritimeBank - NH Hang Hai",
		Bank:    "MaritimeBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "MBBank - NH Quan Doi",
		Bank:    "MBBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "OCB - NH Phuong Dong",
		Bank:    "OCB",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "SacomBank – NH Sai Gon Thuong Tin",
		Bank:    "SacomBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "SCB - NH Sai Gon",
		Bank:    "SCB",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "SHBank - NH Sai Gon Ha Noi",
		Bank:    "SHBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Techcombank - NH Ky Thuong",
		Bank:    "Techcombank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "TPBank - NH Tien Phong",
		Bank:    "TPBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "VIB - NH Quoc Te",
		Bank:    "VIB",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Vietcombank - NH Ngoai Thuong VN",
		Bank:    "Vietcombank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "VietinBank - NH Cong Thuong VN",
		Bank:    "VietinBank",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Bangkok Bank – VN Branch",
		Bank:    "Bangkok Bank – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Bank of China – VN Branch",
		Bank:    "Bank of China – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Bank of Communication – VN Branch",
		Bank:    "Bank of Communication – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Bank of Tokyo Mitsubishi – VN Branch",
		Bank:    "Bank of Tokyo Mitsubishi – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Cathay Bank – VN Branch",
		Bank:    "Cathay Bank – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Chifon Dai Loan – VN Branch",
		Bank:    "Chifon Dai Loan – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "China Construction Bank  – VN Branch",
		Bank:    "China Construction Bank  – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Chinatrust Commercial Bank – VN Branch",
		Bank:    "Chinatrust Commercial Bank – VN Branch",
		Country: "VN",
	})

	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Commonwealth Bank Viet Nam",
		Bank:    "Commonwealth Bank Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Credit Agricole CIB – VN Branch",
		Bank:    "Credit Agricole CIB – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "DBS Bank Ltd – VN Branch",
		Bank:    "DBS Bank Ltd – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Deutsche Bank Viet Nam",
		Bank:    "Deutsche Bank Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Far East National Bank – VN Branch",
		Bank:    "Far East National Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "First Commercial Bank – VN Branch",
		Bank:    "First Commercial Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Hong Leong Bank  – VN Branch",
		Bank:    "Hong Leong Bank  – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Hua Nan Commercial Bank – VN Branch",
		Bank:    "Hua Nan Commercial Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Indovina Bank",
		Bank:    "Indovina Bank",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Industrial Bank Of Korea – VN Branch",
		Bank:    "Industrial Bank Of Korea – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "JP Morgan Chase Bank – VN Branch",
		Bank:    "JP Morgan Chase Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Kho bac nha nuoc Viet Nam",
		Bank:    "Kho bac nha nuoc Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Kookmin Bank – VN Branch",
		Bank:    "Kookmin Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Korea Exchange Bank – VN Branch",
		Bank:    "Korea Exchange Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Malayan Banking Berhad – VN Branch",
		Bank:    "Malayan Banking Berhad – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "May Bank – VN Branch",
		Bank:    "May Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Mega ICBC – VN Branch",
		Bank:    "Mega ICBC – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Mizuho Corporate Bank, Ltd – VN Branch",
		Bank:    "Mizuho Corporate Bank, Ltd – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Natexis Banques Populaires – VN Branch",
		Bank:    "Natexis Banques Populaires – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH An Binh (AnbinhBank)",
		Bank:    "NH An Binh (AnbinhBank)",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Bac A",
		Bank:    "NH Bac A",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Ban Viet (VietCapital)",
		Bank:    "NH Ban Viet (VietCapital)",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Bao Viet",
		Bank:    "NH Bao Viet",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Buu dien Lien Viet (PostBank)",
		Bank:    "NH Buu dien Lien Viet (PostBank)",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Chinh sach Xa Hoi",
		Bank:    "NH Chinh sach Xa Hoi",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Cong Thuong Trung Quoc – VN Branch",
		Bank:    "NH Cong Thuong Trung Quoc – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dai A",
		Bank:    "NH Dai A",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dai Chung",
		Bank:    "NH Dai Chung",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dai Duong (Ocean Bank)",
		Bank:    "NH Dai Duong (Ocean Bank)",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dau khi toan cau",
		Bank:    "NH Dau khi toan cau",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dau Tu va PT Campuchia – VN Branch",
		Bank:    "NH Dau Tu va PT Campuchia – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Dong Nam A",
		Bank:    "NH Dong Nam A",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Hop Tac Xa VN",
		Bank:    "NH Hop Tac Xa VN",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH lien doanh Viet Lao – VN Branch",
		Bank:    "NH lien doanh Viet Lao – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Lien doanh Viet Nga – VN Branch",
		Bank:    "NH Lien doanh Viet Nga – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Lien doanh Viet Thai – VN Branch",
		Bank:    "NH Lien doanh Viet Thai – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Nam A",
		Bank:    "NH Nam A",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Nha Nuoc Viet Nam",
		Bank:    "NH Nha Nuoc Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Phat trien Viet Nam",
		Bank:    "NH Phat trien Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Quoc Dan",
		Bank:    "NH Quoc Dan",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Sai Gon Cong Thuong",
		Bank:    "NH Sai Gon Cong Thuong",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH VID Public",
		Bank:    "NH VID Public",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Viet A",
		Bank:    "NH Viet A",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Viet Nam Thuong Tin",
		Bank:    "NH Viet Nam Thuong Tin",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Xang dau Petrolimex",
		Bank:    "NH Xang dau Petrolimex",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Xay Dung VN",
		Bank:    "NH Xay Dung VN",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Oversea Chinese Banking Corp – VN Branch",
		Bank:    "Oversea Chinese Banking Corp – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "PARIBAS – VN Branch",
		Bank:    "PARIBAS – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Shanghai Commercial & Savings – VN Branch",
		Bank:    "Shanghai Commercial & Savings – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Shinhan Vietnam Bank Limited",
		Bank:    "Shinhan Vietnam Bank Limited",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Standard Chartered Viet Nam",
		Bank:    "Standard Chartered Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Sumitomo Mitsui Bank – VN Branch",
		Bank:    "Sumitomo Mitsui Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Taipei Fubon Bank – VN Branch",
		Bank:    "Taipei Fubon Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "United Overseas Bank – VN Branch",
		Bank:    "United Overseas Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Woori Bank – VN Branch",
		Bank:    "Woori Bank – VN Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Bank Of India HCM Branch",
		Bank:    "Bank Of India HCM Branch",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH Kien Long",
		Bank:    "NH Kien Long",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "NH TNHH MTV CIMB Viet Nam",
		Bank:    "NH TNHH MTV CIMB Viet Nam",
		Country: "VN",
	})
	s.dao.AddCoinBank(&bean.CoinBank{
		Name:    "Nong Hyup Bank",
		Bank:    "Nong Hyup Bank",
		Country: "VN",
	})
}
