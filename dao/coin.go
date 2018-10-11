package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/shopspring/decimal"
	"strings"
)

type CoinDao struct {
}

func (dao CoinDao) ListCoinCenter(country string) (t TransferObject) {
	ListObjects(GetCoinCenterCountryCurrenyPath(country), &t, nil, snapshotToCoinCenter)
	return
}

func (dao CoinDao) GetCoinOrder(id string) (t TransferObject) {
	GetObject(GetCoinOrderItemPath(id), &t, snapshotToCoinOrder)
	return
}

func (dao CoinDao) GetCoinOrderByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCoinOrder)
	return
}

func (dao CoinDao) ListCoinOrders(status string, limit int, startAt interface{}) (t TransferObject) {
	ListPagingObjects(GetCoinOrderPath(), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.OrderBy("created_at", firestore.Desc)
		if status != "" {
			query = query.Where("status", "==", status)
		}

		return query
	}, snapshotToCashOrder)

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
		OrderRef: GetCoinOrderItemPath(order.Id),
		Order:    order.Id,
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
		if order.Type != bean.COIN_ORDER_TYPE_COD {
			txErr = tx.Set(docOrderRefRef, orderRefCode.GetAdd(), firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
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

func (dao CoinDao) CancelCoinOrder(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))
	docOrderRefRef := dbClient.Doc(GetCoinOrderRefCodeItemPath(order.RefCode))

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
		amount := common.StringToDecimal(order.Amount)
		usage = usage.Sub(amount)
		if usage.LessThan(common.Zero) {
			usage = common.Zero
		}

		txErr = tx.Set(docRef, order.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(docUserRef, order.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		if order.Type != bean.COIN_ORDER_TYPE_COD {
			txErr = tx.Delete(docOrderRefRef)
			if txErr != nil {
				return txErr
			}
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

func (dao CoinDao) UpdateCoinOrderReceipt(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))
	docOrderRefRef := dbClient.Doc(GetCoinOrderRefCodeItemPath(order.RefCode))

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetReceiptUpdate(), firestore.MergeAll)
	batch.Set(docUserRef, order.GetReceiptUpdate(), firestore.MergeAll)
	batch.Delete(docOrderRefRef)
	_, err := batch.Commit(context.Background())

	return err
}

func (dao CoinDao) UpdateCoinOrder(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetUpdate(), firestore.MergeAll)
	batch.Set(docUserRef, order.GetUpdate(), firestore.MergeAll)
	_, err := batch.Commit(context.Background())

	return err
}

func (dao CoinDao) FinishCoinOrder(order *bean.CoinOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCoinOrderUserItemPath(order.UID, order.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetUpdate(), firestore.MergeAll)
	batch.Set(docUserRef, order.GetUpdate(), firestore.MergeAll)
	_, err := batch.Commit(context.Background())

	return err
}

func (dao CoinDao) UpdateNotificationCoinOrder(order bean.CoinOrder) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationCoinOrderPath(order.UID, order.Id))
	err := ref.Set(context.Background(), order.GetNotificationUpdate())

	return err
}

func (dao CoinDao) ListCoinOrderRefCode() (t TransferObject) {
	ListObjects(GetCoinOrderRefCodePath(), &t, nil, snapshotToCoinOrderRefCode)
	return
}

func (dao CoinDao) GetCoinOrderRefCode(refCode string) (t TransferObject) {
	GetObject(GetCoinOrderRefCodeItemPath(refCode), &t, snapshotToCoinOrderRefCode)
	return
}

func (dao CoinDao) GetCoinPool(currency string) (t TransferObject) {
	GetObject(GetCoinPoolItemPath(currency), &t, snapshotToCoinPool)
	return
}

func (dao CoinDao) GetCoinPayment(orderId string) (t TransferObject) {
	GetObject(GetCoinPaymentItemPath(orderId), &t, snapshotToCoinPayment)
	return
}

func (dao CoinDao) AddCoinPayment(payment *bean.CoinPayment) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinPaymentItemPath(payment.Order))
	_, err := docRef.Set(context.Background(), payment.GetAdd())

	return err
}

func (dao CoinDao) UpdateCoinPayment(payment *bean.CoinPayment, addAmount decimal.Decimal) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCoinPaymentItemPath(payment.Order))
	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		paymentDoc, txErr := tx.Get(docRef)
		if txErr != nil {
			return txErr
		}
		amount, txErr := common.ConvertToDecimal(paymentDoc, "fiat_amount")
		if txErr != nil {
			return txErr
		}
		amount = amount.Add(addAmount)
		payment.FiatAmount = amount.String()

		txErr = tx.Set(docRef, payment.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		return txErr
	})

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

func GetCoinOrderRefCodePath() string {
	return fmt.Sprintf("coin_order_refs")
}

func GetCoinOrderRefCodeItemPath(refCode string) string {
	return fmt.Sprintf("coin_order_refs/%s", refCode)
}

func GetCoinPoolItemPath(currency string) string {
	return fmt.Sprintf("coin_pools/%s", currency)
}

func GetCoinPaymentItemPath(orderId string) string {
	return fmt.Sprintf("coin_payments/%s", orderId)
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
