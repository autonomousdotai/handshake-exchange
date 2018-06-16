package bean

type CheckoutCardPaymentRequest struct {
	Card     CheckoutCard `json:"card"`
	Currency string       `json:"currency"`
	// CustomerId  string       `json:"customerId"`
	Email       string `json:"email"`
	Value       int64  `json:"value"`
	AutoCapture string `json:"autoCapture"`
	Description string `json:"description"`
	Descriptor  string `json:"descriptor"`
}

type CheckoutCard struct {
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	Number      string `json:"number"`
	CVV         string `json:"CVV"`
}

type CheckoutCardPaymentResponse struct {
	Id                   string               `json:"id"`
	Created              string               `json:"created"`
	Value                int64                `json:"value"`
	ResponseMessage      string               `json:"responseMessage"`
	ResponseAdvancedInfo string               `json:"responseAdvancedInfo"`
	ResponseCode         string               `json:"responseCode"`
	Status               string               `json:"status"`
	AuthCode             string               `json:"authCode"`
	Card                 CheckoutCardResponse `json:"card"`
}

type CheckoutCardResponse struct {
	Id          string `json:"id"`
	CustomerId  string `json:"customerId"`
	Last4       string `json:"last4"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	Fingerprint string `json:"fingerprint"`
	CVVCheck    string `json:"cvvCheck"`
	AVSCheck    string `json:"avsCheck"`
}
