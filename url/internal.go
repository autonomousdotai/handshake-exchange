package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type InternalUrl struct {
}

func (url InternalUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/internal")

	creditApi := api.CreditApi{}
	group.GET("credit/withdraw", creditApi.ListWithdraw)
	group.POST("credit/withdraw/:id", creditApi.UpdateProcessedWithdraw)
	// group.POST("redeem")

	return group
}
