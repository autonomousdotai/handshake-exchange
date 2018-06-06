package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
)

type OfferStoreService struct {
	dao     *dao.OfferStoreDao
	userDao *dao.UserDao
	miscDao *dao.MiscDao
}

func (s OfferStoreService) CreateOfferStore(userId string, offerBody bean.OfferStoreSetup) (offer bean.OfferStoreSetup, ce SimpleContextError) {
	return
}

func (s OfferStoreService) AddOfferStoreItem(userId string, offerStoreId string, offerItem bean.OfferStoreItem) (offer bean.OfferStoreItem, ce SimpleContextError) {
	return
}

func (s OfferStoreService) RemoveOfferStoreItem(userId string, offerStoreId string, currency string) (ce SimpleContextError) {
	return
}

func (s OfferStoreService) CreateOfferStoreShake(userId string, offerStoreId string, offerShakeBody bean.OfferStoreShake) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	return
}

func (s OfferStoreService) RejectOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	return
}

func (s OfferStoreService) CompleteOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	return
}

func (s OfferStoreService) UpdateOnChainOfferStore(offerStoreId string, oldStatus string, newStatus string) (offer bean.OfferStore, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStore(offerStoreId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	if !offerTO.Found {
		return
	}
	offer = offerTO.Object.(bean.OfferStore)
	if offer.Status != oldStatus {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	offer.Status = newStatus
	err := s.dao.UpdateOfferStore(offer, offer.GetChangeStatus())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	// notification.SendOfferNotification(offer)

	return
}

func (s OfferStoreService) UpdateOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string, oldStatus string, newStatus string) (offer bean.OfferStoreShake, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	if !offerTO.Found {
		return
	}
	offer = offerTO.Object.(bean.OfferStoreShake)
	if offer.Status != oldStatus {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	offer.Status = newStatus
	err := s.dao.UpdateOfferStoreShake(offerStoreId, offer, offer.GetChangeStatus())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	// notification.SendOfferNotification(offer)

	return
}

func (s OfferStoreService) ActiveOnChainOfferStore(offerStoreId string, hid int64) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainOfferStore(offerStoreId, bean.OFFER_STORE_STATUS_CREATED, bean.OFFER_STORE_STATUS_ACTIVE)
}

func (s OfferStoreService) CloseOnChainOfferStore(offerStoreId string) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainOfferStore(offerStoreId, bean.OFFER_STORE_STATUS_CLOSING, bean.OFFER_STORE_STATUS_CLOSED)
}

func (s OfferStoreService) ShakeOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, bean.OFFER_STATUS_SHAKING, bean.OFFER_STATUS_SHAKE)
}

func (s OfferStoreService) RejectOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, bean.OFFER_STATUS_REJECTING, bean.OFFER_STATUS_REJECTED)
}

func (s OfferStoreService) CompleteOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, bean.OFFER_STATUS_COMPLETING, bean.OFFER_STATUS_COMPLETED)
}
