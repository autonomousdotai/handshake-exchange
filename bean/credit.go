package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/shopspring/decimal"
	"time"
)

const CREDIT_STATUS_ACTIVE = "active"
const CREDIT_STATUS_INACTIVE = "inactive"

const CREDIT_ITEM_STATUS_CREATE = "create"
const CREDIT_ITEM_STATUS_ACTIVE = "active"
const CREDIT_ITEM_STATUS_INACTIVE = "inactive"

const CREDIT_ITEM_SUB_STATUS_TRANSFERRING = "transferring"
const CREDIT_ITEM_SUB_STATUS_TRANSFERRED = "transferred"

const CREDIT_DEPOSIT_STATUS_CREATED = "created"
const CREDIT_DEPOSIT_STATUS_TRANSFERRING = "transferring"
const CREDIT_DEPOSIT_STATUS_FAILED = "failed"
const CREDIT_DEPOSIT_STATUS_TRANSFERRED = "transferred"

const CREDIT_WITHDRAW_STATUS_PROCESSING = "processing"
const CREDIT_WITHDRAW_STATUS_FAILED = "failed"
const CREDIT_WITHDRAW_STATUS_PROCESSED = "processed"

type Credit struct {
	UID       string                `json:"-" firestore:"uid"`
	Username  string                `json:"username" firestore:"username"`
	Email     string                `json:"email" firestore:"email"`
	Language  string                `json:"-" firestore:"language"`
	FCM       string                `json:"-" firestore:"fcm"`
	ChainId   string                `json:"-" firestore:"chain_id"`
	Status    string                `json:"status" firestore:"status"`
	Items     map[string]CreditItem `json:"items"`
	Revenue   string                `json:"revenue" firestore:"revenue"`
	CreatedAt time.Time             `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" firestore:"updated_at"`
}

func (b Credit) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":        b.UID,
		"username":   b.Username,
		"email":      b.Email,
		"language":   b.Language,
		"fcm":        b.FCM,
		"chain_id":   b.ChainId,
		"revenue":    common.Zero,
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

type CreditItem struct {
	Hid            int64       `json:"hid" firestore:"hid"`
	UID            string      `json:"-" firestore:"uid"`
	Currency       string      `json:"currency" firestore:"currency"`
	Status         string      `json:"status" firestore:"status"`
	SubStatus      string      `json:"sub_status" firestore:"sub_status"`
	LastActionData interface{} `json:"-" firestore:"last_action_data"`
	Balance        string      `json:"balance" firestore:"balance"`
	Sold           string      `json:"sold" firestore:"sold"`
	CreditRevenue  string      `json:"-"`
	Revenue        string      `json:"revenue" firestore:"revenue"`
	Percentage     string      `json:"percentage" firestore:"percentage"`
	UserAddress    string      `json:"user_address" firestore:"user_address"`
	CreatedAt      time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (b CreditItem) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":          b.UID,
		"hid":          b.Hid,
		"currency":     b.Currency,
		"status":       b.Status,
		"sub_status":   b.SubStatus,
		"balance":      common.Zero.String(),
		"sold":         common.Zero.String(),
		"revenue":      common.Zero.String(),
		"percentage":   b.Percentage,
		"user_address": b.UserAddress,
		"created_at":   firestore.ServerTimestamp,
	}
}

func (b CreditItem) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":           b.Status,
		"sub_status":       b.SubStatus,
		"last_action_data": b.LastActionData,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CreditItem) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CreditItem) GetUpdateBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"sold":       b.Sold,
		"revenue":    b.Revenue,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CreditDepositInput struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	UserAddress string `json:"user_address"`
	Percentage  string `json:"percentage"`
}

