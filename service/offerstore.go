package service

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/blockchainio_service"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/service/notification"
	"github.com/shopspring/decimal"
	"time"
)

type OfferStoreService struct {
	dao      *dao.OfferStoreDao
	userDao  *dao.UserDao
	miscDao  *dao.MiscDao
	transDao *dao.TransactionDao
	offerDao *dao.OfferDao
}

func (s OfferStoreService) CreateOfferStore(userId string, offerSetup bean.OfferStoreSetup) (offer bean.OfferStoreSetup, ce SimpleContextError) {
	offerBody := offerSetup.Offer
	offerItemBody := offerSetup.Item

	// Check offer store exists
	// Allow to re-create if offer store exist but all item is closed
	offerTO := s.dao.GetOfferStore(userId)
	if offerTO.Found {
		offerCheck := offerTO.Object.(bean.OfferStore)
		allFalse := true
		for _, v := range offerCheck.ItemFlags {
			if v == true {
				allFalse = false
				break
			}
		}
		if !allFalse {
			ce.SetStatusKey(api_error.OfferStoreExists)
			return
		}
	}

	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}

	s.prepareOfferStore(&offerBody, &offerItemBody, profile, &ce)
	if ce.HasError() {
		return
	}

	offerNew, err := s.dao.AddOfferStore(offerBody, offerItemBody, *profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offerNew.CreatedAt = time.Now().UTC()
	offerNew.ItemSnapshots = offerBody.ItemSnapshots
	notification.SendOfferStoreNotification(offerNew, offerItemBody)

	offer.Offer = offerNew
	offer.Item = offerItemBody

	return
}

func (s OfferStoreService) GetOfferStore(userId string, offerId string) (offer bean.OfferStore, ce SimpleContextError) {
	GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}

	offerTO := s.dao.GetOfferStore(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	notFound := false
	if offerTO.Found {
		offer = offerTO.Object.(bean.OfferStore)
		allFalse := true
		for _, v := range offer.ItemFlags {
			if v == true {
				allFalse = false
				break
			}
		}
		if allFalse {
			notFound = true
		}
	} else {
		notFound = true
	}

	if notFound {
		ce.NotFound = true
	}

	return
}

