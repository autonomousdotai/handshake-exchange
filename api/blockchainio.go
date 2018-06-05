package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/ninjadotorg/handshake-exchange/service"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
)

type BlockChainApi struct {
}

func (api BlockChainApi) ReceiveCallback(context *gin.Context) {
	txHash := context.DefaultQuery("transaction_hash", "")
	address := context.DefaultQuery("address", "")
	confirmationsStr := context.DefaultQuery("confirmations", "")
	valueStr := context.DefaultQuery("value", "")
	offerId := context.DefaultQuery("offer", "")

	if txHash == "" || address == "" || confirmationsStr == "" || valueStr == "" || offerId == "" {
		if api_error.AbortWithValidateErrorSimple(context, api_error.InvalidRequestParam) != nil {
			return
		}
	}
	confirmations, fmtErr := strconv.Atoi(confirmationsStr)
	if fmtErr != nil {
		if api_error.AbortWithValidateErrorSimple(context, api_error.InvalidRequestParam) != nil {
			return
		}
	}
	valueDecimal, fmtErr := decimal.NewFromString(valueStr)
	if fmtErr != nil {
		if api_error.AbortWithValidateErrorSimple(context, api_error.InvalidRequestParam) != nil {
			return
		}
	}
	value := valueDecimal.IntPart()

	offerTO := dao.OfferDaoInst.GetOffer(offerId)
	if offerTO.ContextValidate(context) {
		return
	}

	// Get Firestore client
	dbClient := firebase_service.FirestoreClient

	// ATTENTION: Use bodyNotification which is get from coinbase instead of body
	// To make sure the data is real
	coinbaseItemPath := getBlockChainIoItemPath(txHash)
	coinbaseRef := dbClient.Doc(coinbaseItemPath)
	_, err := coinbaseRef.Get(context)
	if err != nil {
		// This mean this ID does not exist
		_, err = coinbaseRef.Set(context, bean.BlockChainIoCallback{
			TxHash:        txHash,
			Address:       address,
			Confirmations: int64(confirmations),
			Value:         value,
		})
	} else {
		// Already receive this transaction
		bean.SuccessResponse(context, "*ok*")
		return
	}

	amount := valueDecimal.Div(decimal.NewFromBigInt(big.NewInt(100000000), 0))

	offer, ce := service.OfferServiceInst.ActiveOffer(address, amount.String())
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

	// Response
	bean.SuccessResponse(context, "*ok*")
}

func getBlockChainIoItemPath(id string) string {
	return fmt.Sprintf("blockchainio/%s", id)
}
