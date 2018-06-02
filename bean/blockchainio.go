package bean

type BlockChainIoPaymentResponse struct {
	Message string `json:"message" firestore:"message"`
	TxHash  string `json:"tx_hash" firestore:"tx_hash"`
	Notice  string `json:"notice" firestore:"notice"`
}

type BlockChainIoBalance struct {
	Balance int64 `json:"balance" firestore:"balance"`
}

type BlockChainBalanceUpdates struct {
}
