package bitstamp_service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"os"
	"strconv"
	"strings"
	"time"
)

type BitstampClient struct {
	url        string
	apiKey     string
	apiSecret  string
	custonerId string
}

func (c *BitstampClient) Initialize() {
	c.url = os.Getenv("BITSTAMP_URL")
	c.apiKey = os.Getenv("BITSTAMP_API_KEY")
	c.apiSecret = os.Getenv("BITSTAMP_API_SECRET")
	c.custonerId = os.Getenv("BITSTAMP_CUSTOMER_ID")
}

func (c BitstampClient) buildAuthParameters() string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	message := []byte(timestamp + c.custonerId + c.apiKey)
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write(message)
	sign := strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))

	requestParams := fmt.Sprintf("key=%s&signature=%s&nonce=%s", c.apiKey, sign, timestamp)
	fmt.Println(sign)

	return requestParams
}

func (c BitstampClient) Get(uri string) (*grequests.Response, error) {
	url := c.url + uri

	ro := &grequests.RequestOptions{}
	resp, err := grequests.Get(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func (c BitstampClient) Post(uri string, params map[string]string, body interface{}) (*grequests.Response, error) {
	url := c.url + uri

	authParams := c.buildAuthParameters()
	for key, value := range params {
		authParams += fmt.Sprintf("&%s=%s", key, value)
	}
	r := bytes.NewReader([]byte(authParams))
	ro := &grequests.RequestOptions{Headers: map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, RequestBody: r}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func GetBuyPrice(currency string) (TickerResponse, error) {
	client := BitstampClient{}
	client.Initialize()

	var response TickerResponse

	uri := fmt.Sprintf("/v2/ticker/%susd/", strings.ToLower(currency))
	resp, err := client.Get(uri)
	if err == nil {
		resp.JSON(&response)
	}
	response.Amount = response.Ask

	return response, err
}

func GetSellPrice(currency string) (TickerResponse, error) {
	client := BitstampClient{}
	client.Initialize()

	var response TickerResponse

	resp, err := client.Get(fmt.Sprintf("/v2/ticker/%s/", strings.ToLower(fmt.Sprintf(currency, "usd"))))
	if err == nil {
		resp.JSON(&response)
	}
	response.Amount = response.Bid

	return response, err
}

func SendTransaction(address string, amount string, currency string, description string, withdrawId string) (TransferResponse, error) {
	client := BitstampClient{}
	client.Initialize()

	var response TransferResponse

	currencyMapping := map[string]string{
		bean.BTC.Code: "bitcoin_withdrawal",
		bean.ETH.Code: "v2/eth_withdrawal",
		bean.BCH.Code: "v2/bch_withdrawal",
	}

	resp, err := client.Post(fmt.Sprintf("/%s/", currencyMapping[currency]), map[string]string{
		"amount":  amount,
		"address": address,
	}, nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

type TickerResponse struct {
	Last   string `json:"last"`
	Bid    string `json:"bid"`
	Ask    string `json:"ask"`
	Amount string
}

type TransferResponse struct {
	Id string `json:"id"`
}

type WithdrawRequestResponse struct {
	Id            string                 `json:"id"`
	DateTime      string                 `json:"datetime"`
	Type          string                 `json:"type"`
	Currency      string                 `json:"currency"`
	Amount        string                 `json:"amount"`
	Status        string                 `json:"status"`
	Data          map[string]interface{} `json:"data"`
	Address       string                 `json:"address"`
	TransactionId string                 `json:"transaction_id"`
}

func (b WithdrawRequestResponse) GetCurrency() string {
	if b.Currency != "" {
		return b.Currency
	}
	if b.Type == "1" {
		return bean.BTC.Code
	}

	return ""
}
