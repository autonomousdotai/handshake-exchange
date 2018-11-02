package bitpay_service

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/shopspring/decimal"
)

type ScriptPubKey struct {
	Addresses []string `json:"addresses"`
}

type TxVOut struct {
	Value        string       `json:"value"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type TxResponse struct {
	VOut          []TxVOut `json:"vout"`
	Confirmations int      `json:"confirmations"`
}

type AdddressResponse struct {
	Balance      float64  `json:"balance"`
	Transactions []string `json:"transactions"`
}

func GetBCHTransaction(txId string) (decimal.Decimal, string, int, error) {
	url := fmt.Sprintf("https://bch-insight.bitpay.com/api/tx/%s", txId)
	headers := map[string]string{
		"Accept": "application/json",
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return common.Zero, "", 0, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	var data TxResponse
	var value decimal.Decimal
	var address string
	var confirmations int
	if err == nil {
		err = resp.JSON(&data)
		if err == nil {
			value = common.StringToDecimal(data.VOut[0].Value)
			address = data.VOut[0].ScriptPubKey.Addresses[0]
			confirmations = data.Confirmations
		}
	}

	return value, address, confirmations, err
}

func GetBTCAddress(address string) (string, float64, int, error) {
	url := fmt.Sprintf("https://insight.bitpay.com/api/addr/%s", address)
	headers := map[string]string{
		"Accept": "application/json",
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return "", 0, 0, nil
	}

	txCount := 0
	var txHash string
	var data AdddressResponse
	if err == nil {
		err = resp.JSON(&data)
		fmt.Println(data)
		if err == nil {
			txCount = len(data.Transactions)
			if len(data.Transactions) > 0 {
				txHash = data.Transactions[0]
			}
		}
	}

	return txHash, data.Balance, txCount, err
}

func GetBCHAddress(address string) (string, float64, int, error) {
	url := fmt.Sprintf("https://bch-insight.bitpay.com/api/addr/%s", address)
	headers := map[string]string{
		"Accept": "application/json",
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return "", 0, 0, nil
	}

	txCount := 0
	var txHash string
	var data AdddressResponse
	if err == nil {
		err = resp.JSON(&data)
		fmt.Println(data)
		if err == nil {
			txCount = len(data.Transactions)
			if len(data.Transactions) > 0 {
				txHash = data.Transactions[0]
			}
		}
	}

	return txHash, data.Balance, txCount, err
}
