package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/adyen_service"
	"github.com/ninjadotorg/handshake-exchange/service"
	"strconv"
	"time"
)

type CreditCardApi struct {
}

func (api CreditCardApi) GetProposeInstantOffer(context *gin.Context) {
	amount := context.DefaultQuery("amount", "")
	currency := context.DefaultQuery("currency", "")

	offer, ce := service.CreditCardServiceInst.GetProposeInstantOffer(amount, currency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api CreditCardApi) GetCryptoPrice(context *gin.Context) {
	amount := context.DefaultQuery("amount", "1")
	currency := context.DefaultQuery("currency", "")
	fiatCurrency := context.DefaultQuery("fiat_currency", "")

	offer, ce := service.CreditCardServiceInst.GetCryptoPrice(amount, currency, fiatCurrency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api CreditCardApi) PayInstantOffer(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.InstantOffer
	if common.ValidateBody(context, &body) != nil {
		return
	}

	body.UID = userId
	id, _ := strconv.Atoi(chainId)
	body.ChainId = int64(id)
	body.Language = language
	body.FCM = fcm
	offer, ce := service.CreditCardServiceInst.PayInstantOffer(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api CreditCardApi) GetInstantOffers(context *gin.Context) {
	userId := common.GetUserId(context)

	offerId := context.Param("offerId")

	offerTO := dao.CreditCardDaoInst.GetInstantOffer(userId, offerId)
	if offerTO.ContextValidate(context) {
		return
	}
	offer := offerTO.Object.(bean.InstantOffer)

	bean.SuccessResponse(context, offer)
}

func (api CreditCardApi) InitAdyenPayment(context *gin.Context) {
	var body adyen_service.AdyenAuthorise
	if common.ValidateBody(context, &body) != nil {
		return
	}
	body.Reference = fmt.Sprintf("%d", time.Now().UTC().Unix())
	resp, err := adyen_service.Authorise(body)
	fmt.Print(err)
	bean.SuccessResponse(context, resp)
}