func (s OfferStoreService) AddOfferStoreItem(userId string, offerId string, item bean.OfferStoreItem) (offer bean.OfferStore, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer = *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}

	s.prepareOfferStore(&offer, &item, profile, &ce)

	_, err := s.dao.AddOfferStoreItem(offer, item, *profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) RemoveOfferStoreItem(userId string, offerId string, currency string) (offer bean.OfferStore, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer = *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	item := *GetOfferStoreItem(*s.dao, offerId, currency, &ce)
	if ce.HasError() {
		return
	}

	if item.Status != bean.OFFER_STORE_ITEM_STATUS_ACTIVE && item.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSING {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Only BTC, refund the crypto
	hasSell := false
	sellAmount, _ := decimal.NewFromString(item.SellAmount)
	if sellAmount.GreaterThan(common.Zero) {
		hasSell = true
	}
	if item.Currency == bean.BTC.Code {
		// Do Refund
		if hasSell {
			description := fmt.Sprintf("Refund to userId %s due to close the offer", userId)
			response := s.sendTransaction(item.UserAddress,
				item.SellBalance, item.Currency, description, offer.UID, item.WalletProvider, &ce)
			if !ce.HasError() {
				s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
					Provider:         item.WalletProvider,
					ProviderResponse: response,
					DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE,
					DataRef:          dao.GetOfferStoreItemPath(offerId),
					UID:              userId,
					Description:      description,
					Amount:           item.SellBalance,
					Currency:         item.Currency,
				})
			}
		}
	}

	allFalse := true
	// Just for check
	offer.ItemFlags[item.Currency] = false
	for _, v := range offer.ItemFlags {
		if v == true {
			allFalse = false
			break
		}
	}
	waitOnChain := (item.Currency == bean.ETH.Code && hasSell) || offer.Status == bean.OFFER_STORE_STATUS_CLOSING
	if allFalse {
		if waitOnChain {
			// Need to wait for OnChain
			offer.Status = bean.OFFER_STORE_STATUS_CLOSING
		} else {
			offer.Status = bean.OFFER_STORE_STATUS_CLOSED
		}
	}

	if waitOnChain {
		// Only update
		item.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSING
		offer.ItemSnapshots[item.Currency] = item

		err := s.dao.UpdateOfferStoreItemClosing(offer, item)
		if ce.SetError(api_error.UpdateDataFailed, err) {
			return
		}
	} else {
		profile.ActiveOfferStores[item.Currency] = false
		offer.ItemFlags = profile.ActiveOfferStores

		// Really remove the item
		item.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSED
		offer.ItemSnapshots[item.Currency] = item

		err := s.dao.RemoveOfferStoreItem(offer, item, *profile)
		if ce.SetError(api_error.DeleteDataFailed, err) {
			return
		}
	}

	// Assign to correct flag
	offer.ItemFlags[item.Currency] = item.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSED

	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) CreateOfferStoreShake(userId string, offerId string, offerShakeBody bean.OfferStoreShake) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	if UserServiceInst.CheckOfferLocked(*profile) {
		ce.SetStatusKey(api_error.OfferActionLocked)
		return
	}

	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	if profile.UserId == offer.UID {
		ce.SetStatusKey(api_error.OfferPayMyself)
		return
	}

	item := *GetOfferStoreItem(*s.dao, offerId, offerShakeBody.Currency, &ce)
	if ce.HasError() {
		return
	}

	// Make sure shake on the valid item
	if item.Status != bean.OFFER_STORE_ITEM_STATUS_ACTIVE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}
	var balance decimal.Decimal
	amount, _ := decimal.NewFromString(offerShakeBody.Amount)
	if offerShakeBody.Currency == bean.ETH.Code {
		if amount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if offerShakeBody.Currency == bean.BTC.Code {
		if amount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}

	// Total usage from current usage shake, so it's more accuracy than need to wait real balance update
	usageBalance, err := s.getUsageBalance(offerId, offerShakeBody.Type)
	if err != nil {
		ce.SetError(api_error.GetDataFailed, err)
		return
	}
	if offerShakeBody.Type == bean.OFFER_TYPE_SELL {
		balance, _ = decimal.NewFromString(item.SellAmount)
	} else {
		balance, _ = decimal.NewFromString(item.BuyAmount)
	}
	if balance.LessThan(usageBalance.Add(amount)) {
		ce.SetStatusKey(api_error.OfferStoreNotEnoughBalance)
	}

	offerShakeBody.UID = userId
	offerShakeBody.FiatCurrency = offer.FiatCurrency
	offerShakeBody.Latitude = offer.Latitude
	offerShakeBody.Longitude = offer.Longitude

	s.setupOfferShakePrice(&offerShakeBody, &ce)
	s.setupOfferShakeAmount(&offerShakeBody, &ce)
	if ce.HasError() {
		return
	}

	// Status of shake
	if offerShakeBody.Type == bean.OFFER_TYPE_SELL {
		// SHAKE
		offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
		err = s.dao.UpdateOfferStoreShakeBalance(offer, &item, offerShakeBody, true)
		if ce.SetError(api_error.UpdateDataFailed, err) {
			return
		}
	} else {
		if offerShakeBody.Currency == bean.ETH.Code {
			offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING
		} else {
			offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKING
			s.generateSystemAddressForShake(offer, &offerShakeBody, &ce)
			if ce.HasError() {
				return
			}
		}
	}

	offerShake, err = s.dao.AddOfferStoreShake(offer, offerShakeBody)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offerShake.CreatedAt = time.Now().UTC()
	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) RejectOfferStoreShake(userId string, offerId string, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}
	item := *GetOfferStoreItem(*s.dao, offerId, offerShake.Currency, &ce)
	if ce.HasError() {
		return
	}

	if profile.UserId != offer.UID && profile.UserId != offerShake.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerShake.Type == bean.OFFER_TYPE_SELL {
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
		// REJECTED
		err := s.dao.UpdateOfferStoreShakeBalance(offer, &item, offerShake, false)
		if ce.SetError(api_error.UpdateDataFailed, err) {
			return
		}
	} else {
		if offerShake.Currency == bean.ETH.Code {
			// Only ETH
			offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTING
		} else {
			// Only BTC
			offerStoreItemTO := s.dao.GetOfferStoreItem(userId, offerShake.Currency)
			if offerStoreItemTO.HasError() {
				ce.SetStatusKey(api_error.GetDataFailed)
				return
			}
			offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)

			offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
			description := fmt.Sprintf("Refund to userId %s due to reject the offer", offerShake.UID)
			userAddress := offerShake.UserAddress

			response := s.sendTransaction(userAddress,
				offerShake.Amount, offerShake.Currency, description, offer.UID, offerStoreItem.WalletProvider, &ce)
			if !ce.HasError() {
				s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
					Provider:         offerStoreItem.WalletProvider,
					ProviderResponse: response,
					DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
					DataRef:          dao.GetOfferStoreShakeItemPath(offerId, offerShakeId),
					UID:              userId,
					Description:      description,
					Amount:           offerShake.Amount,
					Currency:         offerShake.Currency,
				})
			}
		}
	}

	transCount := s.getFailedTransCount(offer, offerShake, userId)
	err := s.dao.UpdateOfferStoreShakeReject(offer, offerShake, *profile, transCount)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	if userId == offerShake.UID {
		UserServiceInst.UpdateOfferRejectLock(userId)
	}

	offerShake.ActionUID = userId
	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) CancelOfferStoreShake(userId string, offerId string, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}

	if profile.UserId != offer.UID && profile.UserId != offerShake.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerShake.Currency == bean.ETH.Code {
		// Only ETH
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_CANCELLING
	} else {
		// Only BTC
		offerStoreItemTO := s.dao.GetOfferStoreItem(userId, offerShake.Currency)
		if offerStoreItemTO.HasError() {
			ce.SetStatusKey(api_error.GetDataFailed)
			return
		}
		offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)

		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_CANCELLED
		description := fmt.Sprintf("Refund to userId %s due to reject the offer", offerShake.UID)
		userAddress := offerShake.UserAddress
		if offerShake.Type == bean.OFFER_TYPE_BUY {
			description = fmt.Sprintf("Refund to userId %s due to reject the offer", offer.UID)
			userAddress = offerStoreItem.UserAddress

		}
		response := s.sendTransaction(userAddress,
			offerShake.Amount, offerShake.Currency, description, offer.UID, offerStoreItem.WalletProvider, &ce)
		if !ce.HasError() {
			s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
				Provider:         offerStoreItem.WalletProvider,
				ProviderResponse: response,
				DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
				DataRef:          dao.GetOfferStoreShakeItemPath(offerId, offerShakeId),
				UID:              userId,
				Description:      description,
				Amount:           offerShake.Amount,
				Currency:         offerShake.Currency,
			})
		}
	}

	err := s.dao.UpdateOfferStoreShake(offerId, offerShake, offerShake.GetChangeStatus())
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}
	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) AcceptOfferStoreShake(userId string, offerId string, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}

	if profile.UserId != offer.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerShake.Currency == bean.ETH.Code {
		// Only ETH
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKING
	} else {
		// Only BTC
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
	}

	err := s.dao.UpdateOfferStoreShake(offerId, offerShake, offerShake.GetChangeStatus())
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}
	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) CompleteOfferStoreShake(userId string, offerId string, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}

	if offerShake.Type == bean.OFFER_TYPE_SELL {
		if profile.UserId != offer.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	} else {
		if profile.UserId != offerShake.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}

	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerShake.Currency == bean.ETH.Code {
		// Only ETH
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETING
	} else {
		// Only BTC
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETED
		// Do Transfer
		s.transferCrypto(&offer, &offerShake, &ce)
	}

	transCount1, transCount2 := s.getSuccessTransCount(offer, offerShake, userId)

	err := s.dao.UpdateOfferStoreShakeComplete(offer, offerShake, *profile, transCount1, transCount2)

	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	// For onchain processing
	if offerShake.Hid == 0 {
		offerShake.Hid = offer.Hid
	}
	if offerShake.UserAddress == "" {
		offerShake.UserAddress = offer.ItemSnapshots[offerShake.Currency].UserAddress
	}
	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) UpdateOnChainInitOfferStore(offerId string, hid int64, currency string) (offer bean.OfferStore, ce SimpleContextError) {
	offer = *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	item := *GetOfferStoreItem(*s.dao, offerId, currency, &ce)
	if ce.HasError() {
		return
	}
	if item.Status != bean.OFFER_STORE_ITEM_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	if offer.Hid == 0 {
		offer.Hid = hid
	}
	item.Status = bean.OFFER_STORE_ITEM_STATUS_ACTIVE
	item.SellBalance = item.SellAmount
	if offer.Status == bean.OFFER_STORE_STATUS_CREATED {
		offer.Status = bean.OFFER_STORE_STATUS_ACTIVE
	}
	offer.ItemSnapshots[item.Currency] = item
	err := s.dao.UpdateOfferStoreItemActive(offer, item)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) UpdateOnChainCloseOfferStore(offerId string) (offer bean.OfferStore, ce SimpleContextError) {
	offer = *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	profile := *GetProfile(s.userDao, offer.UID, &ce)
	if ce.HasError() {
		return
	}

	itemTO := s.dao.GetOfferStoreItem(offerId, bean.ETH.Code)
	if !itemTO.Found {
		return
	}
	item := itemTO.Object.(bean.OfferStoreItem)
	if item.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSING {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	item.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSED
	if offer.Status == bean.OFFER_STORE_STATUS_CLOSING {
		offer.Status = bean.OFFER_STORE_STATUS_CLOSED
	}

	profile.ActiveOfferStores[item.Currency] = false
	offer.ItemSnapshots[item.Currency] = item
	err := s.dao.UpdateOfferStoreItemClosed(offer, item, profile)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) UpdateOnChainOfferStoreShake(offerId string, offerShakeId string, hid int64, oldStatus string, newStatus string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}

	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}
	if offerShake.Status != oldStatus {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	offerShake.Status = newStatus
	if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
		itemTO := s.dao.GetOfferStoreItem(offerId, offerShake.Currency)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, itemTO) {
			return
		}
		item := itemTO.Object.(bean.OfferStoreItem)
		var err error
		if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
			// SHAKE
			err = s.dao.UpdateOfferStoreShakeBalance(offer, &item, offerShake, true)
		} else if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
			// REJECTED
			err = s.dao.UpdateOfferStoreShakeBalance(offer, &item, offerShake, false)
		}
		offer.ItemSnapshots[item.Currency] = item
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
		}
	}
	if offerShake.Hid == 0 {
		offerShake.Hid = hid
	}

	err := s.dao.UpdateOfferStoreShake(offerId, offerShake, offerShake.GetChangeStatus())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) ActiveOnChainOfferStore(offerId string, hid int64) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainInitOfferStore(offerId, hid, bean.ETH.Code)
}

