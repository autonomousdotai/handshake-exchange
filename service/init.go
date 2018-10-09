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
	dao:      &dao.OfferStoreDaoInst,
	miscDao:  &dao.MiscDaoInst,
	userDao:  &dao.UserDaoInst,
	transDao: &dao.TransactionDaoInst,
	offerDao: &dao.OfferDaoInst,
}

var ReferralServiceInst = ReferralService{
	dao:     &dao.ReferralDao{},
	miscDao: &dao.MiscDaoInst,
}

var CreditServiceInst = CreditService{
	dao:     &dao.CreditDaoInst,
	miscDao: &dao.MiscDaoInst,
	userDao: &dao.UserDaoInst,
}

var CashServiceInst = CashService{
	dao:     &dao.CashDaoInst,
	miscDao: &dao.MiscDaoInst,
	userDao: &dao.UserDaoInst,
}

var CoinServiceInst = CoinService{
	dao:     &dao.CoinDaoInst,
	miscDao: &dao.MiscDaoInst,
	userDao: &dao.UserDaoInst,
}
