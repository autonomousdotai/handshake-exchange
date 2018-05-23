package service

import (
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/go-errors/errors"
	"github.com/shopspring/decimal"
	"os"
	"strconv"
	"time"
)

type UserService struct {
	dao     *dao.UserDao
	miscDao *dao.MiscDao
}

func (s UserService) AddProfile(profile bean.Profile) error {
	to := s.dao.GetProfile(profile.UserId)
	var err error
	if to.Error == nil {
		if to.Found {
			err = errors.New(api_error.ProfileExists)
		} else {
			err = s.dao.AddProfile(profile)
		}
	}

	return err
}

func (s UserService) GetCCLimitLevel(userId string) (limit bean.UserCreditCardLimit, ce SimpleContextError) {
	to := s.dao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	profile := to.Object.(bean.Profile)
	to = s.dao.GetCCLimit(userId, profile.CreditCard.Token)

	if to.Error != nil || !to.Found {
		limit, to.Error = s.GetUserCCLimitFirstLevel()
		if to.Error != nil {
			to.SetError(api_error.GetDataFailed, to.Error)
		}
		to.Found = true

	} else {
		limit = to.Object.(bean.UserCreditCardLimit)
	}

	return
}

func (s UserService) GetUserCCLimitFirstLevel() (limit bean.UserCreditCardLimit, err error) {
	ccLimitTO := s.miscDao.GetCCLimitByLevelFromCache("1")
	if !ccLimitTO.HasError() {
		ccLimit := ccLimitTO.Object.(bean.CCLimit)
		duration := ccLimit.Duration * int64(time.Hour*24)
		limit = bean.UserCreditCardLimit{
			Level:    ccLimit.Level,
			Amount:   common.Zero.String(),
			Limit:    ccLimit.Limit,
			Duration: ccLimit.Duration,
			EndDate:  time.Now().UTC().Add(time.Duration(duration)),
		}
	} else {
		err = ccLimitTO.Error
	}

	return
}

func (s UserService) UpgradeCCLimitLevel(userId string) (ce SimpleContextError) {
	to := s.dao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	profile := to.Object.(bean.Profile)
	to = s.dao.GetCCLimit(userId, profile.CreditCard.Token)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	creditCardLimit := to.Object.(bean.UserCreditCardLimit)
	finalLevel, _ := strconv.Atoi(os.Getenv("MAX_CC_LIMIT_LEVEL"))
	if creditCardLimit.Level < int64(finalLevel) {
		creditCardLimit.Level += 1
	} else {
		// Reset the last limit
	}

	cacheTO := s.miscDao.GetCCLimitByLevelFromCache(strconv.Itoa(int(creditCardLimit.Level)))
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cacheTO) {
		return
	}

	limit := cacheTO.Object.(bean.CCLimit)
	creditCardLimit.Duration = limit.Duration
	creditCardLimit.Amount = common.Zero.String()
	creditCardLimit.Limit = limit.Limit

	duration := creditCardLimit.Duration * int64(time.Hour*24)
	creditCardLimit.EndDate = time.Now().UTC().Add(time.Duration(duration))

	err := s.dao.UpgradeCCLimitLevel(userId, profile.CreditCard.Token, creditCardLimit)
	if ce.SetError(api_error.UpdateDataFailed, err) {
		return
	}

	return
}

func (s UserService) CheckCCLimit(userId string, amountStr string) (ce SimpleContextError) {
	to := s.dao.GetProfile(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, to) {
		return
	}
	profile := to.Object.(bean.Profile)
	to = s.dao.GetCCLimit(userId, profile.CreditCard.Token)

	var err error
	var creditCardLimit bean.UserCreditCardLimit
	if to.Error != nil || !to.Found {
		// First time
		creditCardLimit, err = s.GetUserCCLimitFirstLevel()
		if ce.SetError(api_error.GetDataFailed, err) {
			return
		}
	} else {
		creditCardLimit = to.Object.(bean.UserCreditCardLimit)
	}

	currentAmount, _ := decimal.NewFromString(creditCardLimit.Amount)
	limit, _ := decimal.NewFromString(strconv.Itoa(int(creditCardLimit.Limit)))
	amount, _ := decimal.NewFromString(amountStr)

	if currentAmount.Add(amount).GreaterThan(limit) {
		ce.SetError(api_error.Success, errors.New(api_error.Success))
	}

	return
}

func (s UserService) UpdateUserCCLimitTracks() (ce SimpleContextError) {
	userIds, updateTO := s.dao.UpdateUserCCLimitTracks()
	ce.FeedDaoTransfer(api_error.UpdateDataFailed, updateTO)

	trackTO := s.dao.GetUserCCLimitEndTracks()
	if !ce.FeedDaoTransfer(api_error.GetDataFailed, trackTO) {
		for _, obj := range trackTO.Objects {
			track := obj.(bean.UserCreditCardLimitTrack)
			// Exclude those UID already update left
			if !common.StringInSlice(track.UID, userIds) {
				s.UpgradeCCLimitLevel(track.UID)
			}
		}
	}

	return
}
