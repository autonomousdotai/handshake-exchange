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
}

func (s OfferStoreService) CreateOfferStore(userId string, offerSetup bean.OfferStoreSetup) (offer bean.OfferStoreSetup, ce SimpleContextError) {
	offerBody := offerSetup.Offer
	offerItemBody := offerSetup.Item

	// Check offer store exists
	// Allow to re-create if offer store exist but all item is closed
	offerStoreTO := s.dao.GetOfferStore(userId)
	if offerStoreTO.Found {
		offerCheck := offerStoreTO.Object.(bean.OfferStore)
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

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	s.prepareOfferStore(&offerBody, &offerItemBody, &profile, &ce)
	if ce.HasError() {
		return
	}

	offerNew, err := s.dao.AddOfferStore(offerBody, offerItemBody, profile)
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
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
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

func (s OfferStoreService) AddOfferStoreItem(userId string, offerStoreId string, offerItem bean.OfferStoreItem) (offerStore bean.OfferStore, ce SimpleContextError) {
	// Check offer store exists
	offerStoreTO := s.dao.GetOfferStore(userId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore = offerStoreTO.Object.(bean.OfferStore)

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	s.prepareOfferStore(&offerStore, &offerItem, &profile, &ce)

	_, err := s.dao.AddOfferStoreItem(offerStore, offerItem, profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offerStore, offerItem)

	return
}

func (s OfferStoreService) RemoveOfferStoreItem(userId string, offerStoreId string, currency string) (offerStore bean.OfferStore, ce SimpleContextError) {
	// Check offer store exists
	offerStoreTO := s.dao.GetOfferStore(userId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore = offerStoreTO.Object.(bean.OfferStore)

	// Check offer item exists
	offerStoreItemTO := s.dao.GetOfferStoreItem(userId, currency)
	if !offerStoreItemTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}

	if offerStoreItem.Status != bean.OFFER_STORE_ITEM_STATUS_ACTIVE && offerStoreItem.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSING {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Only BTC, refund the crypto
	hasSell := false
	sellAmount, _ := decimal.NewFromString(offerStoreItem.SellAmount)
	if sellAmount.GreaterThan(common.Zero) {
		hasSell = true
	}
	if offerStoreItem.Currency == bean.BTC.Code {
		// Do Refund
		if hasSell {
			description := fmt.Sprintf("Refund to userId %s due to close the offer", userId)
			response := s.sendTransaction(offerStoreItem.UserAddress,
				offerStoreItem.SellBalance, offerStoreItem.Currency, description, offerStore.UID, offerStoreItem.WalletProvider, &ce)
			if !ce.HasError() {
				s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
					Provider:         offerStoreItem.WalletProvider,
					ProviderResponse: response,
					DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE,
					DataRef:          dao.GetOfferStoreItemPath(offerStoreId),
					UID:              userId,
					Description:      description,
					Amount:           offerStoreItem.SellBalance,
					Currency:         offerStoreItem.Currency,
				})
			}
		}
	}

	allFalse := true
	// Just for check
	offerStore.ItemFlags[offerStoreItem.Currency] = false
	for _, v := range offerStore.ItemFlags {
		if v == true {
			allFalse = false
			break
		}
	}
	waitOnChain := (offerStoreItem.Currency == bean.ETH.Code && hasSell) || offerStore.Status == bean.OFFER_STORE_STATUS_CLOSING
	if allFalse {
		if waitOnChain {
			// Need to wait for OnChain
			offerStore.Status = bean.OFFER_STORE_STATUS_CLOSING
		} else {
			offerStore.Status = bean.OFFER_STORE_STATUS_CLOSED
		}
	}

	if waitOnChain {
		// Only update
		offerStoreItem.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSING
		offerStore.ItemSnapshots[offerStoreItem.Currency] = offerStoreItem

		err := s.dao.UpdateOfferStoreItemClosing(offerStore, offerStoreItem)
		if ce.SetError(api_error.UpdateDataFailed, err) {
			return
		}
	} else {
		profile := profileTO.Object.(bean.Profile)
		profile.ActiveOfferStores[offerStoreItem.Currency] = false
		offerStore.ItemFlags = profile.ActiveOfferStores

		// Really remove the item
		offerStoreItem.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSED
		offerStore.ItemSnapshots[offerStoreItem.Currency] = offerStoreItem

		err := s.dao.RemoveOfferStoreItem(offerStore, offerStoreItem, profile)
		if ce.SetError(api_error.DeleteDataFailed, err) {
			return
		}
	}

	// Assign to correct flag
	offerStore.ItemFlags[offerStoreItem.Currency] = offerStoreItem.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSED

	notification.SendOfferStoreNotification(offerStore, offerStoreItem)

	return
}

func (s OfferStoreService) CreateOfferStoreShake(userId string, offerStoreId string, offerShakeBody bean.OfferStoreShake) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)
	if profile.UserId == offerStore.UID {
		ce.SetStatusKey(api_error.OfferPayMyself)
		return
	}

	// Check offer item exists
	offerStoreItemTO := s.dao.GetOfferStoreItem(offerStoreId, offerShakeBody.Currency)
	if !offerStoreItemTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	// Make sure shake on the valid item
	offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)
	if offerStoreItem.Status != bean.OFFER_STORE_ITEM_STATUS_ACTIVE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}
	var balance decimal.Decimal
	amount, _ := decimal.NewFromString(offerShakeBody.Amount)
	if offerShakeBody.Type == bean.OFFER_TYPE_SELL {
		balance, _ = decimal.NewFromString(offerStoreItem.SellBalance)
	} else {
		balance, _ = decimal.NewFromString(offerStoreItem.BuyBalance)
	}
	if balance.LessThan(amount) {
		ce.SetStatusKey(api_error.UpdateDataFailed)
	}

	offerShakeBody.UID = userId
	offerShakeBody.FiatCurrency = offerStore.FiatCurrency
	offerShakeBody.Latitude = offerStore.Latitude
	offerShakeBody.Longitude = offerStore.Longitude

	s.setupOfferShakePrice(&offerShakeBody, &ce)
	s.setupOfferShakeAmount(&offerShakeBody, &ce)
	if ce.HasError() {
		return
	}

	// Status of shake
	if offerShakeBody.Type == bean.OFFER_TYPE_SELL {
		offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
	} else {
		if offerShakeBody.Currency == bean.ETH.Code {
			offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING
		} else {
			offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKING
			s.generateSystemAddressForShake(offerStore, &offerShakeBody, &ce)
			if ce.HasError() {
				return
			}
		}
	}

	var err error
	offerShake, err = s.dao.AddOfferStoreShake(offerStore, offerShakeBody)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offerShake.CreatedAt = time.Now().UTC()
	notification.SendOfferStoreShakeNotification(offerShake, offerStore)

	return
}

