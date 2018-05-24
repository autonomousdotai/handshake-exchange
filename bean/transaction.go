package bean

import (
	"cloud.google.com/go/firestore"
	"time"
)

const TRANSACTION_TYPE_BUY = "buy"
const TRANSACTION_TYPE_SELL = "sell"
const TRANSACTION_TYPE_INSTANT_BUY = "instant_buy"

const TRANSACTION_STATUS_SUCCESS = "success"
const TRANSACTION_STATUS_PENDING = "pending"
const TRANSACTION_STATUS_FAILED = "failed"

type Transaction struct {
	Id              string    `json:"id"`
	Amount          string    `json:"amount" firestore:"amount"`
	TotalAmount     string    `json:"total_amount" firestore:"total_amount"`
	Currency        string    `json:"currency" firestore:"currency"`
	FiatAmount      string    `json:"fiat_amount" firestore:"fiat_amount"`
	TotalFiatAmount string    `json:"total_fiat_amount" firestore:"total_fiat_amount"`
	FiatCurrency    string    `json:"fiat_currency" firestore:"fiat_currency"`
	Price           string    `json:"price" firestore:"price"`
	OriginalPrice   string    `json:"-" firestore:"original_price"`
	Type            string    `json:"type" firestore:"type"`
	Status          string    `json:"status" firestore:"status"`
	From            string    `json:"from" firestore:"from"`
	FromUsername    string    `json:"from_username" firestore:"from_username"`
	To              string    `json:"to" firestore:"to"`
	ToUsername      string    `json:"to_username" firestore:"to_username"`
	Fee             string    `json:"fee" firestore:"fee"`
	FeePercentage   string    `json:"fee_percentage" firestore:"fee_percentage"`
	OfferHandshake  string    `json:"offer_handshake" firestore:"offer_handshake"`
	IsOriginal      bool      `json:"-" firestore:"is_original"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}

func NewTransactionFromOfferHandshake(offer OfferHandshake) (Transaction, Transaction) {
	fromTransaction := Transaction{
		Amount:          offer.Amount,
		TotalAmount:     "amount+-fee",
		Currency:        offer.Currency,
		FiatAmount:      offer.FiatAmount,
		TotalFiatAmount: offer.FiatAmount,
		FiatCurrency:    offer.FiatCurrency,
		Price:           offer.Price,
		OriginalPrice:   offer.OriginalPrice,
		Type:            offer.Type,
		Status:          TRANSACTION_STATUS_SUCCESS,
		From:            offer.From,
		FromUsername:    offer.FromUsername,
		To:              offer.To,
		ToUsername:      offer.ToUsername,
		Fee:             offer.Fee,
		FeePercentage:   offer.FeePercentage,
		OfferHandshake:  offer.Id,
		IsOriginal:      true,
	}

	reverseType := fromTransaction.Type
	if reverseType == TRANSACTION_TYPE_SELL {
		reverseType = TRANSACTION_TYPE_BUY
	}
	toTransaction := Transaction{
		Amount:          offer.Amount,
		TotalAmount:     "amount+-fee",
		Currency:        offer.Currency,
		FiatAmount:      offer.OriginalFiatAmount,
		TotalFiatAmount: offer.FiatAmount,
		FiatCurrency:    offer.OriginalFiatCurrency,
		Price:           offer.OriginalPrice,
		OriginalPrice:   offer.Price,
		Type:            reverseType,
		Status:          TRANSACTION_STATUS_SUCCESS,
		From:            offer.To,
		FromUsername:    offer.ToUsername,
		To:              offer.From,
		ToUsername:      offer.FromUsername,
		Fee:             offer.OriginalFee,
		FeePercentage:   offer.FeePercentage,
		OfferHandshake:  offer.Id,
		IsOriginal:      false,
	}

	return fromTransaction, toTransaction
}

func NewTransactionFromInstantOffer(offer InstantOffer) Transaction {
	txType := TRANSACTION_TYPE_INSTANT_BUY
	fromTransaction := Transaction{
		Amount:          offer.Amount,
		Currency:        offer.Currency,
		FiatAmount:      offer.RawFiatAmount,
		TotalFiatAmount: offer.FiatAmount,
		FiatCurrency:    offer.FiatCurrency,
		Price:           offer.Price,
		Type:            txType,
		Status:          TRANSACTION_STATUS_PENDING,
		From:            offer.UID,
		Fee:             offer.Fee,
		FeePercentage:   offer.FeePercentage,
		OfferHandshake:  offer.Id,
		IsOriginal:      true,
	}

	return fromTransaction
}

func (transaction Transaction) GetAddTransaction() map[string]interface{} {
	return map[string]interface{}{
		"amount":            transaction.Amount,
		"total_amount":      transaction.TotalAmount,
		"currency":          transaction.Currency,
		"fiat_amount":       transaction.FiatAmount,
		"total_fiat_amount": transaction.TotalFiatAmount,
		"fiat_currency":     transaction.FiatCurrency,
		"price":             transaction.Price,
		"original_price":    transaction.OriginalPrice,
		"type":              transaction.Type,
		"status":            transaction.Status,
		"from":              transaction.From,
		"from_username":     transaction.FromUsername,
		"to":                transaction.To,
		"to_username":       transaction.ToUsername,
		"fee":               transaction.Fee,
		"fee_percentage":    transaction.FeePercentage,
		"offer_handshake":   transaction.OfferHandshake,
		"is_original":       transaction.IsOriginal,
		"created_at":        firestore.ServerTimestamp,
	}
}

func (transaction Transaction) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":     transaction.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (transaction Transaction) GetPageValue() interface{} {
	return transaction.CreatedAt
}
