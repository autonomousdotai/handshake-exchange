package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const COIN_ORDER_TYPE_COD = "cod"
const COIN_ORDER_TYPE_BANK = "bank"

const COIN_ORDER_STATUS_PENDING = "pending"
const COIN_ORDER_STATUS_PROCESSING = "processing"
const COIN_ORDER_STATUS_FIAT_TRANSFERRING = "fiat_transferring"
const COIN_ORDER_STATUS_TRANSFERRING = "transferring"
const COIN_ORDER_STATUS_SUCCESS = "success"
const COIN_ORDER_STATUS_TRANSFER_FAILED = "transfer_failed"
const COIN_ORDER_STATUS_CANCELLED = "cancelled"

type CoinOrderUpdateInput struct {
	ReceiptURL string `json:"receipt_url" firestore:"receipt_url"`
}

type CoinOrder struct {
	Id                        string            `json:"id" firestore:"id"`
	UID                       string            `json:"-" firestore:"uid"`
	UserInfo                  map[string]string `json:"user_info" firestore:"user_info"`
	Amount                    string            `json:"amount" firestore:"amount" validate:"required"`
	Currency                  string            `json:"currency" firestore:"currency" validate:"required"`
	FiatAmount                string            `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	FiatCurrency              string            `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	FiatLocalAmount           string            `json:"fiat_local_amount" firestore:"fiat_local_amount"`
	FiatLocalCurrency         string            `json:"fiat_local_currency" firestore:"fiat_local_currency"`
	RawFiatAmount             string            `json:"-" firestore:"raw_fiat_amount"`
	Price                     string            `json:"price" firestore:"price"`
	Status                    string            `json:"status" firestore:"status"`
	Type                      string            `json:"type" firestore:"type"`
	Duration                  int64             `json:"-" firestore:"duration"`
	FeePercentage             string            `json:"-" firestore:"fee_percentage"`
	Fee                       string            `json:"-" firestore:"fee"`
	ExternalFeePercentage     string            `json:"-" firestore:"external_fee_percentage"`
	ExternalFee               string            `json:"-" firestore:"external_fee"`
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

func (b CoinOrder) GetAdd() map[string]interface{} {
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
		"price":                   b.Price,
		"status":                  b.Status,
		"type":                    b.Type,
		"fee":                     b.Fee,
		"fee_percentage":          b.FeePercentage,
		"external_fee":            b.ExternalFee,
		"external_fee_percentage": b.ExternalFeePercentage,
		"duration":                b.Duration,
		"address":                 b.Address,
		"ref_code":                b.RefCode,
		"language":                b.Language,
		"fcm":                     b.FCM,
		"chain_id":                b.ChainId,
		"created_at":              firestore.ServerTimestamp,
	}
}

func (b CoinOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CoinOrder) GetReceiptUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":      b.Status,
		"receipt_url": b.ReceiptURL,
		"updated_at":  firestore.ServerTimestamp,
	}
}

func (b CoinOrder) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     b.Id,
		"status": b.Status,
		"type":   "coin_order",
	}
}

func (b CoinOrder) GetPageValue() interface{} {
	return b.CreatedAt
}

type CoinCenter struct {
	Id          string                 `json:"id" firestore:"id"`
	Country     string                 `json:"country" firestore:"country"`
	Information map[string]interface{} `json:"information" firestore:"information"`
}

const COIN_PAYMENT_STATUS_MATCHED = "matched"
const COIN_PAYMENT_STATUS_UNDER = "under"
const COIN_PAYMENT_STATUS_OVER = "over"

type CoinPayment struct {
	Order        string    `json:"order" firestore:"order"`
	FiatAmount   string    `json:"fiat_amount" firestore:"fiat_amount"`
	FiatCurrency string    `json:"fiat_currency" firestore:"fiat_currency"`
	OverSpent    string    `json:"over_spent" firestore:"over_spent"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}

func (b CoinPayment) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"order":         b.Order,
		"fiat_amount":   b.FiatAmount,
		"fiat_currency": b.FiatCurrency,
		"over_spent":    b.OverSpent,
		"status":        b.Status,
		"created_at":    firestore.ServerTimestamp,
	}
}

func (b CoinPayment) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"fiat_amount": b.FiatAmount,
		"over_spent":  b.OverSpent,
		"status":      b.Status,
		"updated_at":  firestore.ServerTimestamp,
	}
}

type CoinOrderRefCode struct {
	RefCode  string `json:"ref_code" firestore:"ref_code"`
	OrderRef string `json:"order_ref" firestore:"order_ref"`
}

func (b CoinOrderRefCode) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"ref_code":   b.RefCode,
		"order_ref":  b.OrderRef,
		"created_at": firestore.ServerTimestamp,
	}
}

type CoinQuote struct {
	FiatAmount         string `json:"fiat_amount"`
	FiatCurrency       string `json:"fiat_currency"`
	FiatLocalAmount    string `json:"fiat_local_amount"`
	FiatLocalCurrency  string `json:"fiat_local_currency"`
	FiatAmountCOD      string `json:"fiat_amount_cod"`
	FiatLocalAmountCOD string `json:"fiat_local_amount_cod"`
	Fee                string `json:"-"`
	FeeLocal           string `json:"-"`
	FeePercentage      string `json:"-"`
	FeeCOD             string `json:"-"`
	FeeLocalCOD        string `json:"-"`
	FeePercentageCOD   string `json:"-"`
	RawFiatAmount      string `json:"-"`
	Price              string `json:"-"`
	Limit              string `json:"limit"`
}
