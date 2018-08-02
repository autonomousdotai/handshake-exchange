package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"time"
)

type ReferralService struct {
	dao     *dao.ReferralDao
	miscDao *dao.MiscDao
}

func (s ReferralService) AddReferral(userId string) (ce SimpleContextError) {
	err := s.dao.AddReferral(userId, false)
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}

	return
}

func (s ReferralService) AddReferralRecord(userId string, toUserId string) (ce SimpleContextError) {
	err := s.dao.AddReferralRecord(bean.ReferralRecord{
		UID:       userId,
		ToUID:     toUserId,
		ExpiredAt: time.Now().UTC().Add(6 * 30 * 24 * time.Hour), // 6 * 30 days
	})
	if err != nil {
		ce.SetError(api_error.AddDataFailed, err)
		return
	}
	return
}

func (s ReferralService) AddReferralOfferStoreShake(profile bean.Profile, offer bean.OfferStore, offerShake bean.OfferStoreShake) (ce SimpleContextError) {
	uid := profile.ReferralUser

	if uid != "" {
		systemConfigTO := s.miscDao.GetSystemConfigFromCache(bean.CONFIG_OFFER_STORE_REFERRAL_PERCENTAGE)
		if ce.FeedDaoTransfer(api_error.GetDataFailed, systemConfigTO) {
			return
		}
		systemConfig := systemConfigTO.Object.(bean.SystemConfig)
		rewardPercentage := common.StringToDecimal(systemConfig.Value)
		fee := common.StringToDecimal(offerShake.Fee)
		reward := rewardPercentage.Mul(fee)

		referralOfferShake := bean.ReferralOfferStoreShakeRecord{
			UID:              uid,
			ToUID:            offerShake.UID,
			ToUsername:       offerShake.Username,
			Offer:            offer.Id,
			OfferShake:       offerShake.Id,
			Status:           bean.REFERRAL_OFFER_STORE_SHAKE_PENDING,
			Currency:         offerShake.Currency,
			RewardPercentage: rewardPercentage.String(),
			Reward:           reward.String(),
		}
		err := s.dao.AddReferralOfferStoreShake(profile.CreatedAt, referralOfferShake)
		if err != nil {
			ce.SetError(api_error.AddDataFailed, err)
			return
		}
	}

	return
}
