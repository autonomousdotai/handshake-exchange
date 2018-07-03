package crypto_service

import (
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/blockchainio_service"
	"github.com/ninjadotorg/handshake-exchange/integration/ethereum_service"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func GetBalance(currency string) (decimal.Decimal, error) {
	if currency == bean.ETH.Code {
		client := ethereum_service.EthereumClient{}
		return client.GetBalance()
	} else if currency == bean.BTC.Code {
		client := blockchainio_service.BlockChainIOClient{}
		return client.GetBalance()
	}

	return common.Zero, errors.New("Currency not support")
}

func SendTransaction(address string, amountStr string, currency string) (string, error) {
	amount, _ := decimal.NewFromString(amountStr)
	if currency == bean.ETH.Code {
		client := ethereum_service.EthereumClient{}
		return client.SendTransaction(address, amount)
	} else if currency == bean.BTC.Code {
		client := blockchainio_service.BlockChainIOClient{}
		return client.SendTransaction(address, amount)
	}

	return "", errors.New("Currency not support")
}

func GetTransactionReceipt(txHash string, currency string) (isSuccess bool, isPending bool, err error) {
	if currency == bean.ETH.Code {
		client := ethereum_service.EthereumClient{}
		return client.GetTransactionReceipt(txHash)
	} else if currency == bean.BTC.Code {
		return true, false, nil
	}

	return false, false, nil
}
