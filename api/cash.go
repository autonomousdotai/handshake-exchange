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

type CashApi struct {
}

func (api CashApi) CashStore(context *gin.Context) {
	userId := common.GetUserId(context)

	cash, ce := service.CashServiceInst.GetCashStore(userId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, cash)
}

func (api CashApi) CashStoreCreate(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.CashStore
	if common.ValidateBody(context, &body) != nil {
		return
	}

	id, _ := strconv.Atoi(chainId)
	body.ChainId = int64(id)
	body.Language = language
	body.FCM = fcm

	cash, ce := service.CashServiceInst.AddCashStore(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, cash)
}

func (api CashApi) CashStoreUpdate(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CashStore
	if common.ValidateBody(context, &body) != nil {
		return
	}
	cash, ce := service.CashServiceInst.UpdateCashStore(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, cash)
}

func (api CashApi) CashStorePrice(context *gin.Context) {
	amount := context.DefaultQuery("amount", "")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "")

	order, ce := service.CashServiceInst.GetProposeCashOrder(amount, currency, fiatCurrency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CashApi) ListCashStoreOrders(context *gin.Context) {
	status := context.DefaultQuery("status", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.CashDaoInst.ListCashOrders(status, limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page, 0)
}

func (api CashApi) CashStoreOrder(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CashOrder
	if common.ValidateBody(context, &body) != nil {
		return
	}
	order, ce := service.CashServiceInst.AddOrder(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CashApi) FinishCashOrder(context *gin.Context) {
	id := context.Param("id")
	amount := context.Param("amount")
	order, overSpent, ce := service.CashServiceInst.FinishOrder(id, amount, "USD")
	if ce.ContextValidate(context) {
		return
	}
	fmt.Println(overSpent)

	bean.SuccessResponse(context, order)
}

func (api CashApi) RejectCashOrder(context *gin.Context) {
	id := context.Param("id")

	order, ce := service.CashServiceInst.RejectOrder(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CashApi) UpdateCashOrder(context *gin.Context) {
	id := context.Param("id")

	var body bean.CashOrderUpdateInput
	if common.ValidateBody(context, &body) != nil {
		return
	}
	order, ce := service.CashServiceInst.UpdateOrderReceipt(id, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CashApi) ListCashCenter(context *gin.Context) {
	country := context.Param("country")
	cashCenters, ce := service.CashServiceInst.ListCashCenter(country)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, cashCenters)
}

func (api CashApi) ListProposeQuotes(context *gin.Context) {
	fiatCurrency := context.DefaultQuery("fiat_currency", "")

	type quoteStruct struct {
		Currency          string
		FiatCurrency      string
		Price             string
		FiatLocalCurrency string
		LocalPrice        string
	}

	result := make([]quoteStruct, 0)
	btcOrder, ce := service.CashServiceInst.GetProposeCashAmount(bean.BTC.Code, fiatCurrency)
	if !ce.HasError() {
		result = append(result, quoteStruct{
			Currency:          btcOrder.Currency,
			FiatCurrency:      btcOrder.FiatCurrency,
			Price:             btcOrder.FiatAmount,
			FiatLocalCurrency: btcOrder.FiatLocalCurrency,
			LocalPrice:        btcOrder.FiatLocalAmount,
		})
	}
	bchOrder, ce := service.CashServiceInst.GetProposeCashAmount(bean.ETH.Code, fiatCurrency)
	if !ce.HasError() {
		result = append(result, quoteStruct{
			Currency:          bchOrder.Currency,
			FiatCurrency:      bchOrder.FiatCurrency,
			Price:             bchOrder.FiatAmount,
			FiatLocalCurrency: bchOrder.FiatLocalCurrency,
			LocalPrice:        bchOrder.FiatLocalAmount,
		})
	}
	ethOrder, ce := service.CashServiceInst.GetProposeCashAmount(bean.BCH.Code, fiatCurrency)
	if !ce.HasError() {
		result = append(result, quoteStruct{
			Currency:          ethOrder.Currency,
			FiatCurrency:      ethOrder.FiatCurrency,
			Price:             ethOrder.FiatAmount,
			FiatLocalCurrency: ethOrder.FiatLocalCurrency,
			LocalPrice:        ethOrder.FiatLocalAmount,
		})
	}

	bean.SuccessResponse(context, result)
}
