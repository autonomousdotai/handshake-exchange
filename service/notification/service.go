package notification

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/service/email"
)

func SendOfferNotification(offer bean.Offer) error {
	return nil
}

func SendInstantOfferNotification(language string, offer bean.InstantOffer) error {
	c := make(chan error)
	go SendInstantOfferEmailNotification(language, offer, c)
	go SendInstantOfferFirebaseNotification(offer, c)

	e1, e2 := <-c, <-c
	err := e1
	if err == nil {
		err = e2
	}

	return err
}

func SendOfferEmailNotification(offer bean.Offer) error {
	if offer.Status == bean.OFFER_STATUS_ACTIVE {

	} else if offer.Status == bean.OFFER_STATUS_SHAKE {

	}
	return nil
}

func SendOfferFirebaseNotification() error {
	return nil
}

func SendInstantOfferEmailNotification(language string, offer bean.InstantOffer, c chan error) {
	err := email.SendOrderInstantCCSuccessEmail(language, offer.Email, offer.Amount, offer.Currency)
	c <- err
}

func SendInstantOfferFirebaseNotification(offer bean.InstantOffer, c chan error) {
	c <- nil
}
