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
	"github.com/ninjadotorg/handshake-exchange/integration/exchangecreditatm_service"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strconv"
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

func (s CreditService) DeactivateCredit(userId string, currency string) (credit bean.Credit, ce SimpleContextError) {
	creditTO := s.dao.GetCredit(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
		return
	}
	credit = creditTO.Object.(bean.Credit)

	creditItemTO := s.dao.GetCreditItem(userId, currency)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditItemTO) {
		return
	}
	creditItem := creditItemTO.Object.(bean.CreditItem)

	percentage, _ := strconv.Atoi(creditItem.Percentage)
	poolTO := s.dao.GetCreditPool(currency, percentage)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, poolTO) {
		return
	}
	pool := poolTO.Object.(bean.CreditPool)

	if creditItem.Status == bean.CREDIT_ITEM_STATUS_ACTIVE && creditItem.SubStatus != bean.CREDIT_ITEM_SUB_STATUS_TRANSFERRING && creditItem.LockedSale == false {
		creditItem.Status = bean.CREDIT_ITEM_STATUS_INACTIVE

		itemHistory := bean.CreditBalanceHistory{}
		itemHistory.ModifyType = bean.CREDIT_POOL_MODIFY_TYPE_CLOSE

		poolHistory := bean.CreditPoolBalanceHistory{
			ModifyType: bean.CREDIT_POOL_MODIFY_TYPE_CLOSE,
		}

		poolOrders := make([]bean.CreditPoolOrder, 0)
		poolOrdersTO := s.dao.ListCreditPoolOrderUser(creditItem.Currency, userId)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, poolOrdersTO) {
			return
		}
		for _, poolOrderItem := range poolOrdersTO.Objects {
			poolOrders = append(poolOrders, poolOrderItem.(bean.CreditPoolOrder))
		}
		err := s.dao.RemoveCreditItem(&creditItem, &itemHistory, &pool, poolOrders, &poolHistory)
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
			return
		}
		client := exchangecreditatm_service.ExchangeCreditAtmClient{}
		amount := common.StringToDecimal(itemHistory.Change)

		if currency == bean.ETH.Code {
			txHash, onChainErr := client.ReleasePartialFund(userId, 2, amount, creditItem.UserAddress)
			if onChainErr != nil {
				fmt.Println(onChainErr)
			} else {
			}
			fmt.Println(txHash)
		} else {
			coinbaseTx, errWithdraw := coinbase_service.SendTransaction(creditItem.UserAddress, amount.String(), currency,
				fmt.Sprintf("Refund userId = %s", creditItem.UID), creditItem.UID)
			if errWithdraw != nil {
				fmt.Println(errWithdraw)
			} else {
			}
			fmt.Println(coinbaseTx)
		}
	} else {
		ce.SetStatusKey(api_error.CreditItemStatusInvalid)
	}

	return
}

