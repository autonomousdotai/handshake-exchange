package service

import "github.com/autonomousdotai/handshake-exchange/dao"

var UserServiceInst = UserService{
	dao:     &dao.UserDaoInst,
	miscDao: &dao.MiscDaoInst,
}

var CreditCardServiceInst = CreditCardService{
	dao:     &dao.CreditCardDaoInst,
	miscDao: &dao.MiscDaoInst,
	userDao: &dao.UserDaoInst,
}
