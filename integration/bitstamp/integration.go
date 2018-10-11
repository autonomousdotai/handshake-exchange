package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
)

type BitstampClient struct {
	url        string
	apiKey     string
	apiSecret  string
	custonerId string
}

func (c *BitstampClient) Initialize() {
	c.url = os.Getenv("BITSTAMP_URL")
	c.apiKey = os.Getenv("BITSTAMP_API_KEY")
	c.apiSecret = os.Getenv("BITSTAMP_API_SECRET")
	c.custonerId = os.Getenv("BITSTAMP_CUSTOMER_ID")
}

func (c BitstampClient) buildAuthParameters(method string, url string, body string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	message := []byte(timestamp + c.custonerId + c.apiKey)
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write(message)
	sign := hex.EncodeToString(mac.Sum(nil))

	requestParams := fmt.Sprintf("key=%s&signature=%s&nonce=%s", c.apiKey, sign, timestamp)

	return requestParams
}
