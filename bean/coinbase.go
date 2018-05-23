package bean

import "time"

type CoinbaseNotificationResponse struct {
	Data CoinbaseNotification `json:"data"`
}

type CoinbaseNotification struct {
	Id             string                 `json:"id" firestore:"id"`
	Type           string                 `json:"type" firestore:"type"`
	Data           map[string]interface{} `json:"data" firestore:"data"`
	User           map[string]interface{} `json:"user" firestore:"user"`
	Account        map[string]interface{} `json:"account" firestore:"account"`
	ResourcePath   string                 `json:"resource_path" firestore:"resource_path"`
	CreatedAt      string                 `json:"created_at" firestore:"created_at"`
	Transaction    map[string]interface{} `json:"transaction" firestore:"transaction"`
	AdditionalData map[string]interface{} `json:"additional_data" firestore:"additional_data"`
}

type CoinbaseCurrency struct {
	Code string `json:"code"`
}

type CoinbaseAccount struct {
	Id       string           `json:"id"`
	Currency CoinbaseCurrency `json:"currency"`
	Type     string           `json:"type"`
}

type CoinbaseAccountResponse struct {
	Data []CoinbaseAccount `json:"data"`
}

type CoinbaseAddress struct {
	Id       string `json:"id"`
	Address  string `json:"address"`
	Resource string `json:"resource"`
}

type CoinbaseAddressResponse struct {
	Data CoinbaseAddress `json:"data"`
}

type CoinbaseSendMoneyRequest struct {
	// Type string
	To          string
	Amount      string
	Currency    string
	Description string
	// Fee string
	Idem string
	// ToFinancialInstitution bool
}

func (request CoinbaseSendMoneyRequest) GetRequestBody() map[string]interface{} {
	return map[string]interface{}{
		"type":        "send",
		"to":          request.To,
		"amount":      request.Amount,
		"currency":    request.Currency,
		"description": request.Description,
		"idem":        request.Idem,
		"to_financial_institution": false,
	}
}

type CoinbaseTransactionResponse struct {
	Data CoinbaseTransaction `json:"data"`
}

type CoinbaseAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CoinbaseNetwork struct {
	Status string `json:"status"`
	Hash   string `json:"hash"`
	Name   string `json:"name"`
}

type CoinbaseTransaction struct {
	Id           string          `json:"id"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	Amount       CoinbaseAmount  `json:"amount"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ResourcePath string          `json:"resource_path"`
	Network      CoinbaseNetwork `json:"network"`
	To           CoinbaseAddress `json:"to"`
	Details      interface{}     `json:"details"`
}

type CoinbasePriceResponse struct {
	Data CoinbaseAmount `json:"data"`
}
