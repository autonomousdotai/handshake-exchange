package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
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

func (api CoinApi) CoinOrder(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CoinOrder
	if common.ValidateBody(context, &body) != nil {
		return
	}
	order, ce := service.CoinServiceInst.AddOrder(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
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

func (api CoinApi) UpdateCoinOrder(context *gin.Context) {
	id := context.Param("id")

	var body bean.CoinOrderUpdateInput
	if common.ValidateBody(context, &body) != nil {
		return
	}
	order, ce := service.CoinServiceInst.UpdateOrderReceipt(id, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CoinApi) PickCoinOrder(context *gin.Context) {
	bean.SuccessResponse(context, true)
}
