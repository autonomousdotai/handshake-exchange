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

// created -> active
// active -> shaking
// shaking -> shake, pre_shake
// pre_shake -> shake
// shake -> completing
// shake -> rejected
// completing -> completed
const OFFER_STATUS_CREATED = "created"
const OFFER_STATUS_ACTIVE = "active"
const OFFER_STATUS_PRE_SHAKE = "pre_shake"
const OFFER_STATUS_SHAKING = "shaking"
const OFFER_STATUS_SHAKE = "shake"
const OFFER_STATUS_COMPLETING = "completing"
const OFFER_STATUS_COMPLETED = "completed"
const OFFER_STATUS_CLOSED = "closed"
const OFFER_STATUS_REJECTED = "rejected"

type Offer struct {
	Id             string      `json:"id"`
	Amount         string      `json:"amount" firestore:"amount" validate:"required"`
	AmountNumber   float64     `json:"-" firestore:"amount_number"`
	TotalAmount    string      `json:"total_amount" firestore:"total_amount"`
	Currency       string      `json:"currency" firestore:"currency" validate:"required"`
	PriceNumber    float64     `json:"-" firestore:"price_number"`
	PriceNumberUSD float64     `json:"-" firestore:"price_number_usd"`
	Price          string      `json:"price" firestore:"price"`
	PriceUSD       string      `json:"-" firestore:"price_usd"`
	FiatCurrency   string      `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	FiatAmount     string      `json:"fiat_amount" firestore:"fiat_amount"`
	Type           string      `json:"type" firestore:"type" validate:"required"`
	Status         string      `json:"status" firestore:"status"`
	UID            string      `json:"uid" firestore:"uid"`
	Username       string      `json:"-" firestore:"username"`
	ToUID          string      `json:"to_uid" firestore:"to_uid"`
	ToUsername     string      `json:"to_username" firestore:"to_username"`
	ContactInfo    string      `json:"contact_info" firestore:"contact_info" validate:"required"`
	SystemAddress  string      `json:"system_address" firestore:"system_address"`
	UserAddress    string      `json:"user_address" firestore:"user_address"`
	RefundAddress  string      `json:"refund_address" firestore:"refund_address"`
	Provider       string      `json:"provider" firestore:"provider"`
	ProviderData   interface{} `json:"provider_data" firestore:"provider_data"`
	Fee            string      `json:"-" firestore:"fee"`
	FeePercentage  string      `json:"-" firestore:"fee_percentage"`
	CreatedAt      time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" firestore:"updated_at"`
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
	return map[string]interface{}{
		"id":             offer.Id,
		"amount":         offer.Amount,
		"amount_number":  offer.AmountNumber,
		"currency":       strings.ToUpper(offer.Currency),
		"price_currency": strings.ToUpper(offer.FiatCurrency),
		"type":           strings.ToLower(offer.Type),
		"contact_info":   offer.ContactInfo,
		"system_address": offer.SystemAddress,
		"user_address":   offer.UserAddress,
		"refund_address": offer.RefundAddress,
		"status":         offer.Status,
		"uid":            offer.UID,
		"username":       offer.Username,
		"created_at":     firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferActive() map[string]interface{} {
	return map[string]interface{}{
		"price_number":     0,
		"price_number_usd": 0,
		"price":            "0",
		"price_usd":        "0",
		"fiat_amount":      "0",
		"user_address":     offer.UserAddress,
		"status":           OFFER_STATUS_ACTIVE,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferShaking() map[string]interface{} {
	return map[string]interface{}{
		"price_number":     offer.PriceNumber,
		"price_number_usd": offer.PriceNumberUSD,
		"price":            offer.Price,
		"price_usd":        offer.PriceUSD,
		"fiat_amount":      offer.FiatAmount,
		"user_address":     offer.UserAddress,
		"refund_address":   offer.RefundAddress,
		"status":           OFFER_STATUS_SHAKING,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferCompleting() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        OFFER_STATUS_COMPLETING,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferClose() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        OFFER_STATUS_CLOSED,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOfferReject() map[string]interface{} {
	return map[string]interface{}{
		"provider":      offer.Provider,
		"provider_data": offer.ProviderData,
		"status":        OFFER_STATUS_REJECTED,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (offer Offer) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     strings.ToLower(offer.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer Offer) GetPageValue() interface{} {
	return offer.CreatedAt
}

type OfferShakeRequest struct {
	FiatAmount string `json:"fiat_amount" validate:"required"`
	Address    string `json:"address"`
}

type OfferCloseRequest struct {
	Address string `json:"address"`
}
