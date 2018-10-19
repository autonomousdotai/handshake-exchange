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
const InvalidConfig = "InvalidConfig"

const GetDataFailed = "GetDataFailed"
const AddDataFailed = "AddDataFailed"
const UpdateDataFailed = "UpdateDataFailed"
const DeleteDataFailed = "DeleteDataFailed"

const ProfileExists = "ProfileExists"
const ProfileNotExist = "ProfileNotExist"
const CCOverLimit = "CCOverLimit"
const InvalidCC = "InvalidCC"
const UnsupportedOfferType = "UnsupportedOfferType"
const UnsupportedCurrency = "UnsupportedCurrency"
const OfferStatusInvalid = "OfferStatusInvalid"
const OfferPayMyself = "OfferPayMyself"
const TooManyOffer = "TooManyOffer"
const AmountIsTooSmall = "AmountIsTooSmall"
const InvalidAmount = "InvalidAmount"
const InvalidUserToCompleteHandshake = "InvalidUserToCompleteHandshake"
const OfferActionLocked = "OfferActionLocked"
const OfferStoreExists = "OfferStoreExists"
const OfferStoreNotExist = "OfferStoreNotExist"
const OfferStoreShakeNotExist = "OfferStoreShakeNotExist"
const OfferStoreShakeActiveExist = "OfferStoreShakeActiveExist"
const OfferStoreNotEnoughBalance = "OfferStoreNotEnoughBalance"
const OfferStoreAlreadyReviewed = "OfferStoreAlreadyReviewed"
const InvalidFreeStartAmount = "InvalidFreeStartAmount"
const RegisterFreeStartFailed = "RegisterFreeStartFailed"
const ChargeCCFailed = "ChargeCCFailed"
const CCOverGlobalLimit = "CCOverGlobalLimit"

const CreditExists = "CreditExists"
const CreditPriceChanged = "CreditPriceChanged"
const CreditOutOfStock = "CreditOutOfStock"
const CreditItemStatusInvalid = "CreditItemStatusInvalid"

const CashStoreExists = "CashStoreExists"
const CashOrderStatusInvalid = "CashOrderStatusInvalid"

const CoinOrderStatusInvalid = "CoinOrderStatusInvalid"
const CoinOverLimit = "CoinOverLimit"

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
	InvalidConfig:       {http.StatusBadRequest, -9, "Invalid config"},

	GetDataFailed:    {http.StatusBadRequest, -201, "Get data failed"},
	AddDataFailed:    {http.StatusBadRequest, -202, "Add data failed"},
	UpdateDataFailed: {http.StatusBadRequest, -203, "Update data failed"},
	DeleteDataFailed: {http.StatusBadRequest, -204, "Delete data failed"},

	ProfileExists:                  {http.StatusBadRequest, -301, "Profile exists"},
	ProfileNotExist:                {http.StatusBadRequest, -302, "Profile not exist"},
	CCOverLimit:                    {http.StatusBadRequest, -303, "CC over limit"},
	InvalidCC:                      {http.StatusBadRequest, -304, "CC is invalid"},
	UnsupportedCurrency:            {http.StatusBadRequest, -305, "This currency is not supported"},
	UnsupportedOfferType:           {http.StatusBadRequest, -306, "This offer type is not supported"},
	OfferStatusInvalid:             {http.StatusBadRequest, -307, "This offer status is invalid"},
	OfferPayMyself:                 {http.StatusBadRequest, -308, "You cannot pay offer for yourself"},
	TooManyOffer:                   {http.StatusBadRequest, -309, "Too many offer"},
	AmountIsTooSmall:               {http.StatusBadRequest, -310, "Amount is too small"},
	InvalidUserToCompleteHandshake: {http.StatusBadRequest, -311, "Invalid user to complete handshake"},
	OfferActionLocked:              {http.StatusBadRequest, -312, "Your offer action is locked"},
	OfferStoreExists:               {http.StatusBadRequest, -313, "Offer store exists"},
	OfferStoreNotExist:             {http.StatusBadRequest, -314, "Offer store not exist"},
	OfferStoreShakeNotExist:        {http.StatusBadRequest, -315, "Offer store shake not exist"},
	OfferStoreNotEnoughBalance:     {http.StatusBadRequest, -316, "Offer not enough balance"},
	OfferStoreAlreadyReviewed:      {http.StatusBadRequest, -317, "Offer store already reviewed"},
	InvalidAmount:                  {http.StatusBadRequest, -318, "Invalid amount"},
	InvalidFreeStartAmount:         {http.StatusBadRequest, -319, "Invalid free start amount"},
	RegisterFreeStartFailed:        {http.StatusBadRequest, -320, "Register free start failed"},
	OfferStoreShakeActiveExist:     {http.StatusBadRequest, -321, "There is offer store shake active"},
	ChargeCCFailed:                 {http.StatusBadRequest, -322, "Charge CC failed"},
	CCOverGlobalLimit:              {http.StatusBadRequest, -323, "CC over global limit"},
	CreditExists:                   {http.StatusBadRequest, -324, "Credit exists"},
	CreditPriceChanged:             {http.StatusBadRequest, -325, "Credit price changed"},
	CreditOutOfStock:               {http.StatusBadRequest, -326, "Credit out of stock"},
	CreditItemStatusInvalid:        {http.StatusBadRequest, -327, "Credit item status is invalid"},
	CashStoreExists:                {http.StatusBadRequest, -328, "Cash store exists"},
	CashOrderStatusInvalid:         {http.StatusBadRequest, -329, "Cash order status invalid"},
	CoinOrderStatusInvalid:         {http.StatusBadRequest, -330, "Coin order status invalid"},
	CoinOverLimit:                  {http.StatusBadRequest, -331, "Coin over limit"},
}
