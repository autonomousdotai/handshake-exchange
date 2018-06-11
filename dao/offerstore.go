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
)

type OfferStoreDao struct {
}

func (dao OfferStoreDao) GetOfferStore(offerStoreId string) (t TransferObject) {
	GetObject(GetOfferStoreItemPath(offerStoreId), &t, snapshotToOfferStore)

	return
}

func (dao OfferStoreDao) AddOfferStore(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) (bean.OfferStore, error) {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.UID))
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
			OfferRef: GetOfferStoreItemPath(offer.Id),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}
	batch.Set(itemDocRef, item.GetAddOfferStoreItem())

	batch.Set(docRef, offer.GetAddOfferStore())
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

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

func (dao OfferStoreDao) AddOfferStoreItem(offer bean.OfferStore, item bean.OfferStoreItem, profile bean.Profile) (bean.OfferStoreItem, error) {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))

	batch := dbClient.Batch()
	itemDocRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, item.Currency))
	if item.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  item.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferStoreItemPath(offer.Id),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(item.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	batch.Set(itemDocRef, item.GetAddOfferStoreItem())
	batch.Set(docRef, offer.GetUpdateOfferStoreChangeItem(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferStoreProfile(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return item, err
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

func (dao OfferDao) UpdateOfferStoreItemActive(offer bean.OfferStore, offerItem bean.OfferStoreItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	docItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, offerItem.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferStoreActive(), firestore.MergeAll)
	batch.Set(docItemRef, offerItem.GetUpdateOfferStoreItemActive(), firestore.MergeAll)
	if offerItem.SystemAddress != "" {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offerItem.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreItemClosing(offer bean.OfferStore, offerItem bean.OfferStoreItem) error {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()
	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	docItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offer.Id, offerItem.Currency))

	batch.Set(docRef, offer.GetUpdateOfferStoreChangeItem(), firestore.MergeAll)
	batch.Set(docItemRef, offerItem.GetUpdateOfferStoreItemClosing(), firestore.MergeAll)
	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) GetOfferStoreShake(offerStoreId string, offerStoreShakeId string) (t TransferObject) {
	GetObject(GetOfferStoreShakeItemPath(offerStoreId, offerStoreShakeId), &t, snapshotToOfferStoreShake)

	return
}

func (dao OfferStoreDao) AddOfferStoreShake(offerStore bean.OfferStore, shake bean.OfferStoreShake) (bean.OfferStoreShake, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetOfferStoreShakePath(offerStore.Id)).NewDoc()
	shake.Id = docRef.ID

	batch := dbClient.Batch()
	if shake.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  shake.SystemAddress,
			Offer:    shake.Id,
			OfferRef: GetOfferStoreShakeItemPath(offerStore.Id, shake.Id),
			UID:      shake.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER_STORE,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(shake.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	batch.Set(docRef, shake.GetAddOfferStoreShake())

	_, err := batch.Commit(context.Background())

	return shake, err
}

func (dao OfferStoreDao) UpdateOfferStoreShake(offerStoreId string, offer bean.OfferStoreShake, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreShakeItemPath(offerStoreId, offer.Id))
	_, err := docRef.Set(context.Background(), updateData, firestore.MergeAll)

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeReject(offerStore bean.OfferStore, offer bean.OfferStoreShake, profile bean.Profile, transactionCount bean.TransactionCount) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreShakeItemPath(offerStore.Id, offer.Id))
	transCountDocRef := dbClient.Doc(GetTransactionCountItemPath(profile.UserId, transactionCount.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetChangeStatus(), firestore.MergeAll)
	batch.Set(transCountDocRef, transactionCount.GetUpdateFailed(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeComplete(offerStore bean.OfferStore, offer bean.OfferStoreShake, profile bean.Profile,
	transactionCount1 bean.TransactionCount, transactionCount2 bean.TransactionCount) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreShakeItemPath(offerStore.Id, offer.Id))
	transCountDocRef1 := dbClient.Doc(GetTransactionCountItemPath(profile.UserId, transactionCount1.Currency))
	transCountDocRef2 := dbClient.Doc(GetTransactionCountItemPath(profile.UserId, transactionCount2.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetChangeStatus(), firestore.MergeAll)
	batch.Set(transCountDocRef1, transactionCount1.GetUpdateFailed(), firestore.MergeAll)
	batch.Set(transCountDocRef2, transactionCount2.GetUpdateFailed(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferStoreDao) UpdateOfferStoreShakeBalance(offerStore bean.OfferStore, offerStoreItem *bean.OfferStoreItem, offerStoreShake bean.OfferStoreShake, shakeOrReject bool) error {
	dbClient := firebase_service.FirestoreClient

	offerStoreItemRef := dbClient.Doc(GetOfferStoreItemItemPath(offerStore.Id, offerStoreItem.Currency))
	offerStoreShakeRef := dbClient.Doc(GetOfferStoreShakeItemPath(offerStore.Id, offerStoreShake.Id))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		// Get From Wallet Balance
		walletDoc, err := tx.Get(offerStoreItemRef)
		if err != nil {
			return err
		}
		amount, _ := decimal.NewFromString(offerStoreShake.Amount)

		buyBalance, err := common.ConvertToDecimal(walletDoc, "buy_balance")
		if err != nil {
			return err
		}
		sellBalance, err := common.ConvertToDecimal(walletDoc, "sell_balance")
		if err != nil {
			return err
		}

		if offerStoreShake.Type == bean.OFFER_TYPE_BUY {
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

			offerStoreItem.BuyBalance = buyBalance.String()
		} else {
			if shakeOrReject {
				// Shake, decrease
				sellBalance = buyBalance.Add(amount.Neg())
			} else {
				// Shake, decrease
				sellBalance = buyBalance.Add(amount)
			}

			if sellBalance.LessThan(common.Zero) {
				return errors.New("Not enough balance")
			}

			offerStoreItem.SellBalance = sellBalance.String()
		}

		err = tx.Set(offerStoreShakeRef, offerStoreShake.GetChangeStatus(), firestore.MergeAll)
		err = tx.Set(offerStoreItemRef, offerStoreItem.GetUpdateOfferStoreItemBalance(), firestore.MergeAll)
		return err

	})
	return err
}

// DB path
func GetOfferStorePath() string {
	return "offer_stores"
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
