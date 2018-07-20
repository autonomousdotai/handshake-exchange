package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/shopspring/decimal"
	"google.golang.org/api/iterator"
	"strings"
)

type OfferStoreDao struct {
}

func (dao OfferStoreDao) GetOfferStore(offerId string) (t TransferObject) {
	GetObject(GetOfferStoreItemPath(offerId), &t, snapshotToOfferStore)

	return
}

func (dao OfferStoreDao) ListOfferStore() (t TransferObject) {
	ListObjects(GetOfferStorePath(), &t, nil, snapshotToOfferStore)
	return
}

func (dao OfferStoreDao) AddOfferStore(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) (bean.OfferStore, error) {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferStoreItemPath(offer.UID)
	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(offerPath)
	offer.Id = docRef.ID
	offer.ItemSnapshots = map[string]bean.OfferStoreItem{
		item.Currency: item,
	}

	batch := dbClient.Batch()

	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	if item.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  item.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferStoreItemItemPath(offer.Id, item.Currency),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}
	batch.Set(itemDocRef, item.GetAddOfferStoreItem())

	batch.Set(docRef, offer.GetAddOfferStore())
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

	if item.Currency == bean.ETH.Code && item.Status == bean.OFFER_STORE_ITEM_STATUS_CREATED {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   item.Status,
			Currency: item.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return offer, err
}

func (dao OfferStoreDao) UpdateOfferStore(offer bean.OfferStore, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	_, err := docRef.Set(context.Background(), updateData, firestore.MergeAll)

	return err
}

func (dao OfferStoreDao) GetOfferStoreItem(userId string, currency string) (t TransferObject) {
	GetObject(GetOfferStoreItemItemPath(userId, currency), &t, snapshotToOfferStoreItem)

	return
}

func (dao OfferStoreDao) GetOfferStoreItemByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToOfferStoreItem)

	return
}