func (s OfferStoreService) RejectOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)

	offerStoreShakeTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if !offerStoreShakeTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStoreShake := offerStoreShakeTO.Object.(bean.OfferStoreShake)

	if profile.UserId != offerStore.UID && profile.UserId != offerStoreShake.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerStoreShake.Status != bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerStoreShake.Type == bean.OFFER_TYPE_SELL {
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
	} else {
		if offerStoreShake.Currency == bean.ETH.Code {
			// Only ETH
			offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTING
		} else {
			// Only BTC
			offerStoreItemTO := s.dao.GetOfferStoreItem(userId, offerStoreShake.Currency)
			if offerStoreItemTO.HasError() {
				ce.SetStatusKey(api_error.GetDataFailed)
				return
			}
			offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)

			offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
			description := fmt.Sprintf("Refund to userId %s due to reject the offer", offerStoreShake.UID)
			userAddress := offerStoreShake.UserAddress

			response := s.sendTransaction(userAddress,
				offerStoreShake.Amount, offerStoreShake.Currency, description, offerStore.UID, offerStoreItem.WalletProvider, &ce)
			if !ce.HasError() {
				s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
					Provider:         offerStoreItem.WalletProvider,
					ProviderResponse: response,
					DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
					DataRef:          dao.GetOfferStoreShakeItemPath(offerStoreId, offerStoreShakeId),
					UID:              userId,
					Description:      description,
					Amount:           offerStoreShake.Amount,
					Currency:         offerStoreShake.Currency,
				})
			}
		}
	}

	transCount := s.getFailedTransCount(offerStore, offerStoreShake, userId)
	err := s.dao.UpdateOfferStoreShakeReject(offerStore, offerStoreShake, profile, transCount)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	if userId == offerStoreShake.UID {
		UserServiceInst.UpdateOfferRejectLock(userId)
	}

	offerStoreShake.ActionUID = userId
	notification.SendOfferStoreShakeNotification(offerStoreShake, offerStore)
	offerShake = offerStoreShake

	return
}