type CreditDeposit struct {
	Id            string    `json:"id" firestore:"id"`
	UID           string    `json:"-" firestore:"uid"`
	ItemRef       string    `json:"-" firestore:"item_ref"`
	Status        string    `json:"status" firestore:"status"`
	Currency      string    `json:"currency" firestore:"currency" validator:"oneof=BTC ETH BCH"`
	Amount        string    `json:"amount" firestore:"amount"`
	SystemAddress string    `json:"system_address" firestore:"system_address"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditDeposit) GetAdd() map[string]interface{} {
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

func (b CreditDeposit) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"amount":     b.Amount,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CreditWithdraw struct {
	Id          string    `json:"id" firestore:"id"`
	UID         string    `json:"-" firestore:"uid"`
	Currency    string    `json:"currency" firestore:"currency"`
	ItemRef     string    `json:"-" firestore:"item_ref"`
	Status      string    `json:"status" firestore:"status"`
	Amount      string    `json:"amount" firestore:"amount"`
	Information string    `json:"information" firestore:"information"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditWithdraw) GetAdd() map[string]interface{} {
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

type CreditBalanceHistory struct {
	Id         string    `json:"id" firestore:"id"`
	ItemRef    string    `json:"-" firestore:"item_ref"`
	ModifyRef  string    `json:"-" firestore:"modify_ref"`
	ModifyType string    `json:"-" firestore:"modify_type"`
	Old        string    `json:"old" firestore:"old"`
	Change     string    `json:"change" firestore:"change"`
	New        string    `json:"new" firestore:"new"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditBalanceHistory) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"item_ref":    b.ItemRef,
		"modify_ref":  b.ModifyRef,
		"modify_type": b.ModifyType,
		"old":         b.Old,
		"change":      b.Change,
		"new":         b.New,
		"created_at":  firestore.ServerTimestamp,
	}
}

type CreditOnchain struct {
	Hid      int64
	UID      string
	Currency string
}

type CreditOnChainActionTrackingInput struct {
	Deposit  string `json:"deposit"`
	TxHash   string `json:"tx_hash"`
	Action   string `json:"action" validate:"oneof=deposit close"`
	Reason   string `json:"reason"`
	Currency string `json:"currency"`
}

const CREDIT_ON_CHAIN_ACTION_DEPOSIT = "deposit"
const CREDIT_ON_CHAIN_ACTION_CLOSE = "close"

type CreditOnChainActionTracking struct {
	Id         string    `json:"id" firestore:"id"`
	UID        string    `json:"uid" firestore:"uid"`
	ItemRef    string    `json:"-" firestore:"item_ref"`
	DepositRef string    `json:"-" firestore:"deposit_ref"`
	TxHash     string    `json:"tx_hash" firestore:"tx_hash"`
	Amount     string    `json:"amount" firestore:"amount"`
	Currency   string    `json:"currency" firestore:"currency"`
	Action     string    `json:"action" firestore:"action"`
	Reason     string    `json:"reason" firestore:"reason"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditOnChainActionTracking) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"uid":         b.UID,
		"item_ref":    b.ItemRef,
		"deposit_ref": b.DepositRef,
		"tx_hash":     b.TxHash,
		"amount":      b.Amount,
		"action":      b.Action,
		"reason":      b.Reason,
		"currency":    b.Currency,
		"created_at":  firestore.ServerTimestamp,
	}
}

func (b CreditOnChainActionTracking) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"amount":     b.Amount,
		"reason":     b.Reason,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CreditPool struct {
	Level           string    `json:"level" firestore:"level"`
	Balance         string    `json:"balance" firestore:"balance"`
	CapturedBalance string    `json:"captured_balance" firestore:"captured_balance"`
	Currency        string    `json:"currency" firestore:"currency"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}

func (b CreditPool) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"level":            b.Level,
		"balance":          b.Balance,
		"captured_balance": b.CapturedBalance,
		"currency":         b.Currency,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CreditPool) GetUpdateBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CreditPool) GetUpdateCapturedBalance() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CreditPool) GetUpdateAllBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":          b.Balance,
		"captured_balance": b.CapturedBalance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

const CREDIT_POOL_MODIFY_TYPE_DEPOSIT = "deposit"
const CREDIT_POOL_MODIFY_TYPE_CLOSE = "close"
const CREDIT_POOL_MODIFY_TYPE_PURCHASE = "purchase"

type CreditPoolBalanceHistory struct {
	Id         string    `json:"id" firestore:"id"`
	ItemRef    string    `json:"-" firestore:"item_ref"`
	ModifyRef  string    `json:"-" firestore:"modify_ref"`
	ModifyType string    `json:"-" firestore:"modify_type"`
	Old        string    `json:"old" firestore:"old"`
	Change     string    `json:"change" firestore:"change"`
	New        string    `json:"new" firestore:"new"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

func (b CreditPoolBalanceHistory) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"item_ref":    b.ItemRef,
		"modify_ref":  b.ModifyRef,
		"modify_type": b.ModifyType,
		"old":         b.Old,
		"change":      b.Change,
		"new":         b.New,
		"created_at":  firestore.ServerTimestamp,
	}
}

