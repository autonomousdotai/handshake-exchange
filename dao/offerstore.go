package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
)

type OfferStoreDao struct {
}

func (dao OfferStoreDao) GetOfferStore(offerStoreId string) (t TransferObject) {
	GetObject(GetOfferStoreItemPath(offerStoreId), &t, snapshotToOfferStore)

	return
}

func (dao OfferStoreDao) UpdateOfferStore(offer bean.OfferStore, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreItemPath(offer.Id))
	_, err := docRef.Set(context.Background(), updateData, firestore.MergeAll)

	return err
}

func (dao OfferStoreDao) GetOfferStoreShake(offerStoreId string, offerStoreShakeId string) (t TransferObject) {
	GetObject(GetOfferStoreShakeItemPath(offerStoreId, offerStoreShakeId), &t, snapshotToOfferStoreShake)

	return
}

func (dao OfferStoreDao) UpdateOfferStoreShake(offerStoreId string, offer bean.OfferStoreShake, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferStoreShakeItemPath(offerStoreId, offer.Id))
	_, err := docRef.Set(context.Background(), updateData, firestore.MergeAll)

	return err
}

// DB path
func GetOfferStorePath() string {
	return "offer_stores"
}

func GetOfferStoreItemPath(id string) string {
	return fmt.Sprintf("offer_stores/%s", id)
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

func snapshotToOfferStoreShake(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferStoreShake
	snapshot.DataTo(&obj)
	return obj
}
