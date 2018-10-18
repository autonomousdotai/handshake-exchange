package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/adyen_service"
	"github.com/ninjadotorg/handshake-exchange/integration/bitpay_service"
	"github.com/ninjadotorg/handshake-exchange/integration/bitstamp_service"
	"github.com/ninjadotorg/handshake-exchange/integration/coinapi_service"
	"github.com/ninjadotorg/handshake-exchange/integration/ethereum_service"
	"github.com/ninjadotorg/handshake-exchange/integration/openexchangerates_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/service"
	"github.com/shopspring/decimal"
	"net/http"
	"os"
	"strings"
	"time"
)

type MiscApi struct {
}

// CRON JOB
func (api MiscApi) UpdateCurrencyRates(context *gin.Context) {
	rates, err := openexchangerates_service.GetExchangeRate()
	if api_error.PropagateErrorAndAbort(context, api_error.ExternalApiFailed, err) != nil {
		return
	}

	err = dao.MiscDaoInst.UpdateCurrencyRate(rates)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
func (api MiscApi) UpdateCryptoRates(context *gin.Context) {
	allRates := make([]bean.CryptoRate, 0)
	for _, currency := range []string{bean.BTC.Code, bean.ETH.Code, bean.LTC.Code, bean.BCH.Code} {
		rates := make([]bean.CryptoRate, 0)
		// resp, _ := coinbase_service.GetBuyPrice(currency)
		resp, _ := bitstamp_service.GetBuyPrice(currency)
		buy, _ := decimal.NewFromString(resp.Amount)
		buyFloat, _ := buy.Float64()
		rate := bean.CryptoRate{
			From:     currency,
			To:       bean.USD.Code,
			Buy:      buyFloat,
			Sell:     0,
			Exchange: bean.INSTANT_OFFER_PROVIDER_BITSTAMP,
		}
		rates = append(rates, rate)
		allRates = append(allRates, rate)

		err := dao.MiscDaoInst.UpdateCryptoRates(map[string][]bean.CryptoRate{currency: rates})
		if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
			return
		}
	}

	bean.SuccessResponse(context, allRates)
}

// CRON JOB
func (api MiscApi) UpdateCryptoRatesExtra(context *gin.Context) {
	rates, err := coinapi_service.GetExchangeRate([]string{bean.XRP.Code})
	if api_error.PropagateErrorAndAbort(context, api_error.ExternalApiFailed, err) != nil {
		return
	}
	err = dao.MiscDaoInst.UpdateCryptoRates(rates)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	allRates := make([]bean.CryptoRate, 0)
	bean.SuccessResponse(context, allRates)
}

func (api MiscApi) UpdateSystemFee(context *gin.Context) {
	systemFees, err := dao.MiscDaoInst.LoadSystemFeeToCache()
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, systemFees)
}

func (api MiscApi) UpdateSystemConfig(context *gin.Context) {
	systemFees, err := dao.MiscDaoInst.LoadSystemConfigToCache()
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, systemFees)
}

func (api MiscApi) GetCurrencyRate(context *gin.Context) {
	currency := context.Param("currency")
	to := dao.MiscDaoInst.GetCurrencyRateFromCache(currency[:3], currency[3:])
	if to.ContextValidate(context) {
		return
	}
	rate := to.Object.(bean.CurrencyRate)

	bean.SuccessResponse(context, rate)
}

func (api MiscApi) ListCryptoRates(context *gin.Context) {
	currency := context.Param("currency")
	to := dao.MiscDaoInst.GetCryptoRatesFromCache(currency)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, to.Objects)
}

func (api MiscApi) GetSystemFee(context *gin.Context) {
	feeKey := context.Param("feeKey")
	to := dao.MiscDaoInst.GetSystemFeeFromCache(feeKey)
	if to.ContextValidate(context) {
		return
	}
	systemFee := to.Object.(bean.SystemFee)

	bean.SuccessResponse(context, systemFee)
}

func (api MiscApi) GetSystemConfig(context *gin.Context) {
	feeKey := context.Param("systemKey")
	to := dao.MiscDaoInst.GetSystemConfigFromCache(feeKey)
	if to.ContextValidate(context) {
		return
	}
	systemConfig := to.Object.(bean.SystemConfig)

	bean.SuccessResponse(context, systemConfig)
}

func (api MiscApi) GetCryptoRate(context *gin.Context) {
	currency := context.Param("currency")

	// resp, err := coinbase_service.GetBuyPrice(currency)
	resp, err := bitstamp_service.GetBuyPrice(currency)
	if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, resp)
}

