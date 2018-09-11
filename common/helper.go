package common

import (
	"github.com/shopspring/decimal"
)

func StringToDecimal(value string) decimal.Decimal {
	if value == "" {
		return Zero
	}
	number, _ := decimal.NewFromString(value)
	return number
}

func DecimalToFiatString(value decimal.Decimal) string {
	return value.Round(2).String()
}