func (s OfferStoreService) CloseOnChainOfferStore(offerId string) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainCloseOfferStore(offerId)
}

func (s OfferStoreService) PreShakeOnChainOfferStoreShake(offerId string, offerShakeId string, hid int64) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerId, offerShakeId, hid, bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING, bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE)
}

func (s OfferStoreService) CancelOnChainOfferStoreShake(offerId string, offerShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerId, offerShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_CANCELLING, bean.OFFER_STORE_SHAKE_STATUS_CANCELLED)
}

func (s OfferStoreService) ShakeOnChainOfferStoreShake(offerId string, offerShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerId, offerShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_SHAKING, bean.OFFER_STORE_SHAKE_STATUS_SHAKE)
}

func (s OfferStoreService) RejectOnChainOfferStoreShake(offerId string, offerShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerId, offerShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_REJECTING, bean.OFFER_STORE_SHAKE_STATUS_REJECTED)
}

func (s OfferStoreService) CompleteOnChainOfferStoreShake(offerId string, offerShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerId, offerShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_COMPLETING, bean.OFFER_STORE_SHAKE_STATUS_COMPLETED)
}

func (s OfferStoreService) ActiveOffChainOfferStore(address string, amountStr string) (offer bean.OfferStore, ce SimpleContextError) {
	addressMapTO := s.offerDao.GetOfferAddress(address)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, addressMapTO) {
		return
	}
	if ce.NotFound {
		ce.SetStatusKey(api_error.ResourceNotFound)
		return
	}
	addressMap := addressMapTO.Object.(bean.OfferAddressMap)

	itemTO := s.dao.GetOfferStoreItemByPath(addressMap.OfferRef)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, itemTO) {
		return
	}
	item := itemTO.Object.(bean.OfferStoreItem)
	if item.Status != bean.OFFER_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	inputAmount, _ := decimal.NewFromString(amountStr)
	offerAmount, _ := decimal.NewFromString(item.SellAmount)

	// Check amount need to deposit
	sub := offerAmount.Sub(inputAmount)

	if sub.Equal(common.Zero) {
		// Good
		_, ce = s.UpdateOnChainInitOfferStore(addressMap.Offer, 0, bean.BTC.Code)
		if ce.HasError() {
			return
		}
	} else {
		// TODO Process to refund?
	}

	return
}

