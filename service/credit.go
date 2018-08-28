package service

import (
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
)

type CreditService struct {
	dao     *dao.CreditDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CreditService) Create() (credit bean.Credit, ce SimpleContextError) {
	return
}
