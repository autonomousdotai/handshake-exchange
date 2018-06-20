package checkout_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/shopspring/decimal"
	"os"
)

type CheckoutClient struct {
	url       string
	apiKey    string
	apiSecret string
}

func (c *CheckoutClient) Initialize() {
	c.url = os.Getenv("CHECKOUT_URL")
	c.apiKey = os.Getenv("CHECKOUT_PUBLIC_KEY")
	c.apiSecret = os.Getenv("CHECKOUT_SECRET_KEY")
}

func (c CheckoutClient) buildHeader() map[string]string {
	headers := map[string]string{
		"Authorization": c.apiSecret,
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=UTF-8",
	}

	return headers
}

func (c CheckoutClient) Post(uri string, body interface{}) (*grequests.Response, error) {
	c.Initialize()

	url := c.url + uri

	bodyStr := ""
	if body != nil {
		b, errBody := json.Marshal(&body)
		if errBody != nil {
			return nil, errBody
		}
		bodyStr = string(b)
	}
	r := bytes.NewReader([]byte(bodyStr))

	headers := c.buildHeader()
	ro := &grequests.RequestOptions{Headers: headers, RequestBody: r}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func ChargeCardToken(userId string, cardToken string, amount decimal.Decimal, statement string, description string) (bean.CheckoutCardPaymentResponse, error) {
	client := CheckoutClient{}
	cardPaymentRequest := bean.CheckOutCardIdPaymentRequest{
		CardToken:   cardToken,
		Email:       fmt.Sprintf("user.%s@shake.ninja", userId),
		Currency:    bean.USD.Code,
		Value:       amount.Mul(decimal.NewFromFloat(100).Round(0)).IntPart(),
		AutoCapture: "n",
		Description: description,
		// Descriptor: statement,
	}
	var response bean.CheckoutCardPaymentResponse
	resp, err := client.Post("/v2/charges/token", cardPaymentRequest)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func ChargeCardId(userId string, token string, amount decimal.Decimal, statement string, description string) (bean.CheckoutCardPaymentResponse, error) {
	client := CheckoutClient{}
	cardPaymentRequest := bean.CheckOutCardIdPaymentRequest{
		CardId:      token,
		Email:       fmt.Sprintf("user.%s@shake.ninja", userId),
		Currency:    bean.USD.Code,
		Value:       amount.Mul(decimal.NewFromFloat(100).Round(0)).IntPart(),
		AutoCapture: "n",
		Description: description,
		// Descriptor: statement,
	}
	var response bean.CheckoutCardPaymentResponse
	resp, err := client.Post("/v2/charges/card", cardPaymentRequest)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func ChargeCard(userId string, cardNum string, date string, cvv string, amount decimal.Decimal, statement string, description string) (bean.CheckoutCardPaymentResponse, error) {
	month := date[:2]
	year := "20" + date[3:]

	client := CheckoutClient{}

	cardPaymentRequest := bean.CheckoutCardPaymentRequest{
		Card: bean.CheckoutCard{
			ExpiryMonth: month,
			ExpiryYear:  year,
			Number:      cardNum,
			CVV:         cvv,
		},
		Email:       fmt.Sprintf("user.%s@shake.ninja", userId),
		Currency:    bean.USD.Code,
		Value:       amount.Mul(decimal.NewFromFloat(100).Round(0)).IntPart(),
		AutoCapture: "n",
		Description: description,
		// Descriptor: statement,
	}
	var response bean.CheckoutCardPaymentResponse
	resp, err := client.Post("/v2/charges/card", cardPaymentRequest)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func Capture(chargeId string) (bean.CheckoutCard2ndStepResponse, error) {
	client := CheckoutClient{}
	var response bean.CheckoutCard2ndStepResponse
	resp, err := client.Post(fmt.Sprintf("/v2/charges/%s/capture", chargeId), nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func Void(chargeId string) (bean.CheckoutCard2ndStepResponse, error) {
	client := CheckoutClient{}
	var response bean.CheckoutCard2ndStepResponse
	resp, err := client.Post(fmt.Sprintf("/v2/charges/%s/void", chargeId), nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func Refund(chargeId string) (bean.CheckoutCard2ndStepResponse, error) {
	client := CheckoutClient{}
	var response bean.CheckoutCard2ndStepResponse
	resp, err := client.Post(fmt.Sprintf("/v2/charges/%s/refund", chargeId), nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}
