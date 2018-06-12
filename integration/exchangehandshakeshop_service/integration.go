package exchangehandshakeshop_service

import (
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ninjadotorg/handshake-exchange/abi"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"os"
)

type ExchangeHandshakeShopClient struct {
	client    *ethclient.Client
	address   common.Address
	handshake *abi.ExchangeHandshakeShop
}

func (c *ExchangeHandshakeShopClient) initialize() (err error) {
	c.client, err = ethclient.Dial(os.Getenv("ETH_NETWORK"))
	if err != nil {
		return
	}
	c.address = common.HexToAddress(os.Getenv("ETH_EXCHANGE_HANDSHAKE_OFFER_STORE_ADDRESS"))
	c.handshake, err = abi.NewExchangeHandshakeShop(c.address, c.client)
	if err != nil {
		return
	}

	return
}

func (c *ExchangeHandshakeShopClient) close() {
	c.client.Close()
}

func (c *ExchangeHandshakeShopClient) GetInitOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterInitByShopOwner(opt)
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetCloseOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterCloseByShopOwner(opt)
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetPreShakeOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterInitByCustomer(opt)
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetCancelOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetShakeOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetRejectOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
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
	endBlock = startBlock
	for notEmpty {
		notEmpty = past.Next()
		if notEmpty {
			endBlock = past.Event.Raw.BlockNumber

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetCompleteOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
	c.initialize()

	opt := &bind.FilterOpts{
		Start: startBlock,
	}
	past, errInit := c.handshake.FilterReleasePartialFund(opt)
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

			offerId := string(bytes.Trim(past.Event.OffchainP[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}

func (c *ExchangeHandshakeShopClient) GetCompleteUserOfferStoreEvent(startBlock uint64) (offers []bean.OfferOnchain, endBlock uint64, err error) {
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

			offerId := string(bytes.Trim(past.Event.Offchain[:], "\x00"))
			if offerId != "" {
				offers = append(offers, bean.OfferOnchain{
					Hid:   int64(past.Event.Hid.Uint64()),
					Offer: offerId,
				})
			}
		}
	}
	c.close()

	return
}
