package adyen_service

import (
	"bytes"
	"encoding/json"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"os"
)

type AdyenAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type AdyenAuthorise struct {
	Amount            AdyenAmount            `json:"amount"`
	Reference         string                 `json:"reference"`
	MerchantAccount   string                 `json:"merchantAccount"`
	CaptureDelayHours int                    `json:"captureDelayHours"`
	AdditionalData    map[string]interface{} `json:"additionalData"`
}

type AdyenClient struct {
	url         string
	apiUsername string
	apiPassword string
}

func (c *AdyenClient) Initialize() {
	c.url = os.Getenv("ADYEN_URL")
	c.apiUsername = os.Getenv("ADYEN_USERNAME")
	c.apiPassword = os.Getenv("ADYEN_PASSWORD")
}

func (c AdyenClient) buildHeader() map[string]string {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;charset=UTF-8",
	}

	return headers
}

func (c AdyenClient) Post(uri string, body interface{}) (*grequests.Response, error) {
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
	ro := &grequests.RequestOptions{Headers: headers, RequestBody: r, Auth: []string{c.apiUsername, c.apiPassword}}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}
