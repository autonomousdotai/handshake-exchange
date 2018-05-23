package url

import (
	"github.com/autonomousdotai/handshake-exchange/api"
	"github.com/gin-gonic/gin"
)

type CronJobUrl struct {
}

func (url CronJobUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/cron-job")

	miscApi := api.MiscApi{}

	group.POST("/currency-rates", func(context *gin.Context) {
		miscApi.UpdateCurrencyRates(context)
	})
	group.POST("/crypto-rates", func(context *gin.Context) {
		miscApi.UpdateCryptoRates(context)
	})
	group.POST("/finish-instant-offers", func(context *gin.Context) {
		miscApi.FinishInstantOffers(context)
	})
	//group.POST("/expire-offer-handshakes", func(context *gin.Context) {
	//	miscApi.ExpireOfferHandshakes(context)
	//})
	group.POST("/update-cc-limit-track", func(context *gin.Context) {
		miscApi.UpdateUserCCLimitTracks(context)
	})

	return group
}