func (s CreditService) AddDeposit(userId string, body bean.CreditDepositInput) (deposit bean.CreditDeposit, ce SimpleContextError) {
	var err error

	creditTO := s.dao.GetCredit(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
		return
	}
	credit := creditTO.Object.(bean.Credit)

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
		pNum, _ := strconv.Atoi(body.Percentage)
		if pNum < 0 || pNum > 200 {
			ce.SetStatusKey(api_error.InvalidRequestBody)
			return
		}
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
			if creditItem.Status == bean.CREDIT_ITEM_STATUS_INACTIVE {
				//Reactivate
				creditItem.Status = bean.CREDIT_ITEM_STATUS_ACTIVE
				creditItem.Percentage = body.Percentage
				creditItem.UserAddress = body.UserAddress

				err = s.dao.UpdateCreditItem(&creditItem)
				if err != nil {
					ce.SetError(api_error.AddDataFailed, err)
					return
				}
			} else {
				if creditItem.Percentage != body.Percentage {
					ce.SetStatusKey(api_error.InvalidRequestBody)
					return
				}
			}

			if creditItem.Status == bean.CREDIT_ITEM_STATUS_CREATE || creditItem.SubStatus == bean.CREDIT_ITEM_SUB_STATUS_TRANSFERRING {
				ce.SetStatusKey(api_error.CreditItemStatusInvalid)
				return
			}
			creditItem.UserAddress = body.UserAddress
		}
	}

	deposit = bean.CreditDeposit{
		UID:        userId,
		ItemRef:    dao.GetCreditItemItemPath(userId, body.Currency),
		Status:     bean.CREDIT_DEPOSIT_STATUS_CREATED,
		Currency:   body.Currency,
		Amount:     body.Amount,
		Percentage: body.Percentage,
	}

	if body.Currency != bean.ETH.Code {
		resp, errCoinbase := coinbase_service.GenerateAddress(body.Currency)
		if errCoinbase != nil {
			ce.SetError(api_error.ExternalApiFailed, errCoinbase)
			return
		}
		deposit.SystemAddress = resp.Data.Address
	}

	err = s.dao.AddCreditDeposit(&creditItem, &deposit)
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}

	chainId, _ := strconv.Atoi(credit.ChainId)
	deposit.CreatedAt = time.Now().UTC()
	solr_service.UpdateObject(bean.NewSolrFromCreditDeposit(deposit, int64(chainId)))
	s.dao.UpdateNotificationCreditItem(creditItem)

	return
}

func (s CreditService) AddTracking(userId string, body bean.CreditOnChainActionTrackingInput) (tracking bean.CreditOnChainActionTracking, ce SimpleContextError) {
	creditTO := s.dao.GetCredit(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
		return
	}
	credit := creditTO.Object.(bean.Credit)

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

	chainId, _ := strconv.Atoi(credit.ChainId)
	solr_service.UpdateObject(bean.NewSolrFromCreditDeposit(deposit, int64(chainId)))
	s.dao.UpdateNotificationCreditItem(item)

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
	for percentage <= 200 {
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

	return 0, errors.New("not enough")
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
	userTransList := make([]*bean.CreditTransaction, 0)

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
	var transUserId string
	for k, v := range userTransMap {
		trans.UIDs = append(trans.UIDs, k)
		userTransList = append(userTransList, v)

		transUserId = k
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
			return
		}
	}

	creditTO := s.dao.GetCredit(transUserId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, creditTO) {
		return
	}
	credit := creditTO.Object.(bean.Credit)

	chainId, _ := strconv.Atoi(credit.ChainId)
	for _, userTrans := range userTransList {
		userTrans.CreatedAt = time.Now().UTC()
		solr_service.UpdateObject(bean.NewSolrFromCreditTransaction(*userTrans, int64(chainId)))
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

func (s CreditService) FinishCreditTransaction(currency string, id string, offerRef string,
	revenue decimal.Decimal, fee decimal.Decimal) (ce SimpleContextError) {
	transTO := s.dao.GetCreditTransaction(currency, id)

	if ce.FeedDaoTransfer(api_error.GetDataFailed, transTO) {
		return
	}
	trans := transTO.Object.(bean.CreditTransaction)
	trans.OfferRef = offerRef
	trans.Status = bean.CREDIT_TRANSACTION_STATUS_SUCCESS
	trans.SubStatus = bean.CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED
	trans.Revenue = revenue.RoundBank(2).String()

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
		Change:     amount.Neg().String(),
	}

	items := make([]bean.CreditItem, 0)
	itemHistories := make([]bean.CreditBalanceHistory, 0)
	transList := make([]*bean.CreditTransaction, 0)
	var transUID string
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

		percentageAmount := userAmount.Div(amount)
		userFee := percentageAmount.Mul(fee)
		userRevenue := percentageAmount.Mul(revenue).Sub(userFee)

		userTrans.OfferRef = offerRef
		userTrans.Status = bean.CREDIT_TRANSACTION_STATUS_SUCCESS
		userTrans.SubStatus = bean.CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED
		userTrans.Fee = userFee.RoundBank(2).String()
		userTrans.Revenue = userRevenue.RoundBank(2).String()
		transList = append(transList, &userTrans)

		itemHistory := bean.CreditBalanceHistory{
			ItemRef:    dao.GetCreditItemItemPath(userId, trans.Currency),
			ModifyRef:  dao.GetCreditTransactionItemUserPath(userId, currency, userTrans.Id),
			ModifyType: bean.CREDIT_POOL_MODIFY_TYPE_PURCHASE,
			Change:     userAmount.Neg().String(),
		}
		itemHistories = append(itemHistories, itemHistory)

		transUID = userTrans.UID
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
		return
	}

	creditTO := s.dao.GetCredit(transUID)
	if !creditTO.HasError() {
		credit := creditTO.Object.(bean.Credit)
		chainId, _ := strconv.Atoi(credit.ChainId)
		for _, userTrans := range transList {
			solr_service.UpdateObject(bean.NewSolrFromCreditTransaction(*userTrans, int64(chainId)))
		}
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
	body.Status = bean.CREDIT_WITHDRAW_STATUS_CREATED

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
		return
	}

	withdraw = body
	withdraw.CreatedAt = time.Now().UTC()
	chainId, _ := strconv.Atoi(credit.ChainId)
	solr_service.UpdateObject(bean.NewSolrFromCreditWithdraw(withdraw, int64(chainId)))

	return
}

