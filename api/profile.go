package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/service"
	"strconv"
)

type ProfileApi struct {
}

func (api ProfileApi) AddProfile(context *gin.Context) {
	var body bean.ProfileRequest
	if common.ValidateBody(context, &body) != nil {
		return
	}

	err := service.UserServiceInst.AddProfile(bean.Profile{UserId: strconv.Itoa(body.Id)})
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, body)
}

func (api ProfileApi) UpdateProfileOffline(context *gin.Context) {
	userId := common.GetUserId(context)
	offline := context.Param("offline")

	to := dao.OfferStoreDaoInst.GetOfferStore(userId)
	if to.Found {
		offer := to.Object.(bean.OfferStore)
		offer.Offline = offline
		solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer))
	}
	if to.HasError() {
		if to.ContextValidate(context) {
			return
		}
	}

	bean.SuccessResponse(context, true)
}

func (api ProfileApi) GetProfile(context *gin.Context) {
	userId := common.GetUserId(context)

	to := dao.UserDaoInst.GetProfile(userId)
	if to.ContextValidate(context) {
		return
	}
	profile := to.Object.(bean.Profile)
	if profile.CreditCard.Token != "" {
		profile.CreditCard.Token = "true"
	}

	bean.SuccessResponse(context, profile)
}

func (api ProfileApi) ListTransactions(context *gin.Context) {
	userId := common.GetUserId(context)

	transType := context.DefaultQuery("trans_type", "")
	currency := context.DefaultQuery("currency", "")
	startAt, limit := common.ExtractTimePagingParams(context)

	to := dao.TransactionDaoInst.ListTransactions(userId, transType, currency, limit, startAt)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page)
}

func (api ProfileApi) GetCCLimit(context *gin.Context) {
	userId := common.GetUserId(context)

	userCCLimit, ce := service.UserServiceInst.GetCCLimitLevel(userId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, userCCLimit)
}
