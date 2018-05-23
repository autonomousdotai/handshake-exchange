package api_error

import "net/http"

const Success = "Success"
const UnexpectedError = "UnexpectedError"
const ResourceNotFound = "ResourceNotFound"
const TokenInvalid = "TokenInvalid"
const FirebaseError = "FirebaseError"
const SendEmailError = "SendEmailError"
const InvalidRequestBody = "InvalidRequestBody"
const InvalidRequestParam = "InvalidRequestParam"
const InvalidQueryParam = "InvalidQueryParam"
const ExternalApiFailed = "ExternalApiFailed"
const InvalidNumber = "InvalidNumber"

const GetDataFailed = "GetDataFailed"
const AddDataFailed = "AddDataFailed"
const UpdateDataFailed = "UpdateDataFailed"
const DeleteDataFailed = "DeleteDataFailed"

const ProfileExists = "ProfileExists"
const ProfileNotExist = "ProfileNotExist"

var CodeMessage = map[string]struct {
	StatusCode int
	Code       int
	Message    string
}{
	Success:             {http.StatusOK, 1, "Success"},
	UnexpectedError:     {http.StatusInternalServerError, -1, "Unexpected error"},
	ResourceNotFound:    {http.StatusNotFound, -1, "Resource not found"},
	FirebaseError:       {http.StatusInternalServerError, -1, "Unexpected error"},
	SendEmailError:      {http.StatusInternalServerError, -1, "Unexpected error"},
	TokenInvalid:        {http.StatusUnauthorized, -3, "Token is invalid"},
	InvalidRequestBody:  {http.StatusBadRequest, -4, "Request body is invalid"},
	InvalidRequestParam: {http.StatusBadRequest, -5, "Request param is invalid"},
	InvalidQueryParam:   {http.StatusBadRequest, -6, "Query param is invalid"},
	ExternalApiFailed:   {http.StatusBadRequest, -7, "External API failed"},
	InvalidNumber:       {http.StatusBadRequest, -8, "Invalid number"},

	GetDataFailed:    {http.StatusBadRequest, -201, "Get data failed"},
	AddDataFailed:    {http.StatusBadRequest, -202, "Add data failed"},
	UpdateDataFailed: {http.StatusBadRequest, -203, "Update data failed"},
	DeleteDataFailed: {http.StatusBadRequest, -204, "Delete data failed"},

	ProfileExists:   {http.StatusBadRequest, -205, "Profile exists"},
	ProfileNotExist: {http.StatusBadRequest, -206, "Profile not exist"},
}

//var ErrorSuccess = NewErrorSimple(Success)
//var ErrorUnexpected = NewErrorSimple(UnexpectedError)
//var ErrorResourceNotFound = NewErrorSimple(ResourceNotFound)
//var ErrorFirebase = NewErrorSimple(FirebaseError)
//var ErrorSendEmail = NewErrorSimple(SendEmailError)
//var ErrorTokenInvalid = NewErrorSimple(TokenInvalid)
//var ErrorInvalidRequestBody = NewErrorSimple(InvalidRequestBody)
//var ErrorInvalidRequestParam = NewErrorSimple(InvalidRequestParam)
//var ErrorInvalidQueryParam = NewErrorSimple(InvalidQueryParam)
//var ErrorExternalApiFailed = NewErrorSimple(ExternalApiFailed)
//var ErrorInvalidNumber = NewErrorSimple(InvalidNumber)
//var ErrorGetDataFailed = NewErrorSimple(GetDataFailed)
//var ErrorAddDataFailed = NewErrorSimple(AddDataFailed)
//var ErrorUpdateDataFailed = NewErrorSimple(UpdateDataFailed)
//var ErrorDeleteDataFailed = NewErrorSimple(DeleteDataFailed)
