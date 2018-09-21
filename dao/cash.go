package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
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

func (dao CashDao) FinishCashOrder(order *bean.CashOrder, cash *bean.CashStore) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashOrderItemPath(order.Id))
	docUserRef := dbClient.Doc(GetCashOrderUserItemPath(order.UID, order.Id))
	cashRef := dbClient.Doc(GetCashStorePath(order.UID))

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

		return txErr
	})

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
