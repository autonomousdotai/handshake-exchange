package service

import (
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/coinbase_service"
	"github.com/shopspring/decimal"
)

type OfferService struct {
	dao     *dao.OfferDao
	userDao *dao.UserDao
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
		offerBody.Status = bean.OFFER_STATUS_CREATED
	} else {
		if offerBody.RefundAddress == "" {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		offerBody.Status = bean.OFFER_STATUS_ACTIVE
	}

	currencyInst := bean.CurrencyMapping[offerBody.Currency]
	if currencyInst.Code == "" {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)
	offerBody.UID = userId

	if profile.ActiveOffers[currencyInst.Code] {
		ce.SetStatusKey(api_error.TooManyOffer)
		return
	}
	profile.ActiveOffers[currencyInst.Code] = true

	// Only BTC need to generate address to transfer in
	if offerBody.Currency == bean.BTC.Code {
		addressResponse, err := coinbase_service.GenerateAddress(currencyInst.Code)
		if err != nil {
			ce.SetError(api_error.ExternalApiFailed, err)
		}
		offerBody.SystemAddress = addressResponse.Data.Address
	}

	offer, err := s.dao.AddOffer(offerBody, profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}
	// TODO Add to Algolia

	return
}

func (s OfferService) CloseOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)
	if offer.Status != bean.OFFER_STATUS_ACTIVE && offer.Status != bean.OFFER_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	// Only BTC can refund
	if offer.Type == bean.OFFER_TYPE_SELL && offer.Status == bean.OFFER_STATUS_ACTIVE && offer.Currency == bean.BTC.Code {
		if offer.RefundAddress != "" {
			//Refund
			description := fmt.Sprintf("Refund to userId %s due to close the handshake", userId)
			coinbaseResponse, err := coinbase_service.SendTransaction(offer.RefundAddress, offer.Amount, offer.Currency, description, offer.Id)
			if ce.SetError(api_error.ExternalApiFailed, err) {
				return
			}
			offer.Provider = bean.OFFER_PROVIDER_COINBASE
			offer.ProviderData = coinbaseResponse
		} else {
			ce.SetStatusKey(api_error.InvalidRequestBody)
		}
	}

	offer.Status = bean.OFFER_STATUS_CLOSED
	err := s.dao.UpdateOffer(offer, offer.GetUpdateOfferClose())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	// TODO Remove from Algolia

	return
}

func (s OfferService) ShakeOffer(userId string, offerId string, body bean.OfferShakeRequest) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if profile.UserId == offer.UID {
		ce.SetStatusKey(api_error.OfferPayMyself)
		return
	}

	if offer.Status != bean.OFFER_STATUS_ACTIVE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offer.ToUID = userId
	if offer.Type == bean.OFFER_TYPE_SELL {
		// Only BTC needs to check
		if body.Address == "" && offer.Currency == bean.BTC.Code {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		_, numberErr := decimal.NewFromString(body.FiatAmount)
		if ce.SetError(api_error.InvalidNumber, numberErr) {
			return
		}
		offer.UserAddress = body.Address
		offer.Status = bean.OFFER_STATUS_SHAKE
	} else {
		// Only BTC needs to check
		if body.Address == "" && offer.Currency == bean.BTC.Code {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		offer.RefundAddress = body.Address
		offer.Status = bean.OFFER_STATUS_SHAKING
	}
	offer.FiatAmount = body.FiatAmount

	err := s.dao.UpdateOffer(offer, offer.GetUpdateOfferShaking())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	return
}

func (s OfferService) RejectShakeOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if profile.UserId != offer.UID && profile.UserId != offer.ToUID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if offer.Status != bean.OFFER_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offer.Status = bean.OFFER_STATUS_REJECTED
	if offer.Currency == bean.BTC.Code {
		// Only BTC can refund
		if offer.RefundAddress != "" {
			//Refund
			description := fmt.Sprintf("Refund to userId %s due to reject the handshake", userId)
			coinbaseResponse, err := coinbase_service.SendTransaction(offer.RefundAddress, offer.Amount, offer.Currency, description, offer.Id)
			if ce.SetError(api_error.ExternalApiFailed, err) {
				return
			}
			offer.Provider = bean.OFFER_PROVIDER_COINBASE
			offer.ProviderData = coinbaseResponse

			err = s.dao.UpdateOffer(offer, offer.GetUpdateOfferReject())
			if ce.SetError(api_error.UpdateDataFailed, err) {
				return
			}
		} else {
			ce.SetStatusKey(api_error.InvalidRequestBody)
		}
	}

	return
}

func (s OfferService) CompleteShakeOffer(userId string, offerId string) (offer bean.Offer, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerTO := s.dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.Offer)

	if profile.UserId != offer.UID && profile.UserId != offer.ToUID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}

	if offer.Status != bean.OFFER_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offer.Status = bean.OFFER_STATUS_COMPLETING
	// Only BTC can transfer
	var externalId string
	if offer.Currency == bean.BTC.Code {
		if offer.UserAddress != "" {
			//Transfer
			description := fmt.Sprintf("Transfer to userId %s due to complete the handshake", userId)
			coinbaseResponse, err := coinbase_service.SendTransaction(offer.UserAddress, offer.Amount, offer.Currency, description, offer.Id)
			if ce.SetError(api_error.ExternalApiFailed, err) {
				return
			}
			offer.Provider = bean.OFFER_PROVIDER_COINBASE
			offer.ProviderData = coinbaseResponse
			externalId = coinbaseResponse.Id
		} else {
			ce.SetStatusKey(api_error.InvalidRequestBody)
		}
	}

	err := s.dao.UpdateOfferCompleting(offer, externalId)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	return
}
