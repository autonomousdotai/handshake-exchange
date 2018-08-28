package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CreditUrl struct {
}

func (url CreditUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/atm")

	creditApi := api.CreditApi{}
	group.GET("", creditApi.Dashboard)
	group.GET("/transactions", creditApi.ListTransaction)
	group.POST("", creditApi.Create)
	group.PUT("", creditApi.Deposit)
	group.POST("/transfer", creditApi.Deposit)
	group.DELETE("", creditApi.Deactivate)
	group.POST("/withdraw", creditApi.Withdraw)

	return group
}
