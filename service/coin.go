package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/shopspring/decimal"
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
	price := amount.Mul(decimal.NewFromFloat(cryptoRate.Buy).Round(2))

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
