package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
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
		}

	}

	deposit = bean.CreditDeposit{
		UID:      userId,
		ItemRef:  dao.GetCreditItemItemPath(userId, body.Currency),
		Status:   bean.CREDIT_DEPOSIT_STATUS_CREATED,
		Currency: body.Currency,
		Amount:   body.Amount,
	}

	if body.Currency != bean.ETH.Code {
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

	tracking = bean.CreditOnChainActionTracking{
		UID:        userId,
		ItemRef:    deposit.ItemRef,
		DepositRef: dao.GetCreditDepositItemPath(body.Currency, deposit.Id),
		TxHash:     body.TxHash,
		Action:     body.Action,
		Reason:     body.Reason,
		Currency:   body.Currency,
	}
	s.dao.AddCreditOnChainActionTracking(&tracking)

	return
}
