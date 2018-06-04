package coinbase_service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/levigross/grequests"
	"os"
	"strconv"
	"strings"
	"time"
)

type CoinbaseClient struct {
	url        string
	apiVersion string
	apiKey     string
	apiSecret  string
	accounts   map[string]string
}

func (c *CoinbaseClient) Initialize() {
	c.url = os.Getenv("COINBASE_URL")
	c.apiVersion = os.Getenv("COINBASE_VERSION")
	c.apiKey = os.Getenv("COINBASE_API_KEY")
	c.apiSecret = os.Getenv("COINBASE_API_SECRET")

	// BTC=xxx-----ETH=xxx
	mapAccounts := os.Getenv("COINBASE_ACCOUNTS")
	ss := strings.Split(mapAccounts, "-----")
	c.accounts = make(map[string]string)
	for _, pair := range ss {
		z := strings.Split(pair, "=")
		c.accounts[z[0]] = z[1]
	}
}

func (c CoinbaseClient) buildHeader(method string, url string, body string) map[string]string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	message := []byte(timestamp + method + url + body)
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write(message)
	sign := hex.EncodeToString(mac.Sum(nil))

	headers := map[string]string{
		"CB-ACCESS-SIGN":      sign,
		"CB-ACCESS-TIMESTAMP": timestamp,
		"CB-ACCESS-KEY":       c.apiKey,
		"CB-VERSION":          c.apiVersion,
	}

	return headers
}

func (c CoinbaseClient) Get(uri string) (*grequests.Response, error) {
	url := c.url + uri
	headers := c.buildHeader("GET", uri, "")
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func (c CoinbaseClient) Post(uri string, body interface{}) (*grequests.Response, error) {
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

func (c CoinbaseClient) GetAccount(currency string) string {
	return c.accounts[currency]
}

//func GetAccount() (bean.CoinbaseAccountResponse, error) {
//	client := CoinbaseClient{}
//	client.Initialize()
//
//	var response bean.CoinbaseAccountResponse
//	resp, err := client.Get("/v2/accounts")
//	if err == nil {
//		resp.JSON(&response)
//	}
//
//	return response, err
//}

func GetName() string {
	return "Coinbase"
}

func GenerateAddress(currency string) (bean.CoinbaseAddressResponse, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbaseAddressResponse
	resp, err := client.Post("/v2/accounts/"+client.GetAccount(currency)+"/addresses", nil)
	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func GetTransaction(transId string, currency string) (bean.CoinbaseTransaction, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbaseTransactionResponse
	resp, err := client.Get("/v2/accounts/" + client.GetAccount(currency) + "/transactions/" + transId)
	if err == nil {
		resp.JSON(&response)
	}

	return response.Data, err
}

func SendTransaction(address string, amount string, currency string, description string, withdrawId string) (bean.CoinbaseTransaction, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbaseTransactionResponse

	resp, err := client.Post("/v2/accounts/"+client.GetAccount(currency)+"/transactions", bean.CoinbaseSendMoneyRequest{
		To:          address,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Idem:        withdrawId,
	}.GetRequestBody())

	if err == nil {
		resp.JSON(&response)
	}

	return response.Data, err
}

func GetNotification(resource string) (bean.CoinbaseNotification, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbaseNotificationResponse

	resp, err := client.Get(resource)
	if err == nil {
		resp.JSON(&response)
	}

	return response.Data, err
}

func GetBuyPrice(currency string) (bean.CoinbaseAmount, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbasePriceResponse

	resp, err := client.Get(fmt.Sprintf("/v2/prices/%s-USD/buy", currency))
	if err == nil {
		resp.JSON(&response)
	}

	return response.Data, err
}

func GetSellPrice(currency string) (bean.CoinbaseAmount, error) {
	client := CoinbaseClient{}
	client.Initialize()

	var response bean.CoinbasePriceResponse

	resp, err := client.Get(fmt.Sprintf("/v2/prices/%s-USD/sell", currency))
	if err == nil {
		resp.JSON(&response)
	}

	return response.Data, err
}
