package bean

import (
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	StatusCode int         `json:"status"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type BasePagingResponse struct {
	StatusCode int         `json:"status"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Page       interface{} `json:"page"`
	CanMove    bool        `json:"can_move"`
}

func DefaultSuccessResponse(context *gin.Context) {
	context.JSON(http.StatusOK, BaseResponse{
		http.StatusOK,
		api_error.CodeMessage[api_error.Success].Code,
		"Success",
		nil})
}

func SuccessResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, BaseResponse{
		http.StatusOK,
		api_error.CodeMessage[api_error.Success].Code,
		"Success",
		data})
	context.Set("ResponseData", data)
}

func SuccessPagingResponse(context *gin.Context, data interface{}, canMove bool, nextAt interface{}) {
	context.JSON(http.StatusOK, BasePagingResponse{
		http.StatusOK,
		api_error.CodeMessage[api_error.Success].Code,
		"Success",
		data,
		nextAt,
		canMove,
	})
}

func CustomSuccessResponse(context *gin.Context, statusCode int, code int, message string, data interface{}) {
	context.JSON(http.StatusOK, BaseResponse{
		statusCode,
		code,
		message,
		data})
	context.Set("ResponseData", data)
}

type Paging interface {
	GetPageValue() interface{}
}
