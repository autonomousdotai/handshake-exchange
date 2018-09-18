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

type CashDao struct {
}

func (dao CashDao) GetCashCredit(userId string) (t TransferObject) {
	GetObject(GetCashCreditUserPath(userId), &t, snapshotToCashCredit)
	return
}

func (dao CashDao) GetCashCreditItem(userId string, currency string) (t TransferObject) {
	GetObject(GetCashCreditItemItemPath(userId, currency), &t, snapshotToCashCreditItem)
	return
}

func (dao CashDao) ListCashCreditItem(userId string) (t TransferObject) {
	ListObjects(GetCashCreditItemPath(userId), &t, nil, snapshotToCashCreditItem)
	return
}

func (dao CashDao) AddCashCredit(credit *bean.CashCredit) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCashCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetAdd())

	return err
}

func (dao CashDao) UpdateCashCredit(credit *bean.CashCredit) error {
	dbClient := firebase_service.FirestoreClient

	creditPath := GetCashCreditUserPath(credit.UID)
	docRef := dbClient.Doc(creditPath)
	_, err := docRef.Set(context.Background(), credit.GetUpdate(), firestore.MergeAll)

	return err
}

func (dao CashDao) AddCashCreditItem(item *bean.CashCreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetAdd())

	return err
}

func (dao CashDao) UpdateCashCreditItem(item *bean.CashCreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetUpdateStatus(), firestore.MergeAll)

	return err
}

func (dao CashDao) UpdateCashCreditItemReactivate(item *bean.CashCreditItem) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashCreditItemItemPath(item.UID, item.Currency))
	_, err := docRef.Set(context.Background(), item.GetUpdateReactivate(), firestore.MergeAll)

	return err
}

func (dao CashDao) GetCashCreditDeposit(currency string, depositId string) (t TransferObject) {
	t = dao.GetCashCreditDepositByPath(GetCashCreditDepositItemPath(currency, depositId))
	return
}

func (dao CashDao) GetCashCreditDepositByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCashCreditDeposit)
	return
}

