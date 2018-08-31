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
		amount := decimal.Zero
		isSuccess, isPending, amount, errChain := crypto_service.GetTransactionReceipt(trackingItem.TxHash, trackingItem.Currency)
		fmt.Println(fmt.Sprintf("%s %s %s %s", isSuccess, isPending, amount.String(), errChain))

		fmt.Println(fmt.Sprintf("%s %s %s %s", trackingItem.Id, trackingItem.UID, trackingItem.TxHash, amount.String()))
		if errChain == nil {
			if isSuccess && !isPending && amount.GreaterThan(common.Zero) {
				trackingItem.Amount = amount.String()
				s.finishTrackingItem(trackingItem)
			}
		} else {
			ce.SetError(api_error.ExternalApiFailed, errChain)
		}
	}

	trackingTO = s.dao.ListCreditOnChainActionTracking(bean.BTC.Code)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, trackingTO) {
		return
	}
	for _, item := range trackingTO.Objects {
		trackingItem := item.(bean.CreditOnChainActionTracking)
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
	}

	trackingTO = s.dao.ListCreditOnChainActionTracking(bean.BCH.Code)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, trackingTO) {
		return
	}
	for _, item := range trackingTO.Objects {
		trackingItem := item.(bean.CreditOnChainActionTracking)
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
	if depositTO.HasError() {
		return depositTO.Error
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
			if err != nil {
				ce.SetError(api_error.AddDataFailed, err)
			}
			level += 1
		}
	}

	return
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
