package common

import (
	"github.com/go-playground/validator"
	"github.com/shopspring/decimal"
)

var DataValidator = NewValidator()

func NewValidator() *validator.Validate {
	return validator.New()
}

var Zero = decimal.NewFromFloat(0)
var NegativeOne = decimal.NewFromFloat(-1)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
