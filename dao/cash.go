package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/shopspring/decimal"
	"strings"
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
	GetObject(GetCashOrderItemPath(id), &t, snapshotToCashOrder)
	return
}

func (dao CashDao) GetCashOrderByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCashOrder)
	return
}

func (dao CashDao) ListCashOrders(status string, limit int, startAt interface{}) (t TransferObject) {
	ListPagingObjects(GetCashOrderPath(), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.OrderBy("created_at", firestore.Desc)
		if status != "" {
			query = query.Where("status", "==", status)
		}

		return query
	}, snapshotToCashOrder)

	return
}

func (dao CashDao) AddCashOrder(order *bean.CashOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCashOrderPath()).NewDoc()
	order.Id = docRef.ID
	docUserRef := dbClient.Doc(GetCashOrderUserItemPath(order.UID, order.Id))

	refCode := strings.ToLower(order.Id[:6])
	orderRefCode := bean.CashOrderRefCode{
		RefCode:  refCode,
		OrderRef: GetCashOrderItemPath(order.Id),
	}
	docOrderRefRef := dbClient.Doc(GetCashOrderRefCodeItemPath(refCode))
	order.RefCode = refCode

	batch := dbClient.Batch()
	batch.Set(docRef, order.GetAdd())
	batch.Set(docUserRef, order.GetAdd())
	batch.Set(docOrderRefRef, orderRefCode.GetAdd())
	_, err := batch.Commit(context.Background())

	return err
}

func (dao CashDao) FinishCashOrder(order *bean.CashOrder, cash *bean.CashStore) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCashOrderUserItemPath(order.UID, order.Id))
	cashRef := dbClient.Doc(GetCashStorePath(order.UID))
	docOrderRefRef := dbClient.Doc(GetCashOrderRefCodeItemPath(order.RefCode))

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		cashDoc, txErr := tx.Get(cashRef)
		if txErr != nil {
			return txErr
		}
		profit, txErr := common.ConvertToDecimal(cashDoc, "profit")
		if txErr != nil {
			if strings.Contains(txErr.Error(), "no field") {
				profit = common.Zero
			} else {
				return txErr
			}
		}
		storeProfit := common.StringToDecimal(order.StoreFee)
		profit = profit.Add(storeProfit)
		cash.Profit = profit.String()

		txErr = tx.Set(docRef, order.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(docUserRef, order.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(cashRef, cash.GetUpdateProfit(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Delete(docOrderRefRef)
		if txErr != nil {
			return txErr
		}

		return txErr
	})

	return err
}

func (dao CashDao) UpdateCashStoreReceipt(order *bean.CashOrder) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashOrderItemPath(order.Id))
	_, err := docRef.Set(context.Background(), order.GetReceiptUpdate(), firestore.MergeAll)

	return err
}

func (dao CashDao) UpdateNotificationCashOrder(order bean.CashOrder) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationCashOrderPath(order.UID, order.Id))
	err := ref.Set(context.Background(), order.GetNotificationUpdate())
	if order.ToUID != "" && err == nil {
		ref = dbClient.NewRef(GetNotificationCashOrderPath(order.ToUID, order.Id))
		err = ref.Set(context.Background(), order.GetNotificationUpdate())
	}

	return err
}

func (dao CashDao) ListCashCenter(country string) (t TransferObject) {
	ListObjects(GetCashCenterCountryPath(country), &t, nil, snapshotToCashCenter)
	return
}

func (dao CashDao) GetCashStorePayment(orderId string) (t TransferObject) {
	GetObject(GetCashStorePaymentItemPath(orderId), &t, snapshotToCashStorePayment)
	return
}

func (dao CashDao) AddCashStorePayment(payment *bean.CashStorePayment) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashStorePaymentItemPath(payment.Order))
	_, err := docRef.Set(context.Background(), payment.GetAdd())

	return err
}

func (dao CashDao) UpdateCashStorePayment(payment *bean.CashStorePayment, addAmount decimal.Decimal) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashStorePaymentItemPath(payment.Order))
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

func (dao CashDao) GetCashOrderRefCode(refCode string) (t TransferObject) {
	GetObject(GetCashOrderRefCodeItemPath(refCode), &t, snapshotToCashOrderRefCode)
	return
}

func GetCashStorePath(userId string) string {
	return fmt.Sprintf("cash/%s", userId)
}

func GetCashOrderPath() string {
	return fmt.Sprintf("cash_orders")
}

func GetCashOrderItemPath(id string) string {
	return fmt.Sprintf("cash_orders/%s", id)
}

func GetCashOrderUserItemPath(userId string, id string) string {
	return fmt.Sprintf("cash/%s/orders/%s", userId, id)
}

func GetNotificationCashOrderPath(userId string, id string) string {
	return fmt.Sprintf("users/%s/cash/cash_order_%s", userId, id)
}

func GetCashCenterCountryPath(country string) string {
	return fmt.Sprintf("cash_centers/%s/items", country)
}

func GetCashStorePaymentItemPath(orderId string) string {
	return fmt.Sprintf("cash_payments/%s", orderId)
}

func GetCashOrderRefCodeItemPath(refCode string) string {
	return fmt.Sprintf("cash_order_refs/%s", refCode)
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

func snapshotToCashCenter(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCenter
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCashStorePayment(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashStorePayment
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCashOrderRefCode(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashOrderRefCode
	snapshot.DataTo(&obj)
	return obj
}