func (dao CashDao) AddCashCreditDeposit(item *bean.CashCreditItem, deposit *bean.CashCreditDeposit) (err error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCashCreditDepositPath(deposit.Currency)).NewDoc()
	deposit.Id = docRef.ID
	docUserRef := dbClient.Doc(GetCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, deposit.GetAdd())
	batch.Set(docUserRef, deposit.GetAdd())
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) FinishDepositCashCreditItem(item *bean.CashCreditItem, deposit *bean.CashCreditDeposit,
	itemHistory *bean.CashCreditBalanceHistory,
	pool *bean.CashCreditPool, poolOrder *bean.CashCreditPoolOrder, poolHistory *bean.CashCreditPoolBalanceHistory,
	tracking *bean.CashCreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient
	itemDocRef := dbClient.Doc(GetCashCreditItemItemPath(deposit.UID, deposit.Currency))
	depositUserDocRef := dbClient.Doc(GetCashCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))
	depositDocRef := dbClient.Doc(GetCashCreditDepositItemPath(deposit.Currency, deposit.Id))

	poolDocRef := dbClient.Doc(GetCashCreditPoolItemPath(deposit.Currency, pool.Level))
	poolOrderDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemPath(deposit.Currency, pool.Level, poolOrder.Id))
	poolOrderUserDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemUserPath(deposit.Currency, poolOrder.UID, poolOrder.Id))

	balanceHistoryDocRef := dbClient.Collection(GetCashCreditBalanceHistoryPath(deposit.UID, deposit.Currency)).NewDoc()
	itemHistory.Id = balanceHistoryDocRef.ID
	poolBalanceHistoryDocRef := dbClient.Collection(GetCashCreditPoolBalanceHistoryPath(deposit.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	docLogRef := dbClient.Doc(GetCashCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCashCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

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

	if err == nil {
		dao.SetCashCreditPoolCache(*pool)
	}

	return err
}

func (dao CashDao) RemoveCashCreditItem(item *bean.CashCreditItem, itemHistory *bean.CashCreditBalanceHistory,
	pool *bean.CashCreditPool, poolOrders []bean.CashCreditPoolOrder, poolHistory *bean.CashCreditPoolBalanceHistory) (err error) {

	dbClient := firebase_service.FirestoreClient
	itemDocRef := dbClient.Doc(GetCashCreditItemItemPath(item.UID, item.Currency))

	poolDocRef := dbClient.Doc(GetCashCreditPoolItemPath(item.Currency, pool.Level))

	balanceHistoryDocRef := dbClient.Collection(GetCashCreditBalanceHistoryPath(item.UID, item.Currency)).NewDoc()
	itemHistory.Id = balanceHistoryDocRef.ID
	poolBalanceHistoryDocRef := dbClient.Collection(GetCashCreditPoolBalanceHistoryPath(item.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)
	poolOrderUserDocRefs := make([]*firestore.DocumentRef, 0)

	for _, order := range poolOrders {
		poolOrderDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemPath(item.Currency, pool.Level, order.Id))
		poolOrderUserDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemUserPath(item.Currency, order.UID, order.Id))

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

		item.ReactivateAmount = item.Balance
		item.Balance = zeroStr
		itemHistory.New = item.Balance

		poolBalance = poolBalance.Sub(itemBalance)
		pool.Balance = poolBalance.String()
		poolHistory.New = pool.Balance

		// Update balance
		txErr = tx.Set(itemDocRef, item.GetUpdateDeactivate(), firestore.MergeAll)
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

	if err == nil {
		dao.SetCashCreditPoolCache(*pool)
	}

	return err
}

func (dao CashDao) RemoveCashCreditOnChainActionTracking(tracking bean.CashCreditOnChainActionTracking) error {
	dbClient := firebase_service.FirestoreClient
	docTrackingRef := dbClient.Doc(GetCashCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))
	_, err := docTrackingRef.Delete(context.Background())
	return err
}

func (dao CashDao) ListPendingCashCreditTransaction(currency string) (t TransferObject) {
	ListObjects(GetCashCreditPendingTransactionPath(currency), &t, nil, snapshotToCashCreditTransaction)
	return
}

func (dao CashDao) GetCashCreditTransaction(currency string, id string) (t TransferObject) {
	GetObject(GetCashCreditTransactionItemPath(currency, id), &t, snapshotToCreditTransaction)
	return
}

func (dao CashDao) GetCashCreditTransactionUser(userId string, currency string, id string) (t TransferObject) {
	GetObject(GetCashCreditTransactionItemUserPath(userId, currency, id), &t, snapshotToCashCreditTransaction)
	return
}

func (dao CashDao) AddCashCreditTransaction(pool *bean.CashCreditPool, trans *bean.CashCreditTransaction,
	userTransList []*bean.CashCreditTransaction, selectedOrders []bean.CashCreditPoolOrder) (err error) {

	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCashCreditPoolItemPath(pool.Currency, pool.Level))
	transDocRef := dbClient.Collection(GetCashCreditTransactionPath(pool.Currency)).NewDoc()
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
		orderPath := GetCashCreditPoolItemOrderItemPath(pool.Currency, pool.Level, creditOrder.Id)
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

func (dao CashDao) FinishCashCreditTransaction(pool *bean.CashCreditPool, poolHistory bean.CashCreditPoolBalanceHistory,
	items []bean.CashCreditItem, itemHistories []bean.CashCreditBalanceHistory, poolOrders []bean.CashCreditPoolOrder,
	trans *bean.CashCreditTransaction, transList []*bean.CashCreditTransaction) (err error) {

	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCashCreditPoolItemPath(pool.Currency, pool.Level))
	poolBalanceHistoryDocRef := dbClient.Collection(GetCashCreditPoolBalanceHistoryPath(pool.Currency, pool.Level)).NewDoc()
	poolHistory.Id = poolBalanceHistoryDocRef.ID

	creditDocRefs := make([]*firestore.DocumentRef, 0)
	itemDocRefs := make([]*firestore.DocumentRef, 0)
	itemHistoryDocRefs := make([]*firestore.DocumentRef, 0)

	transDocRef := dbClient.Doc(GetCashCreditTransactionItemPath(trans.Currency, trans.Id))
	transUserDocRefs := make([]*firestore.DocumentRef, 0)
	for itemIndex, item := range items {
		creditDocRef := dbClient.Doc(GetCashCreditUserPath(item.UID))
		itemDocRef := dbClient.Doc(GetCashCreditItemItemPath(item.UID, item.Currency))
		balanceHistoryDocRef := dbClient.Collection(GetCashCreditBalanceHistoryPath(item.UID, item.Currency)).NewDoc()

		itemHistories[itemIndex].Id = balanceHistoryDocRef.ID

		creditDocRefs = append(creditDocRefs, creditDocRef)
		itemDocRefs = append(itemDocRefs, itemDocRef)
		itemHistoryDocRefs = append(itemHistoryDocRefs, balanceHistoryDocRef)

		transUserDocRef := dbClient.Doc(GetCashCreditTransactionItemUserPath(transList[itemIndex].UID, transList[itemIndex].Currency, transList[itemIndex].Id))
		transUserDocRefs = append(transUserDocRefs, transUserDocRef)
	}

	poolOrderDocRefs := make([]*firestore.DocumentRef, 0)
	poolOrderUserDocRefs := make([]*firestore.DocumentRef, 0)

	for _, poolOrder := range poolOrders {
		poolOrderDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemPath(pool.Currency, pool.Level, poolOrder.Id))
		poolOrderUserDocRef := dbClient.Doc(GetCashCreditPoolItemOrderItemUserPath(pool.Currency, poolOrder.UID, poolOrder.Id))

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
			creditDoc, txErr := tx.Get(creditDocRefs[itemIndex])
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
			creditRevenue, txErr := common.ConvertToDecimal(creditDoc, "revenue")
			if txErr != nil {
				return txErr
			}
			itemRevenue, txErr := common.ConvertToDecimal(itemDoc, "revenue")
			if txErr != nil {
				return txErr
			}
			// Revert
			itemAmount := common.StringToDecimal(itemHistories[itemIndex].Change).Neg()
			revenue := common.StringToDecimal(transList[itemIndex].Revenue)
			itemHistories[itemIndex].Old = itemBalance.String()

			itemBalance = itemBalance.Sub(itemAmount)
			items[itemIndex].Balance = itemBalance.String()
			sold = sold.Add(itemAmount)
			items[itemIndex].Sold = sold.String()
			creditRevenue = creditRevenue.Add(revenue)
			itemRevenue = itemRevenue.Add(revenue)
			items[itemIndex].CreditRevenue = creditRevenue.String()
			items[itemIndex].Revenue = itemRevenue.String()

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

	if err == nil {
		dao.SetCashCreditPoolCache(*pool)
	}

	return err
}

