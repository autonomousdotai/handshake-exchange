package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const OFFER_TYPE_BUY = "buy"
const OFFER_TYPE_SELL = "sell"

const OFFER_PROVIDER_COINBASE = "coinbase"

const OFFER_STATUS_CREATED = "created"
const OFFER_STATUS_ACTIVE = "active"
const OFFER_STATUS_PRE_SHAKING = "pre_shaking"
const OFFER_STATUS_PRE_SHAKE = "pre_shake"
const OFFER_STATUS_SHAKING = "shaking"
const OFFER_STATUS_SHAKE = "shake"
const OFFER_STATUS_COMPLETING = "completing"
const OFFER_STATUS_COMPLETED = "completed"
const OFFER_STATUS_CLOSING = "closing"
const OFFER_STATUS_CLOSED = "closed"
const OFFER_STATUS_CANCELLING = "cancelling"
const OFFER_STATUS_CANCELLED = "cancelled"
const OFFER_STATUS_REJECTING = "rejecting"
const OFFER_STATUS_REJECTED = "rejected"

var MIN_ETH = decimal.NewFromFloat(0.01).Round(2)

// TODO Change to 0.001 when going to production
var MIN_BTC = decimal.NewFromFloat(0.0001).Round(4)

type Offer struct {
	Id               string           `json:"id"`
	Hid              int64            `json:"hid" firestore:"hid"`
	Amount           string           `json:"amount" firestore:"amount" validate:"required"`
	AmountNumber     float64          `json:"-" firestore:"amount_number"`
	TotalAmount      string           `json:"total_amount" firestore:"total_amount"`
	Currency         string           `json:"currency" firestore:"currency" validate:"required"`
	PriceNumber      float64          `json:"-" firestore:"price_number"`
	PriceNumberUSD   float64          `json:"-" firestore:"price_number_usd"`
	Price            string           `json:"price" firestore:"price" validate:"required"`
	PriceUSD         string           `json:"-" firestore:"price_usd"`
	Percentage       string           `json:"percentage" firestore:"percentage"`
	FiatCurrency     string           `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	FiatAmount       string           `json:"fiat_amount" firestore:"fiat_amount"`
	PhysicalItem     string           `json:"physical_item" firestore:"physical_item"`
	PhysicalQuantity int64            `json:"physical_quantity" firestore:"physical_quantity"`
	PhysicalItemDocs []string         `json:"physical_item_docs" firestore:"physical_item_docs"`
	Tags             []string         `json:"tags" firestore:"tags"`
	Type             string           `json:"type" firestore:"type" validate:"required,oneof=buy sell"`
	Status           string           `json:"status" firestore:"status"`
	UID              string           `json:"uid" firestore:"uid"`
	Username         string           `json:"username" firestore:"username"`
	ChatUsername     string           `json:"chat_username" firestore:"chat_username"`
	Email            string           `json:"email" firestore:"email"`
	Language         string           `json:"language" firestore:"language"`
	FCM              string           `json:"fcm" firestore:"fcm"`
	ToUID            string           `json:"to_uid" firestore:"to_uid"`
	ToUsername       string           `json:"to_username" firestore:"to_username"`
	ToChatUsername   string           `json:"chat_username" firestore:"to_chat_username"`
	ToEmail          string           `json:"to_email" firestore:"to_email"`
	ToLanguage       string           `json:"to_language" firestore:"to_language"`
	ToFCM            string           `json:"to_fcm" firestore:"to_fcm"`
	ContactPhone     string           `json:"contact_phone" firestore:"contact_phone"`
	ContactInfo      string           `json:"contact_info" firestore:"contact_info" validate:"required"`
	SystemAddress    string           `json:"system_address" firestore:"system_address"`
	UserAddress      string           `json:"user_address" firestore:"user_address"`
	RefundAddress    string           `json:"refund_address" firestore:"refund_address"`
	RewardAddress    string           `json:"reward_address" firestore:"reward_address"`
	Provider         string           `json:"provider" firestore:"provider"`
	ProviderData     interface{}      `json:"provider_data" firestore:"provider_data"`
	WalletProvider   string           `json:"wallet_provider" firestore:"wallet_provider"`
	Fee              string           `json:"-" firestore:"fee"`
	FeePercentage    string           `json:"-" firestore:"fee_percentage"`
	Reward           string           `json:"-" firestore:"reward"`
	RewardPercentage string           `json:"-" firestore:"reward_percentage"`
	Longitude        float64          `json:"longitude" firestore:"longitude" validate:"required"`
	Latitude         float64          `json:"latitude" firestore:"latitude" validate:"required"`
	TransactionCount TransactionCount `json:"transaction_count" firestore:"transaction_count"`
	ChainId          int64            `json:"chain_id" firestore:"chain_id"`
	ActionUID        string           `json:"-"`
	CreatedAt        time.Time        `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" firestore:"updated_at"`
}

