package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const REFERRAL_OFFER_STORE_SHAKE_PENDING = "pending"
const REFERRAL_OFFER_STORE_SHAKE_PAID = "paid"

// referrals/{user_id}
type ReferralCount struct {
	UID   string `json:"uid" firestore:"uid"`
	Count int64  `json:"count" firestore:"count"`
}

func (b ReferralCount) GetAddData() map[string]interface{} {
	return map[string]interface{}{
		"uid":        b.UID,
		"count":      0,
		"created_at": firestore.ServerTimestamp,
	}
}

func (b ReferralCount) GetUpdateReferralCount() map[string]interface{} {
	return map[string]interface{}{
		"count":      b.Count,
		"updated_at": firestore.ServerTimestamp,
	}
}

// referrals/{user_id}/records/{to_user_id}
type ReferralRecord struct {
	UID        string    `json:"uid" firestore:"uid"`
	ToUID      string    `json:"to_uid" firestore:"to_uid"`
	ToUsername string    `json:"to_username" firestore:"to_username"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	ExpiredAt  time.Time `json:"expired_at" firestore:"expired_at"`
}

func (b ReferralRecord) GetAddData() map[string]interface{} {
	return map[string]interface{}{
		"uid":         b.UID,
		"to_uid":      b.ToUID,
		"to_username": b.ToUsername,
		"created_at":  firestore.ServerTimestamp,
		"expired_at":  b.ExpiredAt,
	}
}

// referrals/{user_id}/currencies/{currency}
type ReferralOfferStoreShake struct {
	UID               string    `json:"uid" firestore:"uid"`
	ToUID             string    `json:"to_uid" firestore:"to_uid"`
	ToUsername        string    `json:"to_username" firestore:"to_username"`
	Currency          string    `json:"currency" firestore:"currency"`
	Reward            string    `json:"reward" firestore:"reward"`
	PendingReward     string    `json:"pending_reward" firestore:"pending_reward"`
	TotalReward       string    `json:"total_reward" firestore:"total_reward"`
	ReferralCreatedAt time.Time `json:"referral_created_at" firestore:"referral_created_at"`
}

func (b ReferralOfferStoreShake) GetAddData() map[string]interface{} {
	return map[string]interface{}{
		"uid":            b.UID,
		"to_uid":         b.ToUID,
		"to_username":    b.ToUsername,
		"currency":       b.Currency,
		"reward":         b.Reward,
		"pending_reward": b.PendingReward,
		"total_reward":   b.TotalReward,
		"created_at":     firestore.ServerTimestamp,
	}
}

func (b ReferralOfferStoreShake) GetOverridePendingReward() map[string]interface{} {
	return map[string]interface{}{
		"uid":            b.UID,
		"to_uid":         b.ToUID,
		"to_username":    b.ToUsername,
		"currency":       b.Currency,
		"reward":         b.Reward,
		"pending_reward": b.PendingReward,
		"total_reward":   b.TotalReward,
		"created_at":     firestore.ServerTimestamp,
	}
}

func (b ReferralOfferStoreShake) GetUpdatePendingReward() map[string]interface{} {
	return map[string]interface{}{
		"pending_reward": b.PendingReward,
		"total_reward":   b.TotalReward,
		"updated_at":     firestore.ServerTimestamp,
	}
}

func (b ReferralOfferStoreShake) GetUpdateReward() map[string]interface{} {
	return map[string]interface{}{
		"pending_reward": b.PendingReward,
		"reward":         b.Reward,
		"updated_at":     firestore.ServerTimestamp,
	}
}

// referrals/{user_id}/currencies/{currency}/records/{offer_shake_id}
type ReferralOfferStoreShakeRecord struct {
	UID              string    `json:"uid" firestore:"uid"`
	ToUID            string    `json:"to_uid" firestore:"to_uid"`
	ToUsername       string    `json:"to_username" firestore:"to_username"`
	Offer            string    `json:"offer" firestore:"offer"`
	OfferShake       string    `json:"offer_shake" firestore:"offer_shake"`
	Status           string    `json:"status" firestore:"status"`
	Currency         string    `json:"currency" firestore:"currency"`
	RewardPercentage string    `json:"percentage" firestore:"percentage"`
	Reward           string    `json:"reward" firestore:"reward"`
	CreatedAt        time.Time `json:"created_at" firestore:"created_at"`
}

func (b ReferralOfferStoreShakeRecord) GetAddData() map[string]interface{} {
	return map[string]interface{}{
		"uid":               b.UID,
		"to_uid":            b.ToUID,
		"to_username":       b.ToUsername,
		"offer":             b.Offer,
		"offer_shake":       b.OfferShake,
		"status":            REFERRAL_OFFER_STORE_SHAKE_PENDING,
		"currency":          b.Currency,
		"reward_percentage": b.RewardPercentage,
		"reward":            b.Reward,
		"created_at":        firestore.ServerTimestamp,
	}
}

func (b ReferralOfferStoreShakeRecord) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}
