package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type CronJobUrl struct {
}

func (url CronJobUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/public")

	miscApi := api.MiscApi{}
	coinbaseApi := api.CoinbaseApi{}
	blockchainIoApi := api.BlockChainApi{}

	// CRON JOB
	group.POST("/currency-rates", func(context *gin.Context) {
		miscApi.UpdateCurrencyRates(context)
	})
	// CRON JOB
	group.POST("/crypto-rates", func(context *gin.Context) {
		miscApi.UpdateCryptoRates(context)
	})
	// CRON JOB
	group.POST("/crypto-rates-extra", func(context *gin.Context) {
		miscApi.UpdateCryptoRatesExtra(context)
	})
	// CRON JOB
	group.POST("/finish-instant-offers", func(context *gin.Context) {
		miscApi.FinishInstantOffers(context)
	})

	// CRON JOB
	group.POST("/finish-crypto-transfer", func(context *gin.Context) {
		miscApi.FinishCryptoTransfer(context)
	})
	// CRON JOB
	group.POST("/update-cc-limit-track", func(context *gin.Context) {
		miscApi.UpdateUserCCLimitTracks(context)
	})

	group.POST("/coinbase/callback", func(context *gin.Context) {
		coinbaseApi.ReceiveCallback(context)
	})
	group.POST("/blockchainio/callback", func(context *gin.Context) {
		blockchainIoApi.ReceiveCallback(context)
	})
	group.POST("/system-fees", func(context *gin.Context) {
		miscApi.UpdateSystemFee(context)
	})
	group.POST("/system-configs", func(context *gin.Context) {
		miscApi.UpdateSystemConfig(context)
	})
	group.POST("/cc-limits", func(context *gin.Context) {
		miscApi.UpdateCCLimits(context)
	})
	group.POST("/sync-to-offer-solr/:offerId", func(context *gin.Context) {
		miscApi.SyncOfferToSolr(context)
	})
	group.POST("/sync-to-offer-store-solr/:offerId", func(context *gin.Context) {
		miscApi.SyncOfferStoreToSolr(context)
	})
	group.POST("/sync-to-offer-store-shake-solr/:offerId/:offerShakeId", func(context *gin.Context) {
		miscApi.SyncOfferStoreShakeToSolr(context)
	})

	group.POST("/start-app", func(context *gin.Context) {
		miscApi.StartApp(context)
	})
	group.POST("/script-update-xyz-123", func(context *gin.Context) {
		miscApi.ScriptUpdateAllOfferStoreSolr(context)
	})
	group.POST("/script-check-xyz-123", func(context *gin.Context) {
		miscApi.ScriptCheckFailedTransfer(context)
	})
	group.GET("/btc-confirmations/:txId", func(context *gin.Context) {
		miscApi.GetBTCConfirmation(context)
	})

	return group
}
