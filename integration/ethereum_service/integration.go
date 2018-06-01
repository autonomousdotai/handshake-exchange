package ethereum_service

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"math/big"
	"os"
)

var WeiDecimal = decimal.NewFromBigInt(big.NewInt(1000000000000000000), 0)

type EthereumClient struct {
	client     *ethclient.Client
	address    common.Address
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (c *EthereumClient) initialize() (err error) {
	c.client, err = ethclient.Dial(os.Getenv("ETH_NETWORK"))
	if err != nil {
		return
	}

	c.privateKey, err = crypto.HexToECDSA(os.Getenv("ETH_KEY"))
	if err != nil {
		return
	}

	publicKey := c.privateKey.Public()
	c.publicKey, _ = publicKey.(*ecdsa.PublicKey)

	c.address = crypto.PubkeyToAddress(*c.publicKey)

	return
}

func (c *EthereumClient) close() {
	c.client.Close()
}

func (c *EthereumClient) SendTransaction(address string, amount decimal.Decimal) (string, error) {
	c.initialize()

	nonce, err := c.client.PendingNonceAt(context.Background(), c.address)
	if err == nil {
		value := big.NewInt(amount.Mul(WeiDecimal).IntPart()) // in wei
		gasLimit := uint64(21000)
		gasPrice, err := c.client.SuggestGasPrice(context.Background())
		if err == nil {
			toAddress := common.HexToAddress(address)
			tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

			signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, c.privateKey)
			if err == nil {
				err = c.client.SendTransaction(context.Background(), signedTx)
				if err == nil {
					return signedTx.Hash().Hex(), nil
				}
			}
		}
	}

	c.close()

	return "", err
}

func (c *EthereumClient) GetBalance() (balance decimal.Decimal, err error) {
	c.initialize()

	intBalance, errCli := c.client.BalanceAt(context.Background(), c.address, nil)

	if errCli == nil {
		balance = decimal.NewFromBigInt(intBalance, 0).Div(WeiDecimal)
	}

	err = errCli
	c.close()

	return
}