type CreditPoolOrder struct {
	Id              string          `json:"id" firestore:"id"`
	UID             string          `json:"-" firestore:"uid"`
	DepositRef      string          `json:"deposit_ref" firestore:"deposit_ref"`
	Amount          string          `json:"amount" firestore:"amount"`
	Balance         string          `json:"balance" firestore:"balance"`
	CapturedAmount  decimal.Decimal `json:"-"`
	CapturedBalance string          `json:"captured_balance" firestore:"captured_balance"`
	CapturedFull    bool            `json:"captured_full" firestore:"captured_full"`
	CreatedAt       time.Time       `json:"created_at" firestore:"created_at"`
}

func (b CreditPoolOrder) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"uid":         b.UID,
		"deposit_ref": b.DepositRef,
		"amount":      b.Amount,
		"balance":     b.Balance,
		"created_at":  firestore.ServerTimestamp,
	}
}

func (b CreditPoolOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CreditPoolOrder) GetUpdateCapture() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"captured_full":    b.CapturedFull,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CreditPoolOrder) GetUpdateAllBalance() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"balance":          b.Balance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

const CREDIT_TRANSACTION_STATUS_CREATE = "create"
const CREDIT_TRANSACTION_STATUS_SUCCESS = "success"
const CREDIT_TRANSACTION_STATUS_FAILED = "failed"

const CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSING = "create"
const CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED = "success"

type CreditTransaction struct {
	Id            string         `json:"id" firestore:"id"`
	UID           string         `json:"-" firestore:"uid"`
	UIDs          []string       `json:"-" firestore:"uids"`
	ToUID         string         `json:"-" firestore:"to_uid"`
	Amount        string         `json:"amount" firestore:"amount"`
	Currency      string         `json:"currency" firestore:"currency"`
	Status        string         `json:"status" firestore:"status"`
	SubStatus     string         `json:"sub_status" firestore:"sub_status"`
	Revenue       string         `json:"revenue" firestore:"revenue"`
	Percentage    string         `json:"percentage" firestore:"percentage"`
	OfferRef      string         `json:"-" firestore:"offer_ref"`
	OrderInfoRefs []OrderInfoRef `json:"-" firestore:"order_info_refs"`
	CreatedAt     time.Time      `json:"created_at" firestore:"created_at"`
}

type OrderInfoRef struct {
	OrderRef string `json:"-" firestore:"order_ref"`
	Amount   string `json:"-" firestore:"amount"`
}

func (b CreditTransaction) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":              b.Id,
		"uid":             b.UID,
		"uids":            b.UIDs,
		"to_uid":          b.ToUID,
		"amount":          b.Amount,
		"currency":        b.Currency,
		"status":          b.Status,
		"revenue":         b.Revenue,
		"percentage":      b.Percentage,
		"offer_ref":       b.OfferRef,
		"order_info_refs": b.OrderInfoRefs,
		"created_at":      firestore.ServerTimestamp,
	}
}

func (b CreditTransaction) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"revenue":    b.Revenue,
		"offer_ref":  b.OfferRef,
		"updated_at": firestore.ServerTimestamp,
	}
}
