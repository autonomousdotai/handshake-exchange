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

	group.GET("/quotes", cashApi.ListProposeQuotes)
	group.GET("/store", cashApi.CashStore)
	group.POST("/store", cashApi.CashStoreCreate)
	group.PUT("/store", cashApi.CashStoreUpdate)
	group.GET("/price", cashApi.CashStorePrice)
	group.GET("/order", cashApi.ListCashStoreOrders)
	group.POST("/order", cashApi.CashStoreOrder)
	group.POST("/order/:id/:amount", cashApi.FinishCashOrder)
	group.DELETE("/order/:id", cashApi.RejectCashOrder)
	group.PUT("/order/:id", cashApi.UpdateCashOrder)
	group.GET("/center/:country", cashApi.ListCashCenter)

	return group
}