func (s OfferStoreService) PreShakeOffChainOfferStoreShake(address string, amountStr string) (offer bean.OfferStoreShake, ce SimpleContextError) {
	addressMapTO := s.offerDao.GetOfferAddress(address)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, addressMapTO) {
		return
	}
	if ce.NotFound {
		ce.SetStatusKey(api_error.ResourceNotFound)
		return
	}
	addressMap := addressMapTO.Object.(bean.OfferAddressMap)

	offerShakeTO := s.dao.GetOfferStoreShakeByPath(addressMap.OfferRef)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerShakeTO) {
		return
	}
	offer = offerShakeTO.Object.(bean.OfferStoreShake)
	if offer.Status != bean.OFFER_STATUS_SHAKING {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	inputAmount, _ := decimal.NewFromString(amountStr)
	offerAmount, _ := decimal.NewFromString(offer.Amount)
	totalAmount, _ := decimal.NewFromString(offer.TotalAmount)

	// Check amount need to deposit
	sub := decimal.NewFromFloat(1)
	if offer.Type == bean.OFFER_TYPE_BUY {
		sub = totalAmount.Sub(inputAmount)
	} else {
		sub = offerAmount.Sub(inputAmount)
	}

	// Check amount need to deposit
	sub = offerAmount.Sub(inputAmount)

	if sub.Equal(common.Zero) {
		// Good
		_, ce = s.UpdateOnChainOfferStoreShake(offer.UID, addressMap.Offer, 0, bean.OFFER_STORE_SHAKE_STATUS_SHAKING, bean.OFFER_STORE_SHAKE_STATUS_SHAKE)
		if ce.HasError() {
			return
		}
	} else {
		// TODO Process to refund?
	}

	return
}