func (s OfferStoreService) CancelOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)

	offerStoreShakeTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if !offerStoreShakeTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStoreShake := offerStoreShakeTO.Object.(bean.OfferStoreShake)

	if profile.UserId != offerStore.UID && profile.UserId != offerStoreShake.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerStoreShake.Status != bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerStoreShake.Currency == bean.ETH.Code {
		// Only ETH
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_CANCELLING
	} else {
		// Only BTC
		offerStoreItemTO := s.dao.GetOfferStoreItem(userId, offerStoreShake.Currency)
		if offerStoreItemTO.HasError() {
			ce.SetStatusKey(api_error.GetDataFailed)
			return
		}
		offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)

		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_CANCELLED
		description := fmt.Sprintf("Refund to userId %s due to reject the offer", offerStoreShake.UID)
		userAddress := offerStoreShake.UserAddress
		if offerStoreShake.Type == bean.OFFER_TYPE_BUY {
			description = fmt.Sprintf("Refund to userId %s due to reject the offer", offerStore.UID)
			userAddress = offerStoreItem.UserAddress

		}
		response := s.sendTransaction(userAddress,
			offerStoreShake.Amount, offerStoreShake.Currency, description, offerStore.UID, offerStoreItem.WalletProvider, &ce)
		if !ce.HasError() {
			s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
				Provider:         offerStoreItem.WalletProvider,
				ProviderResponse: response,
				DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
				DataRef:          dao.GetOfferStoreShakeItemPath(offerStoreId, offerStoreShakeId),
				UID:              userId,
				Description:      description,
				Amount:           offerStoreShake.Amount,
				Currency:         offerStoreShake.Currency,
			})
		}
	}

	err := s.dao.UpdateOfferStoreShake(offerStoreId, offerStoreShake, offerStoreShake.GetChangeStatus())
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	offerShake = offerStoreShake
	notification.SendOfferStoreShakeNotification(offerStoreShake, offerStore)

	return
}

func (s OfferStoreService) AcceptOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)

	offerStoreShakeTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if !offerStoreShakeTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStoreShake := offerStoreShakeTO.Object.(bean.OfferStoreShake)

	if profile.UserId != offerStore.UID {
		ce.SetStatusKey(api_error.InvalidRequestBody)
		return
	}
	if offerStoreShake.Status != bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerStoreShake.Currency == bean.ETH.Code {
		// Only ETH
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKING
	} else {
		// Only BTC
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
	}

	err := s.dao.UpdateOfferStoreShake(offerStoreId, offerStoreShake, offerStoreShake.GetChangeStatus())
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	offerShake = offerStoreShake
	notification.SendOfferStoreShakeNotification(offerStoreShake, offerStore)

	return
}

func (s OfferStoreService) CompleteOfferStoreShake(userId string, offerStoreId string, offerStoreShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)

	offerStoreShakeTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if !offerStoreShakeTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStoreShake := offerStoreShakeTO.Object.(bean.OfferStoreShake)

	if offerStoreShake.Type == bean.OFFER_TYPE_SELL {
		if profile.UserId != offerStore.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	} else {
		if profile.UserId != offerStoreShake.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
	}

	if offerStoreShake.Status != bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	if offerStoreShake.Currency == bean.ETH.Code {
		// Only ETH
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETING
	} else {
		// Only BTC
		offerStoreShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETED
		// Do Transfer
		s.transferCrypto(&offerStore, &offerStoreShake, &ce)
	}

	transCount1, transCount2 := s.getSuccessTransCount(offerStore, offerShake, userId)

	offerShake = offerStoreShake
	notification.SendOfferStoreShakeNotification(offerStoreShake, offerStore)
	err := s.dao.UpdateOfferStoreShakeComplete(offerStore, offerShake, profile, transCount1, transCount2)

	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	// For onchain processing
	offerShake = offerStoreShake
	if offerShake.Hid == 0 {
		offerShake.Hid = offerStore.Hid
	}
	if offerShake.UserAddress == "" {
		offerShake.UserAddress = offerStore.ItemSnapshots[offerShake.Currency].UserAddress
	}
	notification.SendOfferStoreShakeNotification(offerStoreShake, offerStore)

	return
}

