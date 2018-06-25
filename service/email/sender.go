package email

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"
	"os"
)

func SendOfferBuyingActiveEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_buying_active_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferBuyingActive,
		data)
}

func SendOfferSellingActiveEmail(language, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_selling_active_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferSellingActive,
		data)
}

func SendOfferClosedEmail(language, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_closed_subject")
	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		CreateOfferUrl string
	}{
		CreateOfferUrl: fmt.Sprintf("%s/create?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferClosed,
		data)
}

func SendOfferMakerBuyShakeEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_maker_buy_shake_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/me", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferMakerBuyShake,
		data)
}

func SendOfferTakerBuyShakeEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_taker_buy_shake_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferTakerBuyShake,
		data)
}

func SendOfferMakerSellShakeEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_maker_sell_shake_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/me", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferMakerSellShake,
		data)
}

func SendOfferTakerSellShakeEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_taker_sell_shake_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferTakerSellShake,
		data)
}

func SendOfferMakerMakerRejectEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_maker_maker_rejected_subject")

	//host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: "",
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferMakerMakerRejected,
		data)
}

func SendOfferTakerMakerRejectEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_taker_maker_rejected_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferTakerMakerRejected,
		data)
}

func SendOfferMakerTakerRejectEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_maker_taker_rejected_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferMakerTakerRejected,
		data)
}

func SendOfferTakerTakerRejectEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_taker_taker_rejected_subject")

	//host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: "",
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferTakerTakerRejected,
		data)
}

func SendOfferBuyCompleteEmail(language string, emailAddress string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_buy_completed_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/discover?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferBuyCompleted,
		data)
}

func SendOfferSellCompleteEmail(language string, emailAddress string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_sell_completed_subject")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/create?id=6", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferSellCompleted,
		data)
}

func SendOfferWithdrawEmail(language string, emailAddress string,
	amount string, currency string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_withdraw_subject", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount   string
		Currency string
	}{
		Amount:   amount,
		Currency: currency,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferWithdraw,
		data)
}

func SendOrderInstantCCSuccessEmail(language string, emailAddress string, amount string, currency string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_order_instant_cc_success_subject")

	data := struct {
		Name     string
		Currency string
		Amount   string
	}{
		Name:     emailAddress,
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

func SendOfferStoreItemAddedEmail(language string, emailAddress string, sellAmount string, buyAmount string, currency string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_item_added")

	data := struct {
		Currency   string
		SellAmount string
		BuyAmount  string
	}{
		Currency:   currency,
		SellAmount: sellAmount,
		BuyAmount:  buyAmount,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreItemAdded,
		data)
}

func SendOfferStoreItemRemovedEmail(language string, emailAddress string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_item_removed")

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url string
	}{
		Url: fmt.Sprintf("%s/create?id=2", host),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreItemRemoved,
		data)
}

func SendOfferStoreMakerSellShakeEmail(language string, emailAddress string, amount string, currency string,
	fiatAmount string, fiatCurrency string, username string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_maker_sell_shake", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount       string
		Currency     string
		FiatAmount   string
		FiatCurrency string
		Username     string
	}{
		Amount:       amount,
		Currency:     currency,
		FiatAmount:   fiatAmount,
		FiatCurrency: fiatCurrency,
		Username:     username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreMakerSellShake,
		data)
}

func SendOfferStoreMakerBuyShakeEmail(language string, emailAddress string, amount string, currency string,
	fiatAmount string, fiatCurrency string, username string) error {
	if emailAddress == "" {
		return nil
	}

	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_maker_buy_shake", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount       string
		Currency     string
		FiatAmount   string
		FiatCurrency string
		Username     string
	}{
		Amount:       amount,
		Currency:     currency,
		FiatAmount:   fiatAmount,
		FiatCurrency: fiatCurrency,
		Username:     username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreMakerBuyShake,
		data)
}

func SendOfferStoreTakerSellShakeEmail(language string, emailAddress string, amount string, currency string,
	fiatAmount string, fiatCurrency string, username string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_taker_sell_shake", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount       string
		Currency     string
		FiatAmount   string
		FiatCurrency string
		Username     string
	}{
		Amount:       amount,
		Currency:     currency,
		FiatAmount:   fiatAmount,
		FiatCurrency: fiatCurrency,
		Username:     username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreTakerSellShake,
		data)
}

func SendOfferStoreTakerBuyShakeEmail(language string, emailAddress string, amount string, currency string,
	fiatAmount string, fiatCurrency string, username string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_taker_buy_shake", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount       string
		Currency     string
		FiatAmount   string
		FiatCurrency string
		Username     string
	}{
		Amount:       amount,
		Currency:     currency,
		FiatAmount:   fiatAmount,
		FiatCurrency: fiatCurrency,
		Username:     username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreTakerBuyShake,
		data)
}

func SendOfferStoreMakerCompleteEmail(language string, emailAddress string, amount string, currency string, username string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_maker_accept", map[string]string{
		"Currency": currency,
	})

	data := struct {
		Amount   string
		Currency string
		Username string
	}{
		Amount:   amount,
		Currency: currency,
		Username: username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreMakerAccept,
		data)
}

func SendOfferStoreTakerCompleteEmail(language string, emailAddress string, amount string, currency string,
	username string, usernameStore string, offerId string, offerShakeId string) error {
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_taker_accept", map[string]string{
		"Currency": currency,
	})

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Amount        string
		Currency      string
		Username      string
		UsernameStore string
		Url           string
	}{
		Amount:        amount,
		Currency:      currency,
		Username:      username,
		UsernameStore: usernameStore,
		Url:           fmt.Sprintf("%s/me?s=%s&sh=%s", host, offerId, offerShakeId),
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreTakerAccept,
		data)
}

func SendOfferStoreMakerRejectEmail(language string, emailAddress string, username string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_maker_reject", map[string]string{
		"Username": username,
	})

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url      string
		Username string
	}{
		Url:      fmt.Sprintf("%s/create?id=2", host),
		Username: username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreMakerReject,
		data)
}

func SendOfferStoreTakerRejectEmail(language string, emailAddress string, username string) error {
	if emailAddress == "" {
		return nil
	}
	T, _ := i18n.Tfunc(language)

	subject := T("email_offer_store_taker_reject", map[string]string{
		"Username": username,
	})

	host := os.Getenv("FRONTEND_HOST")
	data := struct {
		Url      string
		Username string
	}{
		Url:      fmt.Sprintf("%s/discover?id=2", host),
		Username: username,
	}

	return SendSystemEmailWithTemplate(
		"",
		emailAddress,
		language,
		subject,
		OfferStoreTakerReject,
		data)
}
