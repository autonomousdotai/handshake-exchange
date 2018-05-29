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
	group.GET("/:offerId", func(context *gin.Context) {
		offerApi.GetOffer(context)
	})
	group.POST("/:offerId", func(context *gin.Context) {
		offerApi.ShakeOffer(context)
	})
	group.DELETE("/:offerId", func(context *gin.Context) {
		offerApi.CloseOffer(context)
	})
	group.POST("/:offerId/shake", func(context *gin.Context) {
		offerApi.CompleteShakeOffer(context)
	})
	group.DELETE("/:offerId/shake", func(context *gin.Context) {
		offerApi.RejectShakeOffer(context)
	})

	return group
}
