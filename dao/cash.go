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

	cashPath := GetCashStorePath(cash.UID)
	docRef := dbClient.Doc(cashPath)
	_, err := docRef.Set(context.Background(), cash.GetAdd())

	return err
}

func (dao CashDao) UpdateCashStore(cash *bean.CashStore) error {
	dbClient := firebase_service.FirestoreClient

	cashPath := GetCashStorePath(cash.UID)
	docRef := dbClient.Doc(cashPath)
	_, err := docRef.Set(context.Background(), cash.GetUpdate(), firestore.MergeAll)

	return err
}

func (dao CashDao) GetCashOrder(id string) (t TransferObject) {
	GetObject(GetCashOrderItemPath(id), &t, snapshotToCashStore)
	return
}

func (dao CashDao) AddCashOrder(order *bean.CashOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCashOrderPath()).NewDoc()
	order.Id = docRef.ID
	docUserRef := dbClient.Doc(GetCashOrderUserItemPath(order.UID, order.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetAdd())
	batch.Set(docUserRef, order.GetAdd())
	_, err := batch.Commit(context.Background())

	return err
}

func GetCashStorePath(userId string) string {
	return fmt.Sprintf("cash/%s", userId)
}

func GetCashOrderPath() string {
	return fmt.Sprintf("cash_orders")
}

func GetCashOrderItemPath(userId string) string {
	return fmt.Sprintf("cash_orders/%s", userId)
}

func GetCashOrderUserItemPath(userId string, id string) string {
	return fmt.Sprintf("cash/%s/orders/%s", userId, id)
}

func snapshotToCashStore(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashStore
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCashOrder(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashOrder
	snapshot.DataTo(&obj)
	return obj
}
