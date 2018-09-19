package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
)

type CashService struct {
	dao     *dao.CashDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CashService) GetCashStore(userId string) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}

	if cashTO.Found {
		cash = cashTO.Object.(bean.CashStore)
	} else {
		ce.NotFound = true
	}

	return
}

func (s CashService) AddCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}

	var err error
	if cashTO.Found {
		ce.SetStatusKey(api_error.CashStoreExists)
	} else {
		body.UID = userId
		cash = body

		err = s.dao.AddCashStore(&cash)
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
			return
		}
		ce.NotFound = false
	}

	return
}

func (s CashService) UpdateCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}
	cash = cashTO.Object.(bean.CashStore)

	err := s.dao.UpdateCashStore(&cash)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}

	return
}
