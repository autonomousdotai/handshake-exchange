package url

import (
	"github.com/autonomousdotai/handshake-exchange/api"
	"github.com/gin-gonic/gin"
)

type MiscUrl struct {
}

func (url MiscUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/misc")

	miscApi := api.MiscApi{}
	group.POST("/system-fees", func(context *gin.Context) {
		miscApi.UpdateSystemFee(context)
	})
	group.POST("/cc-limits", func(context *gin.Context) {
		miscApi.UpdateCCLimits(context)
	})

	return group
}