func (s CreditService) SetupCreditPool() (ce SimpleContextError) {
	for _, currency := range []string{bean.BTC.Code, bean.ETH.Code, bean.BCH.Code} {
		level := 0
		for level <= 200 {
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

func (s CreditService) SyncCreditTransactionToSolr(currency string, id string) (trans bean.CreditTransaction, ce SimpleContextError) {
	transTO := s.dao.GetCreditTransaction(currency, id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, transTO) {
		return
	}
	trans = transTO.Object.(bean.CreditTransaction)
	for _, userId := range trans.UIDs {
		creditTO := s.dao.GetCredit(userId)
		credit := creditTO.Object.(bean.Credit)

		transUserTO := s.dao.GetCreditTransactionUser(userId, currency, id)
		transUser := transUserTO.Object.(bean.CreditTransaction)

		chainId, _ := strconv.Atoi(credit.ChainId)
		solr_service.UpdateObject(bean.NewSolrFromCreditTransaction(transUser, int64(chainId)))
	}

	return
}

func (s CreditService) SyncCreditDepositToSolr(currency string, id string) (deposit bean.CreditDeposit, ce SimpleContextError) {
	depositTO := s.dao.GetCreditDeposit(currency, id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, depositTO) {
		return
	}
	deposit = depositTO.Object.(bean.CreditDeposit)
	creditTO := s.dao.GetCredit(deposit.UID)
	credit := creditTO.Object.(bean.Credit)
	chainId, _ := strconv.Atoi(credit.ChainId)
	solr_service.UpdateObject(bean.NewSolrFromCreditDeposit(deposit, int64(chainId)))

	return
}

func (s CreditService) SyncCreditWithdrawToSolr(id string) (withdraw bean.CreditWithdraw, ce SimpleContextError) {
	withdrawTO := s.dao.GetCreditWithdraw(id)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, withdrawTO) {
		return
	}
	withdraw = withdrawTO.Object.(bean.CreditWithdraw)
	creditTO := s.dao.GetCredit(withdraw.UID)
	credit := creditTO.Object.(bean.Credit)
	chainId, _ := strconv.Atoi(credit.ChainId)
	solr_service.UpdateObject(bean.NewSolrFromCreditWithdraw(withdraw, int64(chainId)))

	return
}

func (s CreditService) finishTrackingItem(tracking bean.CreditOnChainActionTracking) error {
	var err error

	creditTO := s.dao.GetCredit(tracking.UID)
	if creditTO.HasError() {
		return creditTO.Error
	}
	credit := creditTO.Object.(bean.Credit)

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

	chainId, _ := strconv.Atoi(credit.ChainId)
	solr_service.UpdateObject(bean.NewSolrFromCreditDeposit(deposit, int64(chainId)))
	s.dao.UpdateNotificationCreditItem(item)

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
