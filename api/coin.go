package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/service"
	"net/http"
	"strconv"
)

type CoinApi struct {
}

func (api CoinApi) GetQuote(context *gin.Context) {
	userId := common.GetUserId(context)

	amount := context.DefaultQuery("amount", "1")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "USD")
	level := context.DefaultQuery("level", "1")
	check := context.DefaultQuery("check", "")
	direction := context.DefaultQuery("direction", "buy")

	if direction == "sell" {
		coinQuote, ce := service.CoinServiceInst.GetCoinSellingQuote(userId, amount, currency, fiatCurrency, level, check)
		if ce.ContextValidate(context) {
			return
		}
		bean.SuccessResponse(context, coinQuote)
	} else {
		coinQuote, ce := service.CoinServiceInst.GetCoinQuote(userId, amount, currency, fiatCurrency, level, check)
		if ce.ContextValidate(context) {
			return
		}
		bean.SuccessResponse(context, coinQuote)
	}
}

func (api CoinApi) GetQuoteReverse(context *gin.Context) {
	userId := common.GetUserId(context)

	amount := context.DefaultQuery("fiat_amount", "")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "USD")
	level := context.DefaultQuery("level", "1")
	check := context.DefaultQuery("check", "")
	orderType := context.DefaultQuery("type", bean.COIN_ORDER_TYPE_BANK)
	direction := context.DefaultQuery("direction", "buy")

	if direction == "sell" {
		coinQuote, ce := service.CoinServiceInst.GetCoinSellingQuoteReverse(userId, amount, currency, fiatCurrency, orderType, level, check)
		if ce.ContextValidate(context) {
			return
		}
		bean.SuccessResponse(context, coinQuote)
	} else {
		coinQuote, ce := service.CoinServiceInst.GetCoinQuoteReverse(userId, amount, currency, fiatCurrency, orderType, level, check)
		if ce.ContextValidate(context) {
			return
		}
		bean.SuccessResponse(context, coinQuote)
	}
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

func (api CoinApi) CoinSellingOrder(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.CoinSellingOrder
	if common.ValidateBody(context, &body) != nil {
		return
	}

	id, _ := strconv.Atoi(chainId)
	body.ChainId = int64(id)
	body.Language = language
	body.FCM = fcm

	order, ce := service.CoinServiceInst.AddSellingOrder(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CoinApi) ListCoinOrders(context *gin.Context) {
	status := context.DefaultQuery("status", "")
	orderType := context.DefaultQuery("type", "")
	refCode := context.DefaultQuery("ref_code", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.CoinDaoInst.ListCoinOrders(status, orderType, refCode, limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page)
}

func (api CoinApi) ListCoinSellingOrders(context *gin.Context) {
	status := context.DefaultQuery("status", "")
	refCode := context.DefaultQuery("ref_code", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.CoinDaoInst.ListCoinSellingOrders(status, refCode, limit, startAt)
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

func (api CoinApi) CloseCoinSellingOrder(context *gin.Context) {
	id := context.Param("id")
	// currency := context.Param("currency")
	// amount := context.Param("amount")
	order, ce := service.CoinServiceInst.CloseSellingOrder(id)
	if ce.ContextValidate(context) {
		return
	}

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

func (api CoinApi) CancelCoinSellingOrder(context *gin.Context) {
	id := context.Param("id")

	order, ce := service.CoinServiceInst.CancelSellingOrder(id)
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

func (api CoinApi) PickCoinSellingOrder(context *gin.Context) {
	id := context.Param("id")

	order, ce := service.CoinServiceInst.UpdateSellingOrder(id)
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

func (api CoinApi) ResetCoinUserLimit(context *gin.Context) {
	direction := context.DefaultQuery("direction", "buy")
	if direction == "sell" {
		ce := service.CoinServiceInst.ResetCoinSellingUserLimit()
		if ce.ContextValidate(context) {
			return
		}
	} else {
		ce := service.CoinServiceInst.ResetCoinUserLimit()
		if ce.ContextValidate(context) {
			return
		}
	}

	bean.SuccessResponse(context, true)
}

func (api CoinApi) AddReview(context *gin.Context) {
	userId := common.GetUserId(context)
	var body bean.CoinReview
	if common.ValidateBody(context, &body) != nil {
		return
	}
	body.UID = userId
	ce := service.CoinServiceInst.AddCoinReview(body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, body)
}

func (api CoinApi) ListReview(context *gin.Context) {
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.CoinDaoInst.ListReviews(limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page)
}

func (api CoinApi) VoiceOrderNotification(context *gin.Context) {
	type ResponseData struct {
		Voice   string `xml:"voice,attr"`
		Content string `xml:",chardata"`
	}
	type Response struct {
		Say ResponseData
	}

	context.XML(http.StatusOK, Response{Say: ResponseData{
		Voice:   "alice",
		Content: "hello world!",
	}})
}
