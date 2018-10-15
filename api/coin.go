package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/service"
	"strconv"
)

type CoinApi struct {
}

func (api CoinApi) GetQuote(context *gin.Context) {
	amount := context.DefaultQuery("amount", "1")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "USD")
	check := context.DefaultQuery("check", "")

	coinQuote, ce := service.CoinServiceInst.GetCoinQuote(amount, currency, fiatCurrency, check)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, coinQuote)
}

func (api CoinApi) GetQuoteReverse(context *gin.Context) {
	amount := context.DefaultQuery("fiat_amount", "")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "USD")
	check := context.DefaultQuery("check", "")
	orderType := context.DefaultQuery("type", bean.COIN_ORDER_TYPE_BANK)

	coinQuote, ce := service.CoinServiceInst.GetCoinQuoteReverse(amount, currency, fiatCurrency, orderType, check)
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
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.CoinOrder
	if common.ValidateBody(context, &body) != nil {
		return
	}

	id, _ := strconv.Atoi(chainId)
	body.ChainId = int64(id)
	body.Language = language
	body.FCM = fcm

	order, ce := service.CoinServiceInst.AddOrder(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CoinApi) ListCoinOrders(context *gin.Context) {
	status := context.DefaultQuery("status", "")
	orderType := context.DefaultQuery("type", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.CoinDaoInst.ListCoinOrders(status, orderType, limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page)
}

func (api CoinApi) FinishCoinOrder(context *gin.Context) {
	id := context.Param("id")
	// currency := context.Param("currency")
	// amount := context.Param("amount")
	order, overSpent, ce := service.CoinServiceInst.FinishOrder(id, "", "")
	if ce.ContextValidate(context) {
		return
	}
	fmt.Println(overSpent)

	bean.SuccessResponse(context, order)
}

func (api CoinApi) CancelCoinOrder(context *gin.Context) {
	id := context.Param("id")

	order, ce := service.CoinServiceInst.CancelOrder(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
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
	id := context.Param("id")

	order, ce := service.CoinServiceInst.UpdateOrder(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CoinApi) RejectCoinOrder(context *gin.Context) {
	id := context.Param("id")

	order, ce := service.CoinServiceInst.RejectOrder(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CoinApi) RemoveExpiredOrder(context *gin.Context) {
	ce := service.CoinServiceInst.RemoveExpiredOrder()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}
