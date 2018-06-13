package bean

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"strings"
	"time"
)

const OFFER_STORE_STATUS_CREATED = "created"
const OFFER_STORE_STATUS_ACTIVE = "active"
const OFFER_STORE_STATUS_CLOSING = "closing"
const OFFER_STORE_STATUS_CLOSED = "closed"

const OFFER_STORE_ITEM_STATUS_CREATED = "created"
const OFFER_STORE_ITEM_STATUS_ACTIVE = "active"
const OFFER_STORE_ITEM_STATUS_CLOSING = "closing"
const OFFER_STORE_ITEM_STATUS_CLOSED = "closed"

type OfferStore struct {
	Id               string                    `json:"id" firestore:"id"`
	Hid              int64                     `json:"hid" firestore:"hid"`
	ItemFlags        map[string]bool           `json:"item_flags" firestore:"item_flags"`
	Status           string                    `json:"status" firestore:"status"`
	UID              string                    `json:"-" firestore:"uid"`
	Username         string                    `json:"username" firestore:"username"`
	ChatUsername     string                    `json:"chat_username" firestore:"chat_username"`
	Email            string                    `json:"email" firestore:"email"`
	Language         string                    `json:"language" firestore:"language"`
	ContactPhone     string                    `json:"contact_phone" firestore:"contact_phone"`
	ContactInfo      string                    `json:"contact_info" firestore:"contact_info"`
	FCM              string                    `json:"-" firestore:"fcm"`
	Longitude        float64                   `json:"longitude" firestore:"longitude"`
	Latitude         float64                   `json:"latitude" firestore:"latitude"`
	ChainId          int64                     `json:"-" firestore:"chain_id"`
	FiatCurrency     string                    `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	TransactionCount TransactionCount          `json:"transaction_count" firestore:"transaction_count"`
	ItemSnapshots    map[string]OfferStoreItem `json:"items" firestore:"item_snapshots"`
	Offline          string                    `json:"-"`
	CreatedAt        time.Time                 `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at" firestore:"updated_at"`
}

type OfferStoreItem struct {
	Currency       string `json:"currency" firestore:"currency"`
	Status         string `json:"status" firestore:"status"`
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
	WalletProvider string `json:"-" firestore:"wallet_provider"`
	RewardAddress  string `json:"reward_address" firestore:"reward_address"`
	Hid            string `json:"hid" firestore:"hid"`
	TxHash         string `json:"tx_hash" firestore:"tx_hash"`
}

func (offer OfferStore) GetAddOfferStore() map[string]interface{} {
	return map[string]interface{}{
		"id":                offer.Id,
		"item_flags":        offer.ItemFlags,
		"status":            OFFER_STORE_STATUS_CREATED,
		"uid":               offer.UID,
		"username":          offer.Username,
		"chat_username":     offer.ChatUsername,
		"email":             offer.Email,
		"language":          offer.Language,
		"contact_phone":     offer.ContactPhone,
		"contact_info":      offer.ContactInfo,
		"fcm":               offer.FCM,
		"latitude":          offer.Latitude,
		"longitude":         offer.Longitude,
		"chain_id":          offer.ChainId,
		"fiat_currency":     offer.FiatCurrency,
		"transaction_count": offer.TransactionCount,
		"item_snapshots":    offer.ItemSnapshots,
		"created_at":        firestore.ServerTimestamp,
	}
}

func (offer OfferStore) GetUpdateOfferStoreChangeItem() map[string]interface{} {
	return map[string]interface{}{
		"status":         offer.Status,
		"item_flags":     offer.ItemFlags,
		"item_snapshots": offer.ItemSnapshots,
		"updated_at":     firestore.ServerTimestamp,
	}
}

func (offer OfferStore) GetUpdateOfferStoreChangeSnapshot() map[string]interface{} {
	return map[string]interface{}{
		"item_snapshots": offer.ItemSnapshots,
		"updated_at":     firestore.ServerTimestamp,
	}
}

func (offer OfferStore) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":         strings.ToLower(offer.Status),
		"item_snapshots": offer.ItemSnapshots,
		"updated_at":     firestore.ServerTimestamp,
	}
}

func (offer OfferStore) GetUpdateOfferStoreActive() map[string]interface{} {
	return map[string]interface{}{
		"hid":        offer.Hid,
		"status":     offer.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer OfferStore) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": offer.Status,
		"type":   "offer_store",
	}
}

func (item OfferStoreItem) GetAddOfferStoreItem() map[string]interface{} {
	return map[string]interface{}{
		"currency":        item.Currency,
		"status":          item.Status,
		"sell_amount_min": item.SellAmountMin,
		"sell_amount":     item.SellAmount,
		"sell_balance":    "0",
		"sell_percentage": item.SellPercentage,
		"buy_amount_min":  item.BuyAmountMin,
		"buy_amount":      item.BuyAmount,
		"buy_balance":     item.BuyAmount,
		"buy_percentage":  item.BuyPercentage,
		"system_address":  item.SystemAddress,
		"user_address":    item.UserAddress,
		"reward_address":  item.RewardAddress,
		"wallet_provider": item.WalletProvider,
	}
}

