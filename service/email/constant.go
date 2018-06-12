package email

const OfferBuyingActive = "OfferBuyingActive"
const OfferSellingActive = "OfferSellingActive"
const OfferClosed = "OfferClosed"
const OfferTakerShake = "OfferTakerShake"
const OfferMakerShake = "OfferMakerShake"
const OfferCompleted = "OfferCompleted"
const OfferTakerRejected = "OfferTakerRejected"
const OfferMakerRejected = "OfferMakerRejected"
const OfferWithdraw = "OfferWithdraw"
const OrderInstantCCSuccess = "OrderInstantCCSuccess"
const OfferStoreItemAdded = "OfferStoreItemAdded"
const OfferStoreItemRemoved = "OfferStoreItemRemoved"
const OfferStoreMakerSellShake = "OfferStoreMakerSellShake"
const OfferStoreMakerBuyShake = "OfferStoreMakerBuyShake"
const OfferStoreTakerSellShake = "OfferStoreTakerSellShake"
const OfferStoreTakerBuyShake = "OfferStoreTakerBuyShake"
const OfferStoreAccept = "OfferStoreAccept"
const OfferStoreMakerReject = "OfferStoreMakerReject"
const OfferStoreTakerReject = "OfferStoreTakerReject"

var TemplateName = map[string]string{
	OfferBuyingActive:        "offer-buying-active-",
	OfferSellingActive:       "offer-selling-active-",
	OfferClosed:              "offer-closed-",
	OrderInstantCCSuccess:    "order-instant-cc-success-",
	OfferTakerShake:          "offer-taker-shake-",
	OfferMakerShake:          "offer-maker-shake-",
	OfferCompleted:           "offer-completed-",
	OfferTakerRejected:       "offer-taker-rejected-",
	OfferMakerRejected:       "offer-maker-rejected-",
	OfferWithdraw:            "offer-withdraw-",
	OfferStoreItemAdded:      "offer-store-item-added-",
	OfferStoreItemRemoved:    "offer-store-item-removed-",
	OfferStoreMakerSellShake: "offer-store-maker-sell-shake-",
	OfferStoreMakerBuyShake:  "offer-store-maker-buy-shake-",
	OfferStoreTakerSellShake: "offer-store-taker-sell-shake-",
	OfferStoreTakerBuyShake:  "offer-store-taker-buy-shake-",
	OfferStoreAccept:         "offer-store-accept-",
	OfferStoreMakerReject:    "offer-store-maker-reject-",
	OfferStoreTakerReject:    "offer-store-taker-reject-",
}
