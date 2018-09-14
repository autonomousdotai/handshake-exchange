package email

const OfferBuyingActive = "OfferBuyingActive"
const OfferSellingActive = "OfferSellingActive"
const OfferClosed = "OfferClosed"
const OfferTakerBuyShake = "OfferTakerBuyShake"
const OfferMakerBuyShake = "OfferMakerBuyShake"
const OfferTakerSellShake = "OfferTakerSellShake"
const OfferMakerSellShake = "OfferMakerSellShake"
const OfferBuyCompleted = "OfferBuyCompleted"
const OfferSellCompleted = "OfferSellCompleted"
const OfferMakerTakerRejected = "OfferMakerTakerRejected"
const OfferMakerMakerRejected = "OfferMakerMakerRejected"
const OfferTakerTakerRejected = "OfferTakerTakerRejected"
const OfferTakerMakerRejected = "OfferTakerMakerRejected"
const OfferWithdraw = "OfferWithdraw"
const OrderInstantCCSuccess = "OrderInstantCCSuccess"
const OfferStoreItemAdded = "OfferStoreItemAdded"
const OfferStoreItemRemoved = "OfferStoreItemRemoved"
const OfferStoreMakerSellShake = "OfferStoreMakerSellShake"
const OfferStoreMakerBuyShake = "OfferStoreMakerBuyShake"
const OfferStoreTakerSellShake = "OfferStoreTakerSellShake"
const OfferStoreTakerBuyShake = "OfferStoreTakerBuyShake"
const OfferStoreMakerAccept = "OfferStoreMakerAccept"
const OfferStoreTakerAccept = "OfferStoreTakerAccept"
const OfferStoreMakerReject = "OfferStoreMakerReject"
const OfferStoreTakerReject = "OfferStoreTakerReject"
const CreditWithdraw = "CreditWithdraw"

var TemplateName = map[string]string{
	OfferBuyingActive:        "offer-buying-active-",
	OfferSellingActive:       "offer-selling-active-",
	OfferClosed:              "offer-closed-",
	OrderInstantCCSuccess:    "order-instant-cc-success-",
	OfferTakerBuyShake:       "offer-taker-buy-shake-",
	OfferMakerBuyShake:       "offer-maker-buy-shake-",
	OfferTakerSellShake:      "offer-taker-sell-shake-",
	OfferMakerSellShake:      "offer-maker-sell-shake-",
	OfferBuyCompleted:        "offer-buy-completed-",
	OfferSellCompleted:       "offer-sell-completed-",
	OfferMakerTakerRejected:  "offer-maker-taker-rejected-",
	OfferMakerMakerRejected:  "offer-maker-maker-rejected-",
	OfferTakerTakerRejected:  "offer-taker-taker-rejected-",
	OfferTakerMakerRejected:  "offer-taker-maker-rejected-",
	OfferWithdraw:            "offer-withdraw-",
	OfferStoreItemAdded:      "offer-store-item-added-",
	OfferStoreItemRemoved:    "offer-store-item-removed-",
	OfferStoreMakerSellShake: "offer-store-maker-sell-shake-",
	OfferStoreMakerBuyShake:  "offer-store-maker-buy-shake-",
	OfferStoreTakerSellShake: "offer-store-taker-sell-shake-",
	OfferStoreTakerBuyShake:  "offer-store-taker-buy-shake-",
	OfferStoreMakerAccept:    "offer-store-maker-accept-",
	OfferStoreTakerAccept:    "offer-store-taker-accept-",
	OfferStoreMakerReject:    "offer-store-maker-reject-",
	OfferStoreTakerReject:    "offer-store-taker-reject-",
	CreditWithdraw:           "credit-withdraw-",
}
