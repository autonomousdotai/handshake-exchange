package bean

import "cloud.google.com/go/firestore"

type OfferEventBlock struct {
	LastBlock int64 `json:"last_block" firestore:"last_block"`
}

func (offer OfferEventBlock) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"last_block": offer.LastBlock,
		"updated_at": firestore.ServerTimestamp,
	}
}
