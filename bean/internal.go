package bean

import "time"

const REDEEM_STATUS_SUCCESS = "success"
const REDEEM_STATUS_FAILED = "failed"

type Redeem struct {
	Id           string    `json:"id" firestore:"id"`
	Address      string    `json:"address" firestore:"address" validate:"required"`
	FiatAmount   string    `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	Amount       string    `json:"amount" firestore:"amount"`
	Currency     string    `json:"currency" firestore:"currency" validate:"required"`
	RefData      string    `json:"ref_data" firestore:"ref_data"`
	Provider     string    `json:"provider" firestore:"provider"`
	ProviderData string    `json:"provider_data" firestore:"provider_data"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}

func (b Redeem) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"address":     b.Address,
		"amount":      b.Amount,
		"fiat_amount": b.FiatAmount,
		"currency":    b.Currency,
		"ref_data":    b.RefData,
		"provider":    b.Provider,
	}
}

type RedeemLimit struct {
	Currency  string    `json:"currency" firestore:"currency"`
	Limit     string    `json:"limit" firestore:"limit"`
	Usage     string    `json:"usage" firestore:"usage"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}
