package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
)

type CashDao struct {
}

func (dao CashDao) GetCashStore(userId string) (t TransferObject) {
	GetObject(GetCashStorePath(userId), &t, snapshotToCashStore)
	return
}

func (dao CashDao) AddCashStore(cash *bean.CashStore) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(cash.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), cash.GetAdd())

	return err
}

func (dao CashDao) UpdateCashStore(cash *bean.CashStore) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(cash.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), cash.GetUpdate(), firestore.MergeAll)

	return err
}

func GetCashStorePath(userId string) string {
	return fmt.Sprintf("cash/%s", userId)
}

func snapshotToCashStore(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashStore
	snapshot.DataTo(&obj)
	return obj
}