func (item OfferStoreItem) GetUpdateOfferStoreItemActive() map[string]interface{} {
	return map[string]interface{}{
		"sell_balance": item.SellBalance,
		"status":       OFFER_STORE_ITEM_STATUS_ACTIVE,
		"updated_at":   firestore.ServerTimestamp,
	}
}

func (item OfferStoreItem) GetUpdateOfferStoreItemClosing() map[string]interface{} {
	return map[string]interface{}{
		"status":     OFFER_STORE_ITEM_STATUS_CLOSING,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (item OfferStoreItem) GetUpdateOfferStoreItemClosed() map[string]interface{} {
	return map[string]interface{}{
		"status":     OFFER_STORE_ITEM_STATUS_CLOSED,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (item OfferStoreItem) GetUpdateOfferStoreItemBalance() map[string]interface{} {
	return map[string]interface{}{
		"buy_balance":  item.BuyBalance,
		"sell_balance": item.SellBalance,
		"updated_at":   firestore.ServerTimestamp,
	}
}

func (item OfferStoreItem) GetNotificationUpdate(offer OfferStore) map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": fmt.Sprintf("%s_%s", strings.ToLower(item.Currency), item.Status),
		"type":   "offer_store",
	}
}

type OfferStoreSetup struct {
	Item  OfferStoreItem `json:"item"`
	Offer OfferStore     `json:"offer"`
}

const OFFER_STORE_SHAKE_STATUS_PRE_SHAKING = "pre_shaking"
const OFFER_STORE_SHAKE_STATUS_PRE_SHAKE = "pre_shake"
const OFFER_STORE_SHAKE_STATUS_CANCELLING = "cancelling"
const OFFER_STORE_SHAKE_STATUS_CANCELLED = "cancelled"
const OFFER_STORE_SHAKE_STATUS_SHAKING = "shaking"
const OFFER_STORE_SHAKE_STATUS_SHAKE = "shake"
const OFFER_STORE_SHAKE_STATUS_REJECTING = "rejecting"
const OFFER_STORE_SHAKE_STATUS_REJECTED = "rejected"
const OFFER_STORE_SHAKE_STATUS_COMPLETING = "completing"
const OFFER_STORE_SHAKE_STATUS_COMPLETED = "completed"

type OfferStoreShake struct {
	Id               string      `json:"id" firestore:"id"`
	Hid              int64       `json:"hid" firestore:"hid"`
	OffChainId       string      `json:"off_chain_id" firestore:"off_chain_id"`
	Type             string      `json:"type" firestore:"type" validate:"required"`
	Status           string      `json:"status" firestore:"status"`
	UID              string      `json:"-" firestore:"uid"`
	Username         string      `json:"username" firestore:"username"`
	ChatUsername     string      `json:"chat_username" firestore:"chat_username"`
	Email            string      `json:"email" firestore:"email"`
	Language         string      `json:"language" firestore:"language"`
	FCM              string      `json:"fcm" firestore:"fcm"`
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
	WalletProvider   string      `json:"-" firestore:"wallet_provider"`
	Provider         string      `json:"-" firestore:"provider"`
	ProviderData     interface{} `json:"-" firestore:"provider_data"`
	ChainId          int64       `json:"-" firestore:"chain_id"`
	Longitude        float64     `json:"longitude" firestore:"longitude"`
	Latitude         float64     `json:"latitude" firestore:"latitude"`
	CreatedAt        time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (offer OfferStoreShake) GetAddOfferStoreShake() map[string]interface{} {
	return map[string]interface{}{
		"id":                offer.Id,
		"off_chain_id":      offer.OffChainId,
		"type":              offer.Type,
		"status":            offer.Status,
		"uid":               offer.UID,
		"username":          offer.Username,
		"chat_username":     offer.ChatUsername,
		"email":             offer.Email,
		"language":          offer.Language,
		"contact_phone":     offer.ContactPhone,
		"currency":          offer.Currency,
		"amount":            offer.Amount,
		"total_amount":      offer.TotalAmount,
		"fiat_currency":     offer.FiatCurrency,
		"fiat_amount":       offer.FiatAmount,
		"price":             offer.Price,
		"wallet_provider":   offer.WalletProvider,
		"system_address":    offer.SystemAddress,
		"user_address":      offer.UserAddress,
		"fee":               offer.Fee,
		"fee_percentage":    offer.FeePercentage,
		"reward":            offer.Reward,
		"reward_percentage": offer.RewardPercentage,
		"action_uid":        offer.ActionUID,
		"chain_id":          offer.ChainId,
		"provider":          offer.Provider,
		"provider_data":     offer.ProviderData,
		"latitude":          offer.Latitude,
		"longitude":         offer.Longitude,
		"created_at":        firestore.ServerTimestamp,
	}
}

func (offer OfferStoreShake) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"hid":        offer.Hid,
		"status":     strings.ToLower(offer.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer OfferStoreShake) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": offer.Status,
		"type":   "offer_store_shake",
	}
}