func (dao CashDao) FinishFailedDepositCashCreditItem(item *bean.CashCreditItem, deposit *bean.CashCreditDeposit,
	tracking *bean.CashCreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	itemDocRef := dbClient.Doc(GetCashCreditItemItemPath(deposit.UID, deposit.Currency))
	depositUserDocRef := dbClient.Doc(GetCashCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))
	depositDocRef := dbClient.Doc(GetCashCreditDepositItemPath(deposit.Currency, deposit.Id))

	docLogRef := dbClient.Doc(GetCashCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCashCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	batch.Set(itemDocRef, item.GetUpdate(), firestore.MergeAll)
	batch.Set(depositUserDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Set(depositDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Delete(docTrackingRef)
	batch.Set(docLogRef, tracking.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) ListCashCreditPool(currency string) (t TransferObject) {
	ListObjects(GetCashCreditPoolPath(currency), &t, nil, snapshotToCashCreditPool)
	return
}

func (dao CashDao) GetCashCreditPool(currency string, percentage int) (t TransferObject) {
	level := fmt.Sprintf("%03d", percentage)
	GetObject(GetCashCreditPoolItemPath(currency, level), &t, snapshotToCashCreditPool)
	return
}

func (dao CashDao) AddCashCreditPool(pool *bean.CashCreditPool) error {
	dbClient := firebase_service.FirestoreClient

	poolDocRef := dbClient.Doc(GetCashCreditPoolItemPath(pool.Currency, pool.Level))
	_, err := poolDocRef.Set(context.Background(), pool.GetAdd())

	return err
}

func (dao CashDao) ListCashCreditPoolOrder(currency string, level string) (t TransferObject) {
	ListObjects(GetCashCreditPoolItemOrderPath(currency, level), &t, nil, snapshotToCashCreditPoolOrder)
	return
}

func (dao CashDao) ListCashCreditOnChainActionTracking(currency string) (t TransferObject) {
	ListObjects(GetCashCreditOnChainActionTrackingPath(currency), &t, nil, snapshotToCashCreditOnChainTracking)
	return
}

func (dao CashDao) GetCashCreditOnChainActionTracking(currency string) (t TransferObject) {
	GetObject(GetCashCreditOnChainActionTrackingPath(currency), &t, snapshotToCashCreditOnChainTracking)
	return
}

func (dao CashDao) AddCashCreditOnChainActionTracking(item *bean.CashCreditItem, deposit *bean.CashCreditDeposit,
	tracking *bean.CashCreditOnChainActionTracking) (err error) {

	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetCashCreditOnChainActionLogPath(tracking.Currency)).NewDoc()
	tracking.Id = docRef.ID
	docTrackingRef := dbClient.Doc(GetCashCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	itemDocRef := dbClient.Doc(GetCashCreditItemItemPath(deposit.UID, deposit.Currency))
	depositDocRef := dbClient.Doc(GetCashCreditDepositItemPath(deposit.Currency, deposit.Id))
	depositUserDocRef := dbClient.Doc(GetCashCreditDepositItemUserPath(deposit.UID, deposit.Currency, deposit.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, tracking.GetAdd())
	batch.Set(docTrackingRef, tracking.GetAdd())
	batch.Set(itemDocRef, item.GetUpdateStatus(), firestore.MergeAll)
	batch.Set(depositDocRef, deposit.GetUpdate(), firestore.MergeAll)
	batch.Set(depositUserDocRef, deposit.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) UpdateCashCreditOnChainActionTracking(tracking *bean.CashCreditOnChainActionTracking) (err error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashCreditOnChainActionLogItemPath(tracking.Currency, tracking.Id))
	docTrackingRef := dbClient.Doc(GetCashCreditOnChainActionTrackingItemPath(tracking.Currency, tracking.Id))

	batch := dbClient.Batch()
	batch.Delete(docTrackingRef)
	batch.Set(docRef, tracking.GetUpdate(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) GetCashCreditWithdraw(withdrawId string) (t TransferObject) {
	GetObject(GetCashCreditWithdrawItemPath(withdrawId), &t, snapshotToCashCreditWithdraw)
	return
}

func (dao CashDao) AddCashCreditWithdraw(credit *bean.CashCredit, creditWithdraw *bean.CashCreditWithdraw) (err error) {
	dbClient := firebase_service.FirestoreClient

	creditDocRef := dbClient.Doc(GetCashCreditUserPath(credit.UID))
	creditWithdrawDocRef := dbClient.Collection(GetCashCreditWithdrawPath()).NewDoc()
	creditWithdraw.Id = creditWithdrawDocRef.ID
	creditWithdrawUserDocRef := dbClient.Doc(GetCashCreditWithdrawItemUserPath(credit.UID, creditWithdraw.Id))

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

func (dao CashDao) ListCashCreditPoolOrderUser(currency string, userId string) (t TransferObject) {
	ListObjects(GetCashCreditPoolItemOrderUserPath(currency, userId), &t, nil, snapshotToCashCreditPoolOrder)
	return
}

func (dao CashDao) ListCashCreditWithdraw() (t TransferObject) {
	ListObjects(GetCashCreditWithdrawPath(), &t, nil, snapshotToCreditWithdraw)
	return
}

func (dao CashDao) UpdateProcessingCashWithdraw(withdraw bean.CashCreditWithdraw) (err error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetCashCreditWithdrawItemPath(withdraw.Id))
	docUserRef := dbClient.Doc(GetCashCreditWithdrawItemUserPath(withdraw.UID, withdraw.Id))
	docProcessedRef := dbClient.Doc(GetCashCreditProcessedWithdrawItemPath(withdraw.Id))

	batch := dbClient.Batch()
	batch.Delete(docRef)
	batch.Set(docProcessedRef, withdraw.GetAdd())
	batch.Set(docUserRef, withdraw.GetUpdateStatus(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) UpdateProcessedCashWithdraw(withdraw bean.CashCreditWithdraw) (err error) {
	dbClient := firebase_service.FirestoreClient

	docUserRef := dbClient.Doc(GetCashCreditWithdrawItemUserPath(withdraw.UID, withdraw.Id))
	docProcessedRef := dbClient.Doc(GetCashCreditProcessedWithdrawItemPath(withdraw.Id))

	batch := dbClient.Batch()
	batch.Set(docProcessedRef, withdraw.GetUpdateStatus(), firestore.MergeAll)
	batch.Set(docUserRef, withdraw.GetUpdateStatus(), firestore.MergeAll)
	_, err = batch.Commit(context.Background())

	return err
}

func (dao CashDao) SetCashCreditPoolCache(pool bean.CashCreditPool) {
	b, _ := json.Marshal(&pool)
	key := GetCashCreditPoolCacheKey(pool.Currency, pool.Level)
	cache.RedisClient.Set(key, string(b), 0)
}

func (dao CashDao) GetCashCreditPoolCache(currency string, level string) TransferObject {
	key := GetCashCreditPoolCacheKey(currency, level)
	var to TransferObject
	GetCacheObject(key, &to, func(val string) interface{} {
		var creditPool bean.CashCreditPool
		json.Unmarshal([]byte(val), &creditPool)
		return creditPool
	})

	return to
}

func (dao CashDao) GetCashCreditPoolOrderByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToCashCreditPoolOrder)
	return
}

func (dao CashDao) UpdateNotificationCashCreditItem(creditItem bean.CashCreditItem) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationCashCreditItemPath(creditItem.UID, creditItem.Currency))
	err := ref.Set(context.Background(), creditItem.GetNotificationUpdate())

	return err
}

func GetCashCreditUserPath(userId string) string {
	return fmt.Sprintf("credits/%s", userId)
}

func GetCashCreditItemPath(userId string) string {
	return fmt.Sprintf("credits/%s/items", userId)
}

func GetCashCreditItemItemPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s", userId, currency)
}

func GetCashCreditBalanceHistoryPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s/history", userId, currency)
}

func GetCashCreditBalanceHistoryItemPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/history/%s", userId, currency, id)
}

func GetCashCreditDepositUserPath(userId string, currency string) string {
	return fmt.Sprintf("credits/%s/items/%s/deposits", userId, currency)
}

func GetCashCreditDepositItemUserPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/deposits/%s", userId, currency, id)
}

