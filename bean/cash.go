package bean

import (
	"cloud.google.com/go/firestore"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/shopspring/decimal"
	"time"
)

const CASH_CREDIT_STATUS_ACTIVE = "active"
const CASH_CREDIT_STATUS_INACTIVE = "inactive"

const CASH_CREDIT_ITEM_STATUS_CREATE = "create"
const CASH_CREDIT_ITEM_STATUS_ACTIVE = "active"
const CASH_CREDIT_ITEM_STATUS_INACTIVE = "inactive"

const CASH_CREDIT_ITEM_SUB_STATUS_TRANSFERRING = "transferring"
const CASH_CREDIT_ITEM_SUB_STATUS_TRANSFERRED = "transferred"

const CASH_CREDIT_DEPOSIT_STATUS_CREATED = "created"
const CASH_CREDIT_DEPOSIT_STATUS_TRANSFERRING = "transferring"
const CASH_CREDIT_DEPOSIT_STATUS_FAILED = "failed"
const CASH_CREDIT_DEPOSIT_STATUS_TRANSFERRED = "transferred"

const CASH_CREDIT_WITHDRAW_STATUS_CREATED = "created"
const CASH_CREDIT_WITHDRAW_STATUS_PROCESSING = "processing"
const CASH_CREDIT_WITHDRAW_STATUS_FAILED = "failed"
const CASH_CREDIT_WITHDRAW_STATUS_PROCESSED = "processed"

type CashCredit struct {
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

func (b CashCredit) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":        b.UID,
		"username":   b.Username,
		"email":      b.Email,
		"language":   b.Language,
		"fcm":        b.FCM,
		"chain_id":   b.ChainId,
		"revenue":    common.Zero.String(),
		"status":     b.Status,
		"created_at": firestore.ServerTimestamp,
	}
}

func (b CashCredit) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CashCredit) GetUpdateRevenue() map[string]interface{} {
	return map[string]interface{}{
		"revenue":    b.Revenue,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CashCreditItem struct {
	Hid              int64       `json:"hid" firestore:"hid"`
	UID              string      `json:"-" firestore:"uid"`
	Currency         string      `json:"currency" firestore:"currency"`
	Status           string      `json:"status" firestore:"status"`
	SubStatus        string      `json:"sub_status" firestore:"sub_status"`
	LockedSale       bool        `json:"locked_sale" firestore:"locked_sale"`
	LastActionData   interface{} `json:"-" firestore:"last_action_data"`
	Balance          string      `json:"balance" firestore:"balance"`
	Sold             string      `json:"sold" firestore:"sold"`
	CreditRevenue    string      `json:"-"`
	Revenue          string      `json:"revenue" firestore:"revenue"`
	ReactivateAmount string      `json:"reactivate_amount" firestore:"reactivate_amount"`
	Percentage       string      `json:"percentage" firestore:"percentage"`
	UserAddress      string      `json:"user_address" firestore:"user_address"`
	CreatedAt        time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (b CashCreditItem) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"uid":          b.UID,
		"hid":          b.Hid,
		"currency":     b.Currency,
		"status":       b.Status,
		"sub_status":   b.SubStatus,
		"locked_sale":  false,
		"balance":      common.Zero.String(),
		"sold":         common.Zero.String(),
		"revenue":      common.Zero.String(),
		"percentage":   b.Percentage,
		"user_address": b.UserAddress,
		"created_at":   firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetUpdateReactivate() map[string]interface{} {
	return map[string]interface{}{
		"status":       b.Status,
		"percentage":   b.Percentage,
		"user_address": b.UserAddress,
		"updated_at":   firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":           b.Status,
		"sub_status":       b.SubStatus,
		"percentage":       b.Percentage,
		"last_action_data": b.LastActionData,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetUpdateDeactivate() map[string]interface{} {
	return map[string]interface{}{
		"status":            b.Status,
		"sub_status":        b.SubStatus,
		"balance":           b.Balance,
		"reactivate_amount": b.ReactivateAmount,
		"updated_at":        firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetUpdateLockedSale() map[string]interface{} {
	return map[string]interface{}{
		"locked_sale": true,
	}
}

func (b CashCreditItem) GetUpdateBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":     b.Balance,
		"sold":        b.Sold,
		"revenue":     b.Revenue,
		"locked_sale": b.LockedSale,
		"updated_at":  firestore.ServerTimestamp,
	}
}

func (b CashCreditItem) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"type":       "credit_item",
	}
}

type Center struct {
	Center       string            `json:"center" firestore:"center"`
	Country      string            `json:"country" firestore:"country"`
	FiatCurrency string            `json:"fiat_currency" firestore:"fiat_currency"`
	Information  map[string]string `json:"information" firestore:"information"`
}

const CASH_STORE_BUSINESS_TYPE_PERSONAL = "personal"
const CASH_STORE_BUSINESS_TYPE_STORE = "store"

const CASH_STORE_STATUS_OPEN = "open"
const CASH_STORE_STATUS_CLOSE = "close"

type CashStore struct {
	Name         string            `json:"name" firestore:"name"`
	Address      string            `json:"address" firestore:"address"`
	Phone        string            `json:"phone" firestore:"phone"`
	BusinessType string            `json:"business_type" firestore:"business_type"`
	Status       string            `json:"status" firestore:"status"`
	Center       string            `json:"center" firestore:"center"`
	Information  map[string]string `json:"information" firestore:"information"`
	Longitude    float64           `json:"longitude" firestore:"longitude"`
	Latitude     float64           `json:"latitude" firestore:"latitude"`
	ChainId      int64             `json:"chain_id" firestore:"chain_id"`
	CreatedAt    time.Time         `json:"created_at" firestore:"created_at"`
}

func (b CashStore) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"name":          b.Name,
		"address":       b.Address,
		"phone":         b.Phone,
		"business_type": b.BusinessType,
		"status":        b.Status,
		"center":        b.Center,
		"longitude":     b.Longitude,
		"latitude":      b.Latitude,
		"chain_id":      b.ChainId,
		"created_at":    firestore.ServerTimestamp,
	}
}

type CashCreditDepositInput struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	UserAddress string `json:"user_address"`
	Percentage  string `json:"percentage"`
}

