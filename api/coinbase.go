package api

import (
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/integration/coinbase_service"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
	"github.com/autonomousdotai/handshake-exchange/service"
	"github.com/gin-gonic/gin"
)

type CoinbaseApi struct {
}

func (api CoinbaseApi) ReceiveCallback(context *gin.Context) {
	var body bean.CoinbaseNotification
	err := context.BindJSON(&body)
	if api_error.PropagateErrorAndAbort(context, api_error.InvalidRequestBody, err) != nil {
		return
	}

	if body.Type == "wallet:addresses:new-payment" {
		bodyNotification, err := coinbase_service.GetNotification(body.ResourcePath)
		if api_error.PropagateErrorAndAbort(context, api_error.GetDataFailed, err) != nil {
			// This might be fake coinbase request
			return
		}

		// Do some double check
		if body.Id != bodyNotification.Id || body.ResourcePath != bodyNotification.ResourcePath {
			// This might be fake coinbase request
			api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
			return
		}

		// Get Firestore client
		dbClient := firebase_service.FirestoreClient

		// ATTENTION: Use bodyNotification which is get from coinbase instead of body
		// To make sure the data is real
		coinbaseItemPath := getCoinbaseItemPath(bodyNotification.Id)
		coinbaseRef := dbClient.Doc(coinbaseItemPath)
		_, err = coinbaseRef.Get(context)
		if err != nil {
			// This mean this ID does not exist
			_, err = coinbaseRef.Set(context, bodyNotification)
		} else {
			// Already receive this transaction
			bean.SuccessResponse(context, bodyNotification)
			return
		}

		var address string

		additionalData := bodyNotification.AdditionalData
		if amountNode, ok := additionalData["amount"]; ok {
			amountNodeMap := amountNode.(map[string]interface{})
			amountObj, ok1 := amountNodeMap["amount"]
			currencyObj, ok2 := amountNodeMap["currency"]
			if ok1 && ok2 && amountObj.(string) != "" && currencyObj.(string) != "" {
				if err != nil {
					// Data from Coinbase is not valid
					api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
					return
				}

				data := bodyNotification.Data
				if addressObj, ok := data["address"]; !ok || addressObj.(string) != "" {
					address = addressObj.(string)

					offer, ce := service.OfferServiceInst.ActiveOffer(address, amountObj.(string))
					if ce.HasError() {
						if ce.StatusKey == api_error.OfferStatusInvalid {
							_, ce = service.OfferServiceInst.UpdateShakeOffer(offer)
							if ce.HasError() {
								// TODO Need to do some notification if get error
							}
						} else {
							// TODO Need to do some notification if get error
						}
					}
				} else {
					// Data from Coinbase is not valid
					api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
					return
				}
			} else {
				// Data from Coinbase is not valid
				api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
				return
			}
		} else {
			// Data from Coinbase is not valid
			api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
			return
		}
	}

	// Response
	bean.SuccessResponse(context, body)
}

func getCoinbaseItemPath(id string) string {
	return fmt.Sprintf("coinbase/%s", id)
}
