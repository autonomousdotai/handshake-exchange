package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/ninjadotorg/handshake-exchange/service/cache"
	"strings"
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
	ListObjects(GetCreditItemPath(userId), &t, nil, snapshotToCreditItem)
	return
}

func (dao CreditDao) AddCredit(credit *bean.Credit) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetAdd())

	return err
}

func (dao CreditDao) UpdateCredit(credit *bean.Credit) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetUpdate(), firestore.MergeAll)

	return err
}

func (dao CreditDao) AddCreditItem(item *bean.CreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetAdd())

	return err
}

func (dao CreditDao) UpdateCreditItem(item *bean.CreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetUpdateStatus(), firestore.MergeAll)

	return err
}

func (dao CreditDao) UpdateCreditItemReactivate(item *bean.CreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetUpdateReactivate(), firestore.MergeAll)

	return err
}

func (dao CreditDao) GetCreditDeposit(currency string, depositId string) (t TransferObject) {
	t = dao.GetCreditDepositByPath(GetCreditDepositItemPath(currency, depositId))
	return
}

func (dao CreditDao) GetCreditDepositByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCreditDeposit)
	return
}

func (dao CreditDao) AddCreditDeposit(item *bean.CreditItem, deposit *bean.CreditDeposit) (err error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCreditDepositPath(deposit.Currency)).NewDoc()
	deposit.Id = docRef.ID
	docUserRef := dbClient.Doc(GetCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, deposit.GetAdd())
	batch.Set(docUserRef, deposit.GetAdd())
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditDao) FinishDepositCreditItem(item *bean.CreditItem, deposit *bean.CreditDeposit,
	itemHistory *bean.CreditBalanceHistory,
	pool *bean.CreditPool, poolOrder *bean.CreditPoolOrder, poolHistory *bean.CreditPoolBalanceHistory,
	tracking *bean.CreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient
	itemDocRef := dbClient.Doc(GetCreditItemItemPath(deposit.UID, deposit.Currency))
	depositUserDocRef := dbClient.Doc(GetCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))
	depositDocRef := dbClient.Doc(GetCreditDepositItemPath(deposit.Currency, deposit.Id))

	poolDocRef := dbClient.Doc(GetCreditPoolItemPath(deposit.Currency, pool.Level))
	poolOrderDocRef := dbClient.Doc(GetCreditPoolItemOrderItemPath(deposit.Currency, pool.Level, poolOrder.Id))
	poolOrderUserDocRef := dbClient.Doc(GetCreditPoolItemOrderItemUserPath(deposit.Currency, poolOrder.UID, poolOrder.Id))

	balanceHistoryDocRef := dbClient.Collection(GetCreditBalanceHistoryPath(deposit.UID, deposit.Currency)).NewDoc()
	itemHistory.Id = balanceHistoryDocRef.ID
	poolBalanceHistoryDocRef := dbClient.Collection(GetCreditPoolBalanceHistoryPath(deposit.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	docLogRef := dbClient.Doc(GetCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	amount := common.StringToDecimal(deposit.Amount)
	poolOrderUserDocRefs := make([]*firestore.DocumentRef, 0)
	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)

	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		itemDoc, txErr := tx.Get(itemDocRef)
		if txErr != nil {
			return txErr
		}
		itemBalance, txErr := common.ConvertToDecimal(itemDoc, "balance")
		if txErr != nil {
			return txErr
		}
		itemHistory.Old = itemBalance.String()
		itemHistory.Change = amount.String()

		poolDoc, txErr := tx.Get(poolDocRef)
		if err != nil {
			return txErr
		}
		poolBalance, txErr := common.ConvertToDecimal(poolDoc, "balance")
		if txErr != nil {
			return txErr
		}
		poolHistory.Old = poolBalance.String()
		poolHistory.Change = amount.String()

		itemBalance = itemBalance.Add(amount)
		item.Balance = itemBalance.String()
		itemHistory.New = item.Balance

		poolBalance = poolBalance.Add(amount)
		pool.Balance = poolBalance.String()
		poolHistory.New = pool.Balance

		// Update balance
		txErr = tx.Set(itemDocRef, item.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(poolDocRef, pool.GetUpdateBalance(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		// Update status Deposit
		txErr = tx.Set(depositUserDocRef, deposit.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(depositDocRef, deposit.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		// Insert pool order
		if amount.GreaterThanOrEqual(common.Zero) {
			txErr = tx.Set(poolOrderDocRef, poolOrder.GetAdd())
			if txErr != nil {
				return txErr
			}
			txErr = tx.Set(poolOrderUserDocRef, poolOrder.GetAdd())
			if txErr != nil {
				return txErr
			}
		} else {
			// Remove all order of this user
			for i, itemDocRef := range poolOrderUserDocRefs {
				txErr = tx.Delete(itemDocRef)
				if txErr != nil {
					return txErr
				}
				txErr = tx.Delete(poolOrderDocRefs[i])
				if txErr != nil {
					return txErr
				}
			}
		}

		// Insert history
		txErr = tx.Set(balanceHistoryDocRef, itemHistory.GetAdd())
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(poolBalanceHistoryDocRef, poolHistory.GetAdd())
		if txErr != nil {
			return txErr
		}

		// Update tracking
		tx.Delete(docTrackingRef)
		tx.Set(docLogRef, tracking.GetUpdate(), firestore.MergeAll)

		return txErr
	})

	if err != nil {
		dao.SetCreditPoolCache(*pool)
	}

	return err
}

func (dao CreditDao) RemoveCreditItem(item *bean.CreditItem, itemHistory *bean.CreditBalanceHistory,
	pool *bean.CreditPool, poolOrders []bean.CreditPoolOrder, poolHistory *bean.CreditPoolBalanceHistory) (err error) {

	dbClient := firebase_service.FirestoreClient
	itemDocRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))

	poolDocRef := dbClient.Doc(GetCreditPoolItemPath(item.Currency, pool.Level))

	balanceHistoryDocRef := dbClient.Collection(GetCreditBalanceHistoryPath(item.UID, item.Currency)).NewDoc()
	itemHistory.Id = balanceHistoryDocRef.ID
	poolBalanceHistoryDocRef := dbClient.Collection(GetCreditPoolBalanceHistoryPath(item.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)
	poolOrderUserDocRefs := make([]*firestore.DocumentRef, 0)

	for _, order := range poolOrders {
		poolOrderDocRef := dbClient.Doc(GetCreditPoolItemOrderItemPath(item.Currency, pool.Level, order.Id))
		poolOrderUserDocRef := dbClient.Doc(GetCreditPoolItemOrderItemUserPath(item.Currency, order.UID, order.Id))

		poolOrderDocRefs = append(poolOrderDocRefs, poolOrderDocRef)
		poolOrderUserDocRefs = append(poolOrderUserDocRefs, poolOrderUserDocRef)
	}

	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		zeroStr := common.Zero.String()

		itemDoc, txErr := tx.Get(itemDocRef)
		if txErr != nil {
			return txErr
		}
		itemBalance, txErr := common.ConvertToDecimal(itemDoc, "balance")
		if txErr != nil {
			return txErr
		}
		itemHistory.Old = itemBalance.String()
		itemHistory.Change = itemBalance.Neg().String()

		poolDoc, txErr := tx.Get(poolDocRef)
		if err != nil {
			return txErr
		}
		poolBalance, txErr := common.ConvertToDecimal(poolDoc, "balance")
		if txErr != nil {
			return txErr
		}
		poolHistory.Old = poolBalance.String()
		poolHistory.Change = itemBalance.Neg().String()

		item.Balance = zeroStr
		itemHistory.New = item.Balance

		poolBalance = poolBalance.Add(itemBalance)
		pool.Balance = poolBalance.String()
		poolHistory.New = pool.Balance

		// Update balance
		txErr = tx.Set(itemDocRef, item.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(poolDocRef, pool.GetUpdateBalance(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		// Remove all order of this user
		for i, itemDocRef := range poolOrderUserDocRefs {
			txErr = tx.Set(itemDocRef, map[string]interface{}{
				"deleted": true,
			}, firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
			txErr = tx.Delete(poolOrderDocRefs[i])
			if txErr != nil {
				return txErr
			}
		}

		// Insert history
		txErr = tx.Set(balanceHistoryDocRef, itemHistory.GetAdd())
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(poolBalanceHistoryDocRef, poolHistory.GetAdd())
		if txErr != nil {
			return txErr
		}

		return txErr
	})

	if err != nil {
		dao.SetCreditPoolCache(*pool)
	}

	return err
}

func (dao CreditDao) RemoveCreditOnChainActionTracking(tracking bean.CreditOnChainActionTracking) error {
	dbClient := firebase_service.FirestoreClient
	docTrackingRef := dbClient.Doc(GetCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))
	_, err := docTrackingRef.Delete(context.Background())
	return err
}

func (dao CreditDao) ListPendingCreditTransaction(currency string) (t TransferObject) {
	ListObjects(GetCreditPendingTransactionPath(currency), &t, nil, snapshotToCreditTransaction)
	return
}

func (dao CreditDao) GetCreditTransaction(currency string, id string) (t TransferObject) {
	GetObject(GetCreditTransactionItemPath(currency, id), &t, snapshotToCreditTransaction)
	return
}

func (dao CreditDao) GetCreditTransactionUser(userId string, currency string, id string) (t TransferObject) {
	GetObject(GetCreditTransactionItemUserPath(userId, currency, id), &t, snapshotToCreditTransaction)
	return
}

func (dao CreditDao) AddCreditTransaction(pool *bean.CreditPool, trans *bean.CreditTransaction,
	userTransList []*bean.CreditTransaction, selectedOrders []bean.CreditPoolOrder) (err error) {

	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCreditPoolItemPath(pool.Currency, pool.Level))
	transDocRef := dbClient.Collection(GetCreditTransactionPath(pool.Currency)).NewDoc()
	trans.Id = transDocRef.ID
	//pendingTransDocRef := dbClient.Doc(GetCreditPendingTransactionItemPath(pool.Currency, trans.Id))

	transUserDocRefs := make([]*firestore.DocumentRef, 0)
	for userTransIndex, userTrans := range userTransList {
		userTransPath := GetCreditTransactionItemUserPath(userTrans.UID, pool.Currency, trans.Id)
		transUserDocRefs = append(transUserDocRefs, dbClient.Doc(userTransPath))
		userTransList[userTransIndex].Id = trans.Id
	}

	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)
	for _, creditOrder := range selectedOrders {
		orderPath := GetCreditPoolItemOrderItemPath(pool.Currency, pool.Level, creditOrder.Id)
		poolOrderDocRef := dbClient.Doc(orderPath)
		poolOrderDocRefs = append(poolOrderDocRefs, poolOrderDocRef)
	}

	amount := common.StringToDecimal(trans.Amount)

	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		itemDoc, txErr := tx.Get(poolDocRef)
		if txErr != nil {
			return txErr
		}
		itemBalance, txErr := common.ConvertToDecimal(itemDoc, "balance")
		if txErr != nil {
			return txErr
		}
		itemCaptured, txErr := common.ConvertToDecimal(itemDoc, "captured_balance")
		if txErr != nil {
			if strings.Contains(txErr.Error(), "no field") {
				itemCaptured = common.Zero
			} else {
				return txErr
			}
		}

		for creditOrderIndex, creditOrder := range selectedOrders {
			orderDoc, txErr := tx.Get(poolOrderDocRefs[creditOrderIndex])
			if txErr != nil {
				return txErr
			}
			orderBalance, txErr := common.ConvertToDecimal(orderDoc, "balance")
			if txErr != nil {
				return txErr
			}
			orderCapturedBalance, txErr := common.ConvertToDecimal(orderDoc, "captured_balance")
			if txErr != nil {
				if strings.Contains(txErr.Error(), "no field") {
					orderCapturedBalance = common.Zero
				} else {
					return txErr
				}
			}
			//if orderCapturedBalance.GreaterThan(common.Zero) {
			//	// This is dirty item, stop transaction
			//	return errors.New("out of stock")
			//}
			//if orderBalance.Sub(common.StringToDecimal(creditOrder.CapturedBalance)).LessThan(common.Zero) {
			//	// This is dirty item, stop transaction
			//	return errors.New("out of stock")
			//}

			orderCapturedBalance = orderCapturedBalance.Add(creditOrder.CapturedAmount)
			if orderCapturedBalance.GreaterThan(orderBalance) {
				return errors.New("out of stock")
			}
			creditOrder.CapturedBalance = orderCapturedBalance.String()
			selectedOrders[creditOrderIndex] = creditOrder
		}

		itemCaptured = itemCaptured.Add(amount)
		if itemCaptured.GreaterThan(itemBalance) {
			return errors.New("out of stock")
		}

		pool.CapturedBalance = itemCaptured.String()
		txErr = tx.Set(poolDocRef, pool.GetUpdateCapturedBalance(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		txErr = tx.Set(transDocRef, trans.GetAdd())
		if txErr != nil {
			return txErr
		}
		//txErr = tx.Set(pendingTransDocRef, trans.GetAdd())
		//if txErr != nil {
		//	return txErr
		//}

		for transIndex, userTransDocRef := range transUserDocRefs {
			txErr = tx.Set(userTransDocRef, userTransList[transIndex].GetAdd())
			if txErr != nil {
				return txErr
			}
		}

		for orderIndex, poolOrderDocRef := range poolOrderDocRefs {
			// Miss match ID somehow
			if selectedOrders[orderIndex].Id != poolOrderDocRef.ID {
				// This is dirty item, stop transaction
				return errors.New("out of stock")
			}
			txErr = tx.Set(poolOrderDocRef, selectedOrders[orderIndex].GetUpdateCapture(), firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
		}

		return txErr
	})

	return err
}

func (dao CreditDao) FinishCreditTransaction(pool *bean.CreditPool, poolHistory bean.CreditPoolBalanceHistory,
	items []bean.CreditItem, itemHistories []bean.CreditBalanceHistory, poolOrders []bean.CreditPoolOrder,
	trans *bean.CreditTransaction, transList []*bean.CreditTransaction) (err error) {

	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCreditPoolItemPath(pool.Currency, pool.Level))
	poolBalanceHistoryDocRef := dbClient.Collection(GetCreditPoolBalanceHistoryPath(pool.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	creditDocRefs := make([]*firestore.DocumentRef, 0)
	itemDocRefs := make([]*firestore.DocumentRef, 0)
	itemHistoryDocRefs := make([]*firestore.DocumentRef, 0)

	transDocRef := dbClient.Doc(GetCreditTransactionItemPath(trans.Currency, trans.Id))
	transUserDocRefs := make([]*firestore.DocumentRef, 0)
	for itemIndex, item := range items {
		creditDocRef := dbClient.Doc(GetCreditUserPath(item.UID))
		itemDocRef := dbClient.Doc(GetCreditItemItemPath(item.UID, item.Currency))
		balanceHistoryDocRef := dbClient.Collection(GetCreditBalanceHistoryPath(item.UID, item.Currency)).NewDoc()

		itemHistories[itemIndex].Id = balanceHistoryDocRef.ID

		creditDocRefs = append(creditDocRefs, creditDocRef)
		itemDocRefs = append(itemDocRefs, itemDocRef)
		itemHistoryDocRefs = append(itemHistoryDocRefs, balanceHistoryDocRef)

		transUserDocRef := dbClient.Doc(GetCreditTransactionItemUserPath(transList[itemIndex].UID, transList[itemIndex].Currency, transList[itemIndex].Id))
		transUserDocRefs = append(transUserDocRefs, transUserDocRef)
	}

	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)
	poolOrderUserDocRefs := make([]*firestore.DocumentRef, 0)

	for _, poolOrder := range poolOrders {
		poolOrderDocRef := dbClient.Doc(GetCreditPoolItemOrderItemPath(pool.Currency, pool.Level, poolOrder.Id))
		poolOrderUserDocRef := dbClient.Doc(GetCreditPoolItemOrderItemUserPath(pool.Currency, poolOrder.UID, poolOrder.Id))

		poolOrderDocRefs = append(poolOrderDocRefs, poolOrderDocRef)
		poolOrderUserDocRefs = append(poolOrderUserDocRefs, poolOrderUserDocRef)
	}

	amount := common.StringToDecimal(trans.Amount)

	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		poolDoc, txErr := tx.Get(poolDocRef)
		if err != nil {
			return txErr
		}
		poolBalance, txErr := common.ConvertToDecimal(poolDoc, "balance")
		if txErr != nil {
			return txErr
		}
		poolHistory.Old = poolBalance.String()
		poolBalance = poolBalance.Sub(amount)
		pool.Balance = poolBalance.String()
		if poolBalance.LessThan(common.Zero) {
			return errors.New("invalid balance")
		}
		poolHistory.New = pool.Balance

		poolCapturedBalance, txErr := common.ConvertToDecimal(poolDoc, "captured_balance")
		if txErr != nil {
			return txErr
		}
		poolCapturedBalance = poolCapturedBalance.Sub(amount)
		pool.CapturedBalance = poolCapturedBalance.String()
		if poolCapturedBalance.LessThan(common.Zero) {
			return errors.New("invalid balance")
		}

		transDoc, txErr := tx.Get(transDocRef)
		if err != nil {
			return txErr
		}
		transStatus, txErr := transDoc.DataAt("status")
		if err != nil {
			return txErr
		}
		if transStatus.(string) == bean.CREDIT_TRANSACTION_STATUS_CREATE {
			trans.Status = bean.CREDIT_TRANSACTION_STATUS_SUCCESS
		} else {
			return errors.New("invalid status")
		}
		for itemIndex, itemDocRef := range itemDocRefs {
			itemDoc, txErr := tx.Get(itemDocRef)
			if err != nil {
				return txErr
			}
			itemBalance, txErr := common.ConvertToDecimal(itemDoc, "balance")
			if txErr != nil {
				return txErr
			}
			sold, txErr := common.ConvertToDecimal(itemDoc, "sold")
			if txErr != nil {
				return txErr
			}
			creditRevenue, txErr := common.ConvertToDecimal(itemDoc, "revenue")
			if txErr != nil {
				return txErr
			}
			itemRevenue, txErr := common.ConvertToDecimal(itemDoc, "revenue")
			if txErr != nil {
				return txErr
			}
			itemAmount := common.StringToDecimal(itemHistories[itemIndex].Change)
			revenue := common.StringToDecimal(transList[itemIndex].Revenue)
			itemHistories[itemIndex].Old = itemBalance.String()

			itemBalance = itemBalance.Sub(itemAmount)
			items[itemIndex].Balance = itemBalance.String()
			sold = sold.Add(itemAmount)
			items[itemIndex].Sold = sold.String()
			creditRevenue = creditRevenue.Add(revenue)
			itemRevenue = itemRevenue.Add(revenue)
			items[itemIndex].CreditRevenue = creditRevenue.String()
			items[itemIndex].Revenue = revenue.String()

			if itemBalance.LessThan(common.Zero) {
				return errors.New("invalid balance")
			}
			itemHistories[itemIndex].New = items[itemIndex].Balance
		}

		for orderIndex, orderDocRef := range poolOrderDocRefs {
			orderDoc, txErr := tx.Get(orderDocRef)
			if err != nil {
				return txErr
			}
			orderBalance, txErr := common.ConvertToDecimal(orderDoc, "balance")
			if txErr != nil {
				return txErr
			}
			capturedBalance, txErr := common.ConvertToDecimal(orderDoc, "captured_balance")
			if txErr != nil {
				return txErr
			}
			capturedAmount := poolOrders[orderIndex].CapturedAmount
			orderBalance = orderBalance.Sub(capturedAmount)
			capturedBalance = capturedBalance.Sub(capturedAmount)

			poolOrders[orderIndex].Balance = orderBalance.String()
			if orderBalance.LessThan(common.Zero) {
				return errors.New("invalid balance")
			}
			poolOrders[orderIndex].CapturedBalance = capturedBalance.String()
			if capturedBalance.LessThan(common.Zero) {
				return errors.New("invalid balance")
			}
		}

		txErr = tx.Set(poolDocRef, pool.GetUpdateAllBalance(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(poolBalanceHistoryDocRef, poolHistory.GetAdd())
		if txErr != nil {
			return txErr
		}

		for itemIndex, itemDocRef := range itemDocRefs {
			txErr = tx.Set(creditDocRefs[itemIndex], map[string]string{
				"revenue": items[itemIndex].CreditRevenue,
			}, firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
			txErr = tx.Set(itemDocRef, items[itemIndex].GetUpdateBalance(), firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
			txErr = tx.Set(itemHistoryDocRefs[itemIndex], itemHistories[itemIndex].GetAdd())
			if txErr != nil {
				return txErr
			}
			txErr = tx.Set(transUserDocRefs[itemIndex], transList[itemIndex].GetUpdate(), firestore.MergeAll)
			if txErr != nil {
				return txErr
			}
		}
		txErr = tx.Set(transDocRef, trans.GetUpdate(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		for orderIndex, orderDocRef := range poolOrderDocRefs {
			if common.StringToDecimal(poolOrders[orderIndex].Balance).Equal(common.Zero) {
				txErr = tx.Delete(orderDocRef)
				if txErr != nil {
					return txErr
				}
			} else {
				txErr = tx.Set(orderDocRef, poolOrders[orderIndex].GetUpdateAllBalance(), firestore.MergeAll)
				if txErr != nil {
					return txErr
				}
				txErr = tx.Set(poolOrderUserDocRefs[orderIndex], poolOrders[orderIndex].GetUpdateAllBalance(), firestore.MergeAll)
				if txErr != nil {
					return txErr
				}
			}

		}

		return txErr
	})

	if err != nil {
		dao.SetCreditPoolCache(*pool)
	}

	return err
}

func (dao CreditDao) FinishFailedDepositCreditItem(item *bean.CreditItem, deposit *bean.CreditDeposit,
	tracking *bean.CreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	itemDocRef := dbClient.Doc(GetCreditItemItemPath(deposit.UID, deposit.Currency))
	depositUserDocRef := dbClient.Doc(GetCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))
	depositDocRef := dbClient.Doc(GetCreditDepositItemPath(deposit.Currency, deposit.Id))

	docLogRef := dbClient.Doc(GetCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	batch.Set(itemDocRef, item.GetUpdate(), firestore.MergeAll)
	batch.Set(depositUserDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Set(depositDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Delete(docTrackingRef)
	batch.Set(docLogRef, tracking.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditDao) ListCreditPool(currency string) (t TransferObject) {
	ListObjects(GetCreditPoolPath(currency), &t, nil, snapshotToCreditPool)
	return
}

func (dao CreditDao) GetCreditPool(currency string, percentage int) (t TransferObject) {
	level := fmt.Sprintf("%03d", percentage)
	GetObject(GetCreditPoolItemPath(currency, level), &t, snapshotToCreditPool)
	return
}

func (dao CreditDao) AddCreditPool(pool *bean.CreditPool) error {
	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCreditPoolItemPath(pool.Currency, pool.Level))
	_, err := poolDocRef.Set(context.Background(), pool.GetAdd())

	return err
}

func (dao CreditDao) ListCreditPoolOrder(currency string, level string) (t TransferObject) {
	ListObjects(GetCreditPoolItemOrderPath(currency, level), &t, nil, snapshotToCreditPoolOrder)
	return
}

func (dao CreditDao) ListCreditOnChainActionTracking(currency string) (t TransferObject) {
	ListObjects(GetCreditOnChainActionTrackingPath(currency), &t, nil, snapshotToCreditOnChainTracking)
	return
}

func (dao CreditDao) GetCreditOnChainActionTracking(currency string) (t TransferObject) {
	GetObject(GetCreditOnChainActionTrackingPath(currency), &t, snapshotToCreditOnChainTracking)
	return
}

func (dao CreditDao) AddCreditOnChainActionTracking(item *bean.CreditItem, deposit *bean.CreditDeposit,
	tracking *bean.CreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCreditOnChainActionLogPath(tracking.Currency)).NewDoc()
	tracking.Id = docRef.ID
	docTrackingRef := dbClient.Doc(GetCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	itemDocRef := dbClient.Doc(GetCreditItemItemPath(deposit.UID, deposit.Currency))
	depositDocRef := dbClient.Doc(GetCreditDepositItemPath(deposit.Currency, deposit.Id))
	depositUserDocRef := dbClient.Doc(GetCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, tracking.GetAdd())
	batch.Set(docTrackingRef, tracking.GetAdd())
	batch.Set(itemDocRef, item.GetUpdateStatus(), firestore.MergeAll)
	batch.Set(depositDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Set(depositUserDocRef, deposit.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditDao) UpdateCreditOnChainActionTracking(tracking *bean.CreditOnChainActionTracking) (err error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	batch := dbClient.Batch()
	batch.Delete(docTrackingRef)
	batch.Set(docRef, tracking.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditDao) GetCreditWithdraw(withdrawId string) (t TransferObject) {
	GetObject(GetCreditWithdrawItemPath(withdrawId), &t, snapshotToCreditWithdraw)
	return
}

func (dao CreditDao) AddCreditWithdraw(credit *bean.Credit, creditWithdraw *bean.CreditWithdraw) (err error) {
	dbClient := firebase_service.FirestoreClient

	creditDocRef := dbClient.Doc(GetCreditUserPath(credit.UID))
	creditWithdrawDocRef := dbClient.Collection(GetCreditWithdrawPath()).NewDoc()
	creditWithdraw.Id = creditWithdrawDocRef.ID
	creditWithdrawUserDocRef := dbClient.Doc(GetCreditWithdrawItemUserPath(credit.UID, creditWithdraw.Id))

	amount := common.StringToDecimal(creditWithdraw.Amount)
	err = dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		var txErr error

		creditDoc, txErr := tx.Get(creditDocRef)
		if err != nil {
			return txErr
		}
		creditRevenue, txErr := common.ConvertToDecimal(creditDoc, "revenue")
		if txErr != nil {
			return txErr
		}
		creditRevenue = creditRevenue.Sub(amount)
		if creditRevenue.LessThan(common.Zero) {
			return errors.New("invalid amount")
		}
		credit.Revenue = creditRevenue.String()

		txErr = tx.Set(creditDocRef, credit.GetUpdateRevenue(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(creditWithdrawDocRef, creditWithdraw.GetAdd(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}
		txErr = tx.Set(creditWithdrawUserDocRef, creditWithdraw.GetAdd(), firestore.MergeAll)
		if txErr != nil {
			return txErr
		}

		return txErr
	})

	return err
}

func (dao CreditDao) ListCreditPoolOrderUser(currency string, userId string) (t TransferObject) {
	ListObjects(GetCreditPoolItemOrderUserPath(currency, userId), &t, nil, snapshotToCreditPoolOrder)
	return
}

func (dao CreditDao) SetCreditPoolCache(pool bean.CreditPool) {
	b, _ := json.Marshal(&pool)
	key := GetCreditPoolCacheKey(pool.Currency, pool.Level)
	cache.RedisClient.Set(key, string(b), 0)
}

func (dao CreditDao) GetCreditPoolCache(currency string, level string) TransferObject {
	key := GetCreditPoolCacheKey(currency, level)
	var to TransferObject
	GetCacheObject(key, &to, func(val string) interface{} {
		var creditPool bean.CreditPool
		json.Unmarshal([]byte(val), &creditPool)
		return creditPool
	})

	return to
}

func (dao CreditDao) GetCreditPoolOrderByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCreditPoolOrder)
	return
}

func (dao CreditDao) UpdateNotificationCreditItem(creditItem bean.CreditItem) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationCreditItemPath(creditItem.UID, creditItem.Currency))
	err := ref.Set(context.Background(), creditItem.GetNotificationUpdate())

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

func GetCreditBalanceHistoryPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s/history", userId, currency)
}

func GetCreditBalanceHistoryItemPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/history/%s", userId, currency, id)
}

func GetCreditDepositUserPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s/deposits", userId, currency)
}

func GetCreditDepositItemUserPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/deposits/%s", userId, currency, id)
}

func GetCreditDepositPath(currency string) string {
	return fmt.Sprintf("credit_deposits/%s/deposits", currency)
}

func GetCreditDepositItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_deposits/%s/deposits/%s", currency, id)
}

func GetCreditWithdrawUserPath(userId string) string {
	return fmt.Sprintf("credits/%s/withdraws", userId)
}

func GetCreditWithdrawItemUserPath(userId string, id string) string {
	return fmt.Sprintf("credits/%s/withdraws/%s", userId, id)
}

func GetCreditWithdrawPath() string {
	return fmt.Sprintf("credit_withdraws")
}

func GetCreditWithdrawItemPath(id string) string {
	return fmt.Sprintf("credit_withdraws/%s", id)
}

func GetCreditTransactionUserPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s/transactions", userId, currency)
}

func GetCreditTransactionItemUserPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/transactions/%s", userId, currency, id)
}

func GetCreditTransactionPath(currency string) string {
	return fmt.Sprintf("credit_transactions/%s/transactions", currency)
}

func GetCreditPendingTransactionPath(currency string) string {
	return fmt.Sprintf("credit_transactions/%s/pending_transactions", currency)
}

func GetCreditPostTransactionPath(currency string) string {
	return fmt.Sprintf("credit_transactions/%s/post_transactions", currency)
}

func GetCreditTransactionItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_transactions/%s/transactions/%s", currency, id)
}

func GetCreditPendingTransactionItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_transactions/%s/pending_transactions/%s", currency, id)
}

func GetCreditPostTransactionItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_transactions/%s/post_transactions/%s", currency, id)
}

