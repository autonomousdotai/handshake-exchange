package fcm_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"os"
)

func SendFCM(fcm bean.FCMObject) error {
	host := os.Getenv("FCM_SERVICE_URL")
	url := fmt.Sprintf("%s/send", host)

	bodyStr := ""
	b, errBody := json.Marshal(&fcm)
	if errBody != nil {
		return errBody
	}
	bodyStr = string(b)
	fmt.Println(url)
	fmt.Println(bodyStr)

	r := bytes.NewReader([]byte(bodyStr))

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	ro := &grequests.RequestOptions{RequestBody: r, Headers: headers}
	resp, err := grequests.Post(url, ro)
	fmt.Println(err)
	fmt.Println(resp)

	if resp.Ok != true {
		return api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return err
}
