package notification

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/integration/solr_service"
	"github.com/autonomousdotai/handshake-exchange/service/email"
)

func SendOfferNotification(offer bean.Offer) error {
	return nil
}

func SendInstantOfferNotification(offer bean.InstantOffer) []error {
	c := make(chan error)
	go SendInstantOfferToEmail(offer, c)
	go SendInstantOfferToFirebase(offer, c)
	go SendInstantOfferToSolr(offer, c)

	return []error{<-c, <-c, <-c}
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

func SendInstantOfferToEmail(offer bean.InstantOffer, c chan error) {
	if offer.Status == bean.INSTANT_OFFER_STATUS_SUCCESS {
		err := email.SendOrderInstantCCSuccessEmail(offer.Language, offer.Email, offer.Amount, offer.Currency)
		c <- err
	} else {
		c <- nil
	}
}

func SendInstantOfferToFirebase(offer bean.InstantOffer, c chan error) {
	c <- nil
}

func SendInstantOfferToSolr(offer bean.InstantOffer, c chan error) {
	// Always update
	_, err := solr_service.UpdateObject(bean.NewSolrFromInstantOffer(offer))
	c <- err
}
