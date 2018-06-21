package url

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api"
)

type InfoUrl struct {
}

func (url InfoUrl) Create(router *gin.Engine) *gin.RouterGroup {
	group := router.Group("/info")

	miscApi := api.MiscApi{}
	creditCardAPi := api.CreditCardApi{}

	group.GET("/currency-rates/:currency", func(context *gin.Context) {
		miscApi.GetCurrencyRate(context)
	})
	group.GET("/crypto-rates/:currency", func(context *gin.Context) {
		miscApi.GetCryptoRate(context)
	})
	group.GET("/crypto-rates-all/:currency", func(context *gin.Context) {
		miscApi.GetCryptoRateAll(context)
	})
	group.GET("/system-fees/:feeKey", func(context *gin.Context) {
		miscApi.GetSystemFee(context)
	})
	group.GET("/system-configs/:systemKey", func(context *gin.Context) {
		miscApi.GetSystemConfig(context)
	})
	group.GET("/instant-buy/price", func(context *gin.Context) {
		creditCardAPi.GetProposeInstantOffer(context)
	})
	group.GET("/crypto-quote", func(context *gin.Context) {
		miscApi.GetCryptoQuote(context)
	})
	group.GET("/crypto-quotes", func(context *gin.Context) {
		miscApi.GetAllCryptoQuotes(context)
	})
	group.GET("/cc-limits", func(context *gin.Context) {
		miscApi.GetCCLimits(context)
	})
	group.POST("/cc-limits", func(context *gin.Context) {
		miscApi.GetCCLimits(context)
	})

	return group
}