func (s OfferStoreService) ReviewOfferStore(userId string, offerId string, score int64, offerShakeId string) (offer bean.OfferStore, ce SimpleContextError) {
	if score < 0 && score > 5 {
		ce.SetStatusKey(api_error.InvalidQueryParam)
	}

	GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	offer = *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake := *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}
	if offerShake.UID != userId {
		ce.SetStatusKey(api_error.InvalidUserToCompleteHandshake)
		return
	}
	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_COMPLETING && offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	reviewTO := s.dao.GetOfferStoreReview(offer.Id, offerShakeId)
	if reviewTO.Found {
		ce.SetStatusKey(api_error.OfferStoreExists)
		return
	}

	offer.ReviewCount += 1
	offer.Review += score

	err := s.dao.AddOfferStoreReview(offer, bean.OfferStoreReview{
		Id:    offerShakeId,
		UID:   userId,
		Score: score,
	})
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	notification.SendOfferStoreNotification(offer, bean.OfferStoreItem{})

	return
}

func (s OfferStoreService) GetQuote(quoteType string, amountStr string, currency string, fiatCurrency string) (price decimal.Decimal, fiatPrice decimal.Decimal,
	fiatAmount decimal.Decimal, err error) {
	amount, numberErr := decimal.NewFromString(amountStr)
	to := dao.MiscDaoInst.GetCurrencyRateFromCache(bean.USD.Code, fiatCurrency)
	if numberErr != nil {
		err = numberErr
	}
	rate := to.Object.(bean.CurrencyRate)
	rateNumber := decimal.NewFromFloat(rate.Rate)
	tmpAmount := amount.Mul(rateNumber)

	if quoteType == "buy" {
		resp, errResp := coinbase_service.GetBuyPrice(currency)
		err = errResp
		if err != nil {
			return
		}
		price, _ = decimal.NewFromString(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else if quoteType == "sell" {
		resp, errResp := coinbase_service.GetSellPrice(currency)
		err = errResp
		if err != nil {
			return
		}
		price, _ := decimal.NewFromString(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else {
		err = errors.New(api_error.InvalidQueryParam)
	}

	return
}

func (s OfferStoreService) SyncOfferStoreToSolr(offerId string) (offer bean.OfferStore, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStore(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.OfferStore)
	solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, bean.OfferStoreItem{}))

	return
}

func (s OfferStoreService) SyncOfferStoreShakeToSolr(offerId, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	offerStoreTO := s.dao.GetOfferStore(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerStoreTO) {
		return
	}
	offer := offerStoreTO.Object.(bean.OfferStore)
	offerShakeTO := s.dao.GetOfferStoreShake(offerId, offerShakeId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerShakeTO) {
		return
	}
	offerShake = offerShakeTO.Object.(bean.OfferStoreShake)
	solr_service.UpdateObject(bean.NewSolrFromOfferStoreShake(offerShake, offer))

	return
}

func (s OfferStoreService) prepareOfferStore(offer *bean.OfferStore, item *bean.OfferStoreItem, profile *bean.Profile, ce *SimpleContextError) {
	currencyInst := bean.CurrencyMapping[item.Currency]
	if currencyInst.Code == "" {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}

	if profile.ActiveOfferStores == nil {
		profile.ActiveOfferStores = make(map[string]bool)
	}
	if check, ok := profile.ActiveOfferStores[currencyInst.Code]; ok && check {
		// Has Key and already had setup
		ce.SetStatusKey(api_error.TooManyOffer)
		return
	}

	profile.ActiveOfferStores[currencyInst.Code] = true
	offer.ItemFlags = profile.ActiveOfferStores
	offer.UID = profile.UserId

	s.checkOfferStoreItemAmount(item, ce)
	if ce.HasError() {
		return
	}

	s.generateSystemAddress(*offer, item, ce)

	sellAmount, _ := decimal.NewFromString(item.SellAmount)
	if sellAmount.Equal(common.Zero) {
		// Only the case that shop doesn't sell, so don't need to wait to active
		item.Status = bean.OFFER_STORE_ITEM_STATUS_ACTIVE
		// So active the store as well
		offer.Status = bean.OFFER_STORE_STATUS_ACTIVE
	} else {
		item.Status = bean.OFFER_STORE_ITEM_STATUS_CREATED
	}
	if offer.Status != bean.OFFER_STORE_STATUS_ACTIVE {
		offer.Status = bean.OFFER_STORE_STATUS_CREATED
	}

	minAmount := bean.MIN_ETH
	if item.Currency == bean.BTC.Code {
		minAmount = bean.MIN_BTC
	}
	item.BuyBalance = item.BuyAmount
	item.BuyAmountMin = minAmount.String()
	item.SellBalance = "0"
	item.SellAmountMin = minAmount.String()

	if offer.ItemSnapshots == nil {
		offer.ItemSnapshots = make(map[string]bean.OfferStoreItem)
	}
	offer.ItemSnapshots[item.Currency] = *item
}

func (s OfferStoreService) checkOfferStoreItemAmount(item *bean.OfferStoreItem, ce *SimpleContextError) {
	// Minimum amount
	sellAmount, errFmt := decimal.NewFromString(item.SellAmount)
	if ce.SetError(api_error.InvalidRequestBody, errFmt) {
		return
	}
	if item.Currency == bean.ETH.Code {
		if sellAmount.GreaterThan(common.Zero) && sellAmount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.Currency == bean.BTC.Code {
		if sellAmount.GreaterThan(common.Zero) && sellAmount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.SellPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(item.SellPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		item.SellPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		item.SellPercentage = "0"
	}

	buyAmount, errFmt := decimal.NewFromString(item.BuyAmount)
	if ce.SetError(api_error.InvalidRequestBody, errFmt) {
		return
	}
	if item.Currency == bean.ETH.Code {
		if buyAmount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.Currency == bean.BTC.Code {
		if buyAmount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.BuyPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(item.BuyPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		item.BuyPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		item.BuyPercentage = "0"
	}
}

func (s OfferStoreService) generateSystemAddress(offer bean.OfferStore, item *bean.OfferStoreItem, ce *SimpleContextError) {
	// Only BTC need to generate address to transfer in
	if item.Currency == bean.BTC.Code {
		systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_BTC_WALLET)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
			return
		}
		systemConfig := systemConfigTO.Object.(bean.SystemConfig)
		item.WalletProvider = systemConfig.Value
		if systemConfig.Value == bean.BTC_WALLET_COINBASE {
			addressResponse, err := coinbase_service.GenerateAddress(item.Currency)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			item.SystemAddress = addressResponse.Data.Address

		} else if systemConfig.Value == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			address, err := client.GenerateAddress(offer.Id)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			item.SystemAddress = address
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}
}

// TODO remove func duplicate
func (s OfferStoreService) generateSystemAddressForShake(offer bean.OfferStore, offerShake *bean.OfferStoreShake, ce *SimpleContextError) {
	// Only BTC need to generate address to transfer in
	if offerShake.Currency == bean.BTC.Code {
		systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_BTC_WALLET)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
			return
		}
		systemConfig := systemConfigTO.Object.(bean.SystemConfig)
		offerShake.WalletProvider = systemConfig.Value
		if systemConfig.Value == bean.BTC_WALLET_COINBASE {
			addressResponse, err := coinbase_service.GenerateAddress(offerShake.Currency)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offerShake.SystemAddress = addressResponse.Data.Address

		} else if systemConfig.Value == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			address, err := client.GenerateAddress(offer.Id)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offerShake.SystemAddress = address
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}
}

func (s OfferStoreService) transferCrypto(offer *bean.OfferStore, offerShake *bean.OfferStoreShake, ce *SimpleContextError) {
	offerStoreItemTO := s.dao.GetOfferStoreItem(offer.UID, offerShake.Currency)
	if offerStoreItemTO.HasError() {
		ce.SetStatusKey(api_error.GetDataFailed)
		return
	}
	offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)
	userAddress := offerShake.UserAddress
	actionUID := offerShake.UID

	if offerShake.Type == bean.OFFER_TYPE_BUY {
		userAddress = offerStoreItem.UserAddress
		actionUID = offer.UID
	}

	if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		if userAddress != "" {
			//Transfer
			description := fmt.Sprintf("Transfer to userId %s offerShakeId %s status %s", actionUID, offerShake.Id, offerShake.Status)

			var response1 interface{}
			// var response2 interface{}
			var userId string
			transferAmount := offerShake.Amount
			if offerShake.Type == bean.OFFER_TYPE_BUY {
				// Amount = 1, transfer 0.09 (if fee = 1%)
				response1 = s.sendTransaction(offerStoreItem.UserAddress, offerShake.TotalAmount, offerShake.Currency, description, offerShake.Id, offerStoreItem.WalletProvider, ce)
				transferAmount = offerShake.TotalAmount
				userId = offer.UID
			} else {
				// Amount = 1, transfer 1
				response1 = s.sendTransaction(offerShake.UserAddress, offerShake.Amount, offerShake.Currency, description, offerShake.Id, offerStoreItem.WalletProvider, ce)
				userId = offerShake.UID
			}
			s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
				Provider:         offerStoreItem.WalletProvider,
				ProviderResponse: response1,
				DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
				DataRef:          dao.GetOfferStoreShakeItemPath(offer.Id, offerShake.Id),
				UID:              userId,
				Description:      description,
				Amount:           transferAmount,
				Currency:         offerShake.Currency,
			})

			// Transfer reward
			//if offerStoreItem.RewardAddress != "" {
			//	rewardDescription := fmt.Sprintf("Transfer reward to userId %s offerId %s", offerStore.UID, offerStoreShake.Id)
			//	response2 = s.sendTransaction(offerStoreItem.RewardAddress, offerStoreShake.Reward, offerStoreItem.Currency, rewardDescription,
			//		fmt.Sprintf("%s_reward", offerStoreShake.Id), offerStoreItem.WalletProvider, ce)
			//
			//	s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
			//		Provider:         offerStoreItem.WalletProvider,
			//		ProviderResponse: response2,
			//		DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
			//		DataRef:          dao.GetOfferStoreShakeItemPath(offerStore.Id, offerStoreShake.Id),
			//		UID:              offerStore.UID,
			//		Description:      description,
			//		Amount:           offerStoreShake.Amount,
			//		Currency:         offerStoreShake.Currency,
			//	})
			//}
			// Just logging the error, don't throw it
			//if ce.HasError() {
			//	return
			//}
		} else {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}
}

