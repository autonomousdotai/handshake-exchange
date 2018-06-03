package notification

import (
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/dao"
	"github.com/autonomousdotai/handshake-exchange/integration/solr_service"
	"github.com/autonomousdotai/handshake-exchange/service/email"
)

func SendOfferNotification(offer bean.Offer) []error {
	c := make(chan error)
	go SendOfferToEmail(offer, c)
	go SendOfferToFirebase(offer, c)
	go SendOfferToSolr(offer, c)

	return []error{<-c, <-c, <-c}
}

func SendInstantOfferNotification(offer bean.InstantOffer) []error {
	c := make(chan error)
	go SendInstantOfferToEmail(offer, c)
	go SendInstantOfferToFirebase(offer, c)
	go SendInstantOfferToSolr(offer, c)

	return []error{<-c, <-c, <-c}
}

func SendOfferToEmail(offer bean.Offer, c chan error) {
	var err error
	if offer.Status == bean.OFFER_STATUS_ACTIVE {
		if offer.Type == bean.OFFER_TYPE_BUY {
			err = email.SendBuyingOfferSuccessEmail(offer.Language, offer.Email, offer.Email, offer.Currency, offer.Price)
		} else {
			err = email.SendSellingOfferSuccessEmail(offer.Language, offer.Email, offer.Email, offer.Currency, offer.Price)
		}
	} else if offer.Status == bean.OFFER_STATUS_CLOSED {

	} else if offer.Status == bean.OFFER_STATUS_SHAKE {

	} else if offer.Status == bean.OFFER_STATUS_REJECTED {

	} else if offer.Status == bean.OFFER_STATUS_COMPLETED {

	} else if offer.Status == bean.OFFER_STATUS_WITHDRAW {

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

func SendInstantOfferToEmail(offer bean.InstantOffer, c chan error) {
	if offer.Status == bean.INSTANT_OFFER_STATUS_SUCCESS {
		err := email.SendOrderInstantCCSuccessEmail(offer.Language, offer.Email, offer.Amount, offer.Currency)
		c <- err
	} else {
		c <- nil
	}
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
