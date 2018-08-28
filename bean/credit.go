package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/ninjadotorg/handshake-exchange/common"
	"time"
)

type Credit struct {
	UID       string    `json:"-" firestore:"uid"`
	Username  string    `json:"username" firestore:"username"`
	Email     string    `json:"email" firestore:"email"`
	Language  string    `json:"language" firestore:"language"`
	FCM       string    `json:"-" firestore:"fcm"`
	Status    string    `json:"status" firestore:"status"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

func (b Credit) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":        b.UID,
		"username":   b.Username,
		"email":      b.Email,
		"language":   b.Language,
		"fcm":        b.FCM,
		"status":     b.Status,
		"created_at": firestore.ServerTimestamp,
	}
}

func (b Credit) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

type Item struct {
	Hid        int64     `json:"hid" firestore:"hid"`
	Currency   string    `json:"currency" firestore:"currency"`
	Status     string    `json:"status" firestore:"status"`
	SubStatus  string    `json:"sub_status" firestore:"sub_status"`
	Balance    string    `json:"balance" firestore:"balance"`
	Profit     string    `json:"profit" firestore:"profit"`
	Percentage string    `json:"percentage" firestore:"percentage"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
}

func (b Item) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"hid":        b.Hid,
		"currency":   b.Currency,
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"balance":    common.Zero.String(),
		"profit":     common.Zero.String(),
		"percentage": b.Percentage,
		"created_at": firestore.ServerTimestamp,
	}
}

func (b Item) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"balance":    b.Balance,
	}
}

func (b Item) GetUpdateBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance": b.Balance,
		"profit":  b.Profit,
	}
}

type Deposit struct {
	Id            string    `json:"id" firestore:"id"`
	UID           string    `json:"-" firestore:"uid"`
	Currency      string    `json:"currency" firestore:"currency"`
	ItemRef       string    `json:"item_ref" firestore:"item_ref"`
	Status        string    `json:"status" firestore:"status"`
	Amount        string    `json:"balance" firestore:"balance"`
	SystemAddress string    `json:"system_address" firestore:"system_address"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
}

func (b Deposit) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":             b.Id,
		"uid":            b.UID,
		"currency":       b.Currency,
		"item_ref":       b.ItemRef,
		"status":         b.Status,
		"amount":         b.Amount,
		"system_address": b.SystemAddress,
		"created_at":     firestore.ServerTimestamp,
	}
}

type Withdraw struct {
	Id          string    `json:"id" firestore:"id"`
	UID         string    `json:"-" firestore:"uid"`
	Currency    string    `json:"currency" firestore:"currency"`
	ItemRef     string    `json:"item_ref" firestore:"item_ref"`
	Status      string    `json:"status" firestore:"status"`
	Amount      string    `json:"balance" firestore:"balance"`
	Information string    `json:"information" firestore:"information"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
}

func (b Withdraw) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"uid":         b.UID,
		"currency":    b.Currency,
		"item_ref":    b.ItemRef,
		"status":      b.Status,
		"amount":      b.Amount,
		"information": b.Information,
		"created_at":  firestore.ServerTimestamp,
	}
}

type BalanceHistory struct {
	ItemRef   string    `json:"item_ref" firestore:"item_ref"`
	ModifyRef string    `json:"modify_ref" firestore:"modify_ref"`
	Old       string    `json:"old" firestore:"old"`
	Change    string    `json:"change" firestore:"change"`
	New       string    `json:"new" firestore:"new"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

func (b BalanceHistory) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"item_ref":   b.ItemRef,
		"modify_ref": b.ModifyRef,
		"old":        b.Old,
		"change":     b.Change,
		"new":        b.New,
		"created_at": firestore.ServerTimestamp,
	}
}

type CreditOnchain struct {
	Hid      int64
	UID      string
	Currency string
}

type CreditOnChainTransaction struct {
	TxHash   string `json:"tx_hash"`
	Action   string `json:"action"`
	Reason   string `json:"reason"`
	Currency string `json:"currency"`
}

func (b CreditOnChainTransaction) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"tx_hash":    b.TxHash,
		"action":     b.Action,
		"reason":     b.Reason,
		"currency":   b.Currency,
		"created_at": firestore.ServerTimestamp,
	}
}

type CreditOnChainActionTracking struct {
	Id        string    `json:"id" firestore:"id"`
	UID       string    `json:"uid" firestore:"uid"`
	ItemRef   string    `json:"item_ref" firestore:"item_ref"`
	TxHash    string    `json:"tx_hash" firestore:"tx_hash"`
	Currency  string    `json:"currency" firestore:"currency"`
	Action    string    `json:"action" firestore:"action"`
	Reason    string    `json:"reason" firestore:"reason"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditOnChainActionTracking) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":         b.Id,
		"uid":        b.UID,
		"item_ref":   b.ItemRef,
		"tx_hash":    b.TxHash,
		"action":     b.Action,
		"reason":     b.Reason,
		"currency":   b.Currency,
		"created_at": firestore.ServerTimestamp,
	}
}

type CreditPool struct {
	Level     string    `json:"level" firestore:"level"`
	Balance   string    `json:"balance" firestore:"balance"`
	Currency  string    `json:"currency" firestore:"currency"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

func (b CreditPool) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"level":      b.Level,
		"balance":    b.Balance,
		"currency":   b.Currency,
		"updated_at": firestore.ServerTimestamp,
	}
}

type PoolBalanceHistory struct {
	ItemRef   string    `json:"item_ref" firestore:"item_ref"`
	ModifyRef string    `json:"modify_ref" firestore:"modify_ref"`
	Type      string    `json:"type" firestore:"type"`
	Old       string    `json:"old" firestore:"old"`
	Change    string    `json:"change" firestore:"change"`
	New       string    `json:"new" firestore:"new"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

func (b PoolBalanceHistory) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"item_ref":   b.ItemRef,
		"modify_ref": b.ModifyRef,
		"old":        b.Old,
		"change":     b.Change,
		"new":        b.New,
		"created_at": firestore.ServerTimestamp,
	}
}
