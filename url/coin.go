package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CoinUrl struct {
}

func (url CoinUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/coin")

	cashApi := api.CoinApi{}

	group.GET("/quote", cashApi.GetQuote)
	group.GET("/order-type", cashApi.CoinOrderType)
	group.GET("/order", cashApi.ListCoinOrders)
	group.POST("/order", cashApi.CoinOrder)
	group.POST("/order/:id/:amount", cashApi.FinishCoinOrder)
	group.DELETE("/order/:id", cashApi.RejectCoinOrder)
	group.PUT("/order/:id/picked", cashApi.PickCoinOrder)
	group.PUT("/order/:id", cashApi.UpdateCashOrder)
	group.GET("/center/:country", cashApi.ListCoinCenter)

	return group
}
