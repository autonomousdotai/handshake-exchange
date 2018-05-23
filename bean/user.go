package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const PROFILE_STATUS_NEED_UPDATE = "need_update"
const PROFILE_STATUS_NEED_UPDATE_AGAIN = "need_update_again"
const PROFILE_STATUS_VERIFYING = "verifying"
const PROFILE_STATUS_VERIFIED = "verified"

type Profile struct {
	FirstName      string            `json:"first_name" firestore:"first_name"`
	LastName       string            `json:"last_name" firestore:"last_name"`
	Email          string            `json:"email" firestore:"email"`
	Username       string            `json:"username" firestore:"username"`
	Country        string            `json:"country" firestore:"country"`
	Currency       string            `json:"currency" firestore:"currency"`
	Avatar         string            `json:"avatar" firestore:"avatar"`
	Initialized    bool              `json:"-" firestore:"initialized"`
	Status         string            `json:"status" firestore:"status"`
	EmailVerified  bool              `json:"email_verified" firestore:"email_verified"`
	Verified       bool              `json:"verified" firestore:"verified"`
	TwoFA          bool              `json:"two_fa" firestore:"two_fa"`
	TwoFASecretKey string            `json:"two_fa_secret_key" firestore:"two_fa_secret_key"`
	VerifyInfo     VerifyProfileInfo `json:"verify_info" firestore:"verify_info"`
	CreditCard     UserCreditCard    `json:"credit_card" firestore:"credit_card"`
	BankAccount    UserBankAccount   `json:"bank_account" firestore:"bank_account"`
}

type AddProfile struct {
	Username string `json:"username" firestore:"username"`
	Email    string `json:"email" firestore:"email" validate:"required,email"`
	Country  string `json:"country" firestore:"country" validate:"required"`
	Currency string `json:"currency" firestore:"currency" validate:"required"`
}

func (profile AddProfile) GetAddProfile() map[string]interface{} {
	return map[string]interface{}{
		"username":       profile.Username,
		"email":          profile.Email,
		"country":        profile.Country,
		"currency":       profile.Currency,
		"initialized":    true,
		"status":         PROFILE_STATUS_NEED_UPDATE,
		"verified":       false,
		"email_verified": false,
		"two_fa":         false,
	}
}

func (profile Profile) GetUpdateProfile() map[string]interface{} {
	return map[string]interface{}{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
		"username":   profile.Username,
		"avatar":     profile.Avatar,
		// "country":    profile.Country,  // Country will update by verify info
		"updated_at": firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateProfileCountry() map[string]interface{} {
	return map[string]interface{}{
		"country": profile.Country,
	}
}

func (profile Profile) GetUpdateProfile2FA() map[string]interface{} {
	return map[string]interface{}{
		"two_fa":            profile.TwoFA,
		"two_fa_secret_key": profile.TwoFASecretKey,
		"updated_at":        firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateVerifyingProfile() map[string]interface{} {
	return map[string]interface{}{
		"status":     PROFILE_STATUS_VERIFYING,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateVerifiedProfile() map[string]interface{} {
	return map[string]interface{}{
		"status":      PROFILE_STATUS_VERIFIED,
		"verified":    true,
		"verify_info": profile.VerifyInfo,
		"updated_at":  firestore.ServerTimestamp,
	}
}

func (profile Profile) GetUpdateEmailVerifiedProfile() map[string]interface{} {
	return map[string]interface{}{
		"email_verified": true,
		"updated_at":     firestore.ServerTimestamp,
	}
}

type PublicProfile struct {
	FirstName string `json:"first_name" firestore:"first_name"`
	LastName  string `json:"last_name" firestore:"last_name"`
	Username  string `json:"username" firestore:"username"`
	Country   string `json:"country" firestore:"country"`
	Currency  string `json:"currency" firestore:"currency"`
	Avatar    string `json:"avatar" firestore:"avatar"`
}

type UsernameMap struct {
	UID      string `json:"-" firestore:"uid"`
	Username string `json:"username" firestore:"username"`
}

func (b UsernameMap) GetAddUsernameMap() map[string]interface{} {
	return map[string]interface{}{
		"uid":        b.UID,
		"username":   b.Username,
		"created_at": firestore.ServerTimestamp,
	}
}

type VerifyProfileInfo struct {
	UID                  string    `json:"-" firestore:"uid"`
	LegalName            string    `json:"legal_name" firestore:"legal_name"`
	Birthday             time.Time `json:"birthday" firestore:"birthday"`
	Address1             string    `json:"address1" firestore:"address1"`
	Address2             string    `json:"address2" firestore:"address2"`
	City                 string    `json:"city" firestore:"city"`
	PostalCode           string    `json:"postal_code" firestore:"postal_code"`
	Country              string    `json:"country" firestore:"country"`
	IdentityCard         string    `json:"identity_card" firestore:"identity_card"`
	IdentityCardImgFront string    `json:"identity_card_img_front" firestore:"identity_card_img_front"`
	IdentityCardImgBack  string    `json:"identity_card_img_back" firestore:"identity_card_img_back"`
	SelfieWithIdentity   string    `json:"selfie_with_identity" firestore:"selfie_with_identity"`
	RejectDescription    string    `json:"reject_description" firestore:"reject_description"`
	FixDescription       string    `json:"fix_description" firestore:"fix_description"`
}

func (b VerifyProfileInfo) GetAddUpdateVerifyProfileInfo() map[string]interface{} {
	return map[string]interface{}{
		"uid":                     b.UID,
		"legal_name":              b.LegalName,
		"birthday":                b.Birthday,
		"address1":                b.Address1,
		"address2":                b.Address2,
		"city":                    b.City,
		"postal_code":             b.PostalCode,
		"country":                 b.Country,
		"identity_card":           b.IdentityCard,
		"identity_card_img_front": b.IdentityCardImgFront,
		"identity_card_img_back":  b.IdentityCardImgBack,
		"selfie_with_identity":    b.SelfieWithIdentity,
		"reject_description":      b.RejectDescription,
		"fix_description":         b.FixDescription,
		"updated_at":              firestore.ServerTimestamp,
	}
}

type UserCreditCard struct {
	CCNumber       string `json:"cc_number" firestore:"cc_number"`
	ExpirationDate string `json:"expiration_date" firestore:"expiration_date"`
	Token          string `json:"token" firestore:"token"`
}

type UserBankAccount struct {
	AccountNumber string `json:"account_number" firestore:"account_number"`
	Name          string `json:"name" firestore:"name"`
	BankName      string `json:"bank_name" firestore:"bank_name"`
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

func (user UserBankAccount) GetUpdateProfileBankAccount() map[string]interface{} {
	return map[string]interface{}{
		"bank_account": map[string]interface{}{
			"account_number": user.AccountNumber,
			"name":           user.Name,
			"bank_name":      user.BankName,
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
