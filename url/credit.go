package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CreditUrl struct {
}

func (url CreditUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/credit")

	creditApi := api.CreditApi{}
	group.GET("", creditApi.Dashboard)
	group.POST("", creditApi.Create)
	group.DELETE("", creditApi.Deactivate)
	group.GET("/transaction", creditApi.ListTransaction)
	group.POST("/deposit", creditApi.Deposit)
	group.GET("/deposit", creditApi.ListDeposit)
	group.POST("/tracking", creditApi.Tracking)
	group.POST("/withdraw", creditApi.Withdraw)
	group.GET("/withdraw", creditApi.ListWithdraw)
	group.POST("/nonce", creditApi.Nonce)

	return group
}
