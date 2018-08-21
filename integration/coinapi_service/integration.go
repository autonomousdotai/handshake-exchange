package coinapi_service

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"os"
	"strings"
)

func GetExchangeRate(currencies []string) (map[string][]bean.CryptoRate, error) {
	url := fmt.Sprintf("https://rest.coinapi.io/v1/quotes/current")
	headers := map[string]string{
		"X-CoinAPI-Key": os.Getenv("COINAPI_API_KEY"),
	}
	ro := &grequests.RequestOptions{Headers: headers}
	resp, err := grequests.Get(url, ro)
	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	var data []map[string]interface{}
	result := make(map[string][]bean.CryptoRate)
	if err == nil {
		resp.JSON(&data)
		for _, currency := range currencies {
			dataItems := make([]bean.CryptoRate, 0)
			for _, item := range data {
				symbolId := item["symbol_id"].(string)
				symbol := fmt.Sprintf("BITFINEX_SPOT_%s_USD", currency)
				if strings.Contains(symbolId, symbol) {
					dataItems = append(dataItems, bean.CryptoRate{
						From:     currency,
						To:       bean.USD.Code,
						Buy:      item["bid_price"].(float64),
						Sell:     item["ask_price"].(float64),
						Exchange: bean.INSTANT_OFFER_PROVIDER_COINBASE, // Hard code
						// Exchange: strings.Replace(symbolId, symbol, "", -1),
					})
				}
			}
			result[currency] = dataItems
		}
	}

	return result, err
}
