package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type UserUrl struct {
}

func (url UserUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/user")

	profileApi := api.ProfileApi{}
	group.GET("/profile", func(context *gin.Context) {
		profileApi.GetProfile(context)
	})
	group.POST("/profile", func(context *gin.Context) {
		profileApi.AddProfile(context)
	})
	group.PUT("/profile", func(context *gin.Context) {
		profileApi.UpdateProfile(context)
	})
	group.POST("/profile/offline/:offline", func(context *gin.Context) {
		profileApi.UpdateProfileOffline(context)
	})
	group.GET("/profile/cc-limit", func(context *gin.Context) {
		profileApi.GetCCLimit(context)
	})
	group.GET("/transactions", func(context *gin.Context) {
		profileApi.ListTransactions(context)
	})
	group.GET("/transaction-counts", func(context *gin.Context) {
		profileApi.ListTransactionCounts(context)
	})
	group.GET("/referral-summary", func(context *gin.Context) {
		profileApi.ListReferralSummary(context)
	})

	return group
}
