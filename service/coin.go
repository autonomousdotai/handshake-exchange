package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/crypto_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
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

func (s CoinService) GetCoinQuote(amountStr string, currency string, fiatLocalCurrency string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
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

func (s CoinService) GetCoinQuoteReverse(fiatLocalAmountStr string, currency string, fiatLocalCurrency string,
	orderType string, check string) (coinQuote bean.CoinQuote, ce SimpleContextError) {
	fiatLocalAmount := common.StringToDecimal(fiatLocalAmountStr)

	rateTO := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatLocalCurrency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, rateTO) {
		return
	}
	rate := rateTO.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)

	fiatAmount := fiatLocalAmount.Div(rateNumber)

	cryptoRateTO := dao.MiscDaoInst.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_COINBASE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)
	cryptoPrice := decimal.NewFromFloat(cryptoRate.Buy).Round(2)

	configTO := s.miscDao.GetSystemConfigFromCache(bean.COIN_ORDER_LIMIT)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, configTO) {
		return
	}
	systemConfig := configTO.Object.(bean.SystemConfig)
	// There is no free start on
	limit := common.StringToDecimal(systemConfig.Value)

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
	if fiatAmount.GreaterThan(limit) {
		bankFeeTO = dao.MiscDaoInst.GetSystemFeeFromCache(bean.FEE_COIN_ORDER_BANK_HIGH)
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
		amount := bankValue.Div(cryptoPrice)
		coinQuote.Amount = amount.String()
	} else {
		amount := codValue.Div(cryptoPrice)
		coinQuote.Amount = amount.String()
	}

	coinQuote.Limit = limit.String()

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
	orderTest, testOfferCE := s.GetCoinQuote(orderBody.Amount, orderBody.Currency, orderBody.FiatLocalCurrency, "1")
	if ce.FeedContextError(api_error.GetDataFailed, testOfferCE) {
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

	fmt.Println(orderBody.Currency)
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

	if order.Status != bean.COIN_ORDER_STATUS_FIAT_TRANSFERRING && order.Status != bean.COIN_ORDER_STATUS_PROCESSING {
		ce.SetStatusKey(api_error.CoinOrderStatusInvalid)
		return
	}

	if os.Getenv("ENVIRONMENT") == "dev" {
		if order.Currency == bean.ETH.Code {
			txHash, onChainErr := crypto_service.SendTransaction(order.Address, order.Amount, order.Currency)

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
		coinbaseTx, errWithdraw := coinbase_service.SendTransaction(order.Address, order.Amount, order.Currency,
			fmt.Sprintf("Withdraw tx = %s", order.Id), order.Id)
		if errWithdraw == nil {
			order.ProviderWithdrawData = coinbaseTx.Id
			order.Status = bean.COIN_ORDER_STATUS_TRANSFERRING
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
		provider := bean.BTC_WALLET_BITSTAMP
		s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
			Provider:         provider,
			ProviderResponse: order.ProviderWithdrawData,
			DataType:         bean.OFFER_ADDRESS_MAP_COIN_ORDER,
			DataRef:          dao.GetCoinOrderItemPath(order.Id),
			UID:              order.UID,
			Description:      "",
			Amount:           order.Amount,
			Currency:         order.Currency,
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
