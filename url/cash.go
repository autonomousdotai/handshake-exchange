package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CashUrl struct {
}

func (url CashUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/cash")

	cashApi := api.CashApi{}
	group.GET("", cashApi.Dashboard)
	group.POST("", cashApi.Create)
	group.DELETE("", cashApi.Deactivate)
	group.POST("/deposit", cashApi.Deposit)
	group.POST("/tracking", cashApi.Tracking)
	group.POST("/withdraw", cashApi.Withdraw)
	group.POST("/store", cashApi.CashStore)
	group.PUT("/store", cashApi.CashStoreUpdate)
	group.GET("/price", cashApi.CashStorePrice)
	group.POST("/store/order", cashApi.CashStoreOrder)
	group.DELETE("/store/order/:id", cashApi.CashStoreRemoveOrder)

	return group
}
