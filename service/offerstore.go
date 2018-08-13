package service

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/service/notification"
	"github.com/shopspring/decimal"
	"strconv"
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
	if ce.HasError() {
		return
	}

	_, err := s.dao.AddOfferStoreItem(offer, item, *profile)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) UpdateOfferStore(userId string, offerId string, body bean.OfferStoreSetup) (offer bean.OfferStore, ce SimpleContextError) {
	_ = GetProfile(s.userDao, userId, &ce)
	if ce.HasError() {
		return
	}
	checkOffer := GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offer = *checkOffer
	bodyItem := body.Item
	checkOfferItem := GetOfferStoreItem(*s.dao, offerId, bodyItem.Currency, &ce)
	if ce.HasError() {
		return
	}
	// Copy data
	offer.FiatCurrency = body.Offer.FiatCurrency
	offer.ContactPhone = body.Offer.ContactPhone
	offer.ContactInfo = body.Offer.ContactInfo
	item := *checkOfferItem
	if bodyItem.SellPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(bodyItem.SellPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		item.SellPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		item.SellPercentage = "0"
	}
	if bodyItem.SellAmount != "" {
		item.SellAmount = bodyItem.SellAmount
	}

	if bodyItem.BuyPercentage != "" {
		// Convert to 0.0x
		percentage, errFmt := decimal.NewFromString(bodyItem.BuyPercentage)
		if ce.SetError(api_error.InvalidRequestBody, errFmt) {
			return
		}
		item.BuyPercentage = percentage.Div(decimal.NewFromFloat(100)).String()
	} else {
		item.BuyPercentage = "0"
	}
	if bodyItem.BuyAmount != "" {
		item.BuyAmount = bodyItem.BuyAmount
	}

	offer.ItemSnapshots[bodyItem.Currency] = item

	_, err := s.dao.UpdateOfferStoreItem(offer, item)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	// Only sync to solr
	solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, item))

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

	if item.Status != bean.OFFER_STORE_ITEM_STATUS_ACTIVE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
		return
	}

	count, err := s.countActiveShake(offer.Id, "", item.Currency)
	if err != nil {
		ce.SetError(api_error.GetDataFailed, err)
		return
	}
	// There is still active shake so cannot close
	if count > 0 {
		ce.SetStatusKey(api_error.OfferStoreShakeActiveExist)
		return
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

	if allFalse {
		offer.Status = bean.OFFER_STORE_STATUS_CLOSED
	}

	// Just a time to response
	item.UpdatedAt = time.Now().UTC()

	profile.ActiveOfferStores[item.Currency] = false
	offer.ItemFlags = profile.ActiveOfferStores

	// Really remove the item
	item.Status = bean.OFFER_STORE_ITEM_STATUS_CLOSED
	offer.ItemSnapshots[item.Currency] = item

	err = s.dao.RemoveOfferStoreItem(offer, item, *profile)
	if ce.SetError(api_error.DeleteDataFailed, err) {
		return
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

	amount := common.StringToDecimal(offerShakeBody.Amount)
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
	if offerShakeBody.Currency == bean.BCH.Code {
		if amount.LessThan(bean.MIN_BCH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
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

	var err error
	offerShakeBody.Status = bean.OFFER_STORE_SHAKE_STATUS_SHAKE
	s.updatePendingTransCount(offer, offerShakeBody, userId)
	offerShake, err = s.dao.AddOfferStoreShake(offer, offerShakeBody)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	offerShake.CreatedAt = time.Now().UTC()
	notification.SendOfferStoreShakeNotification(offerShake, offer)
	notification.SendOfferStoreNotification(offer, item)

	return
}

func (s OfferStoreService) RejectOfferStoreShake(userId string, offerId string, offerShakeId string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	profile := *GetProfile(s.userDao, userId, &ce)
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
		return
	}

	offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
	s.updateFailedTransCount(offer, offerShake, userId)
	err := s.dao.UpdateOfferStoreShakeReject(offer, offerShake, profile)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	if userId == offerShake.UID {
		UserServiceInst.UpdateOfferRejectLock(profile)
	}

	offerShake.ActionUID = userId
	notification.SendOfferStoreShakeNotification(offerShake, offer)
	notification.SendOfferStoreNotification(offer, item)

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

	// This is profile of user (not shop) to give referral bonus
	var userProfile bean.Profile
	if offerShake.Type == bean.OFFER_TYPE_SELL {
		if profile.UserId != offer.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		userProfile = *GetProfile(s.userDao, offerShake.UID, &SimpleContextError{})
	} else {
		if profile.UserId != offerShake.UID {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
		userProfile = *profile
	}

	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETED
	s.updateSuccessTransCount(offer, offerShake, userId)
	ReferralServiceInst.AddReferralOfferStoreShake(userProfile, offer, offerShake)
	err := s.dao.UpdateOfferStoreShakeComplete(offer, offerShake, *profile)

	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	notification.SendOfferStoreShakeNotification(offerShake, offer)

	// For onchain processing
	//if offerShake.Hid == 0 {
	//	offerShake.Hid = offer.Hid
	//}
	//if offerShake.UserAddress == "" {
	//	offerShake.UserAddress = offer.ItemSnapshots[offerShake.Currency].UserAddress
	//}

	return
}

func (s OfferStoreService) TransferOfferStoreShake(userId string, offerId string, offerShakeId string, txHash string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	offer := *GetOfferStore(*s.dao, offerId, &ce)
	if ce.HasError() {
		return
	}
	offerShake = *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}

	if offerShake.Status != bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		ce.SetStatusKey(api_error.OfferStatusInvalid)
	}

	offerShake.SubStatus = bean.OFFER_STORE_SHAKE_SUB_STATUS_TRANSFERING
	offerShake.TxHash = txHash

	err := s.dao.UpdateOfferStoreShakeTransfer(offer, offerShake)

	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
		return
	}

	notification.SendOfferStoreShakeNotification(offerShake, offer)

	return
}

