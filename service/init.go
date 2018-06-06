package service

import "github.com/ninjadotorg/handshake-exchange/dao"

var UserServiceInst = UserService{
	dao:     &dao.UserDaoInst,
	miscDao: &dao.MiscDaoInst,
}

var CreditCardServiceInst = CreditCardService{
	dao:      &dao.CreditCardDaoInst,
	miscDao:  &dao.MiscDaoInst,
	userDao:  &dao.UserDaoInst,
	transDao: &dao.TransactionDaoInst,
}

var OfferServiceInst = OfferService{
	dao:      &dao.OfferDaoInst,
	miscDao:  &dao.MiscDaoInst,
	userDao:  &dao.UserDaoInst,
	transDao: &dao.TransactionDaoInst,
}

var OfferStoreServiceInst = OfferStoreService{
	dao:     &dao.OfferStoreDaoInst,
	miscDao: &dao.MiscDaoInst,
	userDao: &dao.UserDaoInst,
}
