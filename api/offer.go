package api

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/service"
	"github.com/gin-gonic/gin"
)

type OfferApi struct {
}

func (api OfferApi) CreateOffer(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.Offer
	if common.ValidateBody(context, &body) != nil {
		return
	}

	// status: buy:active or sell:created
	offer, ce := service.OfferServiceInst.CreateOffer(userId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferApi) ListOffers(context *gin.Context) {
	userId := common.GetUserId(context)
	offerType, currency, status, _, startAt, limit := extractListOfferParams(context)

	to := dao.OfferDaoInst.ListOffers(userId, offerType, currency, status, limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page)
}

func (api OfferApi) CloseOffer(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	// status: created->closed, active->closed
	offer, ce := service.OfferServiceInst.CloseOffer(userId, offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferApi) ShakeOffer(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	var body bean.OfferShakeRequest
	if common.ValidateBody(context, &body) != nil {
		return
	}

	// status: active->shaking
	offer, ce := service.OfferServiceInst.ShakeOffer(userId, offerId, body)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferApi) RejectShakeOffer(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	// status: shake->closed
	offer, ce := service.OfferServiceInst.RejectShakeOffer(userId, offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func (api OfferApi) CompleteShakeOffer(context *gin.Context) {
	userId := common.GetUserId(context)
	offerId := context.Param("offerId")

	// status: shake->completing
	offer, ce := service.OfferServiceInst.CompleteShakeOffer(userId, offerId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, offer)
}

func extractListOfferParams(context *gin.Context) (string, string, string, string, interface{}, int) {
	offerType := context.DefaultQuery("type", "")
	currency := context.DefaultQuery("currency", "")
	status := context.DefaultQuery("status", "")
	amount := context.DefaultQuery("amount", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	return offerType, currency, status, amount, startAt, limit
}
