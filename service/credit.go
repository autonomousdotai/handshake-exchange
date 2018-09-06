package service

import (
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/chainso_service"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/crypto_service"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type CreditService struct {
	dao     *dao.CreditDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CreditService) GetCredit(userId string) (credit bean.Credit, ce SimpleContextError) {
	creditTO := s.dao.GetCredit(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
		return
	}

	if creditTO.Found {
		credit = creditTO.Object.(bean.Credit)
		creditItemsTO := s.dao.ListCreditItem(userId)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
			return
		}
		creditItemMap := map[string]bean.CreditItem{}
		for _, creditItem := range creditItemsTO.Objects {
			item := creditItem.(bean.CreditItem)
			creditItemMap[item.Currency] = item
		}
		credit.Items = creditItemMap
	} else {
		ce.NotFound = true
	}

	return
}

func (s CreditService) AddCredit(userId string, body bean.Credit) (credit bean.Credit, ce SimpleContextError) {
	creditTO := s.dao.GetCredit(userId)

	if creditTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO)
		return
	}

	var err error
	if creditTO.Found {
		ce.SetStatusKey(api_error.CreditExists)
	} else {
		body.UID = userId
		credit = body
		credit.Status = bean.CREDIT_STATUS_ACTIVE
		err = s.dao.AddCredit(&credit)
		if err != nil {
			ce.SetError(api_error.AddDataFailed, err)
			return
		}
		ce.NotFound = false
	}

	return
}

func (s CreditService) AddDeposit(userId string, body bean.CreditDepositInput) (deposit bean.CreditDeposit, ce SimpleContextError) {
	var err error

	// Minimum amount
	amount, _ := decimal.NewFromString(body.Amount)
	if body.Currency == bean.ETH.Code {
		if amount.LessThan(bean.MIN_ETH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if body.Currency == bean.BTC.Code {
		if amount.LessThan(bean.MIN_BTC) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}
	if body.Currency == bean.BCH.Code {
		if amount.LessThan(bean.MIN_BCH) {
			ce.SetStatusKey(api_error.AmountIsTooSmall)
			return
		}
	}

	creditItemTO := s.dao.GetCreditItem(userId, body.Currency)
	var creditItem bean.CreditItem
	if creditItemTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, creditItemTO)
		return
	} else {

		if !creditItemTO.Found {
			creditItem = bean.CreditItem{
				UID:         userId,
				Currency:    body.Currency,
				Status:      bean.CREDIT_ITEM_STATUS_CREATE,
				Percentage:  body.Percentage,
				UserAddress: body.UserAddress,
			}
			err = s.dao.AddCreditItem(&creditItem)
			if err != nil {
				ce.SetError(api_error.AddDataFailed, err)
				return
			}
		} else {
			creditItem = creditItemTO.Object.(bean.CreditItem)
			if creditItem.Percentage != body.Percentage {
				ce.SetStatusKey(api_error.InvalidRequestBody)
				return
			}
			if creditItem.Status == bean.CREDIT_ITEM_STATUS_CREATE || creditItem.SubStatus == bean.CREDIT_ITEM_SUB_STATUS_TRANSFERRING {
				ce.SetStatusKey(api_error.OfferStatusInvalid)
				return
			}
			creditItem.UserAddress = body.UserAddress
		}
	}

	deposit = bean.CreditDeposit{
		UID:      userId,
		ItemRef:  dao.GetCreditItemItemPath(userId, body.Currency),
		Status:   bean.CREDIT_DEPOSIT_STATUS_CREATED,
		Currency: body.Currency,
		Amount:   body.Amount,
	}

	if body.Currency != "" {
		resp, errCoinbase := coinbase_service.GenerateAddress(body.Currency)
		if errCoinbase != nil {
			ce.SetError(api_error.ExternalApiFailed, errCoinbase)
		}
		deposit.SystemAddress = resp.Data.Address
	}

	err = s.dao.AddCreditDeposit(&creditItem, &deposit)
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}

	return
}