func GetCashCreditDepositPath(currency string) string {
	return fmt.Sprintf("credit_deposits/%s/deposits", currency)
}

func GetCashCreditDepositItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_deposits/%s/deposits/%s", currency, id)
}

func GetCashCreditWithdrawItemUserPath(userId string, id string) string {
	return fmt.Sprintf("credits/%s/withdraws/%s", userId, id)
}

func GetCashCreditWithdrawPath() string {
	return fmt.Sprintf("credit_withdraws")
}

func GetCashCreditWithdrawItemPath(id string) string {
	return fmt.Sprintf("credit_withdraws/%s", id)
}

func GetCashCreditProcessedWithdrawItemPath(id string) string {
	return fmt.Sprintf("credit_processed_withdraws/%s", id)
}

func GetCashCreditTransactionItemUserPath(userId string, currency string, id string) string {
	return fmt.Sprintf("credits/%s/items/%s/transactions/%s", userId, currency, id)
}

func GetCashCreditTransactionPath(currency string) string {
	return fmt.Sprintf("credit_transactions/%s/transactions", currency)
}

func GetCashCreditPendingTransactionPath(currency string) string {
	return fmt.Sprintf("credit_transactions/%s/pending_transactions", currency)
}

func GetCashCreditTransactionItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_transactions/%s/transactions/%s", currency, id)
}