func (s OfferStoreService) UpdateOnChainInitOfferStore(offerStoreId string, hid int64) (offer bean.OfferStore, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStore(offerStoreId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	if !offerTO.Found {
		return
	}
	offer = offerTO.Object.(bean.OfferStore)
	offerItemTO := s.dao.GetOfferStoreItem(offerStoreId, bean.ETH.Code)
	if !offerItemTO.Found {
		return
	}
	offerItem := offerItemTO.Object.(bean.OfferStoreItem)
	if offerItem.Status != bean.OFFER_STORE_ITEM_STATUS_CREATED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	if offer.Hid == 0 {
		offer.Hid = hid
	}
	offerItem.Status = bean.OFFER_STORE_ITEM_STATUS_ACTIVE
	offerItem.SellBalance = offerItem.SellAmount
	if offer.Status == bean.OFFER_STORE_STATUS_CREATED {
		offer.Status = bean.OFFER_STORE_STATUS_ACTIVE
	}
	offer.ItemSnapshots[offerItem.Currency] = offerItem
	err := s.dao.UpdateOfferStoreItemActive(offer, offerItem)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, offerItem)

	return
}

func (s OfferStoreService) UpdateOnChainCloseOfferStore(offerStoreId string) (offer bean.OfferStore, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStore(offerStoreId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	if !offerTO.Found {
		return
	}
	offer = offerTO.Object.(bean.OfferStore)

	profileTO := s.userDao.GetProfile(offer.UID)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)

	offerItemTO := s.dao.GetOfferStoreItem(offerStoreId, bean.ETH.Code)
	if !offerItemTO.Found {
		return
	}
	offerItem := offerItemTO.Object.(bean.OfferStoreItem)
	if offerItem.Status != bean.OFFER_STORE_ITEM_STATUS_CLOSING {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	offerItem.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSED
	if offer.Status == bean.OFFER_STORE_STATUS_CLOSING {
		offer.Status = bean.OFFER_STORE_STATUS_CLOSED
	}

	profile.ActiveOfferStores[offerItem.Currency] = false
	offer.ItemSnapshots[offerItem.Currency] = offerItem
	err := s.dao.UpdateOfferStoreItemClosed(offer, offerItem, profile)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, offerItem)

	return
}

func (s OfferStoreService) UpdateOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string, hid int64, oldStatus string, newStatus string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	offerTO := s.dao.GetOfferStore(offerStoreId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	if !offerTO.Found {
		return
	}
	offer := offerTO.Object.(bean.OfferStore)

	offerShakeTO := s.dao.GetOfferStoreShake(offerStoreId, offerStoreShakeId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerShakeTO) {
		return
	}
	if !offerShakeTO.Found {
		return
	}
	offerShake = offerShakeTO.Object.(bean.OfferStoreShake)
	if offerShake.Status != oldStatus {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	// Good
	offerShake.Status = newStatus
	if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
		offerItemTO := s.dao.GetOfferStoreItem(offerStoreId, offerShake.Currency)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, offerItemTO) {
			return
		}
		offerItem := offerItemTO.Object.(bean.OfferStoreItem)
		var err error
		if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
			err = s.dao.UpdateOfferStoreShakeBalance(offer, &offerItem, offerShake, true)
		} else if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
			err = s.dao.UpdateOfferStoreShakeBalance(offer, &offerItem, offerShake, false)
		}
		offer.ItemSnapshots[offerItem.Currency] = offerItem
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
		}
	}
	if offerShake.Hid == 0 {
		offerShake.Hid = hid
	}

	err := s.dao.UpdateOfferStoreShake(offerStoreId, offerShake, offerShake.GetChangeStatus())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) ActiveOnChainOfferStore(offerStoreId string, hid int64) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainInitOfferStore(offerStoreId, hid)
}

func (s OfferStoreService) CloseOnChainOfferStore(offerStoreId string) (bean.OfferStore, SimpleContextError) {
	return s.UpdateOnChainCloseOfferStore(offerStoreId)
}

func (s OfferStoreService) PreShakeOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string, hid int64) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, hid, bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING, bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE)
}

func (s OfferStoreService) CancelOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_CANCELLING, bean.OFFER_STORE_SHAKE_STATUS_CANCELLED)
}

func (s OfferStoreService) ShakeOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_SHAKING, bean.OFFER_STORE_SHAKE_STATUS_SHAKE)
}

func (s OfferStoreService) RejectOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_REJECTING, bean.OFFER_STORE_SHAKE_STATUS_REJECTED)
}

func (s OfferStoreService) CompleteOnChainOfferStoreShake(offerStoreId string, offerStoreShakeId string) (bean.OfferStoreShake, SimpleContextError) {
	return s.UpdateOnChainOfferStoreShake(offerStoreId, offerStoreShakeId, 0, bean.OFFER_STORE_SHAKE_STATUS_COMPLETING, bean.OFFER_STORE_SHAKE_STATUS_COMPLETED)
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
	solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer))

	return
}

