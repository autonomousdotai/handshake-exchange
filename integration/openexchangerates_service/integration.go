package openexchangerates_service

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"os"
)

func GetExchangeRate() (map[string]float64, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", os.Getenv("OPENEXCHANGERATES_API_KEY"))
	resp, err := grequests.Get(url, nil)
	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	var data map[string]interface{}
	exchange := make(map[string]float64)
	if err == nil {
		resp.JSON(&data)
		if rates, ok := data["rates"]; ok {
			m := rates.(map[string]interface{})
			for k := range m {
				rate := m[k].(float64)
				exchange[k] = rate
			}
		}
	}

	return exchange, err
}
