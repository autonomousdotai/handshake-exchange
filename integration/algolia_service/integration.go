package algolia_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/levigross/grequests"
	"os"
)

func AddObject(uri string, body interface{}) (*grequests.Response, error) {
	host := os.Getenv("ALGOLIA_SERVICE_URL")
	url := fmt.Sprintf("%s/objects", host)

	bodyStr := ""
	arrBody := make([]interface{}, 1)
	arrBody[0] = body

	b, errBody := json.Marshal(&arrBody)
	if errBody != nil {
		return nil, errBody
	}
	bodyStr = string(b)
	r := bytes.NewReader([]byte(bodyStr))

	ro := &grequests.RequestOptions{RequestBody: r}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func UpdateObject(uri string, body interface{}) (*grequests.Response, error) {
	host := os.Getenv("ALGOLIA_SERVICE_URL")
	url := fmt.Sprintf("%s/objects", host)

	bodyStr := ""
	arrBody := make([]interface{}, 1)
	arrBody[0] = body

	b, errBody := json.Marshal(&arrBody)
	if errBody != nil {
		return nil, errBody
	}
	bodyStr = string(b)
	r := bytes.NewReader([]byte(bodyStr))

	ro := &grequests.RequestOptions{RequestBody: r}
	resp, err := grequests.Put(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func DeleteObject(uri string, objectId string) (*grequests.Response, error) {
	host := os.Getenv("ALGOLIA_SERVICE_URL")
	url := fmt.Sprintf("%s/objects", host)

	bodyStr := ""
	arrBody := make([]string, 1)
	arrBody[0] = objectId

	b, errBody := json.Marshal(&arrBody)
	if errBody != nil {
		return nil, errBody
	}
	bodyStr = string(b)
	r := bytes.NewReader([]byte(bodyStr))

	ro := &grequests.RequestOptions{RequestBody: r}
	resp, err := grequests.Delete(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}
