package bean

type FCMNotificationObject struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ClickAction string `json:"click_action"`
}

type FCMObject struct {
	To           string                `json:"to"`
	Notification FCMNotificationObject `json:"notification"`
}

type FCMRequest struct {
	Data FCMObject `json:"data"`
}
