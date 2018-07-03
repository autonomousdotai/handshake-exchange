package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
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
	group.POST("/:offerId/complete", func(context *gin.Context) {
		offerApi.CompleteShakeOffer(context)
	})
	group.POST("/:offerId/reject", func(context *gin.Context) {
		offerApi.RejectShakeOffer(context)
	})
	group.POST("/:offerId/accept", func(context *gin.Context) {
		offerApi.AcceptShakeOffer(context)
	})
	group.POST("/:offerId/cancel", func(context *gin.Context) {
		offerApi.CancelShakeOffer(context)
	})
	group.POST("/:offerId/onchain-tracking", func(context *gin.Context) {
		offerApi.OnChainOfferTracking(context)
	})
	return group
}
