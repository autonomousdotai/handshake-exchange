package bean

type FCMNotificationObject struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ClickAction string `json:"click_action"`
}

type FCMObject struct {
	Data         interface{}           `json:"data"`
	To           string                `json:"to"`
	Notification FCMNotificationObject `json:"notification"`
}
