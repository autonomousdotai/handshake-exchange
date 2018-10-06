package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/adyen_service"
	"time"
)

type InternalApi struct {
}

func (api InternalApi) Redeem(context *gin.Context) {

	bean.SuccessResponse(context, bean.Redeem{})
}

func (api InternalApi) Payment(context *gin.Context) {
	var body adyen_service.AdyenAuthorise
	if common.ValidateBody(context, &body) != nil {
		return
	}
	if body.Reference == "" {
		body.Reference = fmt.Sprintf("%d", time.Now().UTC().Unix())
	}
	resp, err := adyen_service.AuthoriseNo3D(body)
	if err != nil {
		api_error.AbortWithError(context, err)
	}

	if resp.ResultCode == "Authorised" {
		adyen_service.Capture(adyen_service.AdyenCapture{
			OriginalReference:  resp.PSPReference,
			ModificationAmount: body.Amount,
			Reference:          body.Reference,
		})
	} else {
		api_error.AbortWithError(context, errors.New(resp.ResultCode))
	}

	bean.SuccessResponse(context, resp)
}
