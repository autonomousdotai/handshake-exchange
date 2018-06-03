package email

const CreateBuyingOfferSuccess = "CreateBuyingOfferSuccess"
const CreateSellingOfferSuccess = "CreateSellingOfferSuccess"
const OrderCancelled = "OrderCancelled"
const OrderSuccess = "OrderSuccess"
const OrderFromSuccess = "OrderFromSuccess"
const OrderInstantCCSuccess = "OrderInstantCCSuccess"

var TemplateName = map[string]string{
	CreateBuyingOfferSuccess:  "create-buying-offer-success-",
	CreateSellingOfferSuccess: "create-selling-offer-success-",
	OrderCancelled:            "order-cancelled-",
	OrderInstantCCSuccess:     "order-instant-cc-success-",
	OrderFromSuccess:          "order-from-success-",
	OrderSuccess:              "order-success-",
}
