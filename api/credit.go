package api

import (
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
	bean.SuccessResponse(context, bean.CreditDeposit{})
}

func (api CreditApi) ListDeposit(context *gin.Context) {
	tx := bean.CreditDeposit{}
	bean.SuccessResponse(context, []bean.CreditDeposit{
		tx,
		tx,
		tx,
	})
}

func (api CreditApi) Transfer(context *gin.Context) {
	bean.SuccessResponse(context, bean.CreditOnChainActionTracking{})
}

func (api CreditApi) Deactivate(context *gin.Context) {
	bean.SuccessResponse(context, bean.Credit{})
}

func (api CreditApi) Withdraw(context *gin.Context) {
	bean.SuccessResponse(context, bean.CreditWithdraw{})
}

func (api CreditApi) ListWithdraw(context *gin.Context) {
	tx := bean.CreditWithdraw{}
	bean.SuccessResponse(context, []bean.CreditWithdraw{
		tx,
		tx,
		tx,
	})
}