func GetCreditPoolPath(currency string) string {
	return fmt.Sprintf("credit_pools/%s/items", currency)
}

func GetCreditPoolItemPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s", currency, level)
}

func GetCreditPoolItemOrderPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/orders", currency, level)
}

func GetCreditPoolItemOrderItemPath(currency string, level string, order string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/orders/%s", currency, level, order)
}

func GetCreditPoolItemOrderUserPath(currency string, userId string) string {
	return fmt.Sprintf("credit_pool_orders/%s/items/%s/orders", currency, userId)
}

func GetCreditPoolItemOrderItemUserPath(currency string, userId string, order string) string {
	return fmt.Sprintf("credit_pool_orders/%s/items/%s/orders/%s", currency, userId, order)
}

func GetCreditPoolBalanceHistoryPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/history", currency, level)
}

func GetCreditPoolBalanceHistoryItemPath(currency string, level string, id string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/history/%s", currency, level, id)
}

func GetCreditOnChainActionTrackingPath(currency string) string {
	return fmt.Sprintf("credit_on_chain_trackings/%s/items", currency)
}

func GetCreditOnChainActionTrackingItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_on_chain_trackings/%s/items/%s", currency, id)
}

func GetCreditOnChainActionLogPath(currency string) string {
	return fmt.Sprintf("credit_on_chain_logs/%s/items", currency)
}

func GetCreditOnChainActionLogItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_on_chain_logs/%s/items/%s", currency, id)
}

func GetCreditPoolCacheKey(currency string, level string) string {
	return fmt.Sprintf("credit_pools.%s.%s", currency, level)
}

func GetNotificationCreditItemPath(userId string, currency string) string {
	return fmt.Sprintf("users/%s/credits/credit_item_%s", userId, currency)
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

func snapshotToCreditTransaction(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditTransaction
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCreditPoolOrder(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditPoolOrder
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCreditPoolBalanceHistory(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditPoolBalanceHistory
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCreditOnChainTracking(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CreditOnChainActionTracking
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}
