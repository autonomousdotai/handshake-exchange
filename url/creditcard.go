package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CreditCardUrl struct {
}

func (url CreditCardUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/instant-buys")

	creditCardApi := api.CreditCardApi{}
	group.POST("", func(context *gin.Context) {
		creditCardApi.PayInstantOffer(context)
	})
	group.GET("/:offerId", func(context *gin.Context) {
		creditCardApi.GetInstantOffers(context)
	})
	group.POST("/init-payment", func(context *gin.Context) {
		creditCardApi.InitAdyenPayment(context)
	})

	return group
}
