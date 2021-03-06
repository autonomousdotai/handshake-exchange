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
	internalApi := api.InternalApi{}
	coinApi := api.CoinApi{}

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
	group.POST("/finish-instant-offer-transfers", func(context *gin.Context) {
		miscApi.FinishInstantOfferTransfers(context)
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
	group.POST("/setup-credit-pool", func(context *gin.Context) {
		miscApi.SetupCreditPool(context)
	})
	group.POST("/setup-credit-pool-cache", func(context *gin.Context) {
		miscApi.SetupCreditPoolCache(context)
	})
	group.POST("/setup-contract-keys", miscApi.SetupContractKeys)

	//CRON JOB
	group.POST("/finish-credit-tracking", func(context *gin.Context) {
		miscApi.FinishCreditTracking(context)
	})
	//CRON JOB
	group.POST("/process-credit-withdraw", func(context *gin.Context) {
		miscApi.ProcessCreditWithdraw(context)
	})
	group.POST("/sync-to-credit-transaction-solr/:currency/:id", func(context *gin.Context) {
		miscApi.SyncCreditTransactionToSolr(context)
	})
	group.POST("/sync-to-credit-deposit-solr/:currency/:id", func(context *gin.Context) {
		miscApi.SyncCreditDepositToSolr(context)
	})
	group.POST("/sync-to-credit-withdraw-solr/:id", func(context *gin.Context) {
		miscApi.SyncCreditWithdrawToSolr(context)
	})
	group.POST("/eth-address", miscApi.GenerateAddress)
	group.POST("/add-address/:address", miscApi.AddAdminAddress)

	group.POST("/sync-to-cash-store-solr/:id", func(context *gin.Context) {
		miscApi.SyncCashStoreToSolr(context)
	})
	group.POST("/sync-to-cash-order-solr/:id", func(context *gin.Context) {
		miscApi.SyncCashOrderToSolr(context)
	})
	group.GET("/server-time", func(context *gin.Context) {
		miscApi.ServerTime(context)
	})
	group.POST("/authorise-receive", func(context *gin.Context) {
		miscApi.AdyenRedirect(context)
	})
	group.GET("/authorise-receive/:id", func(context *gin.Context) {
		miscApi.AdyenData(context)
	})

	//CRON JOB
	group.POST("/load-bitstamp-withdraw-to-cache", func(context *gin.Context) {
		miscApi.LoadBitstampWithdrawToCache(context)
	})
	//CRON JOB
	group.POST("/reset-coin-user-limit", func(context *gin.Context) {
		coinApi.ResetCoinUserLimit(context)
	})
	//CRON JOB
	group.POST("/remove-expired-coin-order", func(context *gin.Context) {
		coinApi.RemoveExpiredOrder(context)
	})
	//CRON JOB
	group.POST("/remove-expired-coin-selling-order", func(context *gin.Context) {
		coinApi.SellingRemoveExpiredOrder(context)
	})
	//CRON JOB
	group.POST("/coin-order-call-notification", func(context *gin.Context) {
		coinApi.OrderCallNotification(context)
	})
	//CRON JOB
	group.POST("/coin-address-tracking-deposit", func(context *gin.Context) {
		coinApi.AddressTrackingDeposit(context)
	})
	group.POST("/coin-init-bank", func(context *gin.Context) {
		coinApi.CoinInitBank(context)
	})
	group.GET("/external-bank-list", func(context *gin.Context) {
		miscApi.ExternalBankList(context)
	})

	group.GET("/voice-order-notification", func(context *gin.Context) {
		coinApi.VoiceOrderNotification(context)
	})

	group.POST("/sync-to-coin-order-solr/:id", func(context *gin.Context) {
		miscApi.SyncCoinOrderToSolr(context)
	})
	group.POST("/script-update-xyz-123", func(context *gin.Context) {
		miscApi.ScriptUpdateAllOfferStoreSolr(context)
	})
	group.POST("/script-check-xyz-123", func(context *gin.Context) {
		miscApi.ScriptCheckFailedTransfer(context)
	})
	group.POST("/test-btc-xyz-123", func(context *gin.Context) {
		miscApi.SendBtc(context)
	})
	group.GET("/btc-confirmations/:txId", func(context *gin.Context) {
		miscApi.GetBTCConfirmation(context)
	})

	group.POST("/test-anything", func(context *gin.Context) {
		miscApi.TestAnything(context)
	})

	// For autonomous
	group.POST("/payment", func(context *gin.Context) {
		internalApi.Payment(context)
	})

	return group
}
