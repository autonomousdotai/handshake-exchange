package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type OfferStoreUrl struct {
}

func (url OfferStoreUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/offer-stores")

	offerApi := api.OfferStoreApi{}
	group.POST("", func(context *gin.Context) {
		offerApi.CreateOfferStore(context)
	})
	group.GET("/:offerId", func(context *gin.Context) {
		offerApi.GetOfferStore(context)
	})
	group.POST("/:offerId", func(context *gin.Context) {
		offerApi.AddOfferStoreItem(context)
	})
	group.PUT("/:offerId", func(context *gin.Context) {
		offerApi.UpdateOfferStoreItem(context)
	})
	group.DELETE("/:offerId", func(context *gin.Context) {
		offerApi.RemoveOfferStoreItem(context)
	})

	group.POST("/:offerId/reviews/:offerShakeId", func(context *gin.Context) {
		offerApi.ReviewOfferStore(context)
	})
	group.POST("/:offerId/shakes", func(context *gin.Context) {
		offerApi.CreateOfferStoreShake(context)
	})
	group.DELETE("/:offerId/shakes/:offerShakeId", func(context *gin.Context) {
		offerApi.RejectOfferStoreShake(context)
	})
	group.POST("/:offerId/shakes/:offerShakeId/complete", func(context *gin.Context) {
		offerApi.CompleteOfferStoreShake(context)
	})

	// Support method
	group.POST("/:offerId/shakes/:offerShakeId/7tHCLp8XpajPJaVh", func(context *gin.Context) {
		offerApi.OfferStoreShakeLocationTracking(context)
	})

	return group
}
