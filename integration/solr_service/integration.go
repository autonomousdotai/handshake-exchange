package solr_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/levigross/grequests"
	"os"
)

func UpdateObject(body interface{}) (*grequests.Response, error) {
	type addBodyStruct struct {
		Add []interface{} `json:"add"`
	}
	host := os.Getenv("SOLR_SERVICE_URL")
	url := fmt.Sprintf("%s/handshake/update", host)

	bodyStr := ""

	arrBody := make([]interface{}, 1)
	arrBody[0] = body
	addBody := addBodyStruct{
		Add: arrBody,
	}

	b, errBody := json.Marshal(&addBody)
	if errBody != nil {
		return nil, errBody
	}
	bodyStr = string(b)

	r := bytes.NewReader([]byte(bodyStr))

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	ro := &grequests.RequestOptions{RequestBody: r, Headers: headers}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}

func DeleteObject(objectId string) (*grequests.Response, error) {
	type deleteBodyStruct struct {
		Delete []string `json:"delete"`
	}

	host := os.Getenv("SOLR_SERVICE_URL")
	url := fmt.Sprintf("%s/handshake/update", host)

	bodyStr := ""
	arrBody := make([]string, 1)
	arrBody[0] = fmt.Sprintf("exchange_%s", objectId)
	deleteBody := deleteBodyStruct{
		Delete: arrBody,
	}

	b, errBody := json.Marshal(&deleteBody)
	if errBody != nil {
		return nil, errBody
	}
	bodyStr = string(b)
	r := bytes.NewReader([]byte(bodyStr))

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	ro := &grequests.RequestOptions{RequestBody: r, Headers: headers}
	resp, err := grequests.Post(url, ro)

	if resp.Ok != true {
		return nil, api_error.NewErrorCustom(api_error.ExternalApiFailed, resp.String(), nil)
	}

	return resp, err
}
