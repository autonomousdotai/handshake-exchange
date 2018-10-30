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
	refUserId := context.DefaultQuery("ref", "")

	userId := strconv.Itoa(body.Id)
	err := service.UserServiceInst.AddProfile(bean.Profile{UserId: userId, ReferralUser: refUserId})
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	service.ReferralServiceInst.AddReferral(userId)
	if refUserId != "" {
		service.ReferralServiceInst.AddReferralRecord(refUserId, userId)
	}

	bean.SuccessResponse(context, body)
}

func (api ProfileApi) UpdateProfile(context *gin.Context) {
	var body bean.ProfileRequest
	if common.ValidateBody(context, &body) != nil {
		return
	}
	userId := strconv.Itoa(body.Id)
	alias := context.DefaultQuery("alias", "")

	to := dao.OfferStoreDaoInst.GetOfferStore(userId)
	if to.Found {
		offer := to.Object.(bean.OfferStore)
		offer.Username = alias
		dao.OfferStoreDaoInst.UpdateOfferStore(offer, map[string]interface{}{
			"username": offer.Username,
		})
		solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, bean.OfferStoreItem{}))
	}
	if to.Error != nil {
		if to.ContextValidate(context) {
			return
		}
	}

	bean.SuccessResponse(context, true)
}

func (api ProfileApi) UpdateProfileOffline(context *gin.Context) {
	userId := common.GetUserId(context)
	offline := context.Param("offline")

	to := dao.OfferStoreDaoInst.GetOfferStore(userId)
	if to.Found {
		offer := to.Object.(bean.OfferStore)
		offer.Offline = offline
		solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, bean.OfferStoreItem{}))
	}
	if to.Error != nil {
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

	bean.SuccessPagingResponse(context, to.Objects, to.CanMove, to.Page, 0)
}

func (api ProfileApi) GetCCLimit(context *gin.Context) {
	userId := common.GetUserId(context)

	userCCLimit, ce := service.UserServiceInst.GetCCLimitLevel(userId)
	if ce.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, userCCLimit)
}

func (api ProfileApi) ListTransactionCounts(context *gin.Context) {
	userId := common.GetUserId(context)
	to := dao.TransactionDaoInst.ListTransactionCounts(userId)
	if to.ContextValidate(context) {
		return
	}
	objs := map[string]bean.TransactionCount{}
	if to.Objects != nil {
		for _, item := range to.Objects {
			countItem := item.(bean.TransactionCount)
			if countItem.Currency != "ALL" {
				objs[countItem.Currency] = countItem
			}

		}
	}

	bean.SuccessResponse(context, objs)
}

func (api ProfileApi) ListReferralSummary(context *gin.Context) {
	userId := common.GetUserId(context)
	to := dao.ReferralDaoInst.ListReferralSummary(userId)
	if to.ContextValidate(context) {
		return
	}

	bean.SuccessResponse(context, to.Objects)
}
