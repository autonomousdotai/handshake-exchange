package gdax_service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/levigross/grequests"
	"os"
	"strconv"
	"time"
)

type GdaxClient struct {
	url           string
	apiPassphrase string
	apiKey        string
	apiSecret     string
}

func (c *GdaxClient) Initialize() {
	c.url = os.Getenv("GDAX_URL")
	c.apiPassphrase = os.Getenv("GDAX_API_PASSPHRASE")
	c.apiKey = os.Getenv("GDAX_API_KEY")
	c.apiSecret = os.Getenv("GDAX_API_SECRET")
}

func (c GdaxClient) buildHeader(method string, url string, body string) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	key, _ := base64.StdEncoding.DecodeString(c.apiSecret)
	message := []byte(timestamp + method + url + body)
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	headers := map[string]string{
		"CB-ACCESS-SIGN":       sign,
		"CB-ACCESS-TIMESTAMP":  timestamp,
		"CB-ACCESS-KEY":        c.apiKey,
		"CB-ACCESS-PASSPHRASE": c.apiPassphrase,
		"Content-Type":         "application/json",
	}

	return headers
}

func (c GdaxClient) Get(uri string) (*grequests.Response, error) {
	url := c.url + uri
	headers := c.buildHeader("GET", uri, "")
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func (c GdaxClient) Post(uri string, body interface{}) (*grequests.Response, error) {
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

	headers := c.buildHeader("POST", uri, bodyStr)
	ro := &grequests.RequestOptions{Headers: headers, RequestBody: r}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func PlaceOrder(amount string, currency string, price string) (bean.GdaxOrderResponse, error) {
	client := GdaxClient{}
	client.Initialize()

	var response bean.GdaxOrderResponse
	resp, err := client.Post("/orders", bean.GdaxPlaceOrderRequest{
		Size:      amount,
		Price:     price,
		Side:      "buy",
		ProductId: fmt.Sprintf("%s-USD", currency),
	}.GetRequestBody())
	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func GetOrder(orderId string) (bean.GdaxOrderResponse, error) {
	client := GdaxClient{}
	client.Initialize()

	var response bean.GdaxOrderResponse

	resp, err := client.Get(fmt.Sprintf("/orders/%s", orderId))
	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}
