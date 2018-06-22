package notification

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/fcm_service"
	"os"
)

func SendOrderInstantCCSuccessFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_order_instant_cc_success")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferMakerShakeFCM(language string, fcm string, offerType string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	role := T("common_buyer")
	if offerType == bean.OFFER_TYPE_BUY {
		role = T("common_seller")
	}
	body := T("notification_offer_maker_shake", map[string]string{
		"Role": role,
	})

	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferTakerShakeFCM(language string, fcm string, offerType string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	role := T("common_seller")
	if offerType == bean.OFFER_TYPE_BUY {
		role = T("common_buyer")
	}
	body := T("notification_offer_taker_shake", map[string]string{
		"Role": role,
	})

	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferMakerRejectedFCM(language string, fcm string, offerType string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	role := T("common_buyer")
	if offerType == bean.OFFER_TYPE_BUY {
		role = T("common_seller")
	}
	body := T("notification_offer_maker_rejected", map[string]string{
		"Role": role,
	})

	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferCompletedFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_completed")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreItemAddedFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_added")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s", frontEndHost, "me")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreMakerSellShakeFCM(language string, fcm string, chatUsername string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_maker_sell_shake")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s/%s", frontEndHost, "chat", chatUsername)

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreTakerSellShakeFCM(language string, fcm string, currency string, chatUsername string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_taker_sell_shake", map[string]string{
		"Currency": currency,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s/%s", frontEndHost, "chat", chatUsername)

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreMakerBuyShakeFCM(language string, fcm string, chatUsername string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_maker_buy_shake")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s/%s", frontEndHost, "chat", chatUsername)

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreTakerBuyShakeFCM(language string, fcm string, currency string, chatUsername string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_taker_buy_shake", map[string]string{
		"Currency": currency,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s/%s", frontEndHost, "chat", chatUsername)

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreMakerRejectFCM(language string, fcm string, username string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_maker_reject", map[string]string{
		"Username": username,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s?id=2", frontEndHost, "discover")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreTakerRejectFCM(language string, fcm string, username string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_taker_reject", map[string]string{
		"Username": username,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s?id=2", frontEndHost, "discover")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreMakerCompleteFCM(language string, fcm string, currency string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_maker_accept", map[string]string{
		"Currency": currency,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s?id=2", frontEndHost, "discover")

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}

func SendOfferStoreTakerCompleteFCM(language string, fcm string, currency string, offerId string, offerShakeId string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_store_taker_accept", map[string]string{
		"Currency": currency,
	})
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/%s?s=%s&sh=%s", frontEndHost, "me", offerId, offerShakeId)

	fcmObj := bean.FCMObject{
		Notification: bean.FCMNotificationObject{
			Title:       title,
			Body:        body,
			ClickAction: url,
		},
		To: fcm,
	}

	return fcm_service.SendFCM(fcmObj)
}
