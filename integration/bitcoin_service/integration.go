package bitcoin_service

import "github.com/shopspring/decimal"

type BitcoinService struct {
}

func (c *BitcoinService) SendTransaction(address string, amount decimal.Decimal) (string, error) {
	return "", nil
}
