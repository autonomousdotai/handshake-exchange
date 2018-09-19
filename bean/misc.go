package bean

import "cloud.google.com/go/firestore"

type CurrencyRate struct {
	From string  `json:"from" firestore:"from"`
	To   string  `json:"to" firestore:"to"`
	Rate float64 `json:"rate" firestore:"rate"`
}

const FEE_TYPE_VALUE = "value"
const FEE_TYPE_PERCENTAGE = "percentage"

const FEE_KEY_EXCHANGE = "exchange"
const FEE_KEY_EXCHANGE_COMMISSION = "exchange_commission"
const FEE_KEY_INSTANT_BUY_CRYPTO = "instant_buy_crypto"
const FEE_KEY_CASH_BUY_CRYPTO = "cash_buy_crypto"
const FEE_KEY_CASH_STORE_SELL_CRYPTO = "cash_store_sell_crypto"

const CONFIG_KEY_CC_MODE = "CC_MODE"
const CC_MODE_GDAX = "gdax"
const CC_MODE_INVENTORY = "inventory"
const CC_MODE_CREDIT = "credit"

const CONFIG_BTC_WALLET = "BTC_WALLET"
const BTC_WALLET_COINBASE = "coinbase"
const BTC_WALLET_BLOCKCHAINIO = "blockchainio"

const CONFIG_OFFER_REJECT_LOCK = "OFFER_REJECT_LOCK"

const CONFIG_OFFER_STORE_FREE_START = "OFFER_STORE_FREE_START"
const OFFER_STORE_FREE_START_ON = "1"
const OFFER_STORE_FREE_START_OFF = "0"

const CONFIG_OFFER_STORE_REFERRAL_PERCENTAGE = "OFFER_STORE_REFERRAL_PERCENTAGE"

type SystemFee struct {
	Key   string  `json:"key" firestore:"key"`
	Value float64 `json:"value" firestore:"value"`
	Type  string  `json:"type" firestore:"type"`
}

type CryptoRate struct {
	From     string  `json:"from" firestore:"from"`
	To       string  `json:"to" firestore:"to"`
	Exchange string  `json:"exchange" firestore:"exchange"`
	Buy      float64 `json:"buy" firestore:"buy"`
	Sell     float64 `json:"sell" firestore:"sell"`
}

type TradingBot struct {
	UID       string `json:"uid" firestore:"uid"`
	Id        string `json:"id" firestore:"id"`
	Enabled   bool   `json:"enabled" firestore:"enabled"`
	Currency  string `json:"currency" firestore:"currency"`
	Type      string `json:"type" firestore:"type"`
	MinAmount string `json:"min_amount" firestore:"min_amount"`
	MaxAmount string `json:"max_amount" firestore:"max_amount"`
	Price     string `json:"price" firestore:"price"`
	Duration  int    `json:"duration" firestore:"duration"`
}

type TradingBotPendingOffer struct {
	Offer     string `json:"-" firestore:"offer"`
	OfferRef  string `json:"-" firestore:"offer_ref"`
	CreatedAt string `json:"-" firestore:"created_at"`
}

type CCLimit struct {
	Level    int64 `json:"level" firestore:"level"`
	Limit    int64 `json:"limit" firestore:"limit"`
	Duration int64 `json:"duration" firestore:"duration"`
}

type SystemConfig struct {
	Key   string `json:"key" firestore:"key"`
	Value string `json:"value" firestore:"value"`
}

type CryptoTransferLog struct {
	Id               string      `json:"id" firestore:"id"`
	Provider         string      `json:"provider" firestore:"provider"`
	ProviderResponse interface{} `json:"provider_response" firestore:"provider_response"`
	ExternalId       string      `json:"external_id" firestore:"external_id"`
	DataType         string      `json:"data_type" firestore:"data_type"`
	DataRef          string      `json:"data_ref" firestore:"data_ref"`
	UID              string      `json:"uid" firestore:"uid"`
	Description      string      `json:"description" firestore:"description"`
	Amount           string      `json:"amount" firestore:"amount"`
	FiatAmountUSD    string      `json:"fiat_amount_usd" firestore:"fiat_amount_usd"`
	Currency         string      `json:"currency" firestore:"currency"`
	TxHash           string      `json:"tx_hash" firestore:"tx_hash"`
}

func (log CryptoTransferLog) GetAddLog() map[string]interface{} {
	return map[string]interface{}{
		"id":                log.Id,
		"provider":          log.Provider,
		"provider_response": log.ProviderResponse,
		"external_id":       log.ExternalId,
		"data_type":         log.DataType,
		"data_ref":          log.DataRef,
		"uid":               log.UID,
		"description":       log.Description,
		"amount":            log.Amount,
		"fiat_amount_usd":   log.FiatAmountUSD,
		"currency":          log.Currency,
		"tx_hash":           log.TxHash,
		"created_at":        firestore.ServerTimestamp,
	}
}

type CryptoPendingTransfer struct {
	Id            string `json:"id" firestore:"id"`
	Provider      string `json:"provider" firestore:"provider"`
	ExternalId    string `json:"external_id" firestore:"external_id"`
	TxHash        string `json:"tx_hash" firestore:"tx_hash"`
	DataType      string `json:"data_type" firestore:"data_type"`
	DataRef       string `json:"data_ref" firestore:"data_ref"`
	UID           string `json:"uid" firestore:"uid"`
	Amount        string `json:"amount" firestore:"amount"`
	FiatAmountUSD string `json:"fiat_amount_usd" firestore:"fiat_amount_usd"`
	Currency      string `json:"currency" firestore:"currency"`
}

func (transfer CryptoPendingTransfer) GetAddCryptoPendingTransfer() map[string]interface{} {
	return map[string]interface{}{
		"id":              transfer.Id,
		"provider":        transfer.Provider,
		"external_id":     transfer.ExternalId,
		"tx_hash":         transfer.TxHash,
		"data_type":       transfer.DataType,
		"data_ref":        transfer.DataRef,
		"uid":             transfer.UID,
		"amount":          transfer.Amount,
		"fiat_amount_usd": transfer.FiatAmountUSD,
		"currency":        transfer.Currency,
		"created_at":      firestore.ServerTimestamp,
	}
}