func (s OfferStoreService) SyncOfferStoreShakeToSolr(offerStoreId, offerId string) (offer bean.OfferStoreShake, ce SimpleContextError) {
	offerStoreTO := s.dao.GetOfferStore(offerStoreId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerStoreTO) {
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)
	offerTO := s.dao.GetOfferStoreShake(offerStoreId, offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, offerTO) {
		return
	}
	offer = offerTO.Object.(bean.OfferStoreShake)
	solr_service.UpdateObject(bean.NewSolrFromOfferStoreShake(offer, offerStore))

	return
}

func (s OfferStoreService) prepareOfferStore(offerStore *bean.OfferStore, offerStoreItem *bean.OfferStoreItem, profile *bean.Profile, ce *SimpleContextError) {
	currencyInst := bean.CurrencyMapping[offerStoreItem.Currency]
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
	offerStore.ItemFlags = profile.ActiveOfferStores
	offerStore.UID = profile.UserId

	s.checkOfferStoreItemAmount(offerStoreItem, ce)
	if ce.HasError() {
		return
	}

	s.generateSystemAddress(*offerStore, offerStoreItem, ce)

	sellAmount, _ := decimal.NewFromString(offerStoreItem.SellAmount)
	if sellAmount.Equal(common.Zero) {
		// Only the case that shop doesn't sell, so don't need to wait to active
		offerStoreItem.Status = bean.OFFER_STORE_STATUS_ACTIVE
	} else {
		offerStoreItem.Status = bean.OFFER_STORE_STATUS_CREATED
	}
	if offerStore.Status != bean.OFFER_STORE_STATUS_ACTIVE {
		offerStore.Status = bean.OFFER_STORE_STATUS_CREATED
	}

	minAmount := bean.MIN_ETH
	if offerStoreItem.Currency == bean.BTC.Code {
		minAmount = bean.MIN_BTC
	}
	offerStoreItem.BuyBalance = offerStoreItem.BuyAmount
	offerStoreItem.BuyAmountMin = minAmount.String()
	offerStoreItem.SellBalance = "0"
	offerStoreItem.SellAmountMin = minAmount.String()

	if offerStore.ItemSnapshots == nil {
		offerStore.ItemSnapshots = make(map[string]bean.OfferStoreItem)
	}
	offerStore.ItemSnapshots[offerStoreItem.Currency] = *offerStoreItem
}

