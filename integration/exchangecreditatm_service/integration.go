package exchangecreditatm_service

import (
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

func (c *ExchangeCreditAtmClient) ReleasePartialFund(offerId string, hid int64, amount decimal.Decimal, address string,
	inNonce uint64, overwriteNonce bool) (txHash string, outNonce uint64, err error) {
	c.initialize()
	c.initializeWrite()

	auth, err := c.writeClient.GetAuth(decimal.NewFromFloat(0))
	if auth.Nonce.Uint64() < inNonce {
		outNonce = inNonce
		auth.Nonce = big.NewInt(int64(inNonce))
	} else {
		outNonce = auth.Nonce.Uint64()
	}

	if overwriteNonce {
		auth.Nonce = big.NewInt(int64(inNonce))
	}

	offChain := [32]byte{}
	copy(offChain[:], []byte(offerId))

	toAddress := common.HexToAddress(address)

	decimalAmount := amount.Sub(amount.Floor())
	intAmount := amount.Sub(decimalAmount)

	weiBigAmount := big.NewInt(WeiDecimal.IntPart())
	intBigAmount := big.NewInt(intAmount.IntPart())
	intWeiAmount := intBigAmount.Mul(intBigAmount, weiBigAmount)
	decimalBigAmount := big.NewInt(decimalAmount.Mul(WeiDecimal).IntPart())

	sendAmount := intWeiAmount.Add(intWeiAmount, decimalBigAmount)

	tx, err := c.creditAtm.ReleasePartialFund(auth, toAddress, sendAmount, offChain)
	if err != nil {
		return
	}
	txHash = tx.Hash().Hex()

	c.closeWrite()
	c.close()

	return
}
