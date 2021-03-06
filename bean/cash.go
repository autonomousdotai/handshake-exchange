package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/ninjadotorg/handshake-exchange/common"
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
	Profit       string            `json:"profit" firestore:"profit"`
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
		"profit":        common.Zero.String(),
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

func (b CashStore) GetUpdateProfit() map[string]interface{} {
	return map[string]interface{}{
		"profit":     b.Profit,
		"updated_at": firestore.ServerTimestamp,
	}
}

const CASH_ORDER_STATUS_PROCESSING = "processing"
const CASH_ORDER_STATUS_SUCCESS = "success"
const CASH_ORDER_STATUS_TRANSFERRING = "transferring"
const CASH_ORDER_STATUS_TRANSFER_FAILED = "transfer_failed"
const CASH_ORDER_STATUS_FIAT_TRANSFERRING = "fiat_transferring"
const CASH_ORDER_STATUS_CANCELLED = "cancelled"

type CashOrderUpdateInput struct {
	ReceiptURL string `json:"receipt_url" firestore:"receipt_url"`
}

type CashOrder struct {
	Id                        string            `json:"id" firestore:"id"`
	UID                       string            `json:"-" firestore:"uid"`
	ToUID                     string            `json:"-" firestore:"to_uid"`
	UserInfo                  map[string]string `json:"user_info" firestore:"user_info"`
	Amount                    string            `json:"amount" firestore:"amount" validate:"required"`
	Currency                  string            `json:"currency" firestore:"currency" validate:"required"`
	FiatAmount                string            `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	FiatCurrency              string            `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	FiatLocalAmount           string            `json:"fiat_local_amount" firestore:"fiat_local_amount"`
	FiatLocalCurrency         string            `json:"fiat_local_currency" firestore:"fiat_local_currency"`
	LocalStoreFee             string            `json:"local_store_fee" firestore:"local_store_fee"`
	RawFiatAmount             string            `json:"-" firestore:"raw_fiat_amount"`
	Price                     string            `json:"price" firestore:"price"`
	Status                    string            `json:"status" firestore:"status"`
	Type                      string            `json:"type" firestore:"type"`
	Duration                  int64             `json:"-" firestore:"duration"`
	FeePercentage             string            `json:"-" firestore:"fee_percentage"`
	Fee                       string            `json:"-" firestore:"fee"`
	StoreFeePercentage        string            `json:"-" firestore:"store_fee_percentage"`
	StoreFee                  string            `json:"store_fee" firestore:"store_fee"`
	ExternalFeePercentage     string            `json:"-" firestore:"external_fee_percentage"`
	ExternalFee               string            `json:"-" firestore:"external_fee"`
	PaymentMethod             string            `json:"-" firestore:"payment_method"`
	PaymentMethodRef          string            `json:"-" firestore:"payment_method_ref"`
	PaymentMethodData         interface{}       `json:"payment_method_data"`
	Provider                  string            `json:"-" firestore:"provider"`
	ProviderData              interface{}       `json:"-" firestore:"provider_data"`
	Center                    string            `json:"center" firestore:"center"`
	Address                   string            `json:"address" firestore:"address" validate:"required"`
	ProviderWithdrawData      interface{}       `json:"provider_withdraw_data" firestore:"provider_withdraw_data"`
	ProviderWithdrawDataExtra interface{}       `json:"-" firestore:"provider_withdraw_data_extra"`
	ReceiptURL                string            `json:"receipt_url" firestore:"receipt_url"`
	RefCode                   string            `json:"ref_code" firestore:"ref_code"`
	FCM                       string            `json:"fcm" firestore:"fcm"`
	Language                  string            `json:"language" firestore:"language"`
	ChainId                   int64             `json:"chain_id" firestore:"chain_id"`
	CreatedAt                 time.Time         `json:"created_at" firestore:"created_at"`
	UpdatedAt                 time.Time         `json:"updated_at" firestore:"updated_at"`
}

func (b CashOrder) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":                      b.Id,
		"uid":                     b.UID,
		"user_info":               b.UserInfo,
		"amount":                  b.Amount,
		"currency":                b.Currency,
		"fiat_amount":             b.FiatAmount,
		"raw_fiat_amount":         b.RawFiatAmount,
		"fiat_currency":           b.FiatCurrency,
		"fiat_local_amount":       b.FiatLocalAmount,
		"fiat_local_currency":     b.FiatLocalCurrency,
		"local_store_fee":         b.LocalStoreFee,
		"price":                   b.Price,
		"status":                  b.Status,
		"type":                    b.Type,
		"fee":                     b.Fee,
		"fee_percentage":          b.FeePercentage,
		"store_fee":               b.StoreFee,
		"store_fee_percentage":    b.StoreFeePercentage,
		"external_fee":            b.ExternalFee,
		"external_fee_percentage": b.ExternalFeePercentage,
		"duration":                b.Duration,
		"payment_method":          b.PaymentMethod,
		"payment_method_ref":      b.PaymentMethodRef,
		"center":                  b.Center,
		"address":                 b.Address,
		"provider":                b.Provider,
		"provider_data":           b.ProviderData,
		"ref_code":                b.RefCode,
		"language":                b.Language,
		"fcm":                     b.FCM,
		"chain_id":                b.ChainId,
		"created_at":              firestore.ServerTimestamp,
	}
}

func (b CashOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"provider_withdraw_data":       b.ProviderWithdrawData,
		"provider_withdraw_data_extra": b.ProviderWithdrawDataExtra,
		"status":                       b.Status,
		"updated_at":                   firestore.ServerTimestamp,
	}
}

func (b CashOrder) GetReceiptUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":      b.Status,
		"receipt_url": b.ReceiptURL,
		"updated_at":  firestore.ServerTimestamp,
	}
}

func (b CashOrder) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     b.Id,
		"status": b.Status,
		"type":   "cash_order",
	}
}

func (b CashOrder) GetPageValue() interface{} {
	return b.CreatedAt
}

type CashCenter struct {
	Id          string                 `json:"id" firestore:"id"`
	Country     string                 `json:"country" firestore:"country"`
	Information map[string]interface{} `json:"information" firestore:"information"`
}

const CASH_STORE_PAYMENT_STATUS_MATCHED = "matched"
const CASH_STORE_PAYMENT_STATUS_UNDER = "under"
const CASH_STORE_PAYMENT_STATUS_OVER = "over"

type CashStorePayment struct {
	Order        string    `json:"order" firestore:"order"`
	FiatAmount   string    `json:"fiat_amount" firestore:"fiat_amount"`
	FiatCurrency string    `json:"fiat_currency" firestore:"fiat_currency"`
	OverSpent    string    `json:"over_spent" firestore:"over_spent"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}

func (b CashStorePayment) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"order":         b.Order,
		"fiat_amount":   b.FiatAmount,
		"fiat_currency": b.FiatCurrency,
		"over_spent":    b.OverSpent,
		"status":        b.Status,
		"created_at":    firestore.ServerTimestamp,
	}
}

func (b CashStorePayment) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"fiat_amount": b.FiatAmount,
		"over_spent":  b.OverSpent,
		"status":      b.Status,
		"updated_at":  firestore.ServerTimestamp,
	}
}

type CashOrderRefCode struct {
	RefCode  string `json:"ref_code" firestore:"ref_code"`
	OrderRef string `json:"order_ref" firestore:"order_ref"`
}

func (b CashOrderRefCode) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"ref_code":   b.RefCode,
		"order_ref":  b.OrderRef,
		"created_at": firestore.ServerTimestamp,
	}
}
