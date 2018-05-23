package url

import (
	"github.com/autonomousdotai/handshake-exchange/api"
	"github.com/gin-gonic/gin"
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
	group.GET("/system-fees/:feeKey", func(context *gin.Context) {
		miscApi.GetSystemFee(context)
	})
	group.GET("/instant-buy/price", func(context *gin.Context) {
		creditCardAPi.GetProposeInstantOffer(context)
	})
	group.GET("/cc-limits", func(context *gin.Context) {
		miscApi.GetCCLimits(context)
	})

	return group
}
