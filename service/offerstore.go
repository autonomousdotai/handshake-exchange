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

	s.prepareOfferStore(&offerBody, &offerItemBody, &profile, &ce)

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

func (s OfferStoreService) GetOfferStore(userId string, offerId string) (offer bean.OfferStore, ce SimpleContextError) {
	profileTO := s.userDao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offerTO := s.dao.GetOfferStore(offerId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, profileTO) {
		return
	}
	offer = offerTO.Object.(bean.OfferStore)
	// price, _ := decimal.NewFromString(offer.Price)
	// percentage, _ := decimal.NewFromString(offer.Percentage)

	//price, fiatPrice, fiatAmount, err := s.GetQuote(offer.Type, offer.Amount, offer.Currency, offer.FiatCurrency)
	//if offer.Type == bean.OFFER_TYPE_SELL && price.Equal(common.Zero) {
	//	if ce.SetError(api_error.GetDataFailed, err) {
	//		return
	//	}
	//	markup := fiatAmount.Mul(percentage)
	//	fiatAmount = fiatAmount.Add(markup)
	//}
	//offer.Price = fiatPrice.Round(2).String()
	//offer.FiatAmount = fiatAmount.Round(2).String()

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

	notification.SendOfferStoreNotification(offerStore)

	return
}

func (s OfferStoreService) RemoveOfferStoreItem(userId string, offerStoreId string, currency string) (ce SimpleContextError) {
	// Check offer store exists
	offerStoreTO := s.dao.GetOfferStore(userId)
	if !offerStoreTO.Found {
		ce.SetStatusKey(api_error.OfferStoreNotExist)
		return
	}
	offerStore := offerStoreTO.Object.(bean.OfferStore)

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

	allFalse := true
	for _, v := range offerStore.ItemFlags {
		if v == true {
			allFalse = false
			break
		}
	}
	waitOnChain := offerStoreItem.Currency == bean.ETH.Code && hasSell
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

	notification.SendOfferStoreNotification(offerStore)

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

	offerShakeBody.UID = userId
	offerShakeBody.FiatCurrency = offerStore.FiatCurrency
	s.setupOfferShakePrice(&offerShakeBody, &ce)
	s.setupOfferShakeAmount(&offerShakeBody, &ce)
	if ce.HasError() {
		return
	}

	if offerShakeBody.Currency == bean.ETH.Code {
		// Only ETH
		offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING
	} else {
		// Only BTC
		if offerShakeBody.Type == bean.OFFER_TYPE_SELL {
			offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
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
	if offerStoreItem.Currency == bean.BTC.Code && sellAmount.Equal(common.Zero) {
		// Only the case that shop doesn't sell BTC, so don't need to wait to active
		offerStoreItem.Status = bean.OFFER_STORE_STATUS_ACTIVE
	} else {
		offerStoreItem.Status = bean.OFFER_STORE_STATUS_CREATED
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
	reward := amount.Mul(exchComm)

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
