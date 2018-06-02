package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
)

type OnChainDao struct {
}

func (dao OnChainDao) GetOfferInitEventBlock() (t TransferObject) {
	GetObject(GetOfferInitEventBlockPath(), &t, snapshotToOfferEventBlock)
	return
}

func (dao OnChainDao) UpdateOfferInitEventBlock(offer bean.OfferEventBlock) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetOfferInitEventBlockPath())

	_, err := docRef.Set(context.Background(), offer.GetUpdate())

	return err
}

func (dao OnChainDao) GetOfferShakeEventBlock() (t TransferObject) {
	GetObject(GetOfferShakeEventBlockPath(), &t, snapshotToOfferEventBlock)
	return
}

func (dao OnChainDao) UpdateOfferShakeEventBlock(offer bean.OfferEventBlock) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetOfferShakeEventBlockPath())

	_, err := docRef.Set(context.Background(), offer.GetUpdate())

	return err
}

func GetOfferInitEventBlockPath() string {
	return "onchain_events/offer_init"
}

func GetOfferShakeEventBlockPath() string {
	return "onchain_events/offer_shake"
}

func snapshotToOfferEventBlock(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferEventBlock
	snapshot.DataTo(&obj)
	return obj
}
