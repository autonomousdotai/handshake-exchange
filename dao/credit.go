package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
)

type CreditDao struct {
}

func (dao CreditDao) GetCredit(userId string) (t TransferObject) {
	GetObject(GetCreditUserPath(userId), &t, snapshotToCredit)
	return
}

func (dao CreditDao) GetCreditItem(userId string, currency string) (t TransferObject) {
	GetObject(GetCreditItemItemPath(userId, currency), &t, snapshotToCreditItem)
	return
}

func (dao CreditDao) ListCreditItem(userId string) (t TransferObject) {
	ListObjects(GetCreditItemPath(userId), &t, nil, snapshotToCredit)
	return
}

func (dao CreditDao) AddCredit(credit bean.Credit) (bean.Credit, error) {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetAdd())

	return credit, err
}

func (dao CreditDao) UpdateCredit(credit bean.Credit) (bean.Credit, error) {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetUpdate())

	return credit, err
}

func (dao CreditDao) AddCreditItem(item bean.CreditItem) (bean.CreditItem, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetAdd())

	return item, err
}

func (dao CreditDao) UpdateCreditItem(item bean.CreditItem) (bean.CreditItem, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetUpdateStatus(), firestore.MergeAll)

	return item, err
}

func (dao CreditDao) DepositCreditItem(item *bean.CreditItem, deposit *bean.CreditDeposit) (err error) {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditDao) FinishDepositCreditItem(item *bean.CreditItem, deposit *bean.CreditDeposit,
	itemHistory *bean.CreditBalanceHistory,
	pool *bean.CreditPool, poolHistory *bean.CreditPoolBalanceHistory) (err error) {

	dbClient := firebase_service.FirestoreClient
	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		return txErr
	})

	return err
}

func GetCreditUserPath(userId string) string {
	return fmt.Sprintf("credits/%s", userId)
}

func GetCreditItemPath(userId string) string {
	return fmt.Sprintf("credits/%s/items", userId)
}

func GetCreditItemItemPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s", userId, currency)
}

func snapshotToCredit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.Credit
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCreditItem(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditItem
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCreditBalanceHistory(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditBalanceHistory
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCreditDeposit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditDeposit
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCreditWithdraw(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditWithdraw
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCreditPool(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditPool
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCreditPoolBalanceHistory(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditPoolBalanceHistory
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}
