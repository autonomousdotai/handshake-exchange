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

	group.GET("/store", cashApi.CashStore)
	group.POST("/store", cashApi.CashStoreCreate)
	group.PUT("/store", cashApi.CashStoreUpdate)
	group.GET("/price", cashApi.CashStorePrice)
	group.POST("/order", cashApi.CashStoreOrder)
	group.POST("/order/:id/:amount", cashApi.FinishCashOrder)
	group.DELETE("/order/:id", cashApi.CashStoreRemoveOrder)

	return group
}