func (s CreditService) AddTracking(userId string, body bean.CreditOnChainActionTrackingInput) (tracking bean.CreditOnChainActionTracking, ce SimpleContextError) {
	depositTO := s.dao.GetCreditDeposit(body.Currency, body.Deposit)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, depositTO) {
		return
	}
	deposit := depositTO.Object.(bean.CreditDeposit)

	itemTO := s.dao.GetCreditItem(userId, body.Currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, itemTO) {
		return
	}
	item := itemTO.Object.(bean.CreditItem)

	tracking = bean.CreditOnChainActionTracking{
		UID:        userId,
		ItemRef:    deposit.ItemRef,
		DepositRef: dao.GetCreditDepositItemPath(body.Currency, deposit.Id),
		TxHash:     body.TxHash,
		Action:     body.Action,
		Reason:     body.Reason,
		Currency:   body.Currency,
	}

	item.SubStatus = bean.CREDIT_ITEM_SUB_STATUS_TRANSFERRING
	deposit.Status = bean.CREDIT_DEPOSIT_STATUS_TRANSFERRING

	s.dao.AddCreditOnChainActionTracking(&item, &deposit, &tracking)

	return
}

func (s CreditService) FinishTracking() (ce SimpleContextError) {
	trackingTO := s.dao.ListCreditOnChainActionTracking(bean.ETH.Code)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, trackingTO) {
		return
	}
	for _, item := range trackingTO.Objects {
		trackingItem := item.(bean.CreditOnChainActionTracking)

		if trackingItem.TxHash != "" {
			amount := decimal.Zero
			isSuccess, isPending, amount, errChain := crypto_service.GetTransactionReceipt(trackingItem.TxHash, trackingItem.Currency)
			if errChain == nil {
				if isSuccess && !isPending && amount.GreaterThan(common.Zero) {
					trackingItem.Amount = amount.String()
					s.finishTrackingItem(trackingItem)
				}
			} else {
				ce.SetError(api_error.ExternalApiFailed, errChain)
			}
		} else {
			s.dao.RemoveCreditOnChainActionTracking(trackingItem)
		}
	}

	trackingTO = s.dao.ListCreditOnChainActionTracking(bean.BTC.Code)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, trackingTO) {
		return
	}
	for _, item := range trackingTO.Objects {
		trackingItem := item.(bean.CreditOnChainActionTracking)

		if trackingItem.TxHash != "" {
			confirmation, errChain := chainso_service.GetConfirmations(trackingItem.TxHash)
			amount := decimal.Zero
			if errChain == nil {
				amount, errChain = chainso_service.GetAmount(trackingItem.TxHash)
			} else {
				ce.SetError(api_error.ExternalApiFailed, errChain)
			}

			fmt.Println(fmt.Sprintf("%s %s %s %s", trackingItem.Id, trackingItem.UID, trackingItem.TxHash, amount.String()))
			confirmationRequired := s.getConfirmationRange(amount)
			if errChain == nil {
				if confirmation >= confirmationRequired && amount.GreaterThan(common.Zero) {
					trackingItem.Amount = amount.String()
					s.finishTrackingItem(trackingItem)
				}
			} else {
				ce.SetError(api_error.ExternalApiFailed, errChain)
			}
		} else {
			s.dao.RemoveCreditOnChainActionTracking(trackingItem)
		}
	}

	trackingTO = s.dao.ListCreditOnChainActionTracking(bean.BCH.Code)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, trackingTO) {
		return
	}
	for _, item := range trackingTO.Objects {
		trackingItem := item.(bean.CreditOnChainActionTracking)

		if trackingItem.TxHash != "" {
			confirmation, errChain := chainso_service.GetConfirmations(trackingItem.TxHash)
			amount := decimal.Zero
			if errChain == nil {
				amount, errChain = chainso_service.GetAmount(trackingItem.TxHash)
			} else {
				ce.SetError(api_error.ExternalApiFailed, errChain)
			}

			fmt.Println(fmt.Sprintf("%s %s %s %s", trackingItem.Id, trackingItem.UID, trackingItem.TxHash, amount.String()))
			confirmationRequired := s.getConfirmationRange(amount)
			if errChain == nil {
				if confirmation >= confirmationRequired && amount.GreaterThan(common.Zero) {
					trackingItem.Amount = amount.String()
					s.finishTrackingItem(trackingItem)
				}
			} else {
				ce.SetError(api_error.ExternalApiFailed, errChain)
			}
		} else {
			s.dao.RemoveCreditOnChainActionTracking(trackingItem)
		}
	}

	return
}

func (s CreditService) GetCreditPoolPercentageByCache(currency string, amount decimal.Decimal) (int, error) {
	percentage := 0
	for percentage <= 100 {
		level := fmt.Sprintf("%03d", percentage)

		creditPoolTO := s.dao.GetCreditPoolCache(currency, level)
		if creditPoolTO.HasError() {
			return 0, creditPoolTO.Error
		}
		if creditPoolTO.Found {
			creditPool := creditPoolTO.Object.(bean.CreditPool)
			if common.StringToDecimal(creditPool.Balance).GreaterThanOrEqual(amount) {
				return percentage, nil
			}
		}

		percentage += 1
	}

	return 0, nil
}