func (s OfferStoreService) sendTransaction(address string, amountStr string, currency string, description string, withdrawId string,
	walletProvider string, ce *SimpleContextError) interface{} {
	// Only BTC
	if currency == bean.BTC.Code {

		if walletProvider == bean.BTC_WALLET_COINBASE {
			response, err := coinbase_service.SendTransaction(address, amountStr, currency, description, withdrawId)
			fmt.Println(response)
			fmt.Println(err)
			if ce.SetError(api_error.ExternalApiFailed, err) {
				return ""
			}
			return response
		} else if walletProvider == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			amount, _ := decimal.NewFromString(amountStr)
			hashTx, err := client.SendTransaction(address, amount)
			if ce.SetError(api_error.ExternalApiFailed, err) {
				return ""
			}
			return hashTx
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}

	return ""
}

func (s OfferStoreService) getSuccessTransCount(offer bean.OfferStore, offerShake bean.OfferStoreShake, actionUID string) (transCount1 bean.TransactionCount, transCount2 bean.TransactionCount) {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, "ALL")
	if !transCountTO.HasError() && transCountTO.Found {
		transCount1 = transCountTO.Object.(bean.TransactionCount)
	}
	transCount1.Currency = "ALL"
	transCount1.Success += 1

	transCountTO = s.transDao.GetTransactionCount(offerShake.UID, offerShake.Currency)
	if !transCountTO.HasError() && transCountTO.Found {
		transCount2 = transCountTO.Object.(bean.TransactionCount)
	}
	transCount2.Currency = offerShake.Currency
	transCount2.Success += 1

	return
}

