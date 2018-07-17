package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/service"
	"strconv"
)

type OfferStoreApi struct {
}

func (api OfferStoreApi) CreateOfferStore(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.OfferStoreSetup
	if common.ValidateBody(context, &body) != nil {
		return
	}

	id, _ := strconv.Atoi(chainId)
	body.Offer.ChainId = int64(id)
	body.Offer.Language = language
	body.Offer.FCM = fcm
	offer, ce := service.OfferStoreServiceInst.CreateOfferStore(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer.Offer)
}

func (api OfferStoreApi) GetOfferStore(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	offer, ce := service.OfferStoreServiceInst.GetOfferStore(userId, offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) AddOfferStoreItem(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	var body bean.OfferStoreItem
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.AddOfferStoreItem(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) UpdateOfferStoreItem(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	var body bean.OfferStoreSetup
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.UpdateOfferStore(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) RemoveOfferStoreItem(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	currency := context.DefaultQuery("currency", "")

	if currency == "" {
		api_error.AbortWithValidateErrorSimple(context, api_error.InvalidQueryParam)
	}

	offer, ce := service.OfferStoreServiceInst.RemoveOfferStoreItem(userId, offerId, currency)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) RefillOfferStoreItem(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	var body bean.OfferStoreSetup
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.RefillOfferStoreItem(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) CreateOfferStoreShake(context *gin.Context) {
	userId := common.GetUserId(context)
	chainId := common.GetChainId(context)
	offerId := context.Param("offerId")
	language := common.GetLanguage(context)
	fcm := common.GetFCM(context)

	var body bean.OfferStoreShake
	if common.ValidateBody(context, &body) != nil {
		return
	}

	id, _ := strconv.Atoi(chainId)
	body.ChainId = int64(id)
	body.Language = language
	body.FCM = fcm
	offerShake, ce := service.OfferStoreServiceInst.CreateOfferStoreShake(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerShake)
}

func (api OfferStoreApi) RejectOfferStoreShake(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")

	offerShake, ce := service.OfferStoreServiceInst.RejectOfferStoreShake(userId, offerId, offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerShake)
}

func (api OfferStoreApi) CompleteOfferStoreShake(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")

	offerShake, ce := service.OfferStoreServiceInst.CompleteOfferStoreShake(userId, offerId, offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerShake)
}

func (api OfferStoreApi) AcceptOfferStoreShake(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")

	offerShake, ce := service.OfferStoreServiceInst.AcceptOfferStoreShake(userId, offerId, offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerShake)
}

func (api OfferStoreApi) CancelOfferStoreShake(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")

	offerShake, ce := service.OfferStoreServiceInst.CancelOfferStoreShake(userId, offerId, offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerShake)
}

func (api OfferStoreApi) ReviewOfferStore(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")
	scoreStr := context.DefaultQuery("score", "0")
	score, _ := strconv.Atoi(scoreStr)

	offerStore, ce := service.OfferStoreServiceInst.ReviewOfferStore(userId, offerId, int64(score), offerShakeId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offerStore)
}

func (api OfferStoreApi) OnChainOfferStoreTracking(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	var body bean.OfferOnChainTransaction
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.OnChainOfferStoreTracking(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) OnChainOfferStoreItemTracking(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	var body bean.OfferOnChainTransaction
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.OnChainOfferStoreItemTracking(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferStoreApi) OnChainOfferStoreShakeTracking(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")
	offerShakeId := context.Param("offerShakeId")
	var body bean.OfferOnChainTransaction
	if common.ValidateBody(context, &body) != nil {
		return
	}

	offer, ce := service.OfferStoreServiceInst.OnChainOfferStoreShakeTracking(userId, offerId, offerShakeId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}
