package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/service"
)

type CoinApi struct {
}

func (api CoinApi) GetQuote(context *gin.Context) {
	amount := context.DefaultQuery("amount", "1")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "USD")

	coinQuote, ce := service.CoinServiceInst.GetCoinQuote(amount, currency, fiatCurrency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, coinQuote)
}

func (api CoinApi) ListCoinCenter(context *gin.Context) {
	country := context.Param("country")
	coinCenters, ce := service.CoinServiceInst.ListCoinCenter(country)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, coinCenters)
}

func (api CoinApi) CoinOrderType(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) CoinOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) ListCoinOrders(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) FinishCoinOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) RejectCoinOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) UpdateCashOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api CoinApi) PickCoinOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}
