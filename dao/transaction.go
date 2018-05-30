package dao

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/bean"
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

func GetTransactionPath(userId string) string {
	return fmt.Sprintf("users/%s/transactions", userId)
}

func GetTransactionItemPath(userId string, id string) string {
	return fmt.Sprintf("users/%s/transactions/%s", userId, id)
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
