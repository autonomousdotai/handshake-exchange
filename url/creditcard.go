package url

import (
	"github.com/autonomousdotai/handshake-exchange/api"
	"github.com/gin-gonic/gin"
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

	return group
}
