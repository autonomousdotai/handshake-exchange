package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/bitstamp_service"
	"github.com/ninjadotorg/handshake-exchange/integration/crypto_service"
	"github.com/ninjadotorg/handshake-exchange/integration/slack_integration"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
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

func (s CoinService) GetCoinQuote(userId string, amountStr string, currency string, fiatLocalCurrency string, level string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
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

	if userLimit.LessThan(userUsage.Add(localPrice)) {
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

func (s CoinService) GetCoinSellingQuote(userId string, amountStr string, currency string, fiatLocalCurrency string, level string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
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

	if userLimit.LessThan(userUsage.Add(localPrice)) {
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

func (s CoinService) AddOrder(userId string, orderBody bean.CoinOrder) (order bean.CoinOrder, ce SimpleContextError) {
	orderTest, testOfferCE := s.GetCoinQuote(userId, orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency, orderBody.Level, "1")
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
	order.Status = bean.COIN_ORDER_STATUS_PROCESSING

	err := s.dao.UpdateCoinOrder(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	s.dao.UpdateNotificationCoinOrder(order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) CancelOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	order = s.cancelCoinOrder(orderId, bean.COIN_ORDER_STATUS_CANCELLED, &ce)
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

func (s CoinService) FinishCoinOrderPendingTransfer(ref string) (order bean.CoinOrder, ce SimpleContextError) {
	cashOrderTO := s.dao.GetCoinOrderByPath(ref)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order = cashOrderTO.Object.(bean.CoinOrder)
	order.Status = bean.COIN_ORDER_STATUS_SUCCESS

	err := s.dao.UpdateCoinOrderReceipt(&order)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

	return
}

func (s CoinService) AddCoinReview(review bean.CoinReview) (ce SimpleContextError) {
	cashOrderTO := s.dao.GetCoinOrder(review.Order)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashOrderTO) {
		return
	}
	order := cashOrderTO.Object.(bean.CoinOrder)
	if order.Reviewed {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	err := s.dao.AddCoinReview(&review)
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}
	order.Reviewed = true

	s.dao.UpdateCoinOrderReview(&order)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(order))

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

func (s CoinService) SyncCoinOrderToSolr(id string) (coinOrder bean.CoinOrder, ce SimpleContextError) {
	coinOrderTO := s.dao.GetCoinOrder(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, coinOrderTO) {
		return
	}
	coinOrder = coinOrderTO.Object.(bean.CoinOrder)

	s.dao.UpdateNotificationCoinOrder(coinOrder)
	solr_service.UpdateObject(bean.NewSolrFromCoinOrder(coinOrder))

	return
}

func (s CoinService) expireOrder(orderId string) (order bean.CoinOrder, ce SimpleContextError) {
	order = s.cancelCoinOrder(orderId, bean.COIN_ORDER_STATUS_EXPIRED, &ce)
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

func (s CoinService) NotifyNewCoinOrder(order bean.CoinOrder) error {
	// os.Getenv("FRONTEND_HOST")
	content := fmt.Sprintf("[%s] [ORDER] You have new order, please check ref code: %s", strings.ToUpper(order.Type), order.RefCode)
	// _, err := twilio_service.SendSMS(os.Getenv("COIN_ORDER_TO_NUMBER"), smsBody)
	if os.Getenv("ENVIRONMENT") == "dev" {
		content = "TEST -- " + content
	}

	slack_integration.SendSlack(content)
	err := email.SendEmail("System", "dojo@ninja.org", "Admin", os.Getenv("COIN_ORDER_TO_EMAIL"), content, " ")

	return err
}
