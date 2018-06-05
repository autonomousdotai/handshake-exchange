package blockchainio_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/levigross/grequests"
	"github.com/shopspring/decimal"
	"math/big"
	"os"
)

type BlockChainIOClient struct {
	address     string
	guid        string
	url         string
	apiKey      string
	apiPassword string
}

var BTC_IN_SATOSHI = decimal.NewFromBigInt(big.NewInt(100000000), 0)

func (c *BlockChainIOClient) initialize() {
	c.address = os.Getenv("BLOCKCHAINIO_ADDRESS")
	c.guid = os.Getenv("BLOCKCHAINIO_GUID")

	c.url = os.Getenv("BLOCKCHAINIO_URL")
	c.apiKey = os.Getenv("BLOCKCHAINIO_KEY")
	c.apiPassword = os.Getenv("BLOCKCHAINIO_PASSWORD")
}

func (c *BlockChainIOClient) post(uri string) (*grequests.Response, error) {
	url := c.url + uri

	ro := &grequests.RequestOptions{}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func (c *BlockChainIOClient) get(uri string) (*grequests.Response, error) {
	url := c.url + uri

	ro := &grequests.RequestOptions{}
	resp, err := grequests.Get(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func (c *BlockChainIOClient) GetBalance() (decimal.Decimal, error) {
	//TODO Enable when BlockchainIO
	//c.initialize()
	//
	//var response bean.BlockChainIoBalance
	//resp, err := c.get(fmt.Sprintf("/merchant/%s/address_balance?password=%s&api_code=%s&address=%s", c.guid, c.apiPassword, c.apiKey, c.address))
	//if err == nil {
	//	resp.JSON(&response)
	//}
	//
	//return decimal.NewFromBigInt(big.NewInt(response.Balance), 0).Div(BTC_IN_SATOSHI), err

	return decimal.NewFromFloat(0), nil
}

func (c *BlockChainIOClient) SendTransaction(address string, amount decimal.Decimal) (string, error) {
	c.initialize()

	var response bean.BlockChainIoPayment
	sendAmount := amount.Mul(BTC_IN_SATOSHI)
	resp, err := c.post(fmt.Sprintf("/merchant/%s/payment?password=%s&api_code=%s&to=%s&amount=%d", c.guid, c.apiPassword, c.apiKey, address, sendAmount.IntPart()))
	if err == nil {
		resp.JSON(&response)
	}

	return response.TxHash, err
}

func (c *BlockChainIOClient) GenerateAddress(offerId string) (string, error) {
	c.initialize()

	var response bean.BlockChainIoAddress
	resp, err := c.post(fmt.Sprintf("/merchant/%s/new_address?password=%s&api_code=%s&label=%s", c.guid, c.apiPassword, c.apiKey, offerId))
	if err == nil {
		resp.JSON(&response)
	}

	monitorRequest := bean.BlockChainIoBalanceUpdates{
		Address:        response.Address,
		Op:             "RECEIVE",
		Confirmations:  1,
		OnNotification: "DELETE",
	}

	callbackClient := BlockChainIOCallbackClient{}
	err = callbackClient.MonitorAddress(&monitorRequest, offerId)

	return response.Address, err
}

type BlockChainIOCallbackClient struct {
	url         string
	apiKey      string
	callbackUrl string
}

func (c *BlockChainIOCallbackClient) initialize() {
	c.url = os.Getenv("BLOCKCHAINIO_CALLBACK_URL")
	c.apiKey = os.Getenv("BLOCKCHAINIO_CALLBACK_KEY")
	c.callbackUrl = os.Getenv("BLOCKCHAINIO_URL_TO_CALL")
}

func (c *BlockChainIOCallbackClient) MonitorAddress(update *bean.BlockChainIoBalanceUpdates, extraData string) error {
	c.initialize()

	url := fmt.Sprintf("%s/receive/balance_update", c.url)

	headers := map[string]string{
		"Content-Type": "text/plain",
	}

	update.Callback = fmt.Sprintf("%s?offer=%s", c.callbackUrl, extraData)

	bodyStr := ""
	b, errBody := json.Marshal(&update)
	if errBody != nil {
		return errBody
	}
	bodyStr = string(b)

	r := bytes.NewReader([]byte(bodyStr))
	ro := &grequests.RequestOptions{Headers: headers, RequestBody: r}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	resp.JSON(&update)

	return err
}
