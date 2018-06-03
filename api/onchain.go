package api

import (
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/exchangehandshake_service"
	"github.com/autonomousdotai/handshake-exchange/service"
	"github.com/gin-gonic/gin"
)

type OnChainApi struct {
}

func (api OnChainApi) UpdateOfferInit(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferInitEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetInitEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.ActiveOnChainOffer(offerOnChain.Offer, offerOnChain.Hid)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferInitEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferClose(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferCloseEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCloseEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.CloseOnChainOffer(offerOnChain.Offer)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferCloseEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferShake(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferShakeEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetShakeEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.ShakeOnChainOffer(offerOnChain.Offer)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferShakeEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferReject(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferRejectEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetRejectEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.RejectOnChainOffer(offerOnChain.Offer)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferRejectEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferComplete(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferCompleteEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCompleteEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.CompleteOnChainOffer(offerOnChain.Offer)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferCompleteEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferWithdraw(context *gin.Context) {
	client := exchangehandshake_service.ExchangeHandshakeClient{}
	to := dao.OnChainDaoInst.GetOfferWithdrawEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetWithdrawEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		service.OfferServiceInst.WithdrawOnChainOffer(offerOnChain.Offer)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferWithdrawEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}