func GetCashCreditPoolPath(currency string) string {
	return fmt.Sprintf("credit_pools/%s/items", currency)
}

func GetCashCreditPoolItemPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s", currency, level)
}

func GetCashCreditPoolItemOrderPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/orders", currency, level)
}

func GetCashCreditPoolItemOrderItemPath(currency string, level string, order string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/orders/%s", currency, level, order)
}

func GetCashCreditPoolItemOrderUserPath(currency string, userId string) string {
	return fmt.Sprintf("credit_pool_orders/%s/items/%s/orders", currency, userId)
}

func GetCashCreditPoolItemOrderItemUserPath(currency string, userId string, order string) string {
	return fmt.Sprintf("credit_pool_orders/%s/items/%s/orders/%s", currency, userId, order)
}

func GetCashCreditPoolBalanceHistoryPath(currency string, level string) string {
	return fmt.Sprintf("credit_pools/%s/items/%s/history", currency, level)
}

func GetCashCreditOnChainActionTrackingPath(currency string) string {
	return fmt.Sprintf("credit_on_chain_trackings/%s/items", currency)
}

func GetCashCreditOnChainActionTrackingItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_on_chain_trackings/%s/items/%s", currency, id)
}

func GetCashCreditOnChainActionLogPath(currency string) string {
	return fmt.Sprintf("credit_on_chain_logs/%s/items", currency)
}

func GetCashCreditOnChainActionLogItemPath(currency string, id string) string {
	return fmt.Sprintf("credit_on_chain_logs/%s/items/%s", currency, id)
}

func GetCashCreditPoolCacheKey(currency string, level string) string {
	return fmt.Sprintf("credit_pools.%s.%s", currency, level)
}

func GetNotificationCashCreditItemPath(userId string, currency string) string {
	return fmt.Sprintf("users/%s/credits/credit_item_%s", userId, currency)
}

func snapshotToCashCredit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCredit
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCashCreditItem(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditItem
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToCashCreditDeposit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditDeposit
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCashCreditWithdraw(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditWithdraw
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCashCreditPool(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditPool
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCashCreditTransaction(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditTransaction
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCashCreditPoolOrder(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditPoolOrder
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToCashCreditPoolBalanceHistory(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditPoolBalanceHistory
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToCashCreditOnChainTracking(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CashCreditOnChainActionTracking
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}
