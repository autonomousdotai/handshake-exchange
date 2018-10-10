package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"strings"
)

type CoinDao struct {
}

func (dao CoinDao) ListCoinCenter(country string) (t TransferObject) {
	ListObjects(GetCoinCenterCountryCurrenyPath(country), &t, nil, snapshotToCoinCenter)
	return
}

func (dao CoinDao) GetCoinOrder(id string) (t TransferObject) {
	GetObject(GetCashOrderItemPath(id), &t, snapshotToCoinOrder)
	return
}

func (dao CoinDao) AddCoinOrder(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCoinOrderPath()).NewDoc()
	order.Id = docRef.ID
	docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))

	refCode := strings.ToLower(order.Id[:6])
	orderRefCode := bean.CoinOrderRefCode{
		RefCode:  refCode,
		OrderRef: GetCashOrderItemPath(order.Id),
		Duration: order.Duration,
	}
	docOrderRefRef := dbClient.Doc(GetCoinOrderRefCodeItemPath(refCode))
	order.RefCode = refCode

	docPoolRef := dbClient.Doc(GetCoinPoolItemPath(order.Currency))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		poolDoc, txErr := tx.Get(docPoolRef)
		if txErr != nil {
			return txErr
		}
		usage, txErr := common.ConvertToDecimal(poolDoc, "usage")
		if txErr != nil {
			return txErr
		}
		limit, txErr := common.ConvertToDecimal(poolDoc, "limit")
		if txErr != nil {
			return txErr
		}
		amount := common.StringToDecimal(order.Amount)
		usage = usage.Add(amount)
		if usage.GreaterThan(limit) {
			return errors.New("out of stock")
		}

		txErr = tx.Set(docRef, order.GetAdd(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(docUserRef, order.GetAdd(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(docOrderRefRef, orderRefCode.GetAdd(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(docPoolRef, bean.CoinPool{
			Usage: usage.String(),
		}.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		return txErr
	})

	return err
}

func (dao CoinDao) UpdateCoinStoreReceipt(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinOrderItemPath(order.Id))
	_, err := docRef.Set(context.Background(), order.GetReceiptUpdate(), firestore.MergeAll)

	return err
}

func (dao CoinDao) UpdateNotificationCoinOrder(order bean.CoinOrder) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationCoinOrderPath(order.UID, order.Id))
	err := ref.Set(context.Background(), order.GetNotificationUpdate())

	return err
}

func (dao CoinDao) GetCoinPool(currency string) (t TransferObject) {
	GetObject(GetCoinPoolItemPath(currency), &t, snapshotToCoinPool)
	return
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

func GetCoinPoolItemPath(currency string) string {
	return fmt.Sprintf("coin_pools/%s", currency)
}

func snapshotToCoinOrder(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinOrder
	snapshot.DataTo(&obj)
	return obj
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

func snapshotToCoinPool(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinPool
	snapshot.DataTo(&obj)
	return obj
}
