package api

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/service"
	"github.com/gin-gonic/gin"
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

func (api CreditCardApi) PayInstantOffer(context *gin.Context) {
	userId := common.GetUserId(context)

	var body bean.InstantOffer
	if common.ValidateBody(context, &body) != nil {
		return
	}

	body.UID = userId
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
