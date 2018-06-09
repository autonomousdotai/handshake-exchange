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
	group.DELETE("/:offerId", func(context *gin.Context) {
		offerApi.RemoveOfferStoreItem(context)
	})
	group.POST("/:offerId/shakes", func(context *gin.Context) {
		offerApi.CreateOfferStoreShake(context)
	})
	group.DELETE("/:offerId/shakes/:offerShakeId", func(context *gin.Context) {
		offerApi.RejectOfferStoreShake(context)
	})
	group.POST("/:offerId/shakes/:offerShakeId", func(context *gin.Context) {
		offerApi.CompleteOfferStoreShake(context)
	})

	return group
}