func (api MiscApi) GetCryptoRateAll(context *gin.Context) {
	currency := context.Param("currency")
	rateType := context.DefaultQuery("type", "buy")
	var resp1, resp2 interface{}
	var err error
	if rateType == "buy" {
		//resp1, err = coinbase_service.GetBuyPrice(currency)
		resp1, err = bitstamp_service.GetBuyPrice(currency)
		if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
			return
		}
	} else {
		//resp2, err = coinbase_service.GetSellPrice(currency)
		resp2, err = bitstamp_service.GetSellPrice(currency)
		if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
			return
		}
	}

	bean.SuccessResponse(context, map[string]interface{}{
		"buy":  resp1,
		"sell": resp2,
	})
}

func (api MiscApi) GetCryptoQuote(context *gin.Context) {
	type quoteStruct struct {
		Type         string
		Amount       string
		Currency     string
		FiatCurrency string
		FiatAmount   string
		Price        string
	}

	quoteType := context.DefaultQuery("type", "")
	amountStr := context.DefaultQuery("amount", "")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "")

	quote := quoteStruct{
		Type:         quoteType,
		Amount:       amountStr,
		Currency:     currency,
		FiatCurrency: fiatCurrency,
	}

	_, fiatPrice, fiatAmount, err := service.OfferServiceInst.GetQuote(quoteType, amountStr, currency, fiatCurrency)
	if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
		return
	}
	quote.Price = fiatPrice.Round(2).String()
	quote.FiatAmount = fiatAmount.Round(2).String()

	bean.SuccessResponse(context, quote)
}

func (api MiscApi) GetAllCryptoQuotes(context *gin.Context) {
	fiatCurrencyStr := context.DefaultQuery("fiat_currency", "RUB,VND,PHP,CAD,USD,EUR,HKD")
	fiatCurrencies := strings.Split(fiatCurrencyStr, ",")

	quotes := make([]interface{}, 0)
	for _, fiatCurrency := range fiatCurrencies {
		quotesTmp := service.OfferServiceInst.GetAllQuotes(fiatCurrency)
		quotes = append(quotes, quotesTmp...)
	}

	bean.SuccessResponse(context, quotes)
}

// CRON JOB
func (api MiscApi) FinishInstantOffers(context *gin.Context) {
	_, ce := service.CreditCardServiceInst.FinishInstantOffers()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
func (api MiscApi) FinishInstantOfferTransfers(context *gin.Context) {
	_, ce := service.CreditCardServiceInst.FinishInstantOfferTransfers()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
func (api MiscApi) FinishOfferConfirmingAddresses(context *gin.Context) {
	_, ce := service.OfferServiceInst.FinishOfferConfirmingAddresses()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
func (api MiscApi) FinishCryptoTransfer(context *gin.Context) {
	_, ce := service.OfferServiceInst.FinishCryptoTransfer()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
//func (api MiscApi) ExpireOfferHandshakes(context *gin.Context) {
//	err := dao.OfferDaoInst.UpdateExpiredHandshake()
//	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
//		return
//	}
//
//	bean.SuccessResponse(context, true)
//}

func (api MiscApi) UpdateCCLimits(context *gin.Context) {
	// country := context.Param("country")
	objs, err := dao.MiscDaoInst.LoadCCLimitToCache()
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, objs)
}

func (api MiscApi) GetCCLimits(context *gin.Context) {
	// country := context.Param("country")
	to := dao.MiscDaoInst.GetCCLimitFromCache()
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, to.Objects)
}

// CRON JOB
func (api MiscApi) UpdateUserCCLimitTracks(context *gin.Context) {
	// country := context.Param("country")
	ce := service.UserServiceInst.UpdateUserCCLimitTracks()
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, true)
}

// CRON JOB
//func (api MiscApi) UpdateTransferTracking(context *gin.Context) {
//	ce := service.OfferServiceInst.EndOffers()
//	if ce.ContextValidate(context) {
//		return
//	}
//
//	bean.SuccessResponse(context, true)
//}

func (api MiscApi) GetOfferStoreFreeStart(context *gin.Context) {
	userId := common.GetUserId(context)
	token := context.Param("token")
	freeStart, ce := service.OfferStoreServiceInst.GetCurrentFreeStart(userId, token)
	freeStart.Id = ""
	freeStart.Level = 0
	freeStart.Count = 0
	freeStart.Limit = 0
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, freeStart)
}

