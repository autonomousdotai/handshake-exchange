package url

import (
	"github.com/autonomousdotai/handshake-exchange/api"
	"github.com/gin-gonic/gin"
)

type OfferUrl struct {
}

func (url OfferUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/offers")

	offerApi := api.OfferApi{}
	group.POST("", func(context *gin.Context) {
		offerApi.CreateOffer(context)
	})
	group.GET("", func(context *gin.Context) {
		offerApi.ListOffers(context)
	})

	return group
}
