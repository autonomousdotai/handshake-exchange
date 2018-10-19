package bean

const CURRENCY_CRYPTO = "crypto"
const CURRENCY_FIAT = "fiat"

type Currency struct {
	Name    string
	Code    string
	Type    string
	Decimal int32
}

var BTC = Currency{
	"Bitcoin",
	"BTC",
	CURRENCY_CRYPTO,
	8,
}

var ETH = Currency{
	"Ethereum",
	"ETH",
	CURRENCY_CRYPTO,
	18,
}

var BCH = Currency{
	"Bitcoin Cash",
	"BCH",
	CURRENCY_CRYPTO,
	8,
}

var LTC = Currency{
	"Litecoin",
	"LTC",
	CURRENCY_CRYPTO,
	8,
}

var XRP = Currency{
	"Ripple",
	"XRP",
	CURRENCY_CRYPTO,
	8,
}

var USD = Currency{
	"US Dollar",
	"USD",
	CURRENCY_FIAT,
	2,
}

var HKD = Currency{
	"Hong Kong Dollar",
	"HKD",
	CURRENCY_FIAT,
	2,
}

var VND = Currency{
	"Vietnam Dong",
	"VND",
	CURRENCY_FIAT,
	0,
}

var CurrencyMapping = map[string]Currency{
	USD.Code: USD,
	// HKD.Code: HKD,

	BTC.Code: BTC,
	// BCH.Code: BCH,
	ETH.Code: ETH,
	// LTC.Code: LTC,
}
