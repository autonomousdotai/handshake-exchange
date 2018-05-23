package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const OFFER_TYPE_BUY = "buy"
const OFFER_TYPE_SELL = "sell"

const OFFER_STATUS_CREATED = "created"
const OFFER_STATUS_ACTIVE = "active"
const OFFER_STATUS_WAITING = "waiting"
const OFFER_STATUS_DONE = "done"
const OFFER_STATUS_CLOSED = "closed"

const OFFER_HANDSHANK_STATUS_WAITING = "waiting"
const OFFER_HANDSHANK_STATUS_TIMEOUT = "timeout"
const OFFER_HANDSHANK_STATUS_DONE = "done"

type Offer struct {
	Id              string    `json:"id"`
	MinAmount       string    `json:"min_amount" firestore:"min_amount" validate:"required"`
	MinAmountNumber float64   `json:"-" firestore:"min_amount_number"`
	MaxAmount       string    `json:"max_amount" firestore:"max_amount" validate:"required"`
	MaxAmountNumber float64   `json:"-" firestore:"max_amount_number"`
	Currency        string    `json:"currency" firestore:"currency"`
	PriceNumber     float64   `json:"-" firestore:"price_number"`
	PriceNumberUSD  float64   `json:"-" firestore:"price_number_usd"`
	Price           string    `json:"price" firestore:"price" validate:"required"`
	PriceCurrency   string    `json:"price_currency" firestore:"price_currency"`
	Type            string    `json:"type" firestore:"type" validate:"required"`
	Balance         string    `json:"balance" firestore:"balance"`
	BalanceNumber   float64   `json:"-" firestore:"balance_number"`
	Status          string    `json:"status" firestore:"status"`
	UID             string    `json:"uid" firestore:"uid"`
	Username        string    `json:"username" firestore:"username"`
	FullName        string    `json:"full_name" firestore:"full_name"`
	Country         string    `json:"country" firestore:"country"`
	Avatar          string    `json:"avatar" firestore:"avatar"`
	HasHandshake    bool      `json:"has_handshake" firestore:"has_handshake"`
	ActiveHandshake int64     `json:"-" firestore:"active_handshake"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}

func (offer Offer) ValidateNumbers() (invalid bool) {
	invalid = true
	if _, err := decimal.NewFromString(offer.MinAmount); err != nil {
		return
	}
	if _, err := decimal.NewFromString(offer.MaxAmount); err != nil {
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
		"min_amount":        offer.MinAmount,
		"min_amount_number": offer.MinAmountNumber,
		"max_amount":        offer.MaxAmount,
		"max_amount_number": offer.MaxAmountNumber,
		"currency":          strings.ToUpper(offer.Currency),
		"price_number":      offer.PriceNumber,
		"price_number_usd":  offer.PriceNumberUSD,
		"price":             offer.Price,
		"price_currency":    strings.ToUpper(offer.PriceCurrency),
		"type":              strings.ToLower(offer.Type),
		"balance":           offer.Balance,
		"balance_number":    offer.BalanceNumber,
		"status":            strings.ToLower(offer.Status),
		"uid":               offer.UID,
		"username":          offer.Username,
		"full_name":         offer.FullName,
		"country":           offer.Country,
		"avatar":            offer.Avatar,
		"has_handshake":     false,
		"created_at":        firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateOffer() map[string]interface{} {
	return map[string]interface{}{
		"min_amount":        offer.MinAmount,
		"min_amount_number": offer.MinAmountNumber,
		"max_amount":        offer.MaxAmount,
		"max_amount_number": offer.MaxAmountNumber,
		"price_number":      offer.PriceNumber,
		"price_number_usd":  offer.PriceNumberUSD,
		"price":             offer.Price,
		"balance":           offer.Balance,
		"balance_number":    offer.BalanceNumber,
		"username":          offer.Username,
		"full_name":         offer.FullName,
		"country":           offer.Country,
		"avatar":            offer.Avatar,
		"updated_at":        firestore.ServerTimestamp,
	}
}

func (offer Offer) GetChangeStatus() map[string]interface{} {
	return map[string]interface{}{
		"active_handshake": offer.ActiveHandshake,
		"status":           strings.ToLower(offer.Status),
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (offer Offer) GetUpdateForHandshake() map[string]interface{} {
	return map[string]interface{}{
		"has_handshake":    true,
		"balance":          offer.Balance,
		"balance_number":   offer.BalanceNumber,
		"active_handshake": offer.ActiveHandshake,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (offer Offer) GetPageValue() interface{} {
	return offer.CreatedAt
}

type OfferHandshakeRequest struct {
	Amount string `json:"amount" validate:"required"`
	// Currency string `json:"currency" validate:"required"`
}

type OfferHandshake struct {
	Id                   string    `json:"id" firestore:"id"`
	Amount               string    `json:"amount" firestore:"amount"`
	Currency             string    `json:"currency" firestore:"currency"`
	FiatAmount           string    `json:"fiat_amount" firestore:"fiat_amount"`
	FiatCurrency         string    `json:"fiat_currency" firestore:"fiat_currency"`
	Price                string    `json:"price" firestore:"price"`
	OriginalFiatAmount   string    `json:"-" firestore:"original_fiat_amount"`
	OriginalFiatCurrency string    `json:"-" firestore:"original_fiat_currency"`
	OriginalPrice        string    `json:"-" firestore:"original_price"`
	Rate                 string    `json:"rate" firestore:"rate"`
	Type                 string    `json:"type" firestore:"type"`
	Status               string    `json:"status" firestore:"status"`
	Duration             int64     `json:"duration" firestore:"duration"`
	PayNow               bool      `json:"pay_now" firestore:"pay_now"`
	From                 string    `json:"from" firestore:"from"`
	FromUsername         string    `json:"from_username" firestore:"from_username"`
	FromFullName         string    `json:"from_full_name" firestore:"from_full_name"`
	FromCountry          string    `json:"from_country" firestore:"from_country"`
	FromAvatar           string    `json:"from_avatar" firestore:"from_avatar"`
	FromWalletRef        string    `json:"-" firestore:"from_wallet_ref"`
	FromFiatWalletRef    string    `json:"-" firestore:"from_fiat_wallet_ref"`
	To                   string    `json:"to" firestore:"to"`
	ToUsername           string    `json:"to_username" firestore:"to_username"`
	ToFullName           string    `json:"to_full_name" firestore:"to_full_name"`
	ToCountry            string    `json:"to_country" firestore:"to_country"`
	ToAvatar             string    `json:"to_avatar" firestore:"to_avatar"`
	ToWalletRef          string    `json:"-" firestore:"to_wallet_ref"`
	ToFiatWalletRef      string    `json:"-" firestore:"to_fiat_wallet_ref"`
	Fee                  string    `json:"fee" firestore:"fee"`
	OriginalFee          string    `json:"original_fee" firestore:"original_fee"`
	FeePercentage        string    `json:"fee_percentage" firestore:"fee_percentage"`
	Offer                string    `json:"offer" firestore:"offer"`
	CreatedAt            time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" firestore:"updated_at"`
}