func (dao OfferStoreDao) AddOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) (bean.OfferStoreItem, error) {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	offerPath := GetOfferStoreItemPath(offer.UID)
	docRef := dbClient.Doc(offerPath)

	batch := dbClient.Batch()
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	if item.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  item.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferStoreItemItemPath(offer.Id, item.Currency),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	batch.Set(itemDocRef, item.GetAddOfferStoreItem())
	batch.Set(docRef, offer.GetUpdateOfferStoreChangeItem(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

	if item.Currency == bean.ETH.Code && item.Status == bean.OFFER_STORE_ITEM_STATUS_CREATED {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   item.Status,
			Currency: item.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return item, err
}

func (dao OfferStoreDao) UpdateOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem) (bean.OfferStoreItem, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.UID))
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	batch := dbClient.Batch()

	// For now only update Percentage, other info will not be updated
	batch.Set(docRef, offer.GetUpdateOfferItemInfo(), firestore.MergeAll)
	batch.Set(itemDocRef, item.GetUpdateOfferStoreItemInfo(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return item, err
}

func (dao OfferStoreDao) UpdateRefillOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem) (bean.OfferStoreItem, error) {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferStoreItemPath(offer.UID)
	docRef := dbClient.Doc(offerPath)

	batch := dbClient.Batch()
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	if item.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  item.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferStoreItemItemPath(offer.Id, item.Currency),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_ITEM,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	batch.Set(itemDocRef, item.GetUpdateOfferStoreItemRefill(), firestore.MergeAll)
	batch.Set(docRef, offer.GetUpdateOfferStoreChangeSnapshot(), firestore.MergeAll)

	if item.Currency == bean.ETH.Code && item.SubStatus == bean.OFFER_STORE_ITEM_STATUS_REFILLING {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   item.SubStatus,
			Currency: item.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_ITEM,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return item, err
}

func (dao OfferStoreDao) UpdateCancelRefillOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem) (bean.OfferStoreItem, error) {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferStoreItemPath(offer.UID)
	docRef := dbClient.Doc(offerPath)
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	batch := dbClient.Batch()

	batch.Set(itemDocRef, item.GetCancelOfferStoreItemRefill(), firestore.MergeAll)
	batch.Set(docRef, offer.GetUpdateOfferStoreChangeSnapshot(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return item, err
}

func (dao OfferStoreDao) RefillBalanceOfferStoreItem(offer bean.OfferStore, item *bean.OfferStoreItem, body bean.OfferStoreItem, offerType string) error {
	dbClient := firebase_service.FirestoreClient

	offerStoreRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	offerStoreItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		// Get From Wallet Balance
		walletDoc, err := tx.Get(offerStoreItemRef)
		if err != nil {
			return err
		}

		sellAmount := common.StringToDecimal(body.SellAmount)
		buyAmount := common.StringToDecimal(body.BuyAmount)

		buyBalance, err := common.ConvertToDecimal(walletDoc, "buy_balance")
		if err != nil {
			return err
		}
		sellBalance, err := common.ConvertToDecimal(walletDoc, "sell_balance")
		if err != nil {
			return err
		}

		if offerType == bean.OFFER_TYPE_BUY {
			buyBalance = buyBalance.Add(buyAmount)
			item.BuyBalance = buyBalance.String()
		} else {
			sellBalance = sellBalance.Add(sellAmount)
			item.SellBalance = sellBalance.String()
		}
		err = tx.Set(offerStoreItemRef, item.GetUpdateOfferStoreItemRefillBalance(), firestore.MergeAll)

		offer.ItemSnapshots[item.Currency] = *item
		err = tx.Set(offerStoreRef, offer.GetUpdateOfferStoreChangeSnapshot(), firestore.MergeAll)
		return err

	})
	return err
}

func (dao OfferStoreDao) RemoveOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))

	batch := dbClient.Batch()
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	batch.Delete(itemDocRef)
	batch.Set(docRef, offer.GetUpdateOfferStoreChangeItem(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreItemActive(offer bean.OfferStore, item bean.OfferStoreItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	docItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferStoreActive(), firestore.MergeAll)
	batch.Set(docItemRef, item.GetUpdateOfferStoreItemActive(), firestore.MergeAll)
	if item.SystemAddress != "" {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreItemClosing(offer bean.OfferStore, item bean.OfferStoreItem) error {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	offerPath := GetOfferStoreItemPath(offer.Id)
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	docItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	batch.Set(docRef, offer.GetUpdateOfferStoreChangeItem(), firestore.MergeAll)
	batch.Set(docItemRef, item.GetUpdateOfferStoreItemClosing(), firestore.MergeAll)

	if item.Currency == bean.ETH.Code && item.Status == bean.OFFER_STORE_ITEM_STATUS_CLOSING {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   item.Status,
			Currency: item.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreItemClosed(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	docItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetChangeStatus(), firestore.MergeAll)
	batch.Set(docItemRef, item.GetUpdateOfferStoreItemClosed(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) GetOfferStoreShake(offerId string, offerShakeId string) (t TransferObject) {
	GetObject(GetOfferStoreShakeItemPath(offerId, offerShakeId), &t, snapshotToOfferStoreShake)

	return
}

func (dao OfferStoreDao) GetOfferStoreShakeByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToOfferStoreShake)

	return
}

func (dao OfferStoreDao) ListOfferStoreShake(offerId string) ([]bean.OfferStoreShake, error) {
	dbClient := firebase_service.FirestoreClient
	collection := dbClient.Collection(GetOfferStoreShakePath(offerId))
	offerShakes := make([]bean.OfferStoreShake, 0)

	if collection != nil {
		iter := collection.Documents(context.Background())
		for {
			var offerShake bean.OfferStoreShake
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return offerShakes, err
			}
			doc.DataTo(&offerShake)
			offerShakes = append(offerShakes, offerShake)
		}
	}

	return offerShakes, nil
}

func (dao OfferStoreDao) AddOfferStoreShake(offer bean.OfferStore, offerShake bean.OfferStoreShake) (bean.OfferStoreShake, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetOfferStoreShakePath(offer.Id)).NewDoc()
	offerShake.Id = docRef.ID
	offerShake.OffChainId = fmt.Sprintf("%s-%s", offer.UID, offerShake.Id)

	batch := dbClient.Batch()
	offerStoreShake := GetOfferStoreShakeItemPath(offer.Id, offerShake.Id)
	if offerShake.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  offerShake.SystemAddress,
			Offer:    offerShake.Id,
			OfferRef: offerStoreShake,
			UID:      offerShake.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offerShake.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	batch.Set(docRef, offerShake.GetAddOfferStoreShake())

	if offerShake.Currency == bean.ETH.Code && (offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING || offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKING) {
		// Store a record to check onchain
		docId := strings.Replace(offerStoreShake, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offerShake.Status,
			Currency: offerShake.Currency,
			Offer:    offerShake.Id,
			OfferRef: offerStoreShake,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return offerShake, err
}

func (dao OfferStoreDao) UpdateOfferStoreShake(offerId string, offerShake bean.OfferStoreShake, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	offerShakePath := GetOfferStoreShakeItemPath(offerId, offerShake.Id)
	docRef := dbClient.Doc(offerShakePath)
	batch.Set(docRef, updateData, firestore.MergeAll)

	if offerShake.Currency == bean.ETH.Code &&
		(offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKING ||
			offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKING ||
			offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_CANCELLING ||
			offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTING ||
			offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETING) {

		// Store a record to check onchain
		docId := strings.Replace(offerShakePath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offerShake.Status,
			Currency: offerShake.Currency,
			Offer:    offerShake.Id,
			OfferRef: offerShakePath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
			UID:      offerId,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeReject(offer bean.OfferStore, offerShake bean.OfferStoreShake, profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	offerShakePath := GetOfferStoreShakeItemPath(offer.Id, offerShake.Id)
	docRef := dbClient.Doc(offerShakePath)

	batch := dbClient.Batch()
	batch.Set(docRef, offerShake.GetChangeStatus(), firestore.MergeAll)

	if offerShake.Currency == bean.ETH.Code && (offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTING) {
		// Store a record to check onchain
		docId := strings.Replace(offerShakePath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offerShake.Status,
			Currency: offerShake.Currency,
			Offer:    offerShake.Id,
			OfferRef: offerShakePath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeComplete(offer bean.OfferStore, offerShake bean.OfferStoreShake, profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	offerShakePath := GetOfferStoreShakeItemPath(offer.Id, offerShake.Id)
	docRef := dbClient.Doc(offerShakePath)

	batch := dbClient.Batch()
	batch.Set(docRef, offerShake.GetChangeStatus(), firestore.MergeAll)

	if offerShake.Currency == bean.ETH.Code && (offerShake.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETING) {
		// Store a record to check onchain
		docId := strings.Replace(offerShakePath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offerShake.Status,
			Currency: offerShake.Currency,
			Offer:    offerShake.Id,
			OfferRef: offerShakePath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeBalance(offer bean.OfferStore, item *bean.OfferStoreItem, offerShake bean.OfferStoreShake, shakeOrReject bool) error {
	dbClient := firebase_service.FirestoreClient

	offerStoreRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	offerStoreItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	offerStoreShakeRef := dbClient.Doc(GetOfferStoreShakeItemPath(offer.Id, offerShake.Id))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		// Get From Wallet Balance
		walletDoc, err := tx.Get(offerStoreItemRef)
		if err != nil {
			return err
		}
		amount, _ := decimal.NewFromString(offerShake.Amount)

		buyBalance, err := common.ConvertToDecimal(walletDoc, "buy_balance")
		if err != nil {
			return err
		}
		sellBalance, err := common.ConvertToDecimal(walletDoc, "sell_balance")
		if err != nil {
			return err
		}

		if offerShake.Type == bean.OFFER_TYPE_BUY {
			if shakeOrReject {
				// Shake, decrease
				buyBalance = buyBalance.Add(amount.Neg())
			} else {
				// Reject, increase
				buyBalance = buyBalance.Add(amount)
			}

			if buyBalance.LessThan(common.Zero) {
				return errors.New("Not enough balance")
			}

			item.BuyBalance = buyBalance.String()
		} else {
			if shakeOrReject {
				// Shake, decrease
				sellBalance = sellBalance.Add(amount.Neg())
			} else {
				// Shake, decrease
				sellBalance = sellBalance.Add(amount)
			}

			if sellBalance.LessThan(common.Zero) {
				return errors.New("Not enough balance")
			}

			item.SellBalance = sellBalance.String()
		}

		err = tx.Set(offerStoreShakeRef, offerShake.GetChangeStatus(), firestore.MergeAll)
		err = tx.Set(offerStoreItemRef, item.GetUpdateOfferStoreItemBalance(), firestore.MergeAll)

		offer.ItemSnapshots[item.Currency] = *item
		err = tx.Set(offerStoreRef, offer.GetUpdateOfferStoreChangeSnapshot(), firestore.MergeAll)
		return err

	})
	return err
}

func (dao OfferStoreDao) UpdateNotificationOfferStore(offer bean.OfferStore, item bean.OfferStoreItem) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationOfferStoreItemPath(offer.UID, offer.Id))
	err := ref.Set(context.Background(), item.GetNotificationUpdate(offer))

	return err
}

func (dao OfferStoreDao) UpdateNotificationOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationOfferStoreItemPath(offer.UID, offer.Id))
	err := ref.Set(context.Background(), item.GetNotificationUpdateItem(offer))

	return err
}

func (dao OfferStoreDao) UpdateNotificationOfferStoreShake(offerShake bean.OfferStoreShake, offer bean.OfferStore) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationOfferStoreShakeItemPath(offerShake.UID, offerShake.Id))
	err := ref.Set(context.Background(), offerShake.GetNotificationUpdate())
	ref2 := dbClient.NewRef(GetNotificationOfferStoreShakeItemPath(offer.UID, offerShake.Id))
	err = ref2.Set(context.Background(), offerShake.GetNotificationUpdate())

	return err
}

func (dao OfferStoreDao) GetOfferStoreReview(offerId string, id string) (t TransferObject) {
	GetObject(GetOfferStoreReviewItemPath(offerId, id), &t, snapshotToOfferStoreReview)

	return
}

func (dao OfferStoreDao) AddOfferStoreReview(offer bean.OfferStore, review bean.OfferStoreReview) error {
	dbClient := firebase_service.FirestoreClient

	offerStoreRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	offerStoreReviewRef := dbClient.Doc(GetOfferStoreReviewItemPath(offer.Id, review.Id))

	batch := dbClient.Batch()
	batch.Set(offerStoreRef, offer.GetUpdateOfferStoreReview(), firestore.MergeAll)
	batch.Set(offerStoreReviewRef, review.GetAddOfferStoreReview())

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) ListOfferStoreFreeStart(token string) ([]bean.OfferStoreFreeStart, error) {
	dbClient := firebase_service.FirestoreClient

	iter := dbClient.Collection(GetOfferStoreFreeStartPath()).Documents(context.Background())
	objs := make([]bean.OfferStoreFreeStart, 0)

	for {
		var obj bean.OfferStoreFreeStart
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objs, err
		}
		doc.DataTo(&obj)
		if obj.Token == token {
			objs = append(objs, obj)
		}
	}

	return objs, nil
}

func (dao OfferStoreDao) AddOfferStoreFreeStartUser(freeStart *bean.OfferStoreFreeStart, freeStartUser *bean.OfferStoreFreeStartUser) error {
	dbClient := firebase_service.FirestoreClient

	freeStartRef := dbClient.Doc(GetOfferStoreFreeStartItemPath(freeStart.Id))
	freeStartUserRef := dbClient.Doc(GetOfferStoreFreeStartUserItemPath(freeStartUser.UID))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		// Get From Wallet Balance
		freeStartDoc, err := tx.Get(freeStartRef)
		if err != nil {
			return err
		}

		value, err := freeStartDoc.DataAt("count")
		if err != nil {
			return err
		}
		count := value.(int64)
		count += 1
		if count > freeStart.Limit {
			return errors.New("Over limit")
		}
		freeStart.Count = count
		freeStartUser.Seq = count

		err = tx.Set(freeStartRef, freeStart.GetUpdateFreeStartCount(), firestore.MergeAll)
		if err != nil {
			return err
		}

		err = tx.Set(freeStartUserRef, freeStartUser.GetAddFreeStartUser(), firestore.MergeAll)
		return err

	})
	return err
}

func (dao OfferStoreDao) UpdateOfferStoreFreeStartUserDone(userId string) error {
	dbClient := firebase_service.FirestoreClient
	freeStartUserRef := dbClient.Doc(GetOfferStoreFreeStartUserItemPath(userId))
	_, err := freeStartUserRef.Set(context.Background(), bean.OfferStoreFreeStartUser{}.GetUpdateFreeStartUserDone(), firestore.MergeAll)

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreFreeStartUserUsing(userId string) error {
	dbClient := firebase_service.FirestoreClient
	freeStartUserRef := dbClient.Doc(GetOfferStoreFreeStartUserItemPath(userId))
	_, err := freeStartUserRef.Set(context.Background(), bean.OfferStoreFreeStartUser{}.GetUpdateFreeStartUserUsing(), firestore.MergeAll)

	return err
}

func (dao OfferStoreDao) GetOfferStoreFreeStart(level string) (t TransferObject) {
	GetObject(GetOfferStoreFreeStartItemPath(level), &t, snapshotToOfferStoreFreeStart)

	return
}

func (dao OfferStoreDao) GetOfferStoreFreeStartUser(userId string) (t TransferObject) {
	GetObject(GetOfferStoreFreeStartUserItemPath(userId), &t, snapshotToOfferStoreFreeStartUser)

	return
}

func (dao OfferStoreDao) UpdateOfferStoreLocationShakeTracking(userId string, offerShake bean.OfferStoreShake,
	offerLocation bean.OfferStoreLocationTracking, offerShakeLocation bean.OfferStoreLocationTracking) error {

	dbClient := firebase_service.FirestoreClient
	offerLocationRef := dbClient.Doc(GetOfferStoreLocationTrackingItemPath(userId, offerShake.Id, true))
	shakeLocationRef := dbClient.Doc(GetOfferStoreLocationTrackingItemPath(offerShake.UID, offerShake.Id, false))

	batch := dbClient.Batch()
	batch.Set(offerLocationRef, offerLocation.GetUpdateOfferStoreLocationShake(), firestore.MergeAll)
	batch.Set(shakeLocationRef, offerShakeLocation.GetUpdateOfferStoreLocationShake(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreLocationCompleteTracking(userId string, offerShake bean.OfferStoreShake,
	offerLocation bean.OfferStoreLocationTracking) error {

	dbClient := firebase_service.FirestoreClient
	offerLocationRef := dbClient.Doc(GetOfferStoreLocationTrackingItemPath(userId, offerShake.Id, true))
	shakeLocationRef := dbClient.Doc(GetOfferStoreLocationTrackingItemPath(offerShake.UID, offerShake.Id, false))

	batch := dbClient.Batch()
	batch.Set(offerLocationRef, offerLocation.GetUpdateOfferStoreLocationComplete(), firestore.MergeAll)
	batch.Delete(shakeLocationRef)

	_, err := batch.Commit(context.Background())

	return err
}

// DB path
//func GetOfferStorePath() string {
//	return "offer_stores"
//}

func GetOfferStorePath() string {
	return fmt.Sprintf("offer_stores")
}

func GetOfferStoreItemPath(id string) string {
	return fmt.Sprintf("offer_stores/%s", id)
}

func GetOfferStoreItemItemPath(id string, currency string) string {
	return fmt.Sprintf("offer_stores/%s/items/%s", id, currency)
}

func GetOfferStoreShakePath(offerStoreId string) string {
	return fmt.Sprintf("offer_stores/%s/shakes", offerStoreId)
}

func GetOfferStoreShakeItemPath(offerStoreId string, id string) string {
	return fmt.Sprintf("offer_stores/%s/shakes/%s", offerStoreId, id)
}

func GetOfferStoreReviewItemPath(offerStoreId string, id string) string {
	return fmt.Sprintf("offer_stores/%s/reviews/%s", offerStoreId, id)
}

func GetOfferStoreFreeStartPath() string {
	return fmt.Sprintf("offer_store_free_starts")
}

func GetOfferStoreFreeStartItemPath(level string) string {
	return fmt.Sprintf("offer_store_free_starts/%s", level)
}

//func GetOfferStoreFreeStartUserPath() string {
//	return fmt.Sprintf("offer_store_free_start_users")
//}

func GetOfferStoreFreeStartUserItemPath(userId string) string {
	return fmt.Sprintf("offer_store_free_start_users/%s", userId)
}

func GetOfferStoreLocationTrackingPath(userId string, offer bool) string {
	path := fmt.Sprintf("offer_store_location_tracking/%s/shake", userId)
	if offer {
		path = fmt.Sprintf("offer_store_location_tracking/%s/offer", userId)
	}
	return path
}

func GetOfferStoreLocationTrackingItemPath(userId string, offerShake string, offer bool) string {
	path := fmt.Sprintf("offer_store_location_tracking/%s/shake/%s", userId, offerShake)
	if offer {
		path = fmt.Sprintf("offer_store_location_tracking/%s/offer/%s", userId, offerShake)
	}
	return path
}

// Firebase
func GetNotificationOfferStoreItemPath(userId string, offerId string) string {
	return fmt.Sprintf("users/%s/offers/offer_store_%s", userId, offerId)
}

func GetNotificationOfferStoreShakeItemPath(userId string, offerId string) string {
	return fmt.Sprintf("users/%s/offers/offer_store_shake_%s", userId, offerId)
}

func snapshotToOfferStore(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStore
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferStoreItem(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreItem
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferStoreShake(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreShake
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferStoreReview(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreReview
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferStoreFreeStart(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreFreeStart
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferStoreFreeStartUser(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreFreeStartUser
	snapshot.DataTo(&obj)
	return obj
}
