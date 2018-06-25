package notification

import (
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/ninjadotorg/handshake-exchange/service/email"
)

func SendOfferNotification(offer bean.Offer) []error {
	c := make(chan error)
	go SendOfferToEmail(offer, c)
	go SendOfferToFirebase(offer, c)
	go SendOfferToSolr(offer, c)
	go SendOfferToFCM(offer, c)

	return []error{<-c, <-c, <-c, <-c}
}

func SendInstantOfferNotification(offer bean.InstantOffer) []error {
	c := make(chan error)
	go SendInstantOfferToEmail(offer, c)
	go SendInstantOfferToFirebase(offer, c)
	go SendInstantOfferToSolr(offer, c)
	go SendInstantOfferToFCM(offer, c)

	return []error{<-c, <-c, <-c, <-c}
}

func SendOfferToEmail(offer bean.Offer, c chan error) {
	var err error
	username := offer.Email
	if username == "" {
		username = offer.ContactPhone
	}
	// toUsername := offer.ToEmail

	// coinUsername := toUsername
	// cashEmail := offer.Email
	// coinEmail := offer.ToEmail
	if offer.Type == bean.OFFER_TYPE_BUY {
		// coinUsername = username
		// cashEmail = offer.ToEmail
		// coinEmail = offer.Email
	}

	if offer.Status == bean.OFFER_STATUS_ACTIVE {
		if offer.Email != "" {
			if offer.Type == bean.OFFER_TYPE_BUY {
				err = email.SendOfferBuyingActiveEmail(offer.Language, offer.Email)
			} else {
				err = email.SendOfferSellingActiveEmail(offer.Language, offer.Email)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_CLOSED {
		//if offer.Email != "" {
		//	err = email.SendOfferClosedEmail(offer.Language, offer.Email)
		//}
	} else if offer.Status == bean.OFFER_STATUS_SHAKE {
		if offer.Email != "" {
			if offer.IsTypeSell() {
				err = email.SendOfferMakerSellShakeEmail(offer.Language, offer.Email)
			} else {
				err = email.SendOfferMakerBuyShakeEmail(offer.Language, offer.Email)
			}
		}
		if offer.ToEmail != "" {
			if offer.IsTypeSell() {
				err = email.SendOfferTakerSellShakeEmail(offer.Language, offer.ToEmail)
			} else {
				err = email.SendOfferTakerBuyShakeEmail(offer.Language, offer.ToEmail)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_REJECTED {
		if offer.UID == offer.ActionUID {
			if offer.Email != "" {
				err = email.SendOfferMakerMakerRejectEmail(offer.Language, offer.Email)
			}
			if offer.ToEmail != "" {
				err = email.SendOfferTakerMakerRejectEmail(offer.Language, offer.ToEmail)
			}
		} else {
			if offer.Email != "" {
				err = email.SendOfferMakerTakerRejectEmail(offer.Language, offer.Email)
			}
			if offer.ToEmail != "" {
				err = email.SendOfferTakerTakerRejectEmail(offer.Language, offer.ToEmail)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
		if offer.IsTypeSell() {
			if offer.Email != "" {
				err = email.SendOfferSellCompleteEmail(offer.Language, offer.Email)
			}
			if offer.ToEmail != "" {
				err = email.SendOfferBuyCompleteEmail(offer.Language, offer.ToEmail)
			}
		} else {
			if offer.Email != "" {
				err = email.SendOfferBuyCompleteEmail(offer.Language, offer.Email)
			}
			if offer.ToEmail != "" {
				err = email.SendOfferSellCompleteEmail(offer.Language, offer.ToEmail)
			}
		}
	}
	c <- err
}

func SendOfferToFirebase(offer bean.Offer, c chan error) {
	err := dao.OfferDaoInst.UpdateNotificationOffer(offer)
	c <- err
}

func SendOfferToSolr(offer bean.Offer, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromOffer(offer))
	c <- err
}

func SendOfferToFCM(offer bean.Offer, c chan error) {
	var err error

	if offer.Status == bean.OFFER_STATUS_ACTIVE {
		// Not yet
	} else if offer.Status == bean.OFFER_STATUS_CLOSED {
		// Not yet
	} else if offer.Status == bean.OFFER_STATUS_SHAKE {
		if offer.FCM != "" {
			if offer.IsTypeSell() {
				err = SendOfferMakerSellShakeFCM(offer.Language, offer.FCM)
			} else {
				err = SendOfferMakerBuyShakeFCM(offer.Language, offer.FCM)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_REJECTED {
		if offer.UID == offer.ActionUID {
			if offer.FCM != "" {
				err = SendOfferMakerMakerRejectFCM(offer.Language, offer.FCM)
			}
			if offer.ToFCM != "" {
				err = SendOfferTakerMakerRejectFCM(offer.Language, offer.FCM)
			}
		} else {
			if offer.FCM != "" {
				err = SendOfferMakerTakerRejectFCM(offer.Language, offer.FCM)
			}
			if offer.ToFCM != "" {
				err = SendOfferTakerTakerRejectFCM(offer.Language, offer.FCM)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
		if offer.IsTypeSell() {
			if offer.FCM != "" {
				err = SendOfferSellCompleteFCM(offer.Language, offer.FCM)
			}
			if offer.ToFCM != "" {
				err = SendOfferBuyCompleteFCM(offer.Language, offer.ToFCM)
			}
		} else {
			if offer.FCM != "" {
				err = SendOfferBuyCompleteFCM(offer.Language, offer.FCM)
			}
			if offer.ToFCM != "" {
				err = SendOfferSellCompleteFCM(offer.Language, offer.ToFCM)
			}
		}
	}
	c <- err
}

func SendInstantOfferToEmail(offer bean.InstantOffer, c chan error) {
	var err error
	if offer.Status == bean.INSTANT_OFFER_STATUS_SUCCESS {
		if offer.Email != "" {
			err = email.SendOrderInstantCCSuccessEmail(offer.Language, offer.Email, offer.Amount, offer.Currency)
		}
	}
	c <- err
}

func SendInstantOfferToFirebase(offer bean.InstantOffer, c chan error) {
	err := dao.CreditCardDaoInst.UpdateNotificationInstantOffer(offer)
	c <- err
}

func SendInstantOfferToSolr(offer bean.InstantOffer, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromInstantOffer(offer))
	c <- err
}

func SendInstantOfferToFCM(offer bean.InstantOffer, c chan error) {
	var err error
	if offer.Status == bean.INSTANT_OFFER_STATUS_SUCCESS {
		if offer.FCM != "" {
			err = SendOrderInstantCCSuccessFCM(offer.Language, offer.FCM)
		}
	}
	c <- err
}

func SendOfferStoreNotification(offer bean.OfferStore, offerItem bean.OfferStoreItem) []error {
	c := make(chan error)
	go SendOfferStoreToEmail(offer, offerItem, c)
	go SendOfferStoreToFirebase(offer, offerItem, c)
	go SendOfferStoreToSolr(offer, offerItem, c)
	go SendOfferStoreToFCM(offer, offerItem, c)

	return []error{<-c, <-c, <-c, <-c}
}

func SendOfferStoreToEmail(offer bean.OfferStore, offerItem bean.OfferStoreItem, c chan error) {
	var err error
	if offer.Email != "" {
		if offerItem.Status == bean.OFFER_STORE_ITEM_STATUS_ACTIVE {
			err = email.SendOfferStoreItemAddedEmail(offer.Language, offer.Email, offerItem.SellAmount, offerItem.BuyAmount, offerItem.Currency)
		} else if offerItem.Status == bean.OFFER_STORE_ITEM_STATUS_CLOSED {
			err = email.SendOfferStoreItemRemovedEmail(offer.Language, offer.Email)
		}
	}
	c <- err
}

func SendOfferStoreToFirebase(offer bean.OfferStore, offerItem bean.OfferStoreItem, c chan error) {
	err := dao.OfferStoreDaoInst.UpdateNotificationOfferStore(offer, offerItem)
	c <- err
}

func SendOfferStoreToSolr(offer bean.OfferStore, offerItem bean.OfferStoreItem, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer, offerItem))
	c <- err
}

func SendOfferStoreToFCM(offer bean.OfferStore, offerItem bean.OfferStoreItem, c chan error) {
	var err error
	if offer.Email != "" {
		if offerItem.Status == bean.OFFER_STORE_ITEM_STATUS_ACTIVE {
			SendOfferStoreItemAddedFCM(offer.Language, offer.FCM)
		} else if offerItem.Status == bean.OFFER_STORE_ITEM_STATUS_CLOSED {
		}
	}
	c <- err
}

func SendOfferStoreShakeNotification(offer bean.OfferStoreShake, offerStore bean.OfferStore) []error {
	c := make(chan error)
	go SendOfferStoreShakeToEmail(offer, offerStore, c)
	go SendOfferStoreShakeToFirebase(offer, offerStore, c)
	go SendOfferStoreShakeToSolr(offer, offerStore, c)
	go SendOfferStoreShakeToFCM(offer, offerStore, c)

	return []error{<-c, <-c, <-c, <-c}
}

func SendOfferStoreShakeToEmail(offer bean.OfferStoreShake, offerStore bean.OfferStore, c chan error) {
	var err error

	username := offerStore.Username
	if username == "" {
		username = offerStore.Email
		if username == "" {
			username = offerStore.ContactPhone
		}
	}
	toUsername := offer.Username
	if toUsername == "" {
		toUsername = offer.Email
		if toUsername == "" {
			toUsername = offer.ContactPhone
		}
	}

	if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		if offer.Type == bean.OFFER_TYPE_BUY {
			err = email.SendOfferStoreMakerBuyShakeEmail(offerStore.Language, offerStore.Email, offer.Amount, offer.Currency, offer.FiatAmount, offer.FiatCurrency, toUsername)
			err = email.SendOfferStoreTakerBuyShakeEmail(offer.Language, offer.Email, offer.Amount, offer.Currency, offer.FiatAmount, offer.FiatCurrency, username)
		} else {
			err = email.SendOfferStoreMakerSellShakeEmail(offerStore.Language, offerStore.Email, offer.Amount, offer.Currency, offer.FiatAmount, offer.FiatCurrency, toUsername)
			err = email.SendOfferStoreTakerSellShakeEmail(offer.Language, offer.Email, offer.Amount, offer.Currency, offer.FiatAmount, offer.FiatCurrency, username)
		}
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_CANCELLED {

	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
		if offer.ActionUID == offer.UID {
			err = email.SendOfferStoreMakerRejectEmail(offerStore.Language, offerStore.Email, toUsername)
		} else {
			err = email.SendOfferStoreTakerRejectEmail(offer.Language, offer.Email, username)
		}
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		if offer.Type == bean.OFFER_TYPE_BUY {
			err = email.SendOfferStoreMakerCompleteEmail(offerStore.Language, offerStore.Email, offer.Amount, offer.Currency, username)
		} else {
			err = email.SendOfferStoreTakerCompleteEmail(offer.Language, offer.Email, offer.Amount, offer.Currency, toUsername, username, offerStore.Id, offer.Id)
		}
	}
	c <- err
}

func SendOfferStoreShakeToFirebase(offer bean.OfferStoreShake, offerStore bean.OfferStore, c chan error) {
	err := dao.OfferStoreDaoInst.UpdateNotificationOfferStoreShake(offer, offerStore)
	c <- err
}

func SendOfferStoreShakeToSolr(offer bean.OfferStoreShake, offerStore bean.OfferStore, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromOfferStoreShake(offer, offerStore))
	c <- err
}

func SendOfferStoreShakeToFCM(offer bean.OfferStoreShake, offerStore bean.OfferStore, c chan error) {
	var err error

	username := offerStore.Username
	if username == "" {
		username = offerStore.Email
		if username == "" {
			username = offerStore.ContactPhone
		}
	}
	toUsername := offer.Username
	if toUsername == "" {
		toUsername = offer.Email
		if toUsername == "" {
			toUsername = offer.ContactPhone
		}
	}

	if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_PRE_SHAKE {
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_SHAKE {
		if offer.Type == bean.OFFER_TYPE_BUY {
			err = SendOfferStoreMakerBuyShakeFCM(offerStore.Language, offerStore.FCM, offerStore.ChatUsername)
			err = SendOfferStoreTakerBuyShakeFCM(offer.Language, offer.FCM, offer.Currency, offer.ChatUsername)
		} else {
			err = SendOfferStoreMakerSellShakeFCM(offerStore.Language, offerStore.FCM, offerStore.ChatUsername)
			err = SendOfferStoreTakerSellShakeFCM(offer.Language, offer.FCM, offer.Currency, offer.ChatUsername)
		}
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_CANCELLED {

	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_REJECTED {
		if offer.ActionUID == offer.UID {
			err = SendOfferStoreMakerRejectFCM(offerStore.Language, offerStore.FCM, toUsername)
		} else {
			err = SendOfferStoreTakerRejectFCM(offer.Language, offer.FCM, username)
		}
	} else if offer.Status == bean.OFFER_STORE_SHAKE_STATUS_COMPLETED {
		if offer.Type == bean.OFFER_TYPE_BUY {
			err = SendOfferStoreMakerCompleteFCM(offerStore.Language, offerStore.FCM, offer.Currency)
		} else {
			err = SendOfferStoreTakerCompleteFCM(offer.Language, offer.FCM, offer.Currency, offerStore.Id, offer.Id)
		}
	}
	c <- err
}
