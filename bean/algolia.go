package bean

type AlgoliaOfferObject struct {
	Id          string `json:"objectID"`
	Type        string `json:"type"`
	Hid         string `json:"hid"`
	ShakeCount  int    `json:"shake_count"`
	ViewCount   int    `json:"view_count"`
	FeedRank    int    `json:"feed_rank"`
	RelatedUids string `json:"related_uids"`

	OfferId      string `json:"offer_id"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	FiatCurrency string `json:"fiat_currency"`
	FiatAmount   string `json:"fiat_amount"`
	Status       string `json:"status"`
}

type AlgoliaOfferShakeObject struct {
	Id          string `json:"objectID"`
	RelatedUids string `json:"related_uids"`

	FiatAmount string `json:"fiat_amount"`
	Status     string `json:"status"`
}

type AlgoliaOfferUpdateObject struct {
	Id          string `json:"objectID"`
	RelatedUids string `json:"related_uids"`

	Status string `json:"status"`
}
