package exchangecreditatm_service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ninjadotorg/handshake-exchange/abi"
	"github.com/ninjadotorg/handshake-exchange/integration/ethereum_service"
	"github.com/shopspring/decimal"
	"math/big"
	"os"
)

var WeiDecimal = decimal.NewFromBigInt(big.NewInt(1000000000000000000), 0)

type ExchangeCreditAtmClient struct {
	client      *ethclient.Client
	address     common.Address
	creditAtm   *abi.CreditATM
	writeClient ethereum_service.EthereumClient
}

func (c *ExchangeCreditAtmClient) initialize() (err error) {
	c.client, err = ethclient.Dial(os.Getenv("ETH_NETWORK"))
	if err != nil {
		return
	}
	c.address = common.HexToAddress(os.Getenv("ETH_EXCHANGE_CREDIT_ATM"))
	c.creditAtm, err = abi.NewCreditATM(c.address, c.client)
	if err != nil {
		return
	}

	return
}

func (c *ExchangeCreditAtmClient) initializeWrite() {
	c.writeClient = ethereum_service.EthereumClient{}
	c.writeClient.Initialize()
}

func (c *ExchangeCreditAtmClient) close() {
	c.client.Close()
}

func (c *ExchangeCreditAtmClient) closeWrite() {
	c.writeClient.Close()
}

func (c *ExchangeCreditAtmClient) ReleasePartialFund(offerId string, hid int64, amount decimal.Decimal, address string) (txHash string, err error) {
	c.initialize()
	c.initializeWrite()

	auth, err := c.writeClient.GetAuth(decimal.NewFromFloat(0))

	offChain := [32]byte{}
	copy(offChain[:], []byte(offerId))

	fmt.Println(fmt.Sprintf("%s %s %s %s", hid, address, amount, offerId))
	toAddress := common.HexToAddress(address)
	sendAmount := big.NewInt(amount.Mul(WeiDecimal).IntPart())

	tx, err := c.creditAtm.ReleasePartialFund(auth, big.NewInt(hid), toAddress, sendAmount, offChain)
	if err != nil {
		return
	}
	txHash = tx.Hash().Hex()

	c.closeWrite()
	c.close()

	return
}
