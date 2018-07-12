package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
)

type TransactionDao struct {
}

func (dao TransactionDao) ListTransactions(userId string, transType string, currency string, limit int, startAt interface{}) (t TransferObject) {
	ListPagingObjects(GetTransactionPath(userId), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.OrderBy("created_at", firestore.Desc)
		return query
	}, snapshotToTransaction)

	return
}

func (dao TransactionDao) GetTransaction(userId string, transId string) TransferObject {
	return dao.GetTransactionByPath(GetTransactionItemPath(userId, transId))
}

func (dao TransactionDao) GetTransactionByPath(path string) (t TransferObject) {
	// users/{uid}/transactions/{id}
	GetObject(path, &t, snapshotToTransaction)
	return
}

func (dao TransactionDao) GetTransactionCount(userId string, currency string) TransferObject {
	to := dao.GetTransactionCountByPath(GetTransactionCountItemPath(userId, currency))
	if !to.Found {
		to.Object = bean.TransactionCount{
			Currency:        currency,
			Success:         0,
			Failed:          0,
			Pending:         0,
			BuyAmount:       common.Zero.String(),
			SellAmount:      common.Zero.String(),
			BuyFiatAmounts:  map[string]bean.TransactionFiatAmount{},
			SellFiatAmounts: map[string]bean.TransactionFiatAmount{},
		}
		to.Found = true
	} else {
		transCount := to.Object.(bean.TransactionCount)
		if transCount.SellFiatAmounts == nil {
			transCount.SellFiatAmounts = map[string]bean.TransactionFiatAmount{}
		}
		if transCount.BuyFiatAmounts == nil {
			transCount.BuyFiatAmounts = map[string]bean.TransactionFiatAmount{}
		}
	}

	return to
}

func (dao TransactionDao) UpdateTransactionCount(userId string, currency string, txCountData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetTransactionCountItemPath(userId, currency))
	_, err := docRef.Set(context.Background(), txCountData, firestore.MergeAll)

	return err
}

func (dao TransactionDao) GetTransactionCountByPath(path string) (t TransferObject) {
	// users/{uid}/transaction_counts/{currency}
	GetObject(path, &t, snapshotToTransactionCount)
	return
}

func (dao TransactionDao) ListTransactionCounts(userId string) (t TransferObject) {
	ListObjects(GetTransactionCountPath(userId), &t, nil, snapshotToTransactionCount)

	return
}

func GetTransactionPath(userId string) string {
	return fmt.Sprintf("users/%s/transactions", userId)
}

func GetTransactionItemPath(userId string, id string) string {
	return fmt.Sprintf("users/%s/transactions/%s", userId, id)
}

func GetTransactionCountPath(userId string) string {
	return fmt.Sprintf("users/%s/transaction_counts", userId)
}

func GetTransactionCountItemPath(userId string, currency string) string {
	return fmt.Sprintf("users/%s/transaction_counts/%s", userId, currency)
}

func snapshotToTransaction(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.Transaction
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID

	return obj
}

func snapshotToTransactionCount(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.TransactionCount
	snapshot.DataTo(&obj)

	return obj
}
