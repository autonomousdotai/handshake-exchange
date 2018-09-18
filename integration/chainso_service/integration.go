package chainso_service

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/shopspring/decimal"
)

func GetConfirmations(txId string, currency string) (int, error) {
	url := fmt.Sprintf("https://chain.so/api/v2/is_tx_confirmed/%s/%s", currency, txId)
	headers := map[string]string{
		"Accept": "application/json",
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return 0, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	var data map[string]interface{}
	result := 0
	if err == nil {
		resp.JSON(&data)
		dataNode := data["data"]
		resultField := dataNode.(map[string]interface{})["confirmations"]
		result = int(resultField.(float64))
	}

	return result, err
}

func GetAmount(txId string) (decimal.Decimal, error) {
	url := fmt.Sprintf("https://chain.so/api/v2/get_tx_outputs/BTC/%s/0", txId)
	headers := map[string]string{
		"Accept": "application/json",
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return common.Zero, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	var data map[string]interface{}
	result := common.Zero
	if err == nil {
		resp.JSON(&data)
		dataNode := data["data"]
		outputsNode := dataNode.(map[string]interface{})["outputs"]
		value := outputsNode.(map[string]interface{})["value"].(string)
		result = common.StringToDecimal(value)
	}

	return result, err
}
