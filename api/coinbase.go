package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	// "github.com/ninjadotorg/handshake-exchange/integration/coinbase_service"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/ninjadotorg/handshake-exchange/service"
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
		// bodyNotification, err := coinbase_service.GetNotification(body.ResourcePath)
		bodyNotification := body
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
			// bean.SuccessResponse(context, bodyNotification)
			// return
		}

		var address string

		additionalData := bodyNotification.AdditionalData
		ok := true
		if amountNode := additionalData.Amount; ok {
			amountObj := amountNode.Amount
			currencyObj := amountNode.Currency
			if amountObj != "" && currencyObj != "" {
				if err != nil {
					// Data from Coinbase is not valid
					api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
					return
				}

				data := bodyNotification.Data
				if addressObj, ok := data["address"]; !ok || addressObj.(string) != "" {
					address = addressObj.(string)
					refCodeTO := dao.CoinDaoInst.GetCoinSellingOrderRefCode(address)
					if refCodeTO.ContextValidate(context) {
						return
					}
					refCodeObj := refCodeTO.Object.(bean.CoinOrderRefCode)
					_, _, _ = service.CoinServiceInst.FinishSellingOrder(refCodeObj.Order, amountObj, currencyObj, additionalData.Hash)

					// Do nothing in this case
					//if errSrv.HasError() {
					//	api_error.AbortWithValidateErrorSimple(context, api_error.UpdateDataFailed)
					//	return
					//}
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

func (api CoinbaseApi) ReceiveCallbackOld(context *gin.Context) {
	var body bean.CoinbaseNotification
	err := context.BindJSON(&body)
	if api_error.PropagateErrorAndAbort(context, api_error.InvalidRequestBody, err) != nil {
		return
	}

	if body.Type == "wallet:addresses:new-payment" {
		// bodyNotification, err := coinbase_service.GetNotification(body.ResourcePath)
		bodyNotification := body
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
			// bean.SuccessResponse(context, bodyNotification)
			// return
		}

		var address string

		additionalData := bodyNotification.AdditionalData
		ok := true
		if amountNode := additionalData.Amount; ok {
			amountObj := amountNode.Amount
			currencyObj := amountNode.Currency
			if amountObj != "" && currencyObj != "" {
				if err != nil {
					// Data from Coinbase is not valid
					api_error.AbortWithValidateErrorSimple(context, api_error.GetDataFailed)
					return
				}

				data := bodyNotification.Data
				if addressObj, ok := data["address"]; !ok || addressObj.(string) != "" {
					address = addressObj.(string)
					offerAddrTO := dao.OfferDaoInst.GetOfferAddress(address)
					if offerAddrTO.ContextValidate(context) {
						return
					}
					offerAddr := offerAddrTO.Object.(bean.OfferAddressMap)
					_, _, fiatAmount, _ := service.OfferStoreServiceInst.GetQuote(bean.OFFER_TYPE_BUY, amountObj, currencyObj, bean.USD.Code)
					err := dao.OfferDaoInst.AddOfferConfirmingAddressMap(bean.OfferConfirmingAddressMap{
						UID:        offerAddr.UID,
						Address:    offerAddr.Address,
						Offer:      offerAddr.Offer,
						OfferRef:   offerAddr.OfferRef,
						Type:       offerAddr.Type,
						Amount:     amountObj,
						FiatAmount: fiatAmount.StringFixed(2),
						TxHash:     bodyNotification.AdditionalData.Hash,
						ExternalId: bodyNotification.AdditionalData.Transaction.Id,
						Currency:   currencyObj,
					})
					if err != nil {
						api_error.AbortWithValidateErrorSimple(context, api_error.UpdateDataFailed)
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
	return fmt.Sprintf("coinbase_coin/%s", id)
}
