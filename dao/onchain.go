package dao

import (
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/service/cache"
	"strconv"
)

type OnChainDao struct {
}

func (dao OnChainDao) GetOfferInitEventBlock() (t TransferObject) {
	// GetObject(GetOfferInitEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferInitEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferInitEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferInitEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferInitEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferCloseEventBlock() (t TransferObject) {
	//GetObject(GetOfferCloseEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferCloseEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})
	return
}

func (dao OnChainDao) UpdateOfferCloseEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferCloseEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferCloseEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferShakeEventBlock() (t TransferObject) {
	//GetObject(GetOfferShakeEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferShakeEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})
	return
}

func (dao OnChainDao) UpdateOfferShakeEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferShakeEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferShakeEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferRejectEventBlock() (t TransferObject) {
	//GetObject(GetOfferRejectEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferRejectEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})
	return
}

func (dao OnChainDao) UpdateOfferRejectEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferRejectEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferRejectEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferCompleteEventBlock() (t TransferObject) {
	//GetObject(GetOfferCompleteEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferCompleteEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})
	return
}

func (dao OnChainDao) UpdateOfferCompleteEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferCompleteEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferCompleteEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferWithdrawEventBlock() (t TransferObject) {
	//GetObject(GetOfferWithdrawEventBlockPath(), &t, snapshotToOfferEventBlock)
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferWithdrawEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})
	return
}

func (dao OnChainDao) UpdateOfferWithdrawEventBlock(offer bean.OfferEventBlock) error {
	//dbClient := firebase_service.FirestoreClient
	//docRef := dbClient.Doc(GetOfferWithdrawEventBlockPath())
	//
	//_, err := docRef.Set(context.Background(), offer.GetUpdate())

	key := GetOfferWithdrawEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

//func GetOfferInitEventBlockPath() string {
//	return "onchain_events/offer_init"
//}

func GetOfferInitEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_init"
}

//func GetOfferCloseEventBlockPath() string {
//	return "onchain_events/offer_close"
//}

func GetOfferCloseEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_close"
}

//func GetOfferShakeEventBlockPath() string {
//	return "onchain_events/offer_shake"
//}

func GetOfferShakeEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_shake"
}

//func GetOfferRejectEventBlockPath() string {
//	return "onchain_events/offer_reject"
//}

func GetOfferRejectEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_reject"
}

//func GetOfferCompleteEventBlockPath() string {
//	return "onchain_events/offer_complete"
//}

func GetOfferCompleteEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_complete"
}

//func GetOfferWithdrawEventBlockPath() string {
//	return "onchain_events/offer_withdraw"
//}

func GetOfferWithdrawEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_withdraw"
}

//func snapshotToOfferEventBlock(snapshot *firestore.DocumentSnapshot) interface{} {
//	var obj bean.OfferEventBlock
//	snapshot.DataTo(&obj)
//	return obj
//}

////////////// OFFER STORE ON CHAIN //////////////

func (dao OnChainDao) GetOfferStoreInitEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreInitEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreInitEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreInitEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreCloseEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreCloseEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreCloseEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreCloseEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStorePreShakeEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStorePreShakeEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStorePreShakeEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreCancelEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreCancelEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreCancelEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreCancelEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreCancelEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreShakeEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreShakeEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreShakeEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreShakeEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreRejectEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreRejectEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreRejectEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreRejectEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreCompleteEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreCompleteEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreCompleteEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreCompleteEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func (dao OnChainDao) GetOfferStoreCompleteUserEventBlock() (t TransferObject) {
	obj := bean.OfferEventBlock{}
	GetCacheObject(GetOfferStoreCompleteUserEventBlockKey(), &t, func(val string) interface{} {
		block, _ := strconv.Atoi(val)
		obj.LastBlock = int64(block)
		return obj
	})

	return
}

func (dao OnChainDao) UpdateOfferStoreCompleteUserEventBlock(offer bean.OfferEventBlock) error {
	key := GetOfferStoreCompleteUserEventBlockKey()
	cache.RedisClient.Set(key, offer.LastBlock, 0)

	return nil
}

func GetOfferStoreInitEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_init"
}

func GetOfferStoreCloseEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_close"
}

func GetOfferStorePreShakeEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_preshake"
}

func GetOfferStoreCancelEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_cancel"
}

func GetOfferStoreShakeEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_shake"
}

func GetOfferStoreRejectEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_reject"
}

func GetOfferStoreCompleteEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_complete"
}

func GetOfferStoreCompleteUserEventBlockKey() string {
	return "handshake_exchange.onchain_events.offer_store_complete_user"
}
