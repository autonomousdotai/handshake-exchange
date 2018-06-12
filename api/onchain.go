package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/exchangehandshake_service"
	"github.com/ninjadotorg/handshake-exchange/integration/exchangehandshakeshop_service"
	"github.com/ninjadotorg/handshake-exchange/service"
	"os"
	"strconv"
	"strings"
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

func (api OnChainApi) StartOnChainOfferBlock(context *gin.Context) {
	blockStr := os.Getenv("ETH_EXCHANGE_HANDSHAKE_BLOCK")
	blockInt, _ := strconv.Atoi(blockStr)
	block := int64(blockInt)

	dao.OnChainDaoInst.UpdateOfferInitEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferShakeEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferCloseEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferRejectEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferCompleteEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferWithdrawEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreInit(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreInitEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetInitOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		service.OfferStoreServiceInst.ActiveOnChainOfferStore(offerOnChain.Offer, offerOnChain.Hid)
	}
	if len(offerOnChains) > 0 {
		lastBlock += 1
	}
	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferStoreInitEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreClose(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreCloseEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCloseOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		service.OfferStoreServiceInst.CloseOnChainOfferStore(offerOnChain.Offer)
	}
	if len(offerOnChains) > 0 {
		lastBlock += 1
	}
	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferStoreCloseEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStorePreShake(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStorePreShakeEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetPreShakeOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.PreShakeOnChainOfferStoreShake(parts[0], parts[1], offerOnChain.Hid)
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferStorePreShakeEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreCancel(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreShakeEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCancelOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.CancelOnChainOfferStoreShake(parts[0], parts[1])
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferShakeEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreShake(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreShakeEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetShakeOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.ShakeOnChainOfferStoreShake(parts[0], parts[1])
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferShakeEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreReject(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreRejectEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetRejectOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.RejectOnChainOfferStoreShake(parts[0], parts[1])
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferRejectEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreComplete(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreCompleteEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCompleteOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.CompleteOnChainOfferStoreShake(parts[0], parts[1])
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferStoreCompleteEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) UpdateOfferStoreCompleteUser(context *gin.Context) {
	client := exchangehandshakeshop_service.ExchangeHandshakeShopClient{}
	to := dao.OnChainDaoInst.GetOfferStoreCompleteUserEventBlock()
	if to.ContextValidate(context) {
		return
	}
	block := to.Object.(bean.OfferEventBlock)

	offerOnChains, lastBlock, err := client.GetCompleteUserOfferStoreEvent(uint64(block.LastBlock))
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}
	for _, offerOnChain := range offerOnChains {
		fmt.Println(offerOnChain)
		parts := strings.Split(offerOnChain.Offer, "-")
		service.OfferStoreServiceInst.CompleteOnChainOfferStoreShake(parts[0], parts[1])
	}

	block.LastBlock = int64(lastBlock)
	err = dao.OnChainDaoInst.UpdateOfferStoreCompleteUserEventBlock(block)
	if api_error.PropagateErrorAndAbort(context, api_error.UpdateDataFailed, err) != nil {
		return
	}

	bean.SuccessResponse(context, true)
}

func (api OnChainApi) StartOnChainOfferStoreBlock(context *gin.Context) {
	blockStr := os.Getenv("ETH_EXCHANGE_HANDSHAKE_OFFER_STORE_BLOCK")
	blockInt, _ := strconv.Atoi(blockStr)
	block := int64(blockInt)

	dao.OnChainDaoInst.UpdateOfferStoreInitEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreCloseEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStorePreShakeEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreCancelEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreShakeEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreRejectEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreCompleteEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})
	dao.OnChainDaoInst.UpdateOfferStoreCompleteUserEventBlock(bean.OfferEventBlock{
		LastBlock: block,
	})

	bean.SuccessResponse(context, true)
}
