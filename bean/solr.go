package bean

import (
	"fmt"
	"strconv"
)

type SolrOfferObject struct {
	Id            string   `json:"id"`
	Type          int      `json:"type_i"`
	State         int      `json:"state_i"`
	Status        int      `json:"status_i"`
	Hid           string   `json:"hid_s"`
	InitUserId    int      `json:"init_user_id_i"`
	ShakedUserIds []int    `json:"shaked_user_ids_is"`
	ShakeCount    int      `json:"shake_count_i"`
	ViewCount     int      `json:"view_count_i"`
	CommentCount  int      `json:"comment_count_i"`
	TextSearch    []string `json:"text_search_ss"`

	OfferId           string `json:"offer_id_s"`
	OfferType         string `json:"offer_type_p"`
	OfferAmount       string `json:"offer_amount_s"`
	OfferCurrency     string `json:"offer_currency_s"`
	OfferFiatCurrency string `json:"offer_fiat_currency_s"`
	OfferFiatAmount   string `json:"offer_fiat_amount_s"`
	OfferPrice        string `json:"offer_price_s"`
	OfferStatus       string `json:"offer_status_s"`
	OfferLocation     string `json:"offer_location_p"`
}

var statusMap = map[string]int{
	OFFER_STATUS_CREATED: 0,
	OFFER_STATUS_ACTIVE:  1,
}

func NewSolrFromOffer(offer Offer) (solr SolrOfferObject) {
	solr.Id = fmt.Sprintf("exchange_%s", offer.Id)
	solr.Type = 2
	if offer.Status == OFFER_STATUS_CREATED {
		solr.State = 0
	} else {
		solr.State = 1
	}
	solr.Status = statusMap[offer.Status]
	solr.Hid = ""
	userId, _ := strconv.Atoi(offer.UID)
	solr.InitUserId = userId
	if offer.ToUID != "" {
		userId, _ := strconv.Atoi(offer.ToUID)
		solr.ShakedUserIds = []int{userId}
	}
	solr.OfferId = offer.Id
	solr.OfferType = offer.Type
	solr.OfferAmount = offer.Amount
	solr.OfferCurrency = offer.Currency
	solr.OfferFiatAmount = offer.FiatAmount
	solr.OfferFiatCurrency = offer.FiatCurrency
	solr.OfferPrice = offer.Price
	solr.OfferStatus = offer.Status
	solr.OfferLocation = fmt.Sprintf("%f,%f", offer.Latitude, offer.Longitude)

	return
}
