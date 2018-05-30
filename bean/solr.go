package bean

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type SolrOfferObject struct {
	Id           string   `json:"id"`
	Type         int      `json:"type_i"`
	State        int      `json:"state_i"`
	Status       int      `json:"status_i"`
	Hid          string   `json:"hid_s"`
	IsPrivate    int      `json:"is_private_i"`
	InitUserId   int      `json:"init_user_id_i"`
	ChainId      int64    `json:"chain_id_i"`
	ShakeUserIds []int    `json:"shake_user_ids_is"`
	ShakeCount   int      `json:"shake_count_i"`
	ViewCount    int      `json:"view_count_i"`
	CommentCount int      `json:"comment_count_i"`
	TextSearch   []string `json:"text_search_ss"`
	ExtraData    string   `json:"extra_data_s"`
	Location     string   `json:"location_p"`
	InitAt       int64    `json:"init_at_i"`
	LastUpdateAt int64    `json:"last_update_at_i"`
}

type SolrOfferExtraData struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	FiatCurrency string `json:"fiat_currency"`
	FiatAmount   string `json:"fiat_amount"`
	Price        string `json:"price"`
	Percentage   string `json:"percentage"`
	ContactPhone string `json:"contact_phone"`
	ContactInfo  string `json:"contact_info"`
	Status       string `json:"status"`
	Success      int64  `json:"success"`
	Failed       int64  `json:"failed"`
}

var offerStatusMap = map[string]int{
	OFFER_STATUS_CREATED:   0,
	OFFER_STATUS_ACTIVE:    1,
	OFFER_STATUS_CLOSED:    2,
	OFFER_STATUS_SHAKING:   3,
	OFFER_STATUS_SHAKE:     4,
	OFFER_STATUS_COMPLETED: 5,
	OFFER_STATUS_WITHDRAW:  6,
}

type SolrInstantOfferExtraData struct {
	Id           string `json:"id"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	FiatCurrency string `json:"fiat_currency"`
	FiatAmount   string `json:"fiat_amount"`
	Status       string `json:"status"`
}

var instantOfferStatusMap = map[string]int{
	INSTANT_OFFER_STATUS_PROCESSING: 0,
	INSTANT_OFFER_STATUS_SUCCESS:    1,
	INSTANT_OFFER_STATUS_CANCELLED:  2,
}

func NewSolrFromOffer(offer Offer) (solr SolrOfferObject) {
	solr.Id = fmt.Sprintf("exchange_%s", offer.Id)
	solr.Type = 2
	if offer.Status == OFFER_STATUS_ACTIVE {
		solr.State = 1
		solr.IsPrivate = 0
	} else {
		solr.State = 0
		solr.IsPrivate = 1
	}
	solr.Status = offerStatusMap[offer.Status]
	solr.Hid = ""
	solr.ChainId = offer.ChainId
	userId, _ := strconv.Atoi(offer.UID)
	solr.InitUserId = userId
	if offer.ToUID != "" {
		userId, _ := strconv.Atoi(offer.ToUID)
		solr.ShakeUserIds = []int{userId}
	} else {
		solr.ShakeUserIds = make([]int, 0)
	}
	solr.TextSearch = make([]string, 0)
	solr.Location = fmt.Sprintf("%f,%f", offer.Latitude, offer.Longitude)
	solr.InitAt = offer.CreatedAt.Unix()
	solr.LastUpdateAt = time.Now().UTC().Unix()

	percentage, _ := decimal.NewFromString(offer.Percentage)
	extraData := SolrOfferExtraData{
		Id:           offer.Id,
		Type:         offer.Type,
		Amount:       offer.Amount,
		Currency:     offer.Currency,
		FiatAmount:   offer.FiatAmount,
		FiatCurrency: offer.FiatCurrency,
		Price:        offer.Price,
		Percentage:   percentage.Mul(decimal.NewFromFloat(100)).String(),
		ContactInfo:  offer.ContactInfo,
		ContactPhone: offer.ContactPhone,
		Status:       offer.Status,
		Success:      offer.TransactionCount.Success,
		Failed:       offer.TransactionCount.Failed,
	}
	b, _ := json.Marshal(&extraData)
	solr.ExtraData = string(b)

	return
}

func NewSolrFromInstantOffer(offer InstantOffer) (solr SolrOfferObject) {
	solr.Id = fmt.Sprintf("exchange_%s", offer.Id)
	solr.Type = 2
	solr.State = 0
	solr.IsPrivate = 1
	solr.Status = instantOfferStatusMap[offer.Status]
	solr.Hid = ""
	userId, _ := strconv.Atoi(offer.UID)
	solr.InitUserId = userId
	solr.ShakeUserIds = make([]int, 0)
	solr.TextSearch = make([]string, 0)
	solr.InitAt = offer.CreatedAt.Unix()
	solr.LastUpdateAt = time.Now().UTC().Unix()

	extraData := SolrInstantOfferExtraData{
		Id:           offer.Id,
		Amount:       offer.Amount,
		Currency:     offer.Currency,
		FiatAmount:   offer.FiatAmount,
		FiatCurrency: offer.FiatCurrency,
		Status:       offer.Status,
	}
	b, _ := json.Marshal(&extraData)
	solr.ExtraData = string(b)

	return
}
