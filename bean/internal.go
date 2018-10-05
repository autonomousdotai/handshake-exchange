package bean

type Redeem struct {
	Address    string `json:"address" firestore:"address" validate:"required"`
	FiatAmount string `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	Amount     string `json:"amount" firestore:"amount"`
	Currency   string `json:"currency" firestore:"currency"`
}
