package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const CREDIT_CARD_STATUS_OK = "ok"
const CREDIT_CARD_STATUS_DISPUTED = "disputed"

type ProfileRequest struct {
	Id int `json:"id" validate:"required"`
}

type Profile struct {
	UserId           string          `json:"-" firestore:"user_id"`
	CreditCardStatus string          `json:"-" firestore:"credit_card_status"`
	CreditCard       UserCreditCard  `json:"credit_card" firestore:"credit_card"`
	ActiveOffers     map[string]bool `json:"-" firestore:"active_offers"`
	OfferRejectLock  OfferRejectLock `json:"offer_reject_lock" firestore:"offer_reject_lock"`
}

type OfferRejectLock struct {
	Duration  int64     `json:"duration" firestore:"duration"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

func (o OfferRejectLock) GetAddOfferRejectLock() map[string]interface{} {
	return map[string]interface{}{
		"duration":   o.Duration,
		"created_at": firestore.ServerTimestamp,
	}
}

func (profile Profile) GetAddProfile() map[string]interface{} {
	offerMap := map[string]bool{
		"ETH": false, "BTC": false,
	}
	return map[string]interface{}{
		"user_id":            profile.UserId,
		"credit_card_status": CREDIT_CARD_STATUS_OK,
		"active_offers":      offerMap,
		"created_at":         firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateOfferProfile() map[string]interface{} {
	return map[string]interface{}{
		"active_offers": profile.ActiveOffers,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateOfferRejectLock() map[string]interface{} {
	return map[string]interface{}{
		"offer_reject_lock": profile.OfferRejectLock.GetAddOfferRejectLock(),
		"updated_at":        firestore.ServerTimestamp,
	}
}

type UserCreditCard struct {
	CCNumber       string `json:"cc_number" firestore:"cc_number"`
	ExpirationDate string `json:"expiration_date" firestore:"expiration_date"`
	Token          string `json:"token" firestore:"token"`
}

func (user UserCreditCard) GetUpdateProfileCreditCard() map[string]interface{} {
	return map[string]interface{}{
		"credit_card": map[string]interface{}{
			"cc_number":       user.CCNumber,
			"expiration_date": user.ExpirationDate,
			"token":           user.Token,
		},
	}
}

type UserCreditCardLimit struct {
	Level    int64     `json:"level" firestore:"level"`
	Amount   string    `json:"amount" firestore:"amount"`
	Limit    int64     `json:"limit" firestore:"limit"`
	EndDate  time.Time `json:"end_date" firestore:"end_date"`
	Duration int64     `json:"duration" firestore:"duration"`
}

func (limit UserCreditCardLimit) GetAddUserCreditCardLimit() map[string]interface{} {
	return map[string]interface{}{
		"level":      limit.Level,
		"amount":     limit.Amount,
		"limit":      limit.Limit,
		"duration":   limit.Duration,
		"end_date":   limit.EndDate,
		"created_at": firestore.ServerTimestamp,
	}
}

func (limit UserCreditCardLimit) GetUpdateLevel() map[string]interface{} {
	return map[string]interface{}{
		"level":      limit.Level,
		"amount":     "0",
		"limit":      limit.Limit,
		"duration":   limit.Duration,
		"end_date":   limit.EndDate,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (limit UserCreditCardLimit) GetUpdateAmount() map[string]interface{} {
	return map[string]interface{}{
		"amount":     limit.Amount,
		"updated_at": firestore.ServerTimestamp,
	}
}

type UserCreditCardLimitTrack struct {
	UID      string `json:"uid" firestore:"uid"`
	Level    int64  `json:"level" firestore:"level"`
	Left     int64  `json:"left" firestore:"left"`
	Duration int64  `json:"duration" firestore:"duration"`
}

func (limit UserCreditCardLimitTrack) GetAddUserCreditCardLimitTrack() map[string]interface{} {
	return map[string]interface{}{
		"uid":        limit.UID,
		"level":      limit.Level,
		"left":       limit.Left,
		"duration":   limit.Duration,
		"created_at": firestore.ServerTimestamp,
	}
}

func (limit UserCreditCardLimitTrack) GetUpdateLeft() map[string]interface{} {
	return map[string]interface{}{
		"left":       limit.Left,
		"updated_at": firestore.ServerTimestamp,
	}
}

type NewUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
	Country  string `json:"country" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type ChangePassword struct {
	Password string `json:"password" validate:"required"`
}

type ResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordToken struct {
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

type VerifyEmailToken struct {
	Token string `json:"token" validate:"required"`
}