func (s CreditService) AddCreditTransaction(trans *bean.CreditTransaction) (ce SimpleContextError) {
	poolTO := s.dao.GetCreditPool(trans.Currency, int(common.StringToDecimal(trans.Percentage).IntPart()))
	if ce.FeedDaoTransfer(api_error.GetDataFailed, poolTO) {
		return
	}
	pool := poolTO.Object.(bean.CreditPool)

	percentage := int(common.StringToDecimal(trans.Percentage).IntPart())
	level := fmt.Sprintf("%03d", percentage)
	orderTO := s.dao.ListCreditPoolOrder(trans.Currency, level)

	if ce.FeedDaoTransfer(api_error.GetDataFailed, orderTO) {
		return
	}
	amount := common.StringToDecimal(trans.Amount)
	selectedOrders := make([]bean.CreditPoolOrder, 0)

	userTransMap := map[string]*bean.CreditTransaction{}
	userTransList := make([]bean.CreditTransaction, 0)

	needBreak := false
	for _, item := range orderTO.Objects {
		order := item.(bean.CreditPoolOrder)
		orderBalance := common.StringToDecimal(order.Balance)
		orderAmountSub := orderBalance.Sub(common.StringToDecimal(order.CapturedBalance))

		if !order.CapturedFull {
			var capturedAmount decimal.Decimal
			sub := amount.Sub(orderAmountSub)

			if sub.LessThan(common.Zero) {
				capturedAmount = amount
				needBreak = true
			} else {
				capturedAmount = orderAmountSub
				order.CapturedFull = true

				// out of amount, stop
				if sub.Equal(common.Zero) {
					needBreak = true
				} else {
					amount = amount.Sub(orderAmountSub)
				}
			}
			order.CapturedAmount = capturedAmount
			selectedOrders = append(selectedOrders, order)

			if userTrans, ok := userTransMap[order.UID]; ok {
				transAmount := common.StringToDecimal(userTrans.Amount)
				transAmount = transAmount.Add(capturedAmount)
				userTrans.Amount = transAmount.String()

				userTrans.OrderInfoRefs = append(userTrans.OrderInfoRefs, bean.OrderInfoRef{
					OrderRef: dao.GetCreditPoolItemOrderItemPath(trans.Currency, level, order.Id),
					Amount:   capturedAmount.String(),
				})

				userTransMap[order.UID] = userTrans
			} else {
				userTrans = &bean.CreditTransaction{}
				userTrans.UID = order.UID
				userTrans.ToUID = trans.ToUID
				userTrans.Status = bean.CREDIT_TRANSACTION_STATUS_CREATE
				userTrans.Currency = trans.Currency
				userTrans.Percentage = trans.Percentage
				userTrans.OfferRef = trans.OfferRef
				userTrans.Amount = capturedAmount.String()

				userTrans.OrderInfoRefs = append(userTrans.OrderInfoRefs, bean.OrderInfoRef{
					OrderRef: dao.GetCreditPoolItemOrderItemPath(trans.Currency, level, order.Id),
					Amount:   capturedAmount.String(),
				})

				userTransMap[order.UID] = userTrans
			}
			trans.OrderInfoRefs = append(trans.OrderInfoRefs, bean.OrderInfoRef{
				OrderRef: dao.GetCreditPoolItemOrderItemPath(trans.Currency, level, order.Id),
				Amount:   capturedAmount.String(),
			})

			if needBreak {
				break
			}
		}
	}

	trans.Status = bean.CREDIT_TRANSACTION_STATUS_CREATE
	for k, v := range userTransMap {
		trans.UIDs = append(trans.UIDs, k)
		userTransList = append(userTransList, *v)
	}

	if len(selectedOrders) == 0 {
		ce.SetStatusKey(api_error.CreditPriceChanged)
	}

	err := s.dao.AddCreditTransaction(&pool, trans, userTransList, selectedOrders)
	if err != nil {
		if strings.Contains(err.Error(), "out of stock") {
			ce.SetStatusKey(api_error.CreditPriceChanged)
			return
		} else {
			ce.SetError(api_error.AddDataFailed, err)
		}
	}

	return
}

