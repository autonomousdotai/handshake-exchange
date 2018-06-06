package bean

import (
	"cloud.google.com/go/firestore"
	"strings"
)

const OFFER_STORE_STATUS_CREATED = "created"
const OFFER_STORE_STATUS_ACTIVE = "active"
const OFFER_STORE_STATUS_CLOSING = "closing"
const OFFER_STORE_STATUS_CLOSED = "closed"

type OfferStore struct {
	Id               string           `json:"id" firestore:"id"`
	Hid              int64            `json:"hid" firestore:"hid"`
	ItemFlags        map[string]bool  `json:"item_flags" firestore:"item_flags"`
	Status           string           `json:"status" firestore:"status"`
	UID              string           `json:"-" firestore:"uid"`
	Username         string           `json:"username" firestore:"username"`
	Email            string           `json:"email" firestore:"email"`
	Language         string           `json:"language" firestore:"language"`
	ContactPhone     string           `json:"contact_phone" firestore:"contact_phone"`
	ContactInfo      string           `json:"contact_info" firestore:"contact_info"`
	FCM              string           `json:"-" firestore:"fcm"`
	Longitude        float64          `json:"longitude" firestore:"longitude"`
	Latitude         float64          `json:"latitude" firestore:"latitude"`
	ChainId          int64            `json:"-" firestore:"chain_id"`
	FiatCurrency     string           `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	WalletProvider   string           `json:"-" firestore:"wallet_provider"`
	TransactionCount TransactionCount `json:"transaction_count" firestore:"transaction_count"`
}

type OfferStoreItem struct {
	Currency       string `json:"currency" firestore:"currency"`
	SellAmountMin  string `json:"sell_amount_min" firestore:"sell_amount_min"`
	SellAmount     string `json:"sell_amount" firestore:"sell_amount" validate:"required"`
	SellBalance    string `json:"sell_balance" firestore:"sell_balance"`
	SellPercentage string `json:"sell_percentage" firestore:"sell_percentage"`
	BuyAmountMin   string `json:"buy_amount_min" firestore:"buy_amount_min"`
	BuyAmount      string `json:"buy_amount" firestore:"buy_amount" validate:"required"`
	BuyBalance     string `json:"buy_balance" firestore:"buy_balance"`
	BuyPercentage  string `json:"buy_percentage" firestore:"buy_percentage"`
	SystemAddress  string `json:"system_address" firestore:"system_address"`
	UserAddress    string `json:"user_address" firestore:"user_address"`
	RewardAddress  string `json:"reward_address" firestore:"reward_address"`
}

func (offer OfferStore) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     strings.ToLower(offer.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}

type OfferStoreSetup struct {
	BTC   OfferStoreItem `json:"btc"`
	ETH   OfferStoreItem `json:"eth"`
	Offer OfferStore     `json:"offer"`
}

const OFFER_STORE_SHAKE_STATUS_SHAKING = "shaking"
const OFFER_STORE_SHAKE_STATUS_SHAKE = "shake"
const OFFER_STORE_SHAKE_STATUS_REJECTING = "rejecting"
const OFFER_STORE_SHAKE_STATUS_REJECTED = "rejected"
const OFFER_STORE_SHAKE_STATUS_COMPLETING = "completing"
const OFFER_STORE_SHAKE_STATUS_COMPLETED = "completed"

type OfferStoreShake struct {
	Id               string      `json:"id" firestore:"id"`
	Type             string      `json:"type" firestore:"type" validate:"required"`
	Status           string      `json:"status" firestore:"status"`
	Username         string      `json:"username" firestore:"username"`
	Email            string      `json:"email" firestore:"email"`
	Language         string      `json:"language" firestore:"language"`
	ContactPhone     string      `json:"contact_phone" firestore:"contact_phone"`
	Currency         string      `json:"currency" firestore:"currency"`
	Amount           string      `json:"amount" firestore:"amount" validate:"required"`
	TotalAmount      string      `json:"total_amount" firestore:"total_amount"`
	FiatCurrency     string      `json:"fiat_currency" firestore:"fiat_currency"`
	FiatAmount       string      `json:"fiat_amount" firestore:"fiat_amount"`
	Price            string      `json:"price" firestore:"price"`
	SystemAddress    string      `json:"system_address" firestore:"system_address"`
	UserAddress      string      `json:"user_address" firestore:"user_address"`
	Fee              string      `json:"-" firestore:"fee"`
	FeePercentage    string      `json:"-" firestore:"fee_percentage"`
	Reward           string      `json:"-" firestore:"reward"`
	RewardPercentage string      `json:"-" firestore:"reward_percentage"`
	ActionUID        string      `json:"-" firestore:"action_uid"`
	Provider         string      `json:"-" firestore:"provider"`
	ProviderData     interface{} `json:"-" firestore:"provider_data"`
}

func (offer OfferStoreShake) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     strings.ToLower(offer.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}
