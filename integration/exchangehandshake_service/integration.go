package exchangehandshake_service

import (
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ninjadotorg/handshake-exchange/abi"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"os"
	// "fmt"
)

type ExchangeHandshakeClient struct {
	client    *ethclient.Client
	address   common.Address
	handshake *abi.ExchangeHandshake
}

func (c *ExchangeHandshakeClient) initialize() (err error) {
	c.client, err = ethclient.Dial(os.Getenv("ETH_NETWORK"))
	if err != nil {
		return
	}
	c.address = common.HexToAddress(os.Getenv("ETH_EXCHANGE_HANDSHAKE_ADDRESS"))
	c.handshake, err = abi.NewExchangeHandshake(c.address, c.client)
	if err != nil {
		return
	}

	return
}

func (c *ExchangeHandshakeClient) close() {
	c.client.Close()
}

func (c *ExchangeHandshakeClient) GetInitEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterInitByCoinOwner(opt)
	if errInit != nil {
		err = errInit
		return
	}

	notEmpty := true
	endBlock = startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock = past.Event.Raw.BlockNumber

			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeClient) GetShakeEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterShake(opt)
	if errInit != nil {
		err = errInit
		return
	}

	notEmpty := true
	endBlock = startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock = past.Event.Raw.BlockNumber

			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeClient) GetRejectEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterCancel(opt)
	if errInit != nil {
		err = errInit
		return
	}

	notEmpty := true
	endBlock = startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock = past.Event.Raw.BlockNumber

			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeClient) GetCompleteEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterAccept(opt)
	if errInit != nil {
		err = errInit
		return
	}

	notEmpty := true
	endBlock = startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock = past.Event.Raw.BlockNumber

			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}
	c.close()

	return
}