type CashCreditDeposit struct {
	Id            string    `json:"id" firestore:"id"`
	UID           string    `json:"-" firestore:"uid"`
	ItemRef       string    `json:"-" firestore:"item_ref"`
	Status        string    `json:"status" firestore:"status"`
	Currency      string    `json:"currency" firestore:"currency" validator:"oneof=BTC ETH BCH"`
	Amount        string    `json:"amount" firestore:"amount"`
	Percentage    string    `json:"percentage"`
	SystemAddress string    `json:"system_address" firestore:"system_address"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
}

func (b CashCreditDeposit) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":             b.Id,
		"uid":            b.UID,
		"currency":       b.Currency,
		"item_ref":       b.ItemRef,
		"status":         b.Status,
		"amount":         b.Amount,
		"percentage":     b.Percentage,
		"system_address": b.SystemAddress,
		"created_at":     firestore.ServerTimestamp,
	}
}

func (b CashCreditDeposit) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"amount":     b.Amount,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CashCreditWithdraw struct {
	Id          string            `json:"id" firestore:"id"`
	UID         string            `json:"-" firestore:"uid"`
	Status      string            `json:"status" firestore:"status"`
	Amount      string            `json:"amount" firestore:"amount"`
	Information map[string]string `json:"information" firestore:"information"`
	ProcessedId string            `json:"processed_id" firestore:"processed_id"`
	CreatedAt   time.Time         `json:"created_at" firestore:"created_at"`
}

func (b CashCreditWithdraw) GetPaypalInformation(email string) map[string]string {
	return map[string]string{
		"email": email,
	}
}

func (b CashCreditWithdraw) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"uid":         b.UID,
		"status":      b.Status,
		"amount":      b.Amount,
		"information": b.Information,
		"created_at":  firestore.ServerTimestamp,
	}
}

func (b CashCreditWithdraw) GetUpdateStatus() map[string]interface{} {
	return map[string]interface{}{
		"processed_id": b.ProcessedId,
		"status":       b.Status,
		"updated_at":   firestore.ServerTimestamp,
	}
}

type CashCreditBalanceHistory struct {
	Id         string    `json:"id" firestore:"id"`
	ItemRef    string    `json:"-" firestore:"item_ref"`
	ModifyRef  string    `json:"-" firestore:"modify_ref"`
	ModifyType string    `json:"-" firestore:"modify_type"`
	Old        string    `json:"old" firestore:"old"`
	Change     string    `json:"change" firestore:"change"`
	New        string    `json:"new" firestore:"new"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

func (b CashCreditBalanceHistory) GetAdd() map[string]interface{} {
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

type CashCreditOnchain struct {
	Hid      int64
	UID      string
	Currency string
}

type CashCreditOnChainActionTrackingInput struct {
	Deposit  string `json:"deposit"`
	TxHash   string `json:"tx_hash"`
	Action   string `json:"action" validate:"oneof=deposit close"`
	Reason   string `json:"reason"`
	Currency string `json:"currency"`
}

const CASH_CREDIT_ON_CHAIN_ACTION_DEPOSIT = "deposit"
const CASH_CREDIT_ON_CHAIN_ACTION_CLOSE = "close"

type CashCreditOnChainActionTracking struct {
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

func (b CashCreditOnChainActionTracking) GetAdd() map[string]interface{} {
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

func (b CashCreditOnChainActionTracking) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"amount":     b.Amount,
		"reason":     b.Reason,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CashCreditPool struct {
	Level           string    `json:"level" firestore:"level"`
	Balance         string    `json:"balance" firestore:"balance"`
	CapturedBalance string    `json:"captured_balance" firestore:"captured_balance"`
	Currency        string    `json:"currency" firestore:"currency"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}

func (b CashCreditPool) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"level":            b.Level,
		"balance":          b.Balance,
		"captured_balance": b.CapturedBalance,
		"currency":         b.Currency,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CashCreditPool) GetUpdateBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CashCreditPool) GetUpdateCapturedBalance() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CashCreditPool) GetUpdateAllBalance() map[string]interface{} {
	return map[string]interface{}{
		"balance":          b.Balance,
		"captured_balance": b.CapturedBalance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

const CASH_CREDIT_POOL_MODIFY_TYPE_DEPOSIT = "deposit"
const CASH_CREDIT_POOL_MODIFY_TYPE_CLOSE = "close"
const CASH_CREDIT_POOL_MODIFY_TYPE_PURCHASE = "purchase"

type CashCreditPoolBalanceHistory struct {
	Id         string    `json:"id" firestore:"id"`
	ItemRef    string    `json:"-" firestore:"item_ref"`
	ModifyRef  string    `json:"-" firestore:"modify_ref"`
	ModifyType string    `json:"-" firestore:"modify_type"`
	Old        string    `json:"old" firestore:"old"`
	Change     string    `json:"change" firestore:"change"`
	New        string    `json:"new" firestore:"new"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

func (b CashCreditPoolBalanceHistory) GetAdd() map[string]interface{} {
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

type CashCreditPoolOrder struct {
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

func (b CashCreditPoolOrder) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":          b.Id,
		"uid":         b.UID,
		"deposit_ref": b.DepositRef,
		"amount":      b.Amount,
		"balance":     b.Balance,
		"created_at":  firestore.ServerTimestamp,
	}
}

func (b CashCreditPoolOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"balance":    b.Balance,
		"updated_at": firestore.ServerTimestamp,
	}
}

func (b CashCreditPoolOrder) GetUpdateCapture() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"captured_full":    b.CapturedFull,
		"updated_at":       firestore.ServerTimestamp,
	}
}

func (b CashCreditPoolOrder) GetUpdateAllBalance() map[string]interface{} {
	return map[string]interface{}{
		"captured_balance": b.CapturedBalance,
		"balance":          b.Balance,
		"updated_at":       firestore.ServerTimestamp,
	}
}

const CASH_CREDIT_TRANSACTION_STATUS_CREATE = "create"
const CASH_CREDIT_TRANSACTION_STATUS_SUCCESS = "success"
const CASH_CREDIT_TRANSACTION_STATUS_FAILED = "failed"

const CASH_CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSING = "create"
const CASH_CREDIT_TRANSACTION_SUB_STATUS_REVENUE_PROCESSED = "success"

type CashCreditTransaction struct {
	Id            string                   `json:"id" firestore:"id"`
	UID           string                   `json:"-" firestore:"uid"`
	UIDs          []string                 `json:"-" firestore:"uids"`
	ToUID         string                   `json:"-" firestore:"to_uid"`
	Amount        string                   `json:"amount" firestore:"amount"`
	Currency      string                   `json:"currency" firestore:"currency"`
	Status        string                   `json:"status" firestore:"status"`
	SubStatus     string                   `json:"sub_status" firestore:"sub_status"`
	Revenue       string                   `json:"revenue" firestore:"revenue"`
	Fee           string                   `json:"fee" firestore:"fee"`
	Percentage    string                   `json:"percentage" firestore:"percentage"`
	OfferRef      string                   `json:"-" firestore:"offer_ref"`
	OrderInfoRefs []CashCreditOrderInfoRef `json:"-" firestore:"order_info_refs"`
	CreatedAt     time.Time                `json:"created_at" firestore:"created_at"`
}

type CashCreditOrderInfoRef struct {
	OrderRef string `json:"-" firestore:"order_ref"`
	Amount   string `json:"-" firestore:"amount"`
}

func (b CashCreditTransaction) GetAdd() map[string]interface{} {
	return map[string]interface{}{
		"id":              b.Id,
		"uid":             b.UID,
		"uids":            b.UIDs,
		"to_uid":          b.ToUID,
		"amount":          b.Amount,
		"currency":        b.Currency,
		"status":          b.Status,
		"revenue":         b.Revenue,
		"fee":             b.Fee,
		"percentage":      b.Percentage,
		"offer_ref":       b.OfferRef,
		"order_info_refs": b.OrderInfoRefs,
		"created_at":      firestore.ServerTimestamp,
	}
}

func (b CashCreditTransaction) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"status":     b.Status,
		"sub_status": b.SubStatus,
		"revenue":    b.Revenue,
		"fee":        b.Fee,
		"offer_ref":  b.OfferRef,
		"updated_at": firestore.ServerTimestamp,
	}
}

