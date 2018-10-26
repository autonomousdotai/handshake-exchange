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
	group.GET("/quote-reverse", coinApi.GetQuoteReverse)
	group.GET("/order", coinApi.ListCoinOrders)
	group.GET("/selling-order", coinApi.ListCoinSellingOrders)
	group.POST("/order", coinApi.CoinOrder)
	group.POST("/selling-order", coinApi.CoinSellingOrder)
	group.POST("/order/:id", coinApi.FinishCoinOrder)
	group.POST("/selling-order/:id", coinApi.CloseCoinSellingOrder)
	group.DELETE("/order/:id", coinApi.CancelCoinOrder)
	group.DELETE("/selling-order/:id", coinApi.CancelCoinSellingOrder)
	group.PUT("/order/:id/pick", coinApi.PickCoinOrder)
	group.PUT("/selling-order/:id/pick", coinApi.PickCoinSellingOrder)
	group.PUT("/order/:id/reject", coinApi.RejectCoinOrder)
	group.PUT("/order/:id", coinApi.UpdateCoinOrder)
	group.GET("/center/:country", coinApi.ListCoinCenter)
	group.GET("/review", coinApi.ListReview)
	group.POST("/review", coinApi.AddReview)

	return group
}
