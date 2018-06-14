package service

import (
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserDaoFake struct {
}

func (dao UserDaoFake) GetProfile(userId string) (t dao.TransferObject) {
	t.Found = false
	return
}
func (dao UserDaoFake) AddProfile(profile bean.Profile) error {
	return nil
}
func (dao UserDaoFake) UpdateProfileCreditCard(userId string, creditCard bean.UserCreditCard, userCCLimit bean.UserCreditCardLimit) error {
	return nil
}
func (dao UserDaoFake) UpdateProfileOfferRejectLock(userId string, lock bean.OfferRejectLock) error {
	return nil
}
func (dao UserDaoFake) UpdateUserCCLimitAmount(userId string, token string, amount decimal.Decimal) error {
	return nil
}
func (dao UserDaoFake) UpdateUserCCLimitTracks() (userIds []string, t dao.TransferObject) {
	return
}
func (dao UserDaoFake) GetCCLimit(userId string, token string) (t dao.TransferObject) {
	return
}
func (dao UserDaoFake) GetUserCCLimitEndTracks() (t dao.TransferObject) {
	return
}
func (dao UserDaoFake) UpgradeCCLimitLevel(userId string, token string, limit bean.UserCreditCardLimit) error {
	return nil
}

func TestAddProfileSuccess(t *testing.T) {
	profile := bean.Profile{
		UserId: "1",
	}

	serviceInst := UserService{
		dao: &UserDaoFake{},
	}

	err := serviceInst.AddProfile(profile)
	assert.Equal(t, nil, err)
}