func (handshake OfferHandshake) GetAddOfferHandshake() map[string]interface{} {
	return map[string]interface{}{
		"id":                     handshake.Id,
		"amount":                 handshake.Amount,
		"currency":               handshake.Currency,
		"fiat_amount":            handshake.FiatAmount,
		"fiat_currency":          handshake.FiatCurrency,
		"price":                  handshake.Price,
		"rate":                   handshake.Rate,
		"original_fiat_amount":   handshake.OriginalFiatAmount,
		"original_fiat_currency": handshake.OriginalFiatCurrency,
		"original_price":         handshake.OriginalPrice,
		"type":                   handshake.Type,
		"status":                 OFFER_HANDSHANK_STATUS_WAITING,
		"duration":               handshake.Duration,
		"from":                   handshake.From,
		"from_username":          handshake.FromUsername,
		"from_full_name":         handshake.FromFullName,
		"from_country":           handshake.FromCountry,
		"from_avatar":            handshake.FromAvatar,
		"from_wallet_ref":        handshake.FromWalletRef,
		"from_fiat_wallet_ref":   handshake.FromFiatWalletRef,
		"to":                 handshake.To,
		"to_username":        handshake.ToUsername,
		"to_full_name":       handshake.ToFullName,
		"to_country":         handshake.ToCountry,
		"to_avatar":          handshake.ToAvatar,
		"to_wallet_ref":      handshake.ToWalletRef,
		"to_fiat_wallet_ref": handshake.ToFiatWalletRef,
		"fee":                handshake.Fee,
		"original_fee":       handshake.OriginalFee,
		"fee_percentage":     handshake.FeePercentage,
		"offer":              handshake.Offer,
		"created_at":         firestore.ServerTimestamp,
	}
}

func (handshake OfferHandshake) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     strings.ToLower(handshake.Status),
		"updated_at": firestore.ServerTimestamp,
	}
}

type OfferBalanceHistory struct {
	Previous  string `json:"-" firestore:"previous"`
	New       string `json:"-" firestore:"new"`
	Change    string `json:"-" firestore:"change"`
	Handshake string `json:"-" firestore:"handshake"`
}

func (history OfferBalanceHistory) GetAddBalanceHistory() map[string]interface{} {
	return map[string]interface{}{
		"previous":   history.Previous,
		"new":        history.New,
		"change":     history.Change,
		"handshake":  history.Handshake,
		"created_at": firestore.ServerTimestamp,
	}
}