type CashStoreOrder struct {
	Id                    string      `json:"id" firestore:"id"`
	UID                   string      `json:"uid" firestore:"uid"`
	Username              string      `json:"username" firestore:"username"`
	Amount                string      `json:"amount" firestore:"amount" validate:"required"`
	Currency              string      `json:"currency" firestore:"currency" validate:"required"`
	FiatAmount            string      `json:"fiat_amount" firestore:"fiat_amount" validate:"required"`
	RawFiatAmount         string      `json:"-" firestore:"raw_fiat_amount"`
	FiatCurrency          string      `json:"fiat_currency" firestore:"fiat_currency" validate:"required"`
	Price                 string      `json:"price" firestore:"price"`
	Status                string      `json:"status" firestore:"status"`
	Type                  string      `json:"type" firestore:"type"`
	Duration              int64       `json:"-" firestore:"duration"`
	FeePercentage         string      `json:"-" firestore:"fee_percentage"`
	Fee                   string      `json:"-" firestore:"fee"`
	StoreFeePercentage    string      `json:"-" firestore:"store_fee_percentage"`
	StoreFee              string      `json:"-" firestore:"store_fee"`
	ExternalFeePercentage string      `json:"-" firestore:"external_fee_percentage"`
	ExternalFee           string      `json:"-" firestore:"external_fee"`
	PaymentMethod         string      `json:"-" firestore:"payment_method"`
	PaymentMethodRef      string      `json:"-" firestore:"payment_method_ref"`
	PaymentMethodData     interface{} `json:"payment_method_data" validate:"required"`
	Center                string      `json:"center" firestore:"center"`
	ProviderWithdrawData  interface{} `json:"-" firestore:"provider_withdraw_data"`
	FCM                   string      `json:"fcm" firestore:"fcm"`
	Language              string      `json:"language" firestore:"language"`
	ChainId               int64       `json:"chain_id" firestore:"chain_id"`
	CreatedAt             time.Time   `json:"created_at" firestore:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at" firestore:"updated_at"`
}

