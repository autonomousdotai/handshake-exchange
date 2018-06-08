package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/blockchainio_service"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/service/notification"
	"github.com/shopspring/decimal"
	"time"
)

type OfferStoreService struct {
	dao     *dao.OfferStoreDao
	userDao *dao.UserDao
	miscDao *dao.MiscDao
}

func (s OfferStoreService) CreateOfferStore(userId string, offerSetup bean.OfferStoreSetup) (offer bean.OfferStoreSetup, ce SimpleContextError) {
	offerBody := offerSetup.Offer
	offerItemBody := offerSetup.Item

	currencyInst := bean.CurrencyMapping[offerItemBody.Currency]
	if currencyInst.Code == "" {
		ce.SetStatusKey(api_error.UnsupportedCurrency)
		return
	}

	// Check offer store exists
	offerStoreTO := s.dao.GetOfferStore(userId)
	if offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreExists)
		return
	}

	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	profile := profileTO.Object.(bean.Profile)
	if check, ok := profile.ActiveOfferStores[currencyInst.Code]; ok && check {
		// Has Key and already had setup
		ce.SetStatusKey(api_error.TooManyOffer)
		return
	}
	profile.ActiveOfferStores[currencyInst.Code] = true
	offerBody.ItemFlags = profile.ActiveOfferStores
	offerBody.UID = userId

	s.checkOfferStoreItemAmount(&offerItemBody, &ce)
	if ce.HasError() {
		return
	}

	s.generateSystemAddress(offerBody, &offerItemBody, &ce)

	sellAmount, _ := decimal.NewFromString(offerItemBody.SellAmount)
	if offerItemBody.Currency == bean.BTC.Code && sellAmount.Equal(common.Zero) {
		// Only the case that shop doesn't sell BTC, so don't need to wait to active
		offerItemBody.Status = bean.OFFER_STORE_STATUS_ACTIVE
	} else {
		offerItemBody.Status = bean.OFFER_STORE_STATUS_CREATED
	}

	minAmount := bean.MIN_ETH
	if offerItemBody.Currency == bean.BTC.Code {
		minAmount = bean.MIN_BTC
	}
	offerItemBody.BuyBalance = offerItemBody.BuyAmount
	offerItemBody.BuyAmountMin = minAmount.String()
	offerItemBody.SellBalance = "0"
	offerItemBody.SellAmountMin = minAmount.String()

	offerNew, err := s.dao.AddOfferStore(offerBody, offerItemBody, profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offerNew.CreatedAt = time.Now().UTC()
	notification.SendOfferStoreNotification(offerNew)

	offer.Offer = offerNew
	offer.Item = offerItemBody

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

		if systemConfig.Value == bean.BTC_WALLET_COINBASE {
			addressResponse, err := coinbase_service.GenerateAddress(offer.Currency)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = addressResponse.Data.Address
			offer.WalletProvider = systemConfig.Value
		} else if systemConfig.Value == bean.BTC_WALLET_BLOCKCHAINIO {
			client := blockchainio_service.BlockChainIOClient{}
			address, err := client.GenerateAddress(offerStore.Id)
			if err != nil {
				ce.SetError(api_error.ExternalApiFailed, err)
				return
			}
			offer.SystemAddress = address
			offer.WalletProvider = systemConfig.Value
		} else {
			ce.SetStatusKey(api_error.InvalidConfig)
		}
	}
}
