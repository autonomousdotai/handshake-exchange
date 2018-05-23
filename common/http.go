package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetHeaderWithDefault(context *gin.Context, key string, defaultValue string) string {
	value := context.GetHeader(key)
	if value == "" {
		value = defaultValue
	}

	return value
}

func GetLanguage(context *gin.Context) string {
	return GetHeaderWithDefault(context, "Custom-Language", "en-US")
}

func GetUserId(context *gin.Context) string {
	return GetHeaderWithDefault(context, "Custom-Userid", "")
}

func GetCurrency(context *gin.Context) string {
	return GetHeaderWithDefault(context, "Custom-Currency", "USD")
}

func ExtractTimePagingParams(context *gin.Context) (interface{}, int) {
	startAtStr := context.DefaultQuery("page", "")
	limitStr := context.DefaultQuery("limit", "10")

	startAt := interface{}(nil)
	if startAtStr != "" {
		startAt, _ = time.Parse(time.RFC3339, startAtStr)
	}
	limit, _ := strconv.Atoi(limitStr)

	return startAt, limit
}