func (offer CashStoreOrder) GetAddInstantOffer() map[string]interface{} {
	return map[string]interface{}{
		"id":                      offer.Id,
		"uid":                     offer.UID,
		"amount":                  offer.Amount,
		"currency":                offer.Currency,
		"fiat_amount":             offer.FiatAmount,
		"raw_fiat_amount":         offer.RawFiatAmount,
		"fiat_currency":           offer.FiatCurrency,
		"price":                   offer.Price,
		"status":                  offer.Status,
		"type":                    offer.Type,
		"fee":                     offer.Fee,
		"external_fee":            offer.ExternalFee,
		"fee_percentage":          offer.FeePercentage,
		"external_fee_percentage": offer.ExternalFeePercentage,
		"duration":                offer.Duration,
		"payment_method":          offer.PaymentMethod,
		"payment_method_ref":      offer.PaymentMethodRef,
		"language":                offer.Language,
		"fcm":                     offer.FCM,
		"chain_id":                offer.ChainId,
		"created_at":              firestore.ServerTimestamp,
	}
}

func (offer CashStoreOrder) GetUpdate() map[string]interface{} {
	return map[string]interface{}{
		"provider_withdraw_data": offer.ProviderWithdrawData,
		"status":                 offer.Status,
		"updated_at":             firestore.ServerTimestamp,
	}
}

func (offer CashStoreOrder) GetNotificationUpdate() map[string]interface{} {
	return map[string]interface{}{
		"id":     offer.Id,
		"status": offer.Status,
		"type":   "instant",
	}
}
