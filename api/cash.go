package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
)

type CashApi struct {
}

func (api CashApi) Dashboard(context *gin.Context) {
	//userId := common.GetUserId(context)

	//credit, ce := service.CreditServiceInst.GetCredit(userId)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCredit{})
}

func (api CashApi) Create(context *gin.Context) {
	//userId := common.GetUserId(context)
	//chainId := common.GetChainId(context)
	//language := common.GetLanguage(context)
	//fcm := common.GetFCM(context)
	//
	//var body bean.Credit
	//if common.ValidateBody(context, &body) != nil {
	//	return
	//}
	//
	//body.ChainId = chainId
	//body.Language = language
	//body.FCM = fcm
	//
	//credit, ce := service.CreditServiceInst.AddCredit(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCredit{})
}

func (api CashApi) Deposit(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	var body bean.CashCreditDepositInput
	if common.ValidateBody(context, &body) != nil {
		return
	}
	//
	//tracking, ce := service.CreditServiceInst.AddDeposit(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCreditDeposit{})
}

func (api CashApi) Tracking(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	var body bean.CashCreditOnChainActionTrackingInput
	if common.ValidateBody(context, &body) != nil {
		return
	}
	//
	//tracking, ce := service.CreditServiceInst.AddTracking(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCreditOnChainActionTracking{})
}

func (api CashApi) Deactivate(context *gin.Context) {
	//userId := common.GetUserId(context)
	//currency := context.DefaultQuery("currency", "")
	//
	//credit, ce := service.CreditServiceInst.DeactivateCredit(userId, currency)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCredit{})
}

func (api CashApi) Withdraw(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	var body bean.CashCreditWithdraw
	if common.ValidateBody(context, &body) != nil {
		return
	}
	//withdraw, ce := service.CreditServiceInst.AddCreditWithdraw(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashCreditWithdraw{})
}

func (api CashApi) CashStore(context *gin.Context) {
	//userId := common.GetUserId(context)
	//
	//withdraw, ce := service.CreditServiceInst.AddCreditWithdraw(userId, body)
	//if ce.ContextValidate(context) {
	//	return
	//}

	bean.SuccessResponse(context, bean.CashStore{})
}

func (api CashApi) CashStoreCreate(context *gin.Context) {
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

	bean.SuccessResponse(context, bean.CashStore{})
}

func (api CashApi) CashStoreUpdate(context *gin.Context) {
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

	bean.SuccessResponse(context, bean.CashStore{})
}

func (api CashApi) CashStorePrice(context *gin.Context) {
	bean.SuccessResponse(context, bean.CashStoreOrder{})
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

	bean.SuccessResponse(context, bean.CashStoreOrder{})
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
