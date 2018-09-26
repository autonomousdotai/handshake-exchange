package adyen_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"os"
)

type AdyenAmount struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

type AdyenAuthorise struct {
	Card              map[string]string      `json:"card"`
	Amount            AdyenAmount            `json:"amount"`
	Reference         string                 `json:"reference"`
	MerchantAccount   string                 `json:"merchantAccount"`
	CaptureDelayHours int                    `json:"captureDelayHours"`
	AdditionalData    map[string]interface{} `json:"additionalData"`
}

type AdyenAuthorise3D struct {
	MD              string `json:"md"`
	PAResponse      string `json:"paResponse"`
	ShopperIP       string `json:"shopperIP"`
	MerchantAccount string `json:"merchantAccount"`
}

type AdyenCapture struct {
	OriginalReference  string      `json:"originalReference"`
	ModificationAmount AdyenAmount `json:"modificationAmount"`
	Reference          string      `json:"reference"`
	MerchantAccount    string      `json:"merchantAccount"`
}

type AdyenCancel struct {
	OriginalReference string `json:"originalReference"`
	Reference         string `json:"reference"`
	MerchantAccount   string `json:"merchantAccount"`
}

type AdyenSimpleResponse struct {
	AdditionalData map[string]interface{} `json:"additionalData"`
	PSPReference   string                 `json:"pspReference"`
	Response       string                 `json:"response"`
}

type AdyenAuthoriseResponse struct {
	AdditionalData map[string]interface{} `json:"additionalData"`
	AuthCode       string                 `json:"authCode"`
	IssueUrl       string                 `json:"issuerUrl"`
	MD             string                 `json:"md"`
	PARequest      string                 `json:"paRequest"`
	PSPReference   string                 `json:"pspReference"`
	ResultCode     string                 `json:"resultCode"`
}

type AdyenClient struct {
	url             string
	apiUsername     string
	apiPassword     string
	MerchantAccount string
}

func (c *AdyenClient) Initialize() {
	c.url = os.Getenv("ADYEN_URL")
	c.apiUsername = os.Getenv("ADYEN_USERNAME")
	c.apiPassword = os.Getenv("ADYEN_PASSWORD")
	c.MerchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")
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
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.String())
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func Authorise(authoriseObj AdyenAuthorise) (AdyenAuthoriseResponse, error) {
	client := AdyenClient{}

	var response AdyenAuthoriseResponse
	authoriseObj.AdditionalData["executeThreeD"] = "true"
	authoriseObj.CaptureDelayHours = 1
	authoriseObj.MerchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")

	fmt.Println(authoriseObj)
	resp, err := client.Post("/authorise", authoriseObj)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func Authorise3D(authorise3DObj AdyenAuthorise3D) (AdyenAuthoriseResponse, error) {
	client := AdyenClient{}

	var response AdyenAuthoriseResponse
	authorise3DObj.MerchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")

	resp, err := client.Post("/authorise3d", authorise3DObj)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func Capture(captureObj AdyenCapture) (AdyenSimpleResponse, error) {
	client := AdyenClient{}

	var response AdyenSimpleResponse

	captureObj.MerchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")
	resp, err := client.Post("/capture", captureObj)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func CancelOrRefund(cancelObject AdyenCancel) (AdyenSimpleResponse, error) {
	client := AdyenClient{}

	var response AdyenSimpleResponse

	cancelObject.MerchantAccount = os.Getenv("ADYEN_MERCHANT_ACCOUNT")
	resp, err := client.Post("/cancelOrRefund", cancelObject)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}