func (s CreditService) ListPendingCreditTransaction(currency string) (trans []bean.CreditTransaction, ce SimpleContextError) {
	transTO := s.dao.ListPendingCreditTransaction(currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, transTO) {
		return
	}
	for _, item := range transTO.Objects {
		transItem := item.(bean.CreditTransaction)
		trans = append(trans, transItem)
	}

	return
}

func (s CreditService) FinishCreditTransaction(currency string, id string, offerRef string, revenue decimal.Decimal) (ce SimpleContextError) {
	transTO := s.dao.GetCreditTransaction(currency, id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, transTO) {
		return
	}
	trans := transTO.Object.(bean.CreditTransaction)
	trans.OfferRef = offerRef
	trans.Status = bean.CREDIT_TRANSACTION_STATUS_SUCCESS
	trans.SubStatus = bean.CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED
	trans.Revenue = revenue.String()

	amount := common.StringToDecimal(trans.Amount)

	poolTO := s.dao.GetCreditPool(trans.Currency, int(common.StringToDecimal(trans.Percentage).IntPart()))
	if ce.FeedDaoTransfer(api_error.GetDataFailed, transTO) {
		return
	}

	pool := poolTO.Object.(bean.CreditPool)
	poolHistory := bean.CreditPoolBalanceHistory{
		ItemRef:    "",
		ModifyRef:  dao.GetCreditTransactionItemPath(currency, id),
		ModifyType: bean.CREDIT_POOL_MODIFY_TYPE_PURCHASE,
		Change:     trans.Amount,
	}

	items := make([]bean.CreditItem, 0)
	itemHistories := make([]bean.CreditBalanceHistory, 0)
	transList := make([]bean.CreditTransaction, 0)
	for _, userId := range trans.UIDs {
		itemTO := s.dao.GetCreditItem(userId, trans.Currency)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, itemTO) {
			return
		}
		item := itemTO.Object.(bean.CreditItem)
		items = append(items, item)

		userTransTO := s.dao.GetCreditTransactionUser(userId, trans.Currency, trans.Id)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, userTransTO) {
			return
		}
		userTrans := userTransTO.Object.(bean.CreditTransaction)
		userAmount := common.StringToDecimal(userTrans.Amount)

		userTrans.OfferRef = offerRef
		userTrans.Status = bean.CREDIT_TRANSACTION_STATUS_SUCCESS
		userTrans.SubStatus = bean.CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED
		userTrans.Revenue = userAmount.Div(amount).Mul(revenue).RoundBank(2).String()
		transList = append(transList, userTrans)

		itemHistory := bean.CreditBalanceHistory{
			ItemRef:    dao.GetCreditItemItemPath(userId, trans.Currency),
			ModifyRef:  dao.GetCreditTransactionItemUserPath(userId, currency, userTrans.Id),
			ModifyType: bean.CREDIT_POOL_MODIFY_TYPE_PURCHASE,
			Change:     userTrans.Amount,
		}
		itemHistories = append(itemHistories, itemHistory)
	}

	orders := make([]bean.CreditPoolOrder, 0)
	for _, orderInfo := range trans.OrderInfoRefs {
		orderTO := s.dao.GetCreditPoolOrderByPath(orderInfo.OrderRef)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, orderTO) {
			return
		}
		order := orderTO.Object.(bean.CreditPoolOrder)
		order.CapturedAmount = common.StringToDecimal(orderInfo.Amount)
		orders = append(orders, order)
	}

	err := s.dao.FinishCreditTransaction(&pool, poolHistory, items, itemHistories, orders, &trans, transList)
	if err != nil {
		ce.SetError(api_error.UpdateDataFailed, err)
	}

	return
}

func (s CreditService) AddCreditWithdraw(userId string, body bean.CreditWithdraw) (withdraw bean.CreditWithdraw, ce SimpleContextError) {
	creditTO := s.dao.GetCredit(userId)

	if creditTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO)
		return
	}
	credit := creditTO.Object.(bean.Credit)
	body.UID = userId

	revenue := common.StringToDecimal(credit.Revenue)
	withdrawAmount := common.StringToDecimal(body.Amount)
	if withdrawAmount.GreaterThan(revenue) {
		ce.SetStatusKey(api_error.InvalidAmount)
		return
	}

	err := s.dao.AddCreditWithdraw(&credit, &body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid amount") {
			ce.SetStatusKey(api_error.InvalidAmount)
			return
		}
		ce.SetError(api_error.UpdateDataFailed, err)
	}

	withdraw = body
	return
}

