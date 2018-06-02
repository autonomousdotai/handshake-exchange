package bean

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

const CONFIG_KEY_CC_MODE = "CC_MODE"
const CC_MODE_GDAX = "gdax"
const CC_MODE_INVENTORY = "inventory"

const CONFIG_BTC_WALLET = "BTC_WALLET"
const BTC_WALLET_COINBASE = "coinbase"
const BTC_WALLET_BLOCKCHAINIO = "blockchainio"

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
