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

var TemplateName = map[string]string{
	OfferBuyingActive:     "offer-buying-active-",
	OfferSellingActive:    "offer-selling-active-",
	OfferClosed:           "offer-closed-",
	OrderInstantCCSuccess: "order-instant-cc-success-",
	OfferTakerShake:       "offer-taker-shake-",
	OfferMakerShake:       "offer-maker-shake-",
	OfferCompleted:        "offer-completed-",
	OfferTakerRejected:    "offer-taker-rejected-",
	OfferMakerRejected:    "offer-maker-rejected-",
	OfferWithdraw:         "offer-withdraw-",
	OfferStoreItemAdded:   "offer-store-item-added-",
	OfferStoreItemRemoved: "offer-store-item-removed-",
}
