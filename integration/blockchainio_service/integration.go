package blockchainio_service

import (
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/shopspring/decimal"
	"os"
)

type BlockChainIOClient struct {
	url         string
	apiKey      string
	apiPassword string
}

func (c *BlockChainIOClient) Initialize() {
	c.url = os.Getenv("BLOCKCHAINIO_URL")
	c.apiKey = os.Getenv("BLOCKCHAINIO_KEY")
	c.apiPassword = os.Getenv("BLOCKCHAINIO_PASSWORD")
}

func (c *BlockChainIOClient) GetBalance() (decimal.Decimal, error) {
	return common.Zero, nil
}

func (c *BlockChainIOClient) SendTransaction(address string, amount decimal.Decimal) (string, error) {
	return "", nil
}
