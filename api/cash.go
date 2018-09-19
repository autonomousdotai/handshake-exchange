package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
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
	fiat_currency := context.DefaultQuery("fiat_currency", "")

	order, ce := service.CashServiceInst.GetProposeCashOrder(amount, currency, fiat_currency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, order)
}

func (api CashApi) CashStoreOrder(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	var body bean.CashStore
	if common.ValidateBody(context, &body) != nil {
		return
	}
	//withdraw, ce := service.CreditServiceInst.AddCreditWithdraw(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashOrder{})
}

func (api CashApi) CashStoreRemoveOrder(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	var body bean.CashStore
	if common.ValidateBody(context, &body) != nil {
		return
	}
	//withdraw, ce := service.CreditServiceInst.AddCreditWithdraw(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, true)
}