func (offer Offer) ValidateNumbers() (invalid bool) {
	invalid = true
	if _, err := decimal.NewFromString(offer.Amount); err != nil {
		return
	}
	if _, err := decimal.NewFromString(offer.Price); err != nil {
		return
	}

	invalid = false
	return
}

func (offer Offer) GetAddOffer() map[string]interface{} {
	priceNumber, _ := decimal.NewFromString(offer.Price)
	priceFloat, _ := priceNumber.Float64()
	return map[string]interface{}{
		"id":                 offer.Id,
		"hid":                offer.Hid,
		"amount":             offer.Amount,
		"amount_number":      offer.AmountNumber,
		"total_amount":       offer.TotalAmount,
		"currency":           strings.ToUpper(offer.Currency),
		"price_currency":     strings.ToUpper(offer.FiatCurrency),
		"type":               strings.ToLower(offer.Type),
		"price":              offer.Price,
		"price_number":       priceFloat,
		"fiat_currency":      offer.FiatCurrency,
		"tags":               offer.Tags,
		"physical_item":      offer.PhysicalItem,
		"physical_quantity":  offer.PhysicalQuantity,
		"physical_item_docs": offer.PhysicalItemDocs,
		"percentage":         offer.Percentage,
		"fee_percentage":     offer.FeePercentage,
		"fee":                offer.Fee,
		"reward_percentage":  offer.RewardPercentage,
		"reward":             offer.Reward,
		"contact_info":       offer.ContactInfo,
		"contact_phone":      offer.ContactPhone,
		"email":              offer.Email,
		"language":           offer.Language,
		"fcm":                offer.FCM,
		"system_address":     offer.SystemAddress,
		"user_address":       offer.UserAddress,
		"refund_address":     offer.RefundAddress,
		"reward_address":     offer.RewardAddress,
		"status":             offer.Status,
		"uid":                offer.UID,
		"latitude":           offer.Latitude,
		"longitude":          offer.Longitude,
		"username":           offer.Username,
		"transaction_count":  offer.TransactionCount,
		"chain_id":           offer.ChainId,
		"wallet_provider":    offer.WalletProvider,
		"created_at":         firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferActive() map[string]interface{} {
	return map[string]interface{}{
		// "user_address": offer.UserAddress,
		"hid":        offer.Hid,
		"status":     OFFER_STATUS_ACTIVE,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferShake() map[string]interface{} {
	return map[string]interface{}{
		"price_number":     offer.PriceNumber,
		"price_number_usd": offer.PriceNumberUSD,
		"price":            offer.Price,
		"price_usd":        offer.PriceUSD,
		"fiat_amount":      offer.FiatAmount,
		"to_email":         offer.ToEmail,
		"to_language":      offer.ToLanguage,
		"to_username":      offer.Username,
		"to_chat_username": offer.ToChatUsername,
		"user_address":     offer.UserAddress,
		"refund_address":   offer.RefundAddress,
		"to_uid":           offer.ToUID,
		"status":           offer.Status,
		"updated_at":       firestore.ServerTimestamp,
	}
}

//func (offer Offer) GetUpdateOfferCompleting() map[string]interface{} {
//	return map[string]interface{}{
//		"provider":      offer.Provider,
//		"provider_data": offer.ProviderData,
//		"status":        OFFER_STATUS_COMPLETING,
//		"updated_at":    firestore.ServerTimestamp,
//	}
//}

func (offer Offer) GetUpdateOfferCompleted() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        offer.Status,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferWithdraw() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        offer.Status,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferClose() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        offer.Status,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferReject() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        offer.Status,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"hid":        offer.Hid,
		"status":     strings.ToLower(offer.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer Offer) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": strings.ToLower(offer.Status),
		"type":   "exchange",
	}
}

func (offer Offer) GetPageValue() interface{} {
	return offer.CreatedAt
}

func (offer Offer) IsTypeSell() bool {
	return offer.Type == OFFER_TYPE_SELL
}
func (offer Offer) IsTypeBuy() bool {
	return offer.Type == OFFER_TYPE_BUY
}

type OfferShakeRequest struct {
	FiatAmount   string `json:"fiat_amount"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	ChatUsername string `json:"chat_username"`
	Language     string `json:"language"`
	FCM          string `json:"fcm"`
}

const OFFER_ADDRESS_MAP_OFFER = "offer"
const OFFER_ADDRESS_MAP_OFFER_STORE = "offer_store"
const OFFER_ADDRESS_MAP_OFFER_STORE_SHAKE = "offer_store_shake"

type OfferAddressMap struct {
	UID      string `json:"uid" firestore:"uid"`
	Address  string `json:"address" firestore:"address"`
	Offer    string `json:"offer" firestore:"offer"`
	OfferRef string `json:"offer_ref" firestore:"offer_ref"`
	Type     string `json:"type" firestore:"type"`
}

func (offer OfferAddressMap) GetAddOfferAddressMap() map[string]interface{} {
	return map[string]interface{}{
		"address":    offer.Address,
		"offer":      offer.Offer,
		"uid":        offer.UID,
		"offer_ref":  offer.OfferRef,
		"type":       offer.Type,
		"created_at": firestore.ServerTimestamp,
	}
}

type OfferConfirmingAddressMap struct {
	UID        string `json:"uid" firestore:"uid"`
	Address    string `json:"address" firestore:"address"`
	Offer      string `json:"offer" firestore:"offer"`
	OfferRef   string `json:"offer_ref" firestore:"offer_ref"`
	Type       string `json:"type" firestore:"type"`
	TxHash     string `json:"tx_hash" firestore:"tx_hash"`
	Amount     string `json:"amount" firestore:"amount"`
	Currency   string `json:"currency" firestore:"currency"`
	ExternalId string `json:"external_id" firestore:"external_id"`
}

func (offer OfferConfirmingAddressMap) GetAddOfferConfirmingAddressMap() map[string]interface{} {
	return map[string]interface{}{
		"address":     offer.Address,
		"offer":       offer.Offer,
		"uid":         offer.UID,
		"offer_ref":   offer.OfferRef,
		"external_id": offer.ExternalId,
		"tx_hash":     offer.TxHash,
		"amount":      offer.Amount,
		"currency":    offer.Currency,
		"type":        offer.Type,
		"created_at":  firestore.ServerTimestamp,
	}
}

type OfferTransferMap struct {
	UID        string `json:"uid" firestore:"uid"`
	Address    string `json:"address" firestore:"address"`
	Offer      string `json:"offer" firestore:"offer"`
	OfferRef   string `json:"offer_ref" firestore:"offer_ref"`
	ExternalId string `json:"external_id" firestore:"external_id"`
	Currency   string `json:"currency" firestore:"currency"`
}

func (offer OfferTransferMap) GetAddOfferTransferMap() map[string]interface{} {
	return map[string]interface{}{
		"address":    offer.Address,
		"offer":      offer.Offer,
		"uid":        offer.UID,
		"offer_ref":  offer.OfferRef,
		"currency":   offer.Currency,
		"created_at": firestore.ServerTimestamp,
	}
}

func (offer OfferTransferMap) GetUpdateTick() map[string]interface{} {
	return map[string]interface{}{
		"updated_at": firestore.ServerTimestamp,
	}
}

type OfferOnchain struct {
	Hid   int64
	Offer string
}
