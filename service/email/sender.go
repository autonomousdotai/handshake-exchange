package email

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"os"
	"strings"
)

func SendBuyingOfferSuccessEmail(language string, emailAddress string, name string, currency string, price string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_create_buying_offer_success_subject", map[string]string{
		"Currency": currency,
		"Price":    price,
	})

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Name     string
		Currency string
		Price    string
		BuyUrl   string
	}{
		Name:     name,
		Currency: currency,
		Price:    price,
		BuyUrl:   strings.Join([]string{host, "instant/buy"}, "/"),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		CreateBuyingOfferSuccess,
		data)
}

func SendSellingOfferSuccessEmail(language, emailAddress string, name string, currency string, price string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_create_selling_offer_success_subject", map[string]string{
		"Currency": currency,
		"Price":    price,
	})

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Name     string
		Currency string
		Price    string
		BuyUrl   string
	}{
		Name:     name,
		Currency: currency,
		Price:    price,
		BuyUrl:   strings.Join([]string{host, "instant/sell"}, "/"),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		CreateSellingOfferSuccess,
		data)
}

func SendOrderSuccessEmail(language string, emailAddress string, name string, currency string, amount string, fromUsername string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_order_success_subject", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Name         string
		Currency     string
		Amount       string
		FromUsername string
	}{
		Name:         name,
		Currency:     currency,
		Amount:       amount,
		FromUsername: fromUsername,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OrderSuccess,
		data)
}

func SendOrderFromSuccessEmail(language string, emailAddress string, name string, currency string, price string, toUsername string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_order_from_success_subject", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Name       string
		Currency   string
		Price      string
		ToUsername string
	}{
		Name:       name,
		Currency:   currency,
		Price:      price,
		ToUsername: toUsername,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OrderFromSuccess,
		data)
}

func SendOrderCancelledEmail(language, emailAddress string, name string, currency string, price string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_order_cancelled_subject")
	data := struct {
		Name         string
		Currency     string
		Price        string
		FromUsername string
	}{
		Name:     name,
		Currency: currency,
		Price:    price,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OrderCancelled,
		data)
}

func SendOrderInstantCCSuccessEmail(language string, emailAddress string, amount string, currency string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_order_instant_cc_success_subject")

	data := struct {
		Currency string
		Amount   string
	}{
		Currency: currency,
		Amount:   amount,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OrderInstantCCSuccess,
		data)
}