func (s OfferStoreService) getFailedTransCount(offer bean.OfferStore, offerShake bean.OfferStoreShake, actionUID string) (transCount bean.TransactionCount) {
	if actionUID == offer.UID {
		transCountTO := s.transDao.GetTransactionCount(offer.UID, "all")
		if !transCountTO.HasError() && transCountTO.Found {
			transCount = transCountTO.Object.(bean.TransactionCount)
		}
		transCount.Currency = "ALL"
		transCount.Failed += 1
	} else {
		transCountTO := s.transDao.GetTransactionCount(offerShake.UID, offerShake.Currency)
		if !transCountTO.HasError() && transCountTO.Found {
			transCount = transCountTO.Object.(bean.TransactionCount)
		}
		transCount.Currency = "ALL"
		transCount.Failed += 1
	}

	return
}

func (s OfferStoreService) setupOfferShakePrice(offer *bean.OfferStoreShake, ce *SimpleContextError) {
	_, fiatPrice, fiatAmount, err := s.GetQuote(offer.Type, offer.Amount, offer.Currency, offer.FiatCurrency)
	if ce.SetError(api_error.GetDataFailed, err) {
		return
	}

	offer.Price = fiatPrice.Round(2).String()
	offer.FiatAmount = fiatAmount.Round(2).String()
}

