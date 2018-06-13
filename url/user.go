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
	group.POST("/profile/offline/:offline", func(context *gin.Context) {
		profileApi.UpdateProfileOffline(context)
	})
	group.GET("/profile/cc-limit", func(context *gin.Context) {
		profileApi.GetCCLimit(context)
	})
	group.GET("/transactions", func(context *gin.Context) {
		profileApi.ListTransactions(context)
	})

	return group
}
