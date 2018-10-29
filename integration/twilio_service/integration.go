package twilio_service

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"os"
)

type TwilioClient struct {
	url          string
	apiSid       string
	apiAuthToken string
	fromNumber   string
}

func (c *TwilioClient) Initialize() {
	c.url = os.Getenv("TWILIO_URL")
	c.apiSid = os.Getenv("TWILIO_API_SID")
	c.apiAuthToken = os.Getenv("TWILIO_API_AUTH_TOKEN")
	c.fromNumber = os.Getenv("TWILIO_FROM_NUMBER")
}

func (c TwilioClient) Post(uri string, params map[string]string, body interface{}) (*grequests.Response, error) {
	url := c.url + fmt.Sprintf("/2010-04-01/Accounts/%s", c.apiSid) + uri
	data := map[string]string{
		"From": c.fromNumber,
	}
	for key, value := range params {
		data[key] = value
	}
	ro := &grequests.RequestOptions{Headers: map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	}, Data: data, Auth: []string{c.apiSid, c.apiAuthToken}}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func SendSMS(toPhone string, body string) (SendSMSResponse, error) {
	client := TwilioClient{}
	client.Initialize()

	var response SendSMSResponse

	resp, err := client.Post("/Messages.json", map[string]string{
		"To":   toPhone,
		"Body": body,
	}, nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

func SendVoice(toPhone string, url string) (SendSMSResponse, error) {
	client := TwilioClient{}
	client.Initialize()

	var response SendSMSResponse

	resp, err := client.Post("/Calls.json", map[string]string{
		"To":  toPhone,
		"Url": url,
	}, nil)

	if err == nil {
		resp.JSON(&response)
	}

	return response, err
}

type SendSMSResponse struct {
	Sid    string
	Status string
}
