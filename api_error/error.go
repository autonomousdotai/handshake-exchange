package api_error

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/go-playground/validator"
	"net/http"
	"reflect"
)

type ErrorDetail struct {
	Code    string
	Message string
}

type Error struct {
	OrgError   error
	StatusCode int
	Code       int
	Message    string
	StatusKey  string
	Details    []ErrorDetail
}

type StackTraceError struct {
	OrgError *errors.Error
}

func (e Error) Error() string {
	return e.OrgError.(*errors.Error).ErrorStack()
}

func (e StackTraceError) Error() string {
	return e.OrgError.ErrorStack()
}

func (e Error) ToJSON() interface{} {
	details := make([]gin.H, len(e.Details))
	for i, detail := range e.Details {
		details[i] = gin.H{
			"code":    detail.Code,
			"message": detail.Message,
		}
	}

	return gin.H{
		"status":  e.StatusCode,
		"code":    e.Code,
		"message": e.Message,
		"details": details,
	}
}

func AbortWithError(context *gin.Context, e error) {
	switch e.(type) {
	case Error:
		context.Error(e)
		context.AbortWithStatusJSON(e.(Error).StatusCode,
			e.(Error).ToJSON())
	default:
		newError := NewError(UnexpectedError, nil)
		context.Error(newError)
		context.AbortWithStatusJSON(http.StatusInternalServerError,
			newError.(Error).ToJSON())
	}
}

func AbortWithValidateError(context *gin.Context, statusKey string, validatorErrors error) error {
	code := CodeMessage[statusKey].Code
	message := CodeMessage[statusKey].Message

	errorDetails := make([]ErrorDetail, len(validatorErrors.(validator.ValidationErrors)))
	for i, err := range validatorErrors.(validator.ValidationErrors) {
		errorDetails[i] = ErrorDetail{
			Message: err.Field(),
			Code:    err.ActualTag(),
		}
	}

	// The duplicate Error new is for Sentry
	newError := NewErrorCustom(statusKey, message, Error{
		errors.New(errors.Errorf("%v: %v", code, message)),
		CodeMessage[statusKey].StatusCode,
		code,
		message,
		statusKey,
		errorDetails,
	})

	AbortWithError(context, newError)

	return newError
}

func AbortWithValidateErrorSimple(context *gin.Context, statusKey string) error {
	code := CodeMessage[statusKey].Code
	message := CodeMessage[statusKey].Message

	// The duplicate Error new is for Sentry
	newError := NewErrorCustom(statusKey, message, Error{
		errors.New(errors.Errorf("%v: %v", code, message)),
		CodeMessage[statusKey].StatusCode,
		code,
		message,
		statusKey,
		make([]ErrorDetail, 0),
	})

	AbortWithError(context, newError)

	return newError
}

func AbortWithQueryParamError(context *gin.Context, validatorErrors error) error {
	if validatorErrors != nil {
		return AbortWithValidateError(context, InvalidQueryParam, validatorErrors)
	}

	return nil
}

func AbortWithRequestParamError(context *gin.Context, validatorErrors error) error {
	if validatorErrors != nil {
		return AbortWithValidateError(context, InvalidRequestParam, validatorErrors)
	}

	return nil
}

func AbortWithRequestBodyError(context *gin.Context, validatorErrors error) error {
	if validatorErrors != nil {
		return AbortWithValidateError(context, InvalidRequestBody, validatorErrors)
	}

	return nil
}

func AbortNotFound(context *gin.Context) {
	statusKey := ResourceNotFound
	code := CodeMessage[statusKey].Code
	message := CodeMessage[statusKey].Message

	newError := NewErrorCustom(statusKey, message, Error{
		errors.New(errors.Errorf("%v: %v", code, message)),
		CodeMessage[statusKey].StatusCode,
		code,
		message,
		statusKey,
		make([]ErrorDetail, 0),
	})

	AbortWithError(context, newError)
}

func NewError(statusKey string, err error) error {
	return NewErrorCustom(statusKey, "", err)
}

func NewErrorCustom(statusKey string, message string, err error) error {
	code := CodeMessage[statusKey].Code
	if message == "" {
		message = CodeMessage[statusKey].Message
	}

	details := make([]ErrorDetail, 0)
	if err != nil && reflect.TypeOf(err).String() == "api_error.Error" {
		details = err.(Error).Details
	}

	if err == nil {
		err = errors.New(errors.Errorf("%v: %v", code, message))
	} else {
		err = errors.Wrap(err, 1)
	}
	return Error{
		err,
		CodeMessage[statusKey].StatusCode,
		code,
		message,
		statusKey,
		details,
	}
}

func PropagateError(statusKey string, err error) error {
	if err != nil {
		err = NewError(statusKey, err)
	}
	return err
}

func PropagateErrorAndAbort(context *gin.Context, statusKey string, err error) error {
	if err != nil {
		err = PropagateError(statusKey, err)
		AbortWithError(context, err)
	}

	return err
}