func (api MiscApi) TestCoinbaseReceive(context *gin.Context) {
	address := context.DefaultQuery("address", "")
	amount := context.DefaultQuery("amount", "")

	offer, ce := service.OfferServiceInst.ActiveOffer(address, amount)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api MiscApi) SyncOfferToSolr(context *gin.Context) {
	offerId := context.Param("offerId")

	offer, ce := service.OfferServiceInst.SyncToSolr(offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api MiscApi) SyncOfferStoreToSolr(context *gin.Context) {
	offerId := context.Param("offerId")

	offer, ce := service.OfferStoreServiceInst.SyncOfferStoreToSolr(offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api MiscApi) SyncOfferStoreShakeToSolr(context *gin.Context) {
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")

	offer, ce := service.OfferStoreServiceInst.SyncOfferStoreShakeToSolr(offerId, offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api MiscApi) StartApp(context *gin.Context) {
	api.UpdateCurrencyRates(context)
	api.UpdateCryptoRates(context)
	api.UpdateSystemFee(context)
	api.UpdateSystemConfig(context)
	api.UpdateCCLimits(context)
}

func (api MiscApi) TestEmail(context *gin.Context) {
	//offerStoreTO := dao.OfferStoreDaoInst.GetOfferStore("708")
	//offerStoreShakeTO := dao.OfferStoreDaoInst.GetOfferStoreShake("708", "wQ4rvvy4hWwDuZ31ThMe")
	//
	//c := make(chan error)
	//go notification.SendOfferStoreShakeToEmail(offerStoreShakeTO.Object.(bean.OfferStoreShake), offerStoreTO.Object.(bean.OfferStore), c)
	//fmt.Println(<-c)

	//c := make(chan error)
	//offerTO := dao.CreditCardDaoInst.GetInstantOffer("708", "SsiAUHmUpfopSdc2l5MY")
	//go notification.SendInstantOfferToFCM(offerTO.Object.(bean.InstantOffer), c)
	//err := <-c
	//fmt.Println(err)
	//if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
	//	return
	//}

	//fcmObj := bean.FCMObject{
	//	Notification: bean.FCMNotificationObject{
	//		Title:       "Hi Exchange",
	//		Body:        "Body Exchange",
	//		ClickAction: "https://staging.ninja.org/me",
	//	},
	//	To: "d-RV0aBxmAc:APA91bEUYloX1TkJ-RkaYIuflBFqaM5fZE3j18PufbBV9NiSmJ2qo5PUbOfYA_8nzngGvz77wO_4VyP4TF16whAvDRR55av2RDr0sTSg4hFyAbvlU4bjryMtwAs5GY8MIGiTvJ5cclHW",
	//}
	//err := fcm_service.SendFCM(fcmObj)
	//if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
	//	return
	//}
	//
	//bean.SuccessResponse(context, true)
}

func (api MiscApi) RemoveSolr(context *gin.Context) {
	id := context.Param("id")
	resp, _ := solr_service.DeleteObject(id)

	bean.SuccessResponse(context, resp)
}

func (api MiscApi) ScriptUpdateTxCount(context *gin.Context) {
	service.OfferStoreServiceInst.ScriptUpdateTransactionCount()
	bean.SuccessResponse(context, "ok")
}

func (api MiscApi) ScriptUpdateAllOfferStoreSolr(context *gin.Context) {
	service.OfferStoreServiceInst.ScriptUpdateOfferStoreSolr()
	bean.SuccessResponse(context, "ok")
}

func (api MiscApi) ScriptCheckFailedTransfer(context *gin.Context) {
	service.CreditCardServiceInst.ScriptCheckFailedTransfer()
	bean.SuccessResponse(context, "ok")
}

func (api MiscApi) GetBTCConfirmation(context *gin.Context) {
	txId := context.Param("txId")
	amount, address, confirmations, err := bitpay_service.GetBCHTransaction(txId)
	fmt.Println(err)
	fmt.Println(amount)
	fmt.Println(address)
	fmt.Println(confirmations)

	bean.SuccessResponse(context, amount)
}

func (api MiscApi) SendBtc(context *gin.Context) {
	//btcService := bitcoin_service.BitcoinService{}
	//tx, err := btcService.SendTransaction("1DrUv69utLBiu5CMCiHyiKNg5A9CxvoMJV", common.StringToDecimal("0.00001"))
	//fmt.Println(err)
	//bean.SuccessResponse(context, tx)

	//amountStr := ""
	//address := ""
	//offchainId := "refund"
	////
	//client := exchangecreditatm_service.ExchangeCreditAtmClient{}
	//amount := common.StringToDecimal(amountStr)
	//txHash, _, onChainErr := client.ReleasePartialFund(offchainId, 1, amount, address, uint64(0), false, "")
	//if onChainErr != nil {
	//	fmt.Println(onChainErr)
	//} else {
	//}
	//fmt.Println(txHash)
	//bean.SuccessResponse(context, txHash)

	//coinbaseTx, errWithdraw := coinbase_service.SendTransaction(address, amountStr, "BTC",
	//	fmt.Sprintf("Withdraw tx = %s", "6044"), "i55dcuEH20bketepu37k")
	//if errWithdraw != nil {
	//	fmt.Println(errWithdraw)
	//}
	//fmt.Println(coinbaseTx.Id)

	//coinbaseTx, _ := coinbase_service.GetTransaction("59baebcb-ad50-5508-a729-4422c8a31ddc", "BTC")
	//fmt.Println(coinbaseTx.Id, coinbaseTx.Status, coinbaseTx.Amount, coinbaseTx.Description, coinbaseTx.CreatedAt)

	//bean.SuccessResponse(context, coinbaseTx.Id)
}

func (api MiscApi) FinishCreditTracking(context *gin.Context) {
	ce := service.CreditServiceInst.FinishTracking()
	if ce.ContextValidate(context) {
		return
	}
	bean.SuccessResponse(context, true)
}

func (api MiscApi) ProcessCreditWithdraw(context *gin.Context) {
	ce := service.CreditServiceInst.ProcessCreditWithdraw()
	if ce.ContextValidate(context) {
		return
	}
	bean.SuccessResponse(context, true)
}

func (api MiscApi) SetupCreditPool(context *gin.Context) {
	ce := service.CreditServiceInst.SetupCreditPool()
	if ce.ContextValidate(context) {
		return
	}
	bean.SuccessResponse(context, true)
}

func (api MiscApi) SetupCreditPoolCache(context *gin.Context) {
	ce := service.CreditServiceInst.SetupCreditPoolCache()
	if ce.ContextValidate(context) {
		return
	}
	bean.SuccessResponse(context, true)
}

func (api MiscApi) SyncCreditTransactionToSolr(context *gin.Context) {
	id := context.Param("id")
	currency := context.Param("currency")

	obj, ce := service.CreditServiceInst.SyncCreditTransactionToSolr(currency, id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) SyncCreditDepositToSolr(context *gin.Context) {
	id := context.Param("id")
	currency := context.Param("currency")

	obj, ce := service.CreditServiceInst.SyncCreditDepositToSolr(currency, id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) SyncCreditWithdrawToSolr(context *gin.Context) {
	id := context.Param("id")

	obj, ce := service.CreditServiceInst.SyncCreditWithdrawToSolr(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) SyncCashStoreToSolr(context *gin.Context) {
	id := context.Param("id")

	obj, ce := service.CashServiceInst.SyncCashStoreToSolr(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) SyncCashOrderToSolr(context *gin.Context) {
	id := context.Param("id")

	obj, ce := service.CashServiceInst.SyncCashOrderToSolr(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) SyncCoinOrderToSolr(context *gin.Context) {
	id := context.Param("id")

	obj, ce := service.CoinServiceInst.SyncCoinOrderToSolr(id)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, obj)
}

func (api MiscApi) GenerateAddress(context *gin.Context) {
	address, key := ethereum_service.GenerateAddress()

	bean.SuccessResponse(context, map[string]string{
		"address": address,
		"key":     key,
	})
}

func (api MiscApi) SetupContractKeys(context *gin.Context) {
	service.CreditServiceInst.SetupContractKey("ETH_HIGH_ADMIN_KEYS")
	service.CreditServiceInst.SetupContractKey("ETH_LOW_ADMIN_KEYS")

	bean.SuccessResponse(context, true)
}

func (api MiscApi) LoadBitstampWithdrawToCache(context *gin.Context) {
	resp, err := bitstamp_service.WithdrawalRequests(0)
	if err == nil {
		dao.MiscDaoInst.BitstampWithdrawRequestToCache(resp)
	}
	bean.SuccessResponse(context, true)
}

func (api MiscApi) AddAdminAddress(context *gin.Context) {
	address := context.Param("address")
	service.CreditServiceInst.AddAdminAddressToContract(address)

	bean.SuccessResponse(context, "ok")
}

func (api MiscApi) ServerTime(context *gin.Context) {
	bean.SuccessResponse(context, time.Now().UTC().Format("2006-01-02T15:04:05.000+00:00"))
}

func (api MiscApi) TestAnything(context *gin.Context) {
	bean.SuccessResponse(context, true)
}

func (api MiscApi) AdyenRedirect(context *gin.Context) {
	id := fmt.Sprintf("%d", time.Now().UTC().Unix())
	dao.CreditCardDaoInst.AddInitInstantOffer(id,
		adyen_service.GetNotificationData(context.PostForm("MD"), context.PostForm("PaRes")))
	urlStr := fmt.Sprintf("%s/cc-payment?MD=%s", os.Getenv("ADYEN_REDIRECT_URL"), id)
	str := `<!DOCTYPE HTML>
<html lang="en-US">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="refresh" content="1;url=%s">
        <script type="text/javascript">
            window.location.href = "%s"
        </script>
        <title>Page Redirection</title>
    </head>
    <body>
    </body>
</html>`

	finalStr := fmt.Sprintf(str, urlStr, urlStr)
	context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(finalStr))
}

func (api MiscApi) AdyenData(context *gin.Context) {
	id := context.Param("id")
	t := dao.CreditCardDaoInst.GetInitInstantOffer(id)

	bean.SuccessResponse(context, t.Object)
}
