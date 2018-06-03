package notification

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/service/email"
)

func SendOfferNotification(offer bean.Offer) error {
	return nil
}

func SendInstantOfferNotification(offer bean.InstantOffer) error {
	return nil
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

func SendInstantOfferEmailNotification(language string, offer bean.InstantOffer) error {
	email.SendOrderInstantCCSuccessEmail(language, offer.Email, offer.Amount, offer.Currency)
	return nil
}

func SendInstantOfferFirebaseNotification() error {
	return nil
}
