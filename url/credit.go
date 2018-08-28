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
	group.GET("/transaction", creditApi.ListTransaction)
	group.POST("", creditApi.Create)
	group.POST("/deposit", creditApi.Deposit)
	group.GET("/deposit", creditApi.ListDeposit)
	group.POST("/transfer", creditApi.Transfer)
	group.DELETE("", creditApi.Deactivate)
	group.POST("/withdraw", creditApi.Withdraw)
	group.GET("/withdraw", creditApi.ListWithdraw)

	return group
}
