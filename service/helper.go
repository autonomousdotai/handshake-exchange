package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
)

func GetProfile(dao dao.UserDaoInterface, userId string, ce *SimpleContextError) (profile *bean.Profile) {
	to := dao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	if to.Found {
		obj := to.Object.(bean.Profile)
		profile = &obj
	} else {
		ce.SetStatusKey(api_error.ResourceNotFound)
	}

	return
}

func GetOffer(dao dao.OfferDao, offerId string, ce *SimpleContextError) (offer *bean.Offer) {
	to := dao.GetOffer(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	if to.Found {
		obj := to.Object.(bean.Offer)
		offer = &obj
	} else {
		ce.SetStatusKey(api_error.ResourceNotFound)
	}

	return
}

func GetOfferStore(dao dao.OfferStoreDao, offerId string, ce *SimpleContextError) (offer *bean.OfferStore) {
	to := dao.GetOfferStore(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	if to.Found {
		obj := to.Object.(bean.OfferStore)
		offer = &obj
	} else {
		ce.SetStatusKey(api_error.ResourceNotFound)
	}

	return
}

func GetOfferStoreItem(dao dao.OfferStoreDao, offerId string, currency string, ce *SimpleContextError) (offer *bean.OfferStoreItem) {
	to := dao.GetOfferStoreItem(offerId, currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	if to.Found {
		obj := to.Object.(bean.OfferStoreItem)
		offer = &obj
	} else {
		ce.SetStatusKey(api_error.ResourceNotFound)
	}

	return
}

func GetOfferStoreShake(dao dao.OfferStoreDao, offerId string, offerShakeId string, ce *SimpleContextError) (offer *bean.OfferStoreShake) {
	to := dao.GetOfferStoreShake(offerId, offerShakeId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	if to.Found {
		obj := to.Object.(bean.OfferStoreShake)
		offer = &obj
	} else {
		ce.SetStatusKey(api_error.ResourceNotFound)
	}

	return
}