func (s OfferStoreService) FinishOfferStorePendingTransfer(ref string) (offer bean.OfferStore, ce SimpleContextError) {
	return
}

func (s OfferStoreService) FinishOfferStoreShakePendingTransfer(ref string) (offerShake bean.OfferStoreShake, ce SimpleContextError) {
	to := s.dao.GetOfferStoreShakeByPath(ref)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	offerShake = to.Object.(bean.OfferStoreShake)
	offer := *GetOfferStore(*s.dao, offerShake.UID, &ce)
	if ce.HasError() {
		return
	}

	if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETING {
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_COMPLETED
	} else if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTING {
		offerShake.Status = bean.OFFER_STORE_SHAKE_STATUS_REJECTED
	}

	err := s.dao.UpdateOfferStoreShake(offer.Id, offerShake, offerShake.GetChangeStatus())
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}
	notification.SendOfferStoreShakeNotification(offerShake, offer)

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

func (s OfferStoreService) UpdateOfferStoreShakeLocation(userId string, offerId string, offerShakeId string, body bean.OfferStoreShakeLocation) (offerLocation bean.OfferStoreShakeLocation, ce SimpleContextError) {
	data := body.Data
	offerShake := *GetOfferStoreShake(*s.dao, offerId, offerShakeId, &ce)
	if ce.HasError() {
		return
	}

	locationType := "GPS"
	if data[0:1] != "G" {
		locationType = "IP"
	}

	lat1n, _ := strconv.Atoi(data[1:2])
	lat1 := data[2 : 2+lat1n]
	lat2n, _ := strconv.Atoi(data[2+lat1n : 2+lat1n+1])
	lat2 := data[2+lat1n+1 : 2+lat1n+1+lat2n]
	startLong := 2 + lat1n + 1 + lat2n
	long1n, _ := strconv.Atoi(data[startLong : startLong+1])
	long1 := data[startLong+1 : startLong+1+long1n]
	long2n, _ := strconv.Atoi(data[startLong+1+long1n : startLong+1+long1n+1])
	long2 := data[startLong+1+long1n+1 : startLong+1+long1n+1+long2n]

	lat, _ := decimal.NewFromString(fmt.Sprintf("%s.%s", lat1, lat2))
	long, _ := decimal.NewFromString(fmt.Sprintf("%s.%s", long1, long2))

	offerLocation = body
	offerLocation.ActionUID = userId
	offerLocation.Offer = offerId
	offerLocation.OfferShake = offerShakeId
	offerLocation.LocationType = locationType
	offerLocation.Latitude, _ = lat.Float64()
	offerLocation.Longitude, _ = long.Float64()

	s.dao.UpdateOfferStoreShakeLocation(userId, offerShake, offerLocation)

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
		price := common.StringToDecimal(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else if quoteType == "sell" {
		resp, errResp := coinbase_service.GetSellPrice(currency)
		err = errResp
		if err != nil {
			return
		}
		price := common.StringToDecimal(resp.Amount)
		fiatPrice = price.Mul(rateNumber)
		fiatAmount = tmpAmount.Mul(price)
	} else {
		err = errors.New(api_error.InvalidQueryParam)
	}

	return
}

func (s OfferStoreService) GetCurrentFreeStart(userId string, token string) (freeStart bean.OfferStoreFreeStart, ce SimpleContextError) {
	systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_OFFER_STORE_FREE_START)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
		return
	}
	systemConfig := systemConfigTO.Object.(bean.SystemConfig)
	// There is no free start on
	if systemConfig.Value == bean.OFFER_STORE_FREE_START_OFF {
		return
	}
	to := s.dao.GetOfferStoreFreeStartUser(userId)
	if to.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, to)
		return
	}
	if to.Found {
		freeStartTest := to.Object.(bean.OfferStoreFreeStartUser)
		if freeStartTest.Status == bean.OFFER_STORE_FREE_START_STATUS_DONE {
			return
		}
	}

	freeStarts, err := s.dao.ListOfferStoreFreeStart(token)
	if err != nil {
		ce.SetError(api_error.GetDataFailed, err)
	}

	for _, item := range freeStarts {
		if item.Count < item.Limit {
			freeStart = item
			break
		}
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

	item.Status = bean.OFFER_STORE_ITEM_STATUS_ACTIVE
	offer.Status = bean.OFFER_STORE_STATUS_ACTIVE

	minAmount := bean.MIN_ETH
	if item.Currency == bean.BTC.Code {
		minAmount = bean.MIN_BTC
	}
	if item.Currency == bean.BCH.Code {
		minAmount = bean.MIN_BCH
	}
	item.BuyBalance = common.Zero.String()
	item.BuyAmountMin = minAmount.String()
	item.SellBalance = common.Zero.String()
	item.SellAmountMin = minAmount.String()
	item.SellTotalAmount = common.Zero.String()
	item.CreatedAt = time.Now().UTC()

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
	if item.Currency == bean.BCH.Code {
		if sellAmount.GreaterThan(common.Zero) && sellAmount.LessThan(bean.MIN_BCH) {
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
		if buyAmount.GreaterThan(common.Zero) && buyAmount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.Currency == bean.BTC.Code {
		if buyAmount.GreaterThan(common.Zero) && buyAmount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if item.Currency == bean.BCH.Code {
		if buyAmount.GreaterThan(common.Zero) && buyAmount.LessThan(bean.MIN_BCH) {
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

func (s OfferStoreService) updateSuccessTransCount(offer bean.OfferStore, offerShake bean.OfferStoreShake, actionUID string) (transCount1 bean.TransactionCount, transCount2 bean.TransactionCount) {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, offerShake.Currency)
	if transCountTO.HasError() {
		return
	}
	transCount1 = transCountTO.Object.(bean.TransactionCount)
	transCount1.Currency = offerShake.Currency
	transCount1.Success += 1
	transCount1.Pending -= 1
	if transCount1.Pending < 0 {
		// Just for prevent weird number
		transCount1.Pending = 0
	}
	if offerShake.IsTypeSell() {
		sellAmount := common.StringToDecimal(transCount1.SellAmount)
		amount := common.StringToDecimal(offerShake.Amount)
		transCount1.SellAmount = sellAmount.Add(amount).String()

		if fiatAmountObj, ok := transCount1.SellFiatAmounts[offerShake.FiatCurrency]; ok {
			fiatAmount := common.StringToDecimal(fiatAmountObj.Amount)
			newFiatAmount := common.StringToDecimal(offerShake.FiatAmount)
			fiatAmountObj.Amount = fiatAmount.Add(newFiatAmount).String()
		} else {
			transCount1.SellFiatAmounts[offerShake.FiatCurrency] = bean.TransactionFiatAmount{
				Currency: offerShake.FiatCurrency,
				Amount:   offerShake.FiatAmount,
			}
		}
	} else {
		buyAmount := common.StringToDecimal(transCount1.BuyAmount)
		amount := common.StringToDecimal(offerShake.Amount)
		transCount1.BuyAmount = buyAmount.Add(amount).String()

		if fiatAmountObj, ok := transCount1.BuyFiatAmounts[offerShake.FiatCurrency]; ok {
			fiatAmount := common.StringToDecimal(fiatAmountObj.Amount)
			newFiatAmount := common.StringToDecimal(offerShake.FiatAmount)
			fiatAmountObj.Amount = fiatAmount.Add(newFiatAmount).String()
		} else {
			transCount1.BuyFiatAmounts[offerShake.FiatCurrency] = bean.TransactionFiatAmount{
				Currency: offerShake.FiatCurrency,
				Amount:   offerShake.FiatAmount,
			}
		}
	}

	//transCountTO = s.transDao.GetTransactionCount(offerShake.UID, offerShake.Currency)
	//if transCountTO.HasError() {
	//	return
	//}
	//transCount2 = transCountTO.Object.(bean.TransactionCount)
	//transCount2.Currency = offerShake.Currency
	//transCount2.Success += 1

	s.transDao.UpdateTransactionCount(offer.UID, offerShake.Currency, transCount1.GetUpdateSuccess())
	//s.transDao.UpdateTransactionCount(offerShake.UID, offerShake.Currency, transCount2.GetUpdateSuccess())

	return
}

func (s OfferStoreService) updatePendingTransCount(offer bean.OfferStore, offerShake bean.OfferStoreShake, actionUID string) (transCount bean.TransactionCount) {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, offerShake.Currency)
	if transCountTO.HasError() {
		return
	}
	transCount = transCountTO.Object.(bean.TransactionCount)
	transCount.Currency = offerShake.Currency
	transCount.Pending += 1

	s.transDao.UpdateTransactionCount(offer.UID, offerShake.Currency, transCount.GetUpdatePending())

	return
}

func (s OfferStoreService) updateFailedTransCount(offer bean.OfferStore, offerShake bean.OfferStoreShake, actionUID string) (transCount bean.TransactionCount) {
	transCountTO := s.transDao.GetTransactionCount(offer.UID, offerShake.Currency)
	if transCountTO.HasError() {
		return
	}
	transCount = transCountTO.Object.(bean.TransactionCount)
	transCount.Currency = offerShake.Currency
	if actionUID == offer.UID {
		transCount.Failed += 1
	}
	transCount.Pending -= 1
	if transCount.Pending < 0 {
		// Just for prevent weird number
		transCount.Pending = 0
	}
	s.transDao.UpdateTransactionCount(offer.UID, offerShake.Currency, transCount.GetUpdateFailed())

	return
}

func (s OfferStoreService) setupOfferShakePrice(offer *bean.OfferStoreShake, ce *SimpleContextError) {
	userOfferType := bean.OFFER_TYPE_SELL
	if offer.IsTypeSell() {
		userOfferType = bean.OFFER_TYPE_BUY
	}
	_, fiatPrice, fiatAmount, err := s.GetQuote(userOfferType, offer.Amount, offer.Currency, offer.FiatCurrency)
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
	amount := common.StringToDecimal(offerShake.Amount)
	fee := amount.Mul(exchFee)
	// reward := amount.Mul(exchComm)
	// For now
	reward := decimal.NewFromFloat(0)

	offerShake.FeePercentage = exchFee.String()
	offerShake.RewardPercentage = exchComm.String()
	offerShake.Fee = fee.String()
	offerShake.Reward = reward.String()
	if offerShake.Type == bean.OFFER_TYPE_SELL {
		offerShake.TotalAmount = amount.Add(fee.Add(reward)).String()
	} else if offerShake.Type == bean.OFFER_TYPE_BUY {
		offerShake.TotalAmount = amount.Sub(fee.Add(reward)).String()
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

func (s OfferStoreService) countActiveShake(offerId string, offerType string, currency string) (int, error) {
	offerShakes, err := s.dao.ListOfferStoreShake(offerId)
	count := 0
	countInactive := 0
	if err == nil {
		for _, offerShake := range offerShakes {
			if offerType != "" && offerShake.Type == offerType && offerShake.Currency == currency {
				if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_CANCELLED {
					countInactive += 1
				}
				count += 1
			} else if offerShake.Currency == currency {
				if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_CANCELLED {
					countInactive += 1
				}
				count += 1
			}
		}
	}
	return count - countInactive, err
}

func (s OfferStoreService) ScriptUpdateTransactionCount() error {
	t := s.dao.ListOfferStore()
	if !t.HasError() {
		for _, item := range t.Objects {
			offer := item.(bean.OfferStore)
			fmt.Printf("Updating store %s", offer.UID)
			fmt.Println("")

			shakes, err := s.dao.ListOfferStoreShake(offer.UID)
			if err == nil {
				btcTxCount := bean.TransactionCount{
					Currency:        bean.BTC.Code,
					Success:         0,
					Failed:          0,
					Pending:         0,
					BuyAmount:       common.Zero.String(),
					SellAmount:      common.Zero.String(),
					BuyFiatAmounts:  map[string]bean.TransactionFiatAmount{},
					SellFiatAmounts: map[string]bean.TransactionFiatAmount{},
				}
				ethTxCount := bean.TransactionCount{
					Currency:        bean.ETH.Code,
					Success:         0,
					Failed:          0,
					Pending:         0,
					BuyAmount:       common.Zero.String(),
					SellAmount:      common.Zero.String(),
					BuyFiatAmounts:  map[string]bean.TransactionFiatAmount{},
					SellFiatAmounts: map[string]bean.TransactionFiatAmount{},
				}

				txCountMap := map[string]*bean.TransactionCount{
					bean.ETH.Code: &ethTxCount,
					bean.BTC.Code: &btcTxCount,
				}

				for _, offerShake := range shakes {
					currency := offerShake.Currency
					fmt.Printf("Processing shake %s %s", offerShake.Id, offerShake.Currency)
					fmt.Println("")

					if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
						txCountMap[currency].Success += 1

						if offerShake.IsTypeSell() {
							sellAmount := common.StringToDecimal(txCountMap[currency].SellAmount)
							amount := common.StringToDecimal(offerShake.Amount)
							txCountMap[currency].SellAmount = sellAmount.Add(amount).String()

							if fiatAmountObj, ok := txCountMap[currency].SellFiatAmounts[offerShake.FiatCurrency]; ok {
								fiatAmount := common.StringToDecimal(fiatAmountObj.Amount)
								newFiatAmount := common.StringToDecimal(offerShake.FiatAmount)
								fiatAmountObj.Amount = fiatAmount.Add(newFiatAmount).String()
							} else {
								txCountMap[currency].SellFiatAmounts[offerShake.FiatCurrency] = bean.TransactionFiatAmount{
									Currency: offerShake.FiatCurrency,
									Amount:   offerShake.FiatAmount,
								}
							}
						} else {
							buyAmount := common.StringToDecimal(txCountMap[currency].BuyAmount)
							amount := common.StringToDecimal(offerShake.Amount)
							txCountMap[currency].BuyAmount = buyAmount.Add(amount).String()

							if fiatAmountObj, ok := txCountMap[currency].BuyFiatAmounts[offerShake.FiatCurrency]; ok {
								fiatAmount := common.StringToDecimal(fiatAmountObj.Amount)
								newFiatAmount := common.StringToDecimal(offerShake.FiatAmount)
								fiatAmountObj.Amount = fiatAmount.Add(newFiatAmount).String()
							} else {
								txCountMap[currency].BuyFiatAmounts[offerShake.FiatCurrency] = bean.TransactionFiatAmount{
									Currency: offerShake.FiatCurrency,
									Amount:   offerShake.FiatAmount,
								}
							}
						}
					} else if offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
						txCountMap[currency].Failed += 1
					} else {
						txCountMap[currency].Pending += 1
					}
				}
				if btcTxCount.Pending > 0 || btcTxCount.Success > 0 || btcTxCount.Failed > 0 {
					fmt.Printf("Making update BTC tx count %s", offer.UID)
					fmt.Println("")
					s.transDao.UpdateTransactionCountForce(offer.UID, bean.BTC.Code, btcTxCount.GetUpdateOverride())
				}
				if ethTxCount.Pending > 0 || ethTxCount.Success > 0 || ethTxCount.Failed > 0 {
					fmt.Printf("Making update ETH tx count %s", offer.UID)
					fmt.Println("")
					s.transDao.UpdateTransactionCountForce(offer.UID, bean.ETH.Code, ethTxCount.GetUpdateOverride())
				}
			}
		}
	}

	return t.Error
}

func (s OfferStoreService) ScriptUpdateOfferStoreSolr() error {
	t := s.dao.ListOfferStore()
	if !t.HasError() {
		for _, item := range t.Objects {
			offer := item.(bean.OfferStore)
			fmt.Printf("Updating store %s", offer.UID)
			fmt.Println("")
			solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, bean.OfferStoreItem{}))
		}
	}
	return t.Error
}
