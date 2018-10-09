package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"strings"
)

type CoinDao struct {
}

func (dao CoinDao) ListCoinCenter(country string) (t TransferObject) {
	ListObjects(GetCoinCenterCountryCurrenyPath(country), &t, nil, snapshotToCoinCenter)
	return
}

func (dao CoinDao) AddCoinOrder(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCoinOrderPath()).NewDoc()
	order.Id = docRef.ID
	// docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))

	refCode := strings.ToLower(order.Id[:6])
	orderRefCode := bean.CoinOrderRefCode{
		RefCode:  refCode,
		OrderRef: GetCashOrderItemPath(order.Id),
	}
	docOrderRefRef := dbClient.Doc(GetCoinOrderRefCodeItemPath(refCode))
	order.RefCode = refCode

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetAdd())
	// batch.Set(docUserRef, order.GetAdd())
	batch.Set(docOrderRefRef, orderRefCode.GetAdd())
	_, err := batch.Commit(context.Background())

	return err
}

func GetCoinCenterCountryCurrenyPath(country string) string {
	return fmt.Sprintf("coin_centers/%s/currency", country)
}

func GetCoinOrderPath() string {
	return fmt.Sprintf("coin_orders")
}

func GetCoinOrderItemPath(id string) string {
	return fmt.Sprintf("coin_orders/%s", id)
}

func GetCoinOrderUserItemPath(userId string, id string) string {
	return fmt.Sprintf("coin/%s/orders/%s", userId, id)
}

func GetNotificationCoinOrderPath(userId string, id string) string {
	return fmt.Sprintf("users/%s/coin/coin_order_%s", userId, id)
}

func GetCoinOrderRefCodeItemPath(refCode string) string {
	return fmt.Sprintf("coin_order_refs/%s", refCode)
}

func snapshotToCoinCenter(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinCenter
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCoinPayment(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinPayment
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCoinOrderRefCode(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinOrderRefCode
	snapshot.DataTo(&obj)
	return obj
}
