package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const CC_TRANSACTION_STATUS_PURCHASED = "purchased"
const CC_TRANSACTION_STATUS_CAPTURED = "captured"
const CC_TRANSACTION_STATUS_REFUNDED = "refunded"
const CC_TRANSACTION_TYPE = "instant_buy"

const CC_PROVIDER_STRIPE = "stripe"

type CCTransaction struct {
	Id           string      `json:"id" firestore:"id"`
	UID          string      `json:"uid" firestore:"uid"`
	Amount       string      `json:"amount" firestore:"amount"`
	Currency     string      `json:"currency" firestore:"currency"`
	Status       string      `json:"-" firestore:"status"`
	Provider     string      `json:"-" firestore:"provider"`
	ProviderData interface{} `json:"-" firestore:"provider_data"`
	ExternalId   string      `json:"-" firestore:"external_id"`
	Type         string      `json:"-" firestore:"type"`
	DataRef      string      `json:"-" firestore:"data_ref"`
	CreatedAt    time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (cc CCTransaction) GetAddCCTransaction() map[string]interface{} {
	return map[string]interface{}{
		"id":            cc.Id,
		"uid":           cc.UID,
		"amount":        cc.Amount,
		"currency":      cc.Currency,
		"status":        cc.Status,
		"provider":      cc.Provider,
		"provider_data": cc.ProviderData,
		"external_id":   cc.ExternalId,
		"type":          cc.Type,
		"data_ref":      cc.DataRef,
		"created_at":    firestore.ServerTimestamp,
	}
}

func (cc CCTransaction) GetUpdateCCTransaction() map[string]interface{} {
	return map[string]interface{}{
		"data_ref":   cc.DataRef,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (cc CCTransaction) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":        cc.Status,
		"provider_data": cc.ProviderData,
		"updated_at":    firestore.ServerTimestamp,
	}
}

func (cc CCTransaction) GetPageValue() interface{} {
	return cc.CreatedAt
}

const INSTANT_OFFER_STATUS_CREATED = "created"
const INSTANT_OFFER_STATUS_PROCESSING = "processing"
const INSTANT_OFFER_STATUS_SUCCESS = "success"
const INSTANT_OFFER_STATUS_CANCELLED = "cancelled"

const INSTANT_OFFER_TYPE_BUY = "buy"

const INSTANT_OFFER_PROVIDER_COINBASE = "coinbase"
const INSTANT_OFFER_PROVIDER_GDAX = "gdax"

const INSTANT_OFFER_PAYMENT_METHOD_CC = "creditcard"

type InstantOffer struct {
	Id                string      `json:"id" firestore:"id"`
	UID               string      `json:"uid" firestore:"uid"`
	Amount            string      `json:"amount" firestore:"amount" validate:"required"`
	Currency          string      `json:"currency" firestore:"currency" validate:"required"`
	FiatAmount        string      `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	TotalFiatAmount   string      `json:"total_fiat_amount" firestore:"total_fiat_amount"`
	FiatCurrency      string      `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	Price             string      `json:"price" firestore:"price"`
	Status            string      `json:"status" firestore:"status"`
	Type              string      `json:"type" firestore:"type"`
	Duration          int         `json:"duration" firestore:"duration"`
	Fee               string      `json:"fee" firestore:"fee"`
	ExternalFee       string      `json:"-" firestore:"external_fee"`
	PaymentMethod     string      `json:"-" firestore:"payment_method"`
	PaymentMethodRef  string      `json:"-" firestore:"payment_method_ref"`
	PaymentMethodData interface{} `json:"payment_method_data" validate:"required"`
	FeePercentage     string      `json:"fee_percentage" firestore:"fee_percentage"`
	Provider          string      `json:"-" firestore:"provider"`
	ProviderData      interface{} `json:"-" firestore:"provider_data"`
	CreatedAt         time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" firestore:"updated_at"`
}

type CreditCardInfo struct {
	CCNum          string `json:"cc_num"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
	Token          string `json:"token"`
	Save           bool   `json:"save"`
}

type InstantOfferRequest struct {
	Currency string `json:"currency" firestore:"currency"`
	Amount   string `json:"amount" firestore:"amount" validate:"required"`
}

func (offer InstantOffer) GetAddInstantOffer() map[string]interface{} {
	return map[string]interface{}{
		"id":                 offer.Id,
		"uid":                offer.UID,
		"amount":             offer.Amount,
		"currency":           offer.Currency,
		"fiat_amount":        offer.FiatAmount,
		"total_fiat_amount":  offer.TotalFiatAmount,
		"fiat_currency":      offer.FiatCurrency,
		"price":              offer.Price,
		"status":             offer.Status,
		"type":               offer.Type,
		"fee":                offer.Fee,
		"external_fee":       offer.ExternalFee,
		"fee_percentage":     offer.FeePercentage,
		"payment_method":     offer.PaymentMethod,
		"payment_method_ref": offer.PaymentMethodRef,
		"provider":           offer.Provider,
		"provider_data":      offer.ProviderData,
		"created_at":         firestore.ServerTimestamp,
	}
}

func (offer InstantOffer) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     offer.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (offer InstantOffer) GetPageValue() interface{} {
	return offer.CreatedAt
}

type PendingInstantOffer struct {
	Id              string    `json:"id" firestore:"id"`
	UID             string    `json:"uid" firestore:"uid"`
	InstantOffer    string    `json:"instant_offer" firestore:"instant_offer"`
	InstantOfferRef string    `json:"instant_offer_ref" firestore:"instant_offer_ref"`
	Provider        string    `json:"provider" firestore:"provider"`
	ProviderId      string    `json:"provider_id" firestore:"provider_id"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
}

func (offer PendingInstantOffer) GetAddInstantOffer() map[string]interface{} {
	return map[string]interface{}{
		"id":                offer.Id,
		"uid":               offer.UID,
		"instant_offer":     offer.InstantOffer,
		"instant_offer_ref": offer.InstantOfferRef,
		"provider":          offer.Provider,
		"provider_id":       offer.ProviderId,
		"created_at":        firestore.ServerTimestamp,
	}
}
