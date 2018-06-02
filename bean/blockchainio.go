package bean

type BlockChainIoPayment struct {
	Message string `json:"message" firestore:"message"`
	TxHash  string `json:"tx_hash" firestore:"tx_hash"`
	Notice  string `json:"notice" firestore:"notice"`
}

type BlockChainIoBalance struct {
	Balance int64 `json:"balance" firestore:"balance"`
}

type BlockChainIoAddress struct {
	Address string `json:"address" firestore:"address"`
	Label   string `json:"label" firestore:"label"`
}

type BlockChainIoBalanceUpdates struct {
	Id             int64  `json:"id" firestore:"id"`
	Address        string `json:"addr" firestore:"address"`
	Op             string `json:"op" firestore:"op"`
	Confirmations  int64  `json:"confs" firestore:"confirmations"`
	Callback       string `json:"callback" firestore:"callback"`
	OnNotification string `json:"on_notification" firestore:"on_notification"`
}

type BlockChainIoCallback struct {
	TxHash        string `json:"tx_hash" firestore:"tx_hash"`
	Address       string `json:"addr" firestore:"address"`
	Confirmations int64  `json:"confs" firestore:"confirmations"`
	Value         int64  `json:"value" firestore:"value"`
	Offer         string `json:"offer" firestore:"offer"`
}
