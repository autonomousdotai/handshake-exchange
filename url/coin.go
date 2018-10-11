package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CoinUrl struct {
}

func (url CoinUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/coin")

	coinApi := api.CoinApi{}

	group.GET("/quote", coinApi.GetQuote)
	group.GET("/order", coinApi.ListCoinOrders)
	group.POST("/order", coinApi.CoinOrder)
	group.POST("/order/:id/:currency/:amount", coinApi.FinishCoinOrder)
	group.DELETE("/order/:id", coinApi.CancelCoinOrder)
	group.PUT("/order/:id/pick", coinApi.PickCoinOrder)
	group.PUT("/order/:id/reject", coinApi.RejectCoinOrder)
	group.PUT("/order/:id", coinApi.UpdateCoinOrder)
	group.GET("/center/:country", coinApi.ListCoinCenter)

	return group
}
