package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/service"
)

type CreditApi struct {
}

func (api CreditApi) Dashboard(context *gin.Context) {
	userId := common.GetUserId(context)

	credit, ce := service.CreditServiceInst.GetCredit(userId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, credit)
}

func (api CreditApi) Create(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.Credit
	if common.ValidateBody(context, &body) != nil {
		return
	}

	body.ChainId = chainId
	body.Language = language
	body.FCM = fcm

	credit, ce := service.CreditServiceInst.AddCredit(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, credit)
}

func (api CreditApi) ListTransaction(context *gin.Context) {
	tx := bean.CreditTransaction{}
	bean.SuccessResponse(context, []bean.CreditTransaction{
		tx,
		tx,
		tx,
	})
}

func (api CreditApi) Deposit(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CreditDepositInput
	if common.ValidateBody(context, &body) != nil {
		return
	}

	tracking, ce := service.CreditServiceInst.AddDeposit(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, tracking)
}

func (api CreditApi) ListDeposit(context *gin.Context) {
	tx := bean.CreditDeposit{}
	bean.SuccessResponse(context, []bean.CreditDeposit{
		tx,
		tx,
		tx,
	})
}

func (api CreditApi) Tracking(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CreditOnChainActionTrackingInput
	if common.ValidateBody(context, &body) != nil {
		return
	}

	tracking, ce := service.CreditServiceInst.AddTracking(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, tracking)
}

func (api CreditApi) Deactivate(context *gin.Context) {
	userId := common.GetUserId(context)
	currency := context.DefaultQuery("currency", "")

	credit, ce := service.CreditServiceInst.DeactivateCredit(userId, currency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, credit)
}

func (api CreditApi) Withdraw(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.CreditWithdraw
	if common.ValidateBody(context, &body) != nil {
		return
	}
	withdraw, ce := service.CreditServiceInst.AddCreditWithdraw(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, withdraw)
}

func (api CreditApi) ListWithdraw(context *gin.Context) {
	status := context.DefaultQuery("status", "processing")

	fmt.Println(status)
	if status == "processing" {
		withdraws, ce := service.CreditServiceInst.ListCreditProcessingWithdraw()
		if ce.ContextValidate(context) {
			return
		}
		fmt.Println(withdraws)
		bean.SuccessResponse(context, withdraws)
		return
	} else if status == "processed" {
		withdraws, ce := service.CreditServiceInst.ListCreditProcessedWithdraw()
		if ce.ContextValidate(context) {
			return
		}
		bean.SuccessResponse(context, withdraws)
		return
	}

	withdraws := make([]bean.CreditWithdraw, 0)
	bean.SuccessResponse(context, withdraws)
}

func (api CreditApi) UpdateProcessedWithdraw(context *gin.Context) {
	withdrawId := context.Param("id")

	var body bean.CreditWithdraw
	if common.ValidateBody(context, &body) != nil {
		return
	}
	withdraw, ce := service.CreditServiceInst.FinishCreditWithdraw(withdrawId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, withdraw)
}
