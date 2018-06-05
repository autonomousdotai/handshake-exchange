package url

import (
	"github.com/ninjadotorg/handshake-exchange/api"
	"github.com/gin-gonic/gin"
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
	group.GET("/profile/cc-limit", func(context *gin.Context) {
		profileApi.GetCCLimit(context)
	})
	group.GET("/transactions", func(context *gin.Context) {
		profileApi.ListTransactions(context)
	})

	return group
}
