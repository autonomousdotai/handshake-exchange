package bean

type SolrOfferObject struct {
	Id          string `json:"id"`
	Type        string `json:"type_s"`
	Hid         string `json:"hid_s"`
	ShakeCount  int    `json:"shake_count_i"`
	ViewCount   int    `json:"view_count_i"`
	FeedRank    int    `json:"feed_rank_i"`
	RelatedUids string `json:"related_uids_s"`

	OfferId      string `json:"offer_id_s"`
	Amount       string `json:"amount_s"`
	Currency     string `json:"currency_s"`
	FiatCurrency string `json:"fiat_currency_s"`
	FiatAmount   string `json:"fiat_amount_s"`
	Status       string `json:"status_s"`
	Location     string `json:"location_p"`
}
