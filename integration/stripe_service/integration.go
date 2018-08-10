package stripe_service

import (
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"os"
)

type StripeClient struct {
}

func CreateToken(cardNum string, date string, cvv string) (string, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	month := date[:2]
	year := "20" + date[3:]

	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: cardNum,
			Month:  month,
			Year:   year,
			CVC:    cvv,
		},
	}

	token, err := sc.Tokens.New(tokenParams)

	return token.ID, err
}

func CreateCard(token string, customerId string) (string, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	tokenParams := &stripe.CardParams{
		Token:    token,
		Customer: customerId,
	}

	stripeCard, err := sc.Cards.New(tokenParams)

	return stripeCard.ID, err
}

func CreateCustomer(description string, token string) (string, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	customerParams := &stripe.CustomerParams{
		Desc: description,
	}
	customerParams.SetSource(token) // obtained with Stripe.js
	c, err := sc.Customers.New(customerParams)

	return c.ID, err
}

func CreateCustomerRaw(description string) (string, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	customerParams := &stripe.CustomerParams{
		Desc: description,
	}
	c, err := sc.Customers.New(customerParams)

	return c.ID, err
}

func GetSourceChargeable(token string) (bool, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	source, err := sc.Sources.Get(token, nil)
	return source.Status == stripe.SourceStatusChargeable, err
}

func Charge(token string, customerId string, amount decimal.Decimal, statement string, description string) (*stripe.Charge, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	stripeAmount := amount.Round(2).Mul(decimal.NewFromFloat(100)).IntPart()
	chargeParams := &stripe.ChargeParams{
		Amount:    uint64(stripeAmount),
		Currency:  "usd",
		Desc:      description,
		Statement: statement,
		NoCapture: true,
	}
	if token != "" {
		chargeParams.SetSource(token)
	} else {
		chargeParams.Customer = customerId
	}

	ch, err := sc.Charges.New(chargeParams)

	return ch, err
}

func Refund(chargeId string) (*stripe.Refund, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	r, err := sc.Refunds.New(&stripe.RefundParams{Charge: chargeId})
	return r, err
}

func Capture(chargeId string) (*stripe.Charge, error) {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)

	ch, err := sc.Charges.Capture(chargeId, nil)
	return ch, err
}
