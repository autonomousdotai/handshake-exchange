package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const CASH_STORE_BUSINESS_TYPE_PERSONAL = "personal"
const CASH_STORE_BUSINESS_TYPE_STORE = "store"

const CASH_STORE_STATUS_OPEN = "open"
const CASH_STORE_STATUS_CLOSE = "close"

type CashStore struct {
	UID          string            `json:"uid" firestore:"uid"`
	Name         string            `json:"name" firestore:"name"`
	Address      string            `json:"address" firestore:"address"`
	Phone        string            `json:"phone" firestore:"phone"`
	BusinessType string            `json:"business_type" firestore:"business_type"`
	Status       string            `json:"status" firestore:"status"`
	Center       string            `json:"center" firestore:"center"`
	Information  map[string]string `json:"information" firestore:"information"`
	Longitude    float64           `json:"longitude" firestore:"longitude"`
	Latitude     float64           `json:"latitude" firestore:"latitude"`
	ChainId      int64             `json:"chain_id" firestore:"chain_id"`
	Language     string            `json:"-" firestore:"language"`
	FCM          string            `json:"-" firestore:"fcm"`
	CreatedAt    time.Time         `json:"created_at" firestore:"created_at"`
}

func (b CashStore) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":           b.UID,
		"name":          b.Name,
		"address":       b.Address,
		"phone":         b.Phone,
		"business_type": b.BusinessType,
		"status":        b.Status,
		"center":        b.Center,
		"information":   b.Information,
		"longitude":     b.Longitude,
		"latitude":      b.Latitude,
		"chain_id":      b.ChainId,
		"created_at":    firestore.ServerTimestamp,
	}
}

func (b CashStore) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"name":          b.Name,
		"address":       b.Address,
		"phone":         b.Phone,
		"business_type": b.BusinessType,
		"status":        b.Status,
		"information":   b.Information,
		"center":        b.Center,
		"longitude":     b.Longitude,
		"latitude":      b.Latitude,
		"updated_at":    firestore.ServerTimestamp,
	}
}

type CashStoreOrder struct {
	Id                    string      `json:"id" firestore:"id"`
	UID                   string      `json:"uid" firestore:"uid"`
	Username              string      `json:"username" firestore:"username"`
	Amount                string      `json:"amount" firestore:"amount" validate:"required"`
	Currency              string      `json:"currency" firestore:"currency" validate:"required"`
	FiatAmount            string      `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	FiatCurrency          string      `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	FiatLocalAmount       string      `json:"fiat_local_amount" firestore:"fiat_local_amount" validate:"required"`
	FiatLocalCurrency     string      `json:"fiat_local_currency" firestore:"fiat_local_currency" validate:"required"`
	RawFiatAmount         string      `json:"-" firestore:"raw_fiat_amount"`
	Price                 string      `json:"price" firestore:"price"`
	Status                string      `json:"status" firestore:"status"`
	Type                  string      `json:"type" firestore:"type"`
	Duration              int64       `json:"-" firestore:"duration"`
	FeePercentage         string      `json:"-" firestore:"fee_percentage"`
	Fee                   string      `json:"-" firestore:"fee"`
	StoreFeePercentage    string      `json:"-" firestore:"store_fee_percentage"`
	StoreFee              string      `json:"-" firestore:"store_fee"`
	ExternalFeePercentage string      `json:"-" firestore:"external_fee_percentage"`
	ExternalFee           string      `json:"-" firestore:"external_fee"`
	PaymentMethod         string      `json:"-" firestore:"payment_method"`
	PaymentMethodRef      string      `json:"-" firestore:"payment_method_ref"`
	PaymentMethodData     interface{} `json:"payment_method_data" validate:"required"`
	Center                string      `json:"center" firestore:"center"`
	ProviderWithdrawData  interface{} `json:"-" firestore:"provider_withdraw_data"`
	FCM                   string      `json:"fcm" firestore:"fcm"`
	Language              string      `json:"language" firestore:"language"`
	ChainId               int64       `json:"chain_id" firestore:"chain_id"`
	CreatedAt             time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (offer CashStoreOrder) GetAddInstantOffer() map[string]interface{} {
	return map[string]interface{}{
		"id":                      offer.Id,
		"uid":                     offer.UID,
		"amount":                  offer.Amount,
		"currency":                offer.Currency,
		"fiat_amount":             offer.FiatAmount,
		"raw_fiat_amount":         offer.RawFiatAmount,
		"fiat_currency":           offer.FiatCurrency,
		"price":                   offer.Price,
		"status":                  offer.Status,
		"type":                    offer.Type,
		"fee":                     offer.Fee,
		"external_fee":            offer.ExternalFee,
		"fee_percentage":          offer.FeePercentage,
		"external_fee_percentage": offer.ExternalFeePercentage,
		"duration":                offer.Duration,
		"payment_method":          offer.PaymentMethod,
		"payment_method_ref":      offer.PaymentMethodRef,
		"language":                offer.Language,
		"fcm":                     offer.FCM,
		"chain_id":                offer.ChainId,
		"created_at":              firestore.ServerTimestamp,
	}
}

func (offer CashStoreOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"provider_withdraw_data": offer.ProviderWithdrawData,
		"status":                 offer.Status,
		"updated_at":             firestore.ServerTimestamp,
	}
}

func (offer CashStoreOrder) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": offer.Status,
		"type":   "instant",
	}
}
