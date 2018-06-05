package exchangehandshake_service

import (
	"bytes"
	"github.com/ninjadotorg/handshake-exchange/abi"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
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
	endBlock1 := startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock1 = past.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}

	past2, errInit := c.handshake.FilterInitByCashOwner(opt)
	if errInit != nil {
		err = errInit
		return
	}
	notEmpty = true
	endBlock2 := startBlock
	for notEmpty {
		notEmpty = past2.Next()
		if notEmpty {
			endBlock2 = past2.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past2.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past2.Event.Offchain[:], "\x00")),
			})
		}
	}

	if endBlock1 > endBlock2 {
		endBlock = endBlock1
	} else {
		endBlock = endBlock2
	}

	c.close()

	return
}

func (c *ExchangeHandshakeClient) GetCloseEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterCloseByCoinOwner(opt)
	if errInit != nil {
		err = errInit
		return
	}
	notEmpty := true
	endBlock1 := startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock1 = past.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}

	past2, errInit := c.handshake.FilterCloseByCashOwner(opt)
	if errInit != nil {
		err = errInit
		return
	}
	notEmpty = true
	endBlock2 := startBlock
	for notEmpty {
		notEmpty = past2.Next()
		if notEmpty {
			endBlock2 = past2.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past2.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past2.Event.Offchain[:], "\x00")),
			})
		}
	}

	if endBlock1 > endBlock2 {
		endBlock = endBlock1
	} else {
		endBlock = endBlock2
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
	past, errInit := c.handshake.FilterReject(opt)
	if errInit != nil {
		err = errInit
		return
	}
	notEmpty := true
	endBlock1 := startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock1 = past.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past.Event.Offchain[:], "\x00")),
			})
		}
	}

	past2, errInit := c.handshake.FilterCancel(opt)
	if errInit != nil {
		err = errInit
		return
	}
	notEmpty = true
	endBlock2 := startBlock
	for notEmpty {
		notEmpty = past2.Next()
		if notEmpty {
			endBlock2 = past2.Event.Raw.BlockNumber
			offers = append(offers, bean.OfferOnchain{
				Hid:   int64(past2.Event.Hid.Uint64()),
				Offer: string(bytes.Trim(past2.Event.Offchain[:], "\x00")),
			})
		}
	}

	if endBlock1 > endBlock2 {
		endBlock = endBlock1
	} else {
		endBlock = endBlock2
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

func (c *ExchangeHandshakeClient) GetWithdrawEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterWithdraw(opt)
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
