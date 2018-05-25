package service

import (
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/coinbase_service"
)

type OfferService struct {
	dao *dao.OfferDao
}

func (s OfferService) CreateOffer(userId string, offerBody bean.Offer) (offer bean.Offer, ce SimpleContextError) {
	// Offer type
	if offerBody.Type != bean.OFFER_TYPE_BUY && offerBody.Type != bean.OFFER_TYPE_SELL {
		ce.SetStatusKey(api_error.UnsupportedOfferType)
		return
	}
	if offerBody.Type == bean.OFFER_TYPE_BUY {
		// Need to set address to receive crypto
		if offerBody.UserAddress == "" {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}
	currencyInst := bean.CurrencyMapping[offerBody.Currency]
	if currencyInst.Code == "" {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}
	profileTO := dao.UserDaoInst.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offerBody.UID = userId
	offerBody.Status = bean.OFFER_STATUS_CREATED

	addressResponse, err := coinbase_service.GenerateAddress(currencyInst.Code)
	if err != nil {
		ce.SetError(api_error.ExternalApiFailed, err)
	}
	offerBody.SystemAddress = addressResponse.Data.Address

	offer, err = s.dao.AddOffer(offerBody)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	return
}
