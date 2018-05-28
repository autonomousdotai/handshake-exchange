package api

import (
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/coinbase_service"
	"github.com/autonomousdotai/handshake-exchange/integration/openexchangerates_service"
	"github.com/autonomousdotai/handshake-exchange/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MiscApi struct {
}

// CRON JOB
func (api MiscApi) UpdateCurrencyRates(context *gin.Context) {
	rates, err := openexchangerates_service.GetExchangeRate()
	if api_error.PropagateErrorAndAbort(context, api_error.ExternalApiFailed, err) != nil {
		return
	}

	err = dao.MiscDaoInst.UpdateCurrencyRate(rates)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
func (api MiscApi) UpdateCryptoRates(context *gin.Context) {
	//rates, err := coinapi_service.GetExchangeRate()
	//if api_error.PropagateErrorAndAbort(context, api_error.ExternalApiFailed, err) != nil {
	//	return
	//}

	allRates := make([]bean.CryptoRate, 0)
	for _, currency := range []string{bean.BTC.Code, bean.ETH.Code, bean.LTC.Code, bean.BCH.Code} {
		rates := make([]bean.CryptoRate, 0)
		resp, _ := coinbase_service.GetBuyPrice(currency)
		buy, _ := decimal.NewFromString(resp.Amount)
		buyFloat, _ := buy.Float64()
		rate := bean.CryptoRate{
			From:     currency,
			To:       bean.USD.Code,
			Buy:      buyFloat,
			Sell:     0,
			Exchange: bean.INSTANT_OFFER_PROVIDER_COINBASE,
		}
		rates = append(rates, rate)
		allRates = append(allRates, rate)

		err := dao.MiscDaoInst.UpdateCryptoRates(map[string][]bean.CryptoRate{currency: rates})
		if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
			return
		}
	}

	bean.SuccessResponse(context, allRates)
}

func (api MiscApi) UpdateSystemFee(context *gin.Context) {
	systemFees, err := dao.MiscDaoInst.LoadSystemFeeToCache()
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, systemFees)
}

func (api MiscApi) GetCurrencyRate(context *gin.Context) {
	currency := context.Param("currency")
	to := dao.MiscDaoInst.GetCurrencyRateFromCache(currency[:3], currency[3:])
	if to.ContextValidate(context) {
		return
	}
	rate := to.Object.(bean.CurrencyRate)

	bean.SuccessResponse(context, rate)
}

func (api MiscApi) ListCryptoRates(context *gin.Context) {
	currency := context.Param("currency")
	to := dao.MiscDaoInst.GetCryptoRatesFromCache(currency)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, to.Objects)
}

func (api MiscApi) GetSystemFee(context *gin.Context) {
	feeKey := context.Param("feeKey")
	to := dao.MiscDaoInst.GetSystemFeeFromCache(feeKey)
	if to.ContextValidate(context) {
		return
	}
	systemFee := to.Object.(bean.SystemFee)

	bean.SuccessResponse(context, systemFee)
}

func (api MiscApi) GetCryptoRate(context *gin.Context) {
	currency := context.Param("currency")
	resp, err := coinbase_service.GetBuyPrice(currency)
	if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, resp)
}

func (api MiscApi) GetCryptoQuote(context *gin.Context) {
	type quoteStruct struct {
		Type         string
		Amount       string
		Currency     string
		FiatCurrency string
		FiatAmount   string
	}

	quoteType := context.DefaultQuery("type", "")
	amountStr := context.DefaultQuery("amount", "")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "")

	quote := quoteStruct{
		Type:         quoteType,
		Amount:       amountStr,
		Currency:     currency,
		FiatCurrency: fiatCurrency,
	}

	amount, numberErr := decimal.NewFromString(amountStr)
	if api_error.PropagateErrorAndAbort(context, api_error.InvalidQueryParam, numberErr) != nil {
		return
	}
	to := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatCurrency)
	if to.ContextValidate(context) {
		return
	}
	rate := to.Object.(bean.CurrencyRate)
	fiatAmount := amount.Mul(decimal.NewFromFloat(rate.Rate))

	if quoteType == "buy" {
		resp, err := coinbase_service.GetBuyPrice(currency)
		if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
			return
		}
		price, _ := decimal.NewFromString(resp.Amount)
		fiatAmount = fiatAmount.Mul(price)
	} else if quoteType == "sell" {
		resp, err := coinbase_service.GetSellPrice(currency)
		if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
			return
		}
		price, _ := decimal.NewFromString(resp.Amount)
		fiatAmount = fiatAmount.Mul(price)
	} else {
		if api_error.AbortWithValidateErrorSimple(context, api_error.InvalidQueryParam) != nil {
			return
		}
	}

	quote.FiatAmount = fiatAmount.Round(2).String()

	bean.SuccessResponse(context, quote)
}

// CRON JOB
func (api MiscApi) FinishInstantOffers(context *gin.Context) {
	offers, ce := service.CreditCardServiceInst.FinishInstantOffers()
	if ce.ContextValidate(context) {
		return
	}
	for range offers {
		//profileTO := dao.UserDaoInst.GetProfile(offer.UID)
		//if !profileTO.HasError() {
		//	profile := profileTO.Object.(bean.Profile)
		//	// TODO Do some notification
		//}
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
//func (api MiscApi) ExpireOfferHandshakes(context *gin.Context) {
//	err := dao.OfferDaoInst.UpdateExpiredHandshake()
//	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
//		return
//	}
//
//	bean.SuccessResponse(context, true)
//}

func (api MiscApi) UpdateCCLimits(context *gin.Context) {
	// country := context.Param("country")
	objs, err := dao.MiscDaoInst.LoadCCLimitToCache()
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, objs)
}

func (api MiscApi) GetCCLimits(context *gin.Context) {
	// country := context.Param("country")
	to := dao.MiscDaoInst.GetCCLimitFromCache()
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, to.Objects)
}

// CRON JOB
func (api MiscApi) UpdateUserCCLimitTracks(context *gin.Context) {
	// country := context.Param("country")
	ce := service.UserServiceInst.UpdateUserCCLimitTracks()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}
