package common

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

func ValidateBody(context *gin.Context, body interface{}) error {
	err := context.BindJSON(body)
	if api_error.PropagateErrorAndAbort(context, api_error.InvalidRequestBody, err) != nil {
		return err
	}

	// Validate data
	err = api_error.AbortWithRequestBodyError(context, DataValidator.Struct(body))
	if err != nil {
		return err
	}

	return nil
}

// If there is not found error, error will be returns nil
func CheckNotFound(err error) error {
	if strings.Contains(fmt.Sprintf("%s", err), "code = NotFound") {
		err = nil
	}

	return err
}

func ConvertToDecimal(doc *firestore.DocumentSnapshot, field string) (decimal.Decimal, error) {
	zero := decimal.NewFromFloat(0)
	value, err := doc.DataAt(field)
	if err != nil {
		return zero, err
	}
	valueStr := value.(string)
	realValue := zero
	if valueStr != "" {
		realValue, err = decimal.NewFromString(valueStr)
		if err != nil {
			return zero, err
		}
	}

	return realValue, nil
}

func CheckETHAddress(s string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(s)
}