func (s OfferStoreService) setupOfferShakeAmount(offerShake *bean.OfferStoreShake, ce *SimpleContextError) {
	exchFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_EXCHANGE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, exchFeeTO) {
		return
	}
	exchFeeObj := exchFeeTO.Object.(bean.SystemFee)
	exchCommTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_EXCHANGE_COMMISSION)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, exchCommTO) {
		return
	}
	exchCommObj := exchCommTO.Object.(bean.SystemFee)

	exchFee := decimal.NewFromFloat(exchFeeObj.Value).Round(6)
	exchComm := decimal.NewFromFloat(exchCommObj.Value).Round(6)
	amount, _ := decimal.NewFromString(offerShake.Amount)
	fee := amount.Mul(exchFee)
	// reward := amount.Mul(exchComm)
	// For now
	reward := decimal.NewFromFloat(0)

	offerShake.FeePercentage = exchFee.String()
	offerShake.RewardPercentage = exchComm.String()
	offerShake.Fee = fee.String()
	offerShake.Reward = reward.String()
	if offerShake.Type == bean.OFFER_TYPE_SELL {
		offerShake.TotalAmount = amount.Sub(fee.Add(reward)).String()
	} else if offerShake.Type == bean.OFFER_TYPE_BUY {
		offerShake.TotalAmount = amount.Add(fee.Add(reward)).String()
	}
}

func (s OfferStoreService) getOfferProfile(offer bean.OfferStore, offerShake bean.OfferStoreShake, profile bean.Profile, ce *SimpleContextError) (offerProfile bean.Profile) {
	if profile.UserId == offer.UID {
		offerProfile = profile
	} else {
		offerProfileTO := s.userDao.GetProfile(offerShake.UID)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, offerProfileTO) {
			return
		}
		offerProfile = offerProfileTO.Object.(bean.Profile)
	}

	return
}

func (s OfferStoreService) getUsageBalance(offerId string, offerType string) (decimal.Decimal, error) {
	offerShakes, err := s.dao.ListUsageOfferStoreShake(offerId, offerType)
	usage := common.Zero
	if err == nil {
		for _, offerShake := range offerShakes {
			if offerType == bean.OFFER_TYPE_SELL {
				amount, _ := decimal.NewFromString(offerShake.TotalAmount)
				usage = usage.Add(amount)
			} else {
				amount, _ := decimal.NewFromString(offerShake.Amount)
				usage = usage.Add(amount)
			}
		}
	}
	return usage, err
}
