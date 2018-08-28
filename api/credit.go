package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
)

type CreditApi struct {
}

func (api CreditApi) Create(context *gin.Context) {
	bean.SuccessResponse(context, bean.Credit{})
}

func (api CreditApi) Dashboard(context *gin.Context) {
	item := bean.CreditItem{}
	bean.SuccessResponse(context, bean.Credit{
		Items: map[string]bean.CreditItem{
			bean.BTC.Code: item,
			bean.ETH.Code: item,
			bean.BCH.Code: item,
		},
	})
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
