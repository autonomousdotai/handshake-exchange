package bean

import (
	"encoding/json"
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
	ExtraData     string   `json:"extra_data_s"`
	Location      string   `json:"location_p"`
}

type SolrOfferExtraData struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	FiatCurrency string `json:"fiat_currency"`
	FiatAmount   string `json:"fiat_amount"`
	Price        string `json:"price"`
	Status       string `json:"status"`
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
	} else {
		solr.ShakedUserIds = make([]int, 0)
	}
	solr.TextSearch = make([]string, 0)
	solr.Location = fmt.Sprintf("%f,%f", offer.Latitude, offer.Longitude)

	extraData := SolrOfferExtraData{
		Id:           offer.Id,
		Type:         offer.Type,
		Amount:       offer.Amount,
		Currency:     offer.Currency,
		FiatAmount:   offer.FiatAmount,
		FiatCurrency: offer.FiatCurrency,
		Price:        offer.Price,
		Status:       offer.Status,
	}
	b, _ := json.Marshal(&extraData)
	solr.ExtraData = string(b)

	return
}
