package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/exchangehandshake_service"
	"github.com/ninjadotorg/handshake-exchange/service"
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
		fmt.Println(offerOnChain)
		service.OfferServiceInst.ActiveOnChainOffer(offerOnChain.Offer, offerOnChain.Hid)
	}

	if len(offerOnChains) > 0 {
		lastBlock += 1
	}
	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferInitEventBlock(block)
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
		fmt.Println(offerOnChain)
		service.OfferServiceInst.ShakeOnChainOffer(offerOnChain.Offer)
	}

	if len(offerOnChains) > 0 {
		lastBlock += 1
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
		fmt.Println(offerOnChain)
		service.OfferServiceInst.RejectOnChainOffer(offerOnChain.Offer)
	}

	if len(offerOnChains) > 0 {
		lastBlock += 1
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
		fmt.Println(offerOnChain)
		service.OfferServiceInst.CompleteOnChainOffer(offerOnChain.Offer)
	}

	if len(offerOnChains) > 0 {
		lastBlock += 1
	}
	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferCompleteEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}
