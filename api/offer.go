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

func extractListOfferParams(context *gin.Context) (string, string, string, string, interface{}, int) {
	offerType := context.DefaultQuery("type", "")
	currency := context.DefaultQuery("currency", "")
	status := context.DefaultQuery("status", "")
	amount := context.DefaultQuery("amount", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	return offerType, currency, status, amount, startAt, limit
}
