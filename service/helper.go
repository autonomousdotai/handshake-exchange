package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/ethereum_service"
	"github.com/ninjadotorg/handshake-exchange/integration/exchangecreditatm_service"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
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

func ReleaseContractFund(dao dao.MiscDao, address string, amountStr string, refId string, hid int64, keySet string) (txHash string, outNonce uint64, outAddress string, err error) {
	indexTO := dao.GetCreditKeyIndexFromCache(keySet)
	if indexTO.HasError() {
		err = errors.New("cache error")
		return
	}

	keyStr := os.Getenv(keySet)
	keys := strings.Split(keyStr, ";")

	valStr := indexTO.Object.(string)
	index, _ := strconv.Atoi(valStr)
	index += 1
	if index >= len(keys) {
		index = 0
	}
	dao.CreditKeyIndexToCache(keySet, fmt.Sprintf("%d", index))

	writeClient := ethereum_service.EthereumClient{}
	writeClient.InitializeWithKey(keys[index])
	nonce, clientErr := writeClient.GetNonce()
	if clientErr != nil {
		err = errors.New("network error")
		return
	}

	keyDataTO := dao.GetCreditContractKeyDataFromCache(writeClient.GetAddress())
	if keyDataTO.HasError() {
		err = errors.New("cache error")
		return
	}
	keyData := keyDataTO.Object.(bean.CreditContractKeyData)
	outAddress = keyData.Address

	fmt.Println(nonce)
	fmt.Println(outAddress)
	if int64(nonce) <= keyData.Nonce && nonce > 0 {
		// need to retry
		err = errors.New("retry later")
		return
	}

	client := exchangecreditatm_service.ExchangeCreditAtmClient{}
	amount := common.StringToDecimal(amountStr)
	fmt.Println(refId, hid, amount, address)
	txHash, outNonce, err = client.ReleasePartialFund(refId, hid, amount, address, 0, false, keys[index])

	keyData.Nonce = int64(outNonce)
	dao.CreditContractKeyDataToCache(keyData)

	return
}