func (s CreditService) SetupCreditPool() (ce SimpleContextError) {
	for _, currency := range []string{bean.BTC.Code, bean.ETH.Code, bean.BCH.Code} {
		level := 0
		for level <= 100 {
			pool := bean.CreditPool{
				Level:    fmt.Sprintf("%03d", level),
				Balance:  common.Zero.String(),
				Currency: currency,
			}
			err := s.dao.AddCreditPool(&pool)
			s.dao.SetCreditPoolCache(pool)
			if err != nil {
				ce.SetError(api_error.AddDataFailed, err)
			}
			level += 1
		}
	}

	return
}

func (s CreditService) SetupCreditPoolCache() (ce SimpleContextError) {
	for _, currency := range []string{bean.BTC.Code, bean.ETH.Code, bean.BCH.Code} {
		poolTO := s.dao.ListCreditPool(currency)
		if !poolTO.HasError() {
			for _, item := range poolTO.Objects {
				creditPool := item.(bean.CreditPool)
				s.dao.SetCreditPoolCache(creditPool)
			}
		}
	}

	return
}

func (s CreditService) finishTrackingItem(tracking bean.CreditOnChainActionTracking) error {
	var err error

	depositTO := s.dao.GetCreditDepositByPath(tracking.DepositRef)
	if depositTO.HasError() {
		return depositTO.Error
	}
	deposit := depositTO.Object.(bean.CreditDeposit)

	itemTO := s.dao.GetCreditItem(tracking.UID, tracking.Currency)
	if itemTO.HasError() {
		return itemTO.Error
	}
	item := itemTO.Object.(bean.CreditItem)

	if item.Status == bean.CREDIT_ITEM_STATUS_CREATE || item.Status == bean.CREDIT_ITEM_STATUS_INACTIVE {
		item.Status = bean.CREDIT_ITEM_STATUS_ACTIVE
	}
	item.SubStatus = bean.CREDIT_ITEM_SUB_STATUS_TRANSFERRED
	item.LastActionData = deposit
	deposit.Status = bean.CREDIT_DEPOSIT_STATUS_TRANSFERRED

	poolTO := s.dao.GetCreditPool(item.Currency, int(common.StringToDecimal(item.Percentage).IntPart()))
	if poolTO.HasError() {
		return poolTO.Error
	}
	pool := poolTO.Object.(bean.CreditPool)
	itemHistory := bean.CreditBalanceHistory{
		ItemRef:    tracking.ItemRef,
		ModifyRef:  tracking.DepositRef,
		ModifyType: tracking.Action,
	}
	poolHistory := bean.CreditPoolBalanceHistory{
		ItemRef:    tracking.ItemRef,
		ModifyRef:  tracking.DepositRef,
		ModifyType: tracking.Action,
	}
	poolOrder := bean.CreditPoolOrder{
		Id:         time.Now().UTC().Format("2006-01-02T15:04:05.000000000"),
		UID:        tracking.UID,
		DepositRef: tracking.DepositRef,
		Amount:     tracking.Amount,
		Balance:    tracking.Amount,
	}

	s.dao.FinishDepositCreditItem(&item, &deposit, &itemHistory, &pool, &poolOrder, &poolHistory, &tracking)

	return err
}

func (s CreditService) finishFailedTrackingItem(tracking bean.CreditOnChainActionTracking) error {
	var err error

	depositTO := s.dao.GetCreditDepositByPath(tracking.DepositRef)
	if depositTO.HasError() {
		return depositTO.Error
	}
	deposit := depositTO.Object.(bean.CreditDeposit)

	itemTO := s.dao.GetCreditItem(tracking.UID, tracking.Currency)
	if itemTO.HasError() {
		return itemTO.Error
	}
	item := itemTO.Object.(bean.CreditItem)

	if item.Status == bean.CREDIT_ITEM_STATUS_CREATE || item.Status == bean.CREDIT_ITEM_STATUS_INACTIVE {
		item.Status = bean.CREDIT_ITEM_STATUS_INACTIVE
	}
	item.SubStatus = ""
	deposit.Status = bean.CREDIT_DEPOSIT_STATUS_FAILED

	s.dao.FinishFailedDepositCreditItem(&item, &deposit, &tracking)

	return err
}

func (s CreditService) getConfirmationRange(amount decimal.Decimal) int {
	if amount.LessThan(decimal.NewFromFloat(0.5)) {
		return 1
	} else if amount.LessThan(decimal.NewFromFloat(1)) {
		return 3
	}

	return 6
}
