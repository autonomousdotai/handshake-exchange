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

func SendOfferMakerBuyShakeFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_maker_buy_shake")

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

func SendOfferMakerSellShakeFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_maker_sell_shake")

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

func SendOfferMakerMakerRejectFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_maker_maker_reject")

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

func SendOfferTakerMakerRejectFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_taker_maker_reject")

	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/discover?id=6", frontEndHost)

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

func SendOfferMakerTakerRejectFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_maker_taker_reject")

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

func SendOfferTakerTakerRejectFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_taker_taker_reject")

	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/discover?id=6", frontEndHost)

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

func SendOfferSellCompleteFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_sell_completed")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/discover?id=6", frontEndHost)

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

func SendOfferBuyCompleteFCM(language string, fcm string) error {
	T, _ := i18n.Tfunc(language)

	title := T("common_notification_title")
	body := T("notification_offer_buy_completed")
	frontEndHost := os.Getenv("FRONTEND_HOST")
	url := fmt.Sprintf("%s/discover?id=6", frontEndHost)

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
