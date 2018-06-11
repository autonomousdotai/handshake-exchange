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
	onChainApi := api.OnChainApi{}
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
	group.POST("/finish-instant-offers", func(context *gin.Context) {
		miscApi.FinishInstantOffers(context)
	})
	// CRON JOB
	group.POST("/finish-offer-confirming-addresses", func(context *gin.Context) {
		miscApi.FinishOfferConfirmingAddresses(context)
	})
	// CRON JOB
	//group.POST("/transfer-tracking", func(context *gin.Context) {
	//	miscApi.UpdateTransferTracking(context)
	//})
	// CRON JOB
	group.POST("/update-cc-limit-track", func(context *gin.Context) {
		miscApi.UpdateUserCCLimitTracks(context)
	})
	// CRON JOB
	group.POST("/update-offer-init-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferInit(context)
	})
	// CRON JOB
	group.POST("/update-offer-close-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferClose(context)
	})
	// CRON JOB
	group.POST("/update-offer-shake-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferShake(context)
	})
	// CRON JOB
	group.POST("/update-offer-reject-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferReject(context)
	})
	// CRON JOB
	group.POST("/update-offer-complete-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferComplete(context)
	})
	// CRON JOB
	group.POST("/update-offer-withdraw-on-chain", func(context *gin.Context) {
		onChainApi.UpdateOfferWithdraw(context)
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
	group.POST("/init-handshake-block", func(context *gin.Context) {
		onChainApi.StartOnChainOfferStoreBlock(context)
	})
	group.POST("/start-app", func(context *gin.Context) {
		miscApi.StartApp(context)
	})

	// Internal
	//group.POST("/test-coinbase-receive", func(context *gin.Context) {
	//	miscApi.TestCoinbaseReceive(context)
	//})

	return group
}
