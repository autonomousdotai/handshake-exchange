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
	toUsername := offer.ToEmail

	coinUsername := toUsername
	cashEmail := offer.Email
	coinEmail := offer.ToEmail
	if offer.Type == bean.OFFER_TYPE_BUY {
		coinUsername = username
		cashEmail = offer.ToEmail
		coinEmail = offer.Email
	}

	if offer.Status == bean.OFFER_STATUS_ACTIVE {
		if offer.Email != "" {
			if offer.Type == bean.OFFER_TYPE_BUY {
				err = email.SendOfferBuyingActiveEmail(offer.Language, offer.Email, offer.Currency, offer.Price, offer.FiatCurrency)
			} else {
				err = email.SendOfferSellingActiveEmail(offer.Language, offer.Email, offer.Currency, offer.Price, offer.FiatCurrency)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_CLOSED {
		if offer.Email != "" {
			err = email.SendOfferClosedEmail(offer.Language, offer.Email)
		}
	} else if offer.Status == bean.OFFER_STATUS_SHAKE {
		if offer.Email != "" {
			err = email.SendOfferMakerShakeEmail(offer.Language, offer.Email, toUsername, offer.Amount, offer.Currency, offer.Price, offer.FiatCurrency)
		}
		if offer.ToEmail != "" {
			err = email.SendOfferTakerShakeEmail(offer.Language, offer.ToEmail, username, offer.Amount, offer.Currency, offer.Price, offer.FiatCurrency)
		}
	} else if offer.Status == bean.OFFER_STATUS_REJECTED {
		if offer.UID == offer.ActionUID {
			if offer.ToEmail != "" {
				err = email.SendOfferTakerRejectEmail(offer.Language, offer.ToEmail, username)
			}
		} else {
			if offer.Email != "" {
				err = email.SendOfferMakerRejectEmail(offer.Language, offer.Email, toUsername)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
		if cashEmail != "" {
			err = email.SendOfferCompleteEmail(offer.Language, cashEmail, coinUsername, offer.Amount, offer.Currency)
		}
	} else if offer.Status == bean.OFFER_STATUS_WITHDRAW {
		if coinEmail != "" {
			err = email.SendOfferWithdrawEmail(offer.Language, coinEmail, offer.Amount, offer.Currency)
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
			err = SendOfferMakerShakeFCM(offer.Language, offer.FCM, offer.Type)
		}
		if offer.ToFCM != "" {
			err = SendOfferTakerShakeFCM(offer.ToLanguage, offer.ToFCM, offer.Type)
		}
	} else if offer.Status == bean.OFFER_STATUS_REJECTED {
		if offer.UID == offer.ActionUID {

		} else {
			if offer.FCM != "" {
				err = SendOfferMakerRejectedFCM(offer.Language, offer.FCM, offer.Type)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_COMPLETED {
		if offer.Type == bean.OFFER_TYPE_BUY {
			if offer.FCM != "" {
				err = SendOfferCompletedFCM(offer.Language, offer.FCM)
			}
		} else {
			if offer.ToFCM != "" {
				err = SendOfferCompletedFCM(offer.Language, offer.ToFCM)
			}
		}
	} else if offer.Status == bean.OFFER_STATUS_WITHDRAW {
		// Not yet
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

func SendOfferStoreNotification(offer bean.OfferStore) []error {
	c := make(chan error)
	// go SendOfferToEmail(offer, c)
	// go SendOfferToFirebase(offer, c)
	go SendOfferStoreToSolr(offer, c)
	// go SendOfferToFCM(offer, c)

	// return []error{<-c, <-c, <-c, <-c}
	return []error{<-c}
	// return nil
}

func SendOfferStoreToSolr(offer bean.OfferStore, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromOfferStore(offer))
	c <- err
}

func SendOfferStoreShakeNotification(offer bean.OfferStoreShake, offerStore bean.OfferStore) []error {
	c := make(chan error)
	// go SendOfferToEmail(offer, c)
	// go SendOfferToFirebase(offer, c)
	go SendOfferStoreShakeToSolr(offer, offerStore, c)
	// go SendOfferToFCM(offer, c)

	// return []error{<-c, <-c, <-c, <-c}
	return []error{<-c}
	// return nil
}

func SendOfferStoreShakeToSolr(offer bean.OfferStoreShake, offerStore bean.OfferStore, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromOfferStoreShake(offer, offerStore))
	c <- err
}