func (s OfferStoreService) checkOfferStoreItemAmount(offerStoreItem *bean.OfferStoreItem, ce *SimpleContextError) {
	// Minimum amount
	sellAmount, errFmt := decimal.NewFromString(offerStoreItem.SellAmount)
	if ce.SetError(api_error.InvalidRequestBody, errFmt) {
		return
	}
	if offerStoreItem.Currency == bean.ETH.Code {
		if sellAmount.LessThan(decimal.NewFromFloat(0.1).Round(1)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if offerStoreItem.Currency == bean.BTC.Code {
		if sellAmount.LessThan(decimal.NewFromFloat(0.01).Round(4)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if offerStoreItem.SellPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(offerStoreItem.SellPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		offerStoreItem.SellPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		offerStoreItem.SellPercentage = "0"
	}

	buyAmount, errFmt := decimal.NewFromString(offerStoreItem.BuyAmount)
	if ce.SetError(api_error.InvalidRequestBody, errFmt) {
		return
	}
	if offerStoreItem.Currency == bean.ETH.Code {
		if buyAmount.LessThan(decimal.NewFromFloat(0.1).Round(1)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if offerStoreItem.Currency == bean.BTC.Code {
		if buyAmount.LessThan(decimal.NewFromFloat(0.01).Round(4)) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if offerStoreItem.BuyPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(offerStoreItem.BuyPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		offerStoreItem.BuyPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		offerStoreItem.BuyPercentage = "0"
	}
}

func (s OfferStoreService) generateSystemAddress(offerStore bean.OfferStore, offer *bean.OfferStoreItem, ce *SimpleContextError) {
	// Only BTC need to generate address to transfer in
	if offer.Currency == bean.BTC.Code {
		systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_BTC_WALLET)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
			return
		}
		systemConfig := systemConfigTO.Object.(bean.SystemConfig)
		offer.WalletProvider = systemConfig.Value
		if systemConfig.Value == bean.BTC_WALLET_COINBASE {
			addressResponse, err := coinbase_service.GenerateAddress(offer.Currency)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = addressResponse.Data.Address

		} else if systemConfig.Value == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			address, err := client.GenerateAddress(offerStore.Id)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = address
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}
}

// TODO remove func duplicate
func (s OfferStoreService) generateSystemAddressForShake(offerStore bean.OfferStore, offer *bean.OfferStoreShake, ce *SimpleContextError) {
	// Only BTC need to generate address to transfer in
	if offer.Currency == bean.BTC.Code {
		systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_BTC_WALLET)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
			return
		}
		systemConfig := systemConfigTO.Object.(bean.SystemConfig)
		offer.WalletProvider = systemConfig.Value
		if systemConfig.Value == bean.BTC_WALLET_COINBASE {
			addressResponse, err := coinbase_service.GenerateAddress(offer.Currency)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = addressResponse.Data.Address

		} else if systemConfig.Value == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			address, err := client.GenerateAddress(offerStore.Id)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = address
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}
}

func (s OfferStoreService) transferCrypto(offerStore *bean.OfferStore, offerStoreShake *bean.OfferStoreShake, ce *SimpleContextError) {
	offerStoreItemTO := s.dao.GetOfferStoreItem(offerStore.UID, offerStoreShake.Currency)
	if offerStoreItemTO.HasError() {
		ce.SetStatusKey(api_error.GetDataFailed)
		return
	}
	offerStoreItem := offerStoreItemTO.Object.(bean.OfferStoreItem)
	userAddress := offerStoreShake.UserAddress
	actionUID := offerStoreShake.UID

	if offerStoreShake.Type == bean.OFFER_TYPE_BUY {
		userAddress = offerStoreItem.UserAddress
		actionUID = offerStore.UID
	}

	if offerStoreShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		if userAddress != "" {
			//Transfer
			description := fmt.Sprintf("Transfer to userId %s offerShakeId %s status %s", actionUID, offerStoreShake.Id, offerStoreShake.Status)

			var response1 interface{}
			// var response2 interface{}
			var userId string
			if offerStoreShake.Type == bean.OFFER_TYPE_BUY {
				// Amount = 1, transfer 1
				response1 = s.sendTransaction(offerStoreItem.UserAddress, offerStoreShake.Amount, offerStoreShake.Currency, description, offerStoreShake.Id, offerStoreItem.WalletProvider, ce)
				userId = offerStore.UID
			} else {
				// Amount = 1, transfer 0.09 (if fee = 1%)
				response1 = s.sendTransaction(offerStoreShake.UserAddress, offerStoreShake.TotalAmount, offerStoreShake.Currency, description, offerStoreShake.Id, offerStoreItem.WalletProvider, ce)
				userId = offerStoreShake.UID
			}
			s.miscDao.AddCryptoTransferLog(bean.CryptoTransferLog{
				Provider:         offerStoreItem.WalletProvider,
				ProviderResponse: response1,
				DataType:         bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
				DataRef:          dao.GetOfferStoreShakeItemPath(offerStore.Id, offerStoreShake.Id),
				UID:              userId,
				Description:      description,
				Amount:           offerStoreShake.Amount,
				Currency:         offerStoreShake.Currency,
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

func (s OfferStoreService) setupOfferShakeAmount(offer *bean.OfferStoreShake, ce *SimpleContextError) {
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
	amount, _ := decimal.NewFromString(offer.Amount)
	fee := amount.Mul(exchFee)
	// reward := amount.Mul(exchComm)
	// For now
	reward := decimal.NewFromFloat(0)

	offer.FeePercentage = exchFee.String()
	offer.RewardPercentage = exchComm.String()
	offer.Fee = fee.String()
	offer.Reward = reward.String()
	if offer.Type == bean.OFFER_TYPE_SELL {
		offer.TotalAmount = amount.Sub(fee.Add(reward)).String()
	} else if offer.Type == bean.OFFER_TYPE_BUY {
		offer.TotalAmount = amount.Add(fee.Add(reward)).String()
	}
}

func (s OfferStoreService) getOfferProfile(offerStore bean.OfferStore, offerStoreShake bean.OfferStoreShake, profile bean.Profile, ce *SimpleContextError) (offerProfile bean.Profile) {
	if profile.UserId == offerStore.UID {
		offerProfile = profile
	} else {
		offerProfileTO := s.userDao.GetProfile(offerStoreShake.UID)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, offerProfileTO) {
			return
		}
		offerProfile = offerProfileTO.Object.(bean.Profile)
	}

	return
}
