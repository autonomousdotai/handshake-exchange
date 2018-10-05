package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/bean"
)

type InternalApi struct {
}

func (api InternalApi) Redeem(context *gin.Context) {

	bean.SuccessResponse(context, bean.Redeem{})
}
