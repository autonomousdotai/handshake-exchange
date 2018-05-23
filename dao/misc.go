package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
	"github.com/autonomousdotai/handshake-exchange/service/cache"
	"github.com/shopspring/decimal"
	"google.golang.org/api/iterator"
	"strconv"
)

type MiscDao struct {
}

func (dao MiscDao) UpdateCurrencyRate(rates map[string]float64) error {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	for k := range rates {
		// currency_rates/{USDHKD}
		docRef := dbClient.Doc(GetCurrencyRateItemPath(fmt.Sprintf("USD%s", k)))
		batch.Set(docRef, bean.CurrencyRate{
			From: bean.USD.Code,
			To:   k,
			Rate: rates[k],
		})
		key := GetCurrencyRateItemCacheKey(fmt.Sprintf("USD%s", k))
		cache.RedisClient.Set(key, rates[k], 0)
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao MiscDao) GetCurrencyRate(from string, to string) (t TransferObject) {
	GetObject(GetCurrencyRateItemPath(fmt.Sprintf("%s%s", from, to)), &t, func(snapshot *firestore.DocumentSnapshot) interface{} {
		var obj bean.CurrencyRate
		snapshot.DataTo(&obj)
		return obj
	})

	return
}

func (dao MiscDao) GetCurrencyRateFromCache(from string, to string) (t TransferObject) {
	currencyRate := bean.CurrencyRate{
		From: from,
		To:   to,
	}

	GetCacheObject(GetCurrencyRateItemCacheKey(fmt.Sprintf("%s%s", from, to)), &t, func(val string) interface{} {
		rate, _ := strconv.ParseFloat(val, 64)
		currencyRate.Rate = rate

		return currencyRate
	})

	return
}

func (dao MiscDao) UpdateCryptoRates(rates map[string][]bean.CryptoRate) error {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	for k := range rates {
		// crypto_rates/{BTC}
		for _, item := range rates[k] {
			docRef := dbClient.Doc(GetCryptoRateItemPath(k, item.Exchange))
			batch.Set(docRef, item)
			b, _ := json.Marshal(&item)
			key := GetCryptoRateItemCacheKey(fmt.Sprintf("%s.%s", k, item.Exchange))
			cache.RedisClient.Set(key, string(b), 0)
		}
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao MiscDao) GetCryptoRatesFromCache(from string) (t TransferObject) {
	keys, err := cache.RedisClient.Keys(GetCryptoRateItemCacheKey(from) + "*").Result()
	if err != nil {
		t.SetError(api_error.GetDataFailed, err)
		return
	}

	t.Found = true
	for _, key := range keys {
		var tTemp TransferObject
		GetCacheObject(key, &tTemp, func(val string) interface{} {
			var cryptoRate bean.CryptoRate
			json.Unmarshal([]byte(val), &cryptoRate)
			return cryptoRate
		})
		t.Objects = append(t.Objects, tTemp.Object.(bean.CryptoRate))
	}

	return
}

func (dao MiscDao) GetCryptoRateFromCache(currency string, exchange string) (t TransferObject) {
	GetCacheObject(GetCryptoRateItemCacheKey(fmt.Sprintf("%s.%s", currency, exchange)), &t, func(val string) interface{} {
		var cryptoRate bean.CryptoRate
		json.Unmarshal([]byte(val), &cryptoRate)
		return cryptoRate
	})

	return
}

func (dao MiscDao) LoadSystemFeeToCache() ([]bean.SystemFee, error) {
	dbClient := firebase_service.FirestoreClient

	ClearCache(GetSystemFeeCacheKey("*"))

	// system_fees/
	addressesIter := dbClient.Collection(GetSystemFeePath()).Documents(context.Background())
	systemFees := make([]bean.SystemFee, 0)

	for {
		var systemFee bean.SystemFee
		doc, err := addressesIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return systemFees, err
		}
		doc.DataTo(&systemFee)
		systemFees = append(systemFees, systemFee)

		// To cache
		key := GetSystemFeeCacheKey(systemFee.Key)
		cache.RedisClient.Set(key, systemFee.Value, 0)
	}

	return systemFees, nil
}

func (dao MiscDao) GetSystemFeeFromCache(feeKey string) (t TransferObject) {
	// Warning: Don't use Fee type yet
	systemFee := bean.SystemFee{
		Key: feeKey,
	}

	GetCacheObject(GetSystemFeeCacheKey(feeKey), &t, func(val string) interface{} {
		testVal, _ := decimal.NewFromString(val)
		value, _ := testVal.Float64()
		systemFee.Value = value

		return systemFee
	})

	return
}

func (dao MiscDao) LoadCCLimitToCache() ([]bean.CCLimit, error) {
	dbClient := firebase_service.FirestoreClient

	ClearCache(GetCCLimitCacheKey("*"))

	// cc_limits/
	iter := dbClient.Collection(GetCCLimitPath()).Documents(context.Background())
	objs := make([]bean.CCLimit, 0)

	for {
		var obj bean.CCLimit
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objs, err
		}
		doc.DataTo(&obj)
		objs = append(objs, obj)

		// To cache
		b, _ := json.Marshal(&obj)
		key := GetCCLimitCacheKey(fmt.Sprintf("%d", obj.Level))
		cache.RedisClient.Set(key, string(b), 0)
	}

	return objs, nil
}

func (dao MiscDao) GetCCLimitFromCache() (t TransferObject) {
	keys, err := cache.RedisClient.Keys(GetCCLimitCacheKey("*")).Result()
	if err != nil {
		t.SetError(api_error.GetDataFailed, err)
		return
	}

	t.Found = true
	for _, key := range keys {
		var tTemp TransferObject
		GetCacheObject(key, &tTemp, func(val string) interface{} {
			var obj bean.CCLimit
			json.Unmarshal([]byte(val), &obj)
			return obj
		})
		t.Objects = append(t.Objects, tTemp.Object.(bean.CCLimit))
	}

	return
}

func (dao MiscDao) GetCCLimitByLevelFromCache(level string) (t TransferObject) {
	GetCacheObject(GetCCLimitCacheKey(level), &t, func(val string) interface{} {
		var obj bean.CCLimit
		json.Unmarshal([]byte(val), &obj)
		return obj
	})

	return
}

func GetCurrencyRateItemPath(currency string) string {
	return fmt.Sprintf("currency_rates/%s", currency)
}

func GetCurrencyRateItemCacheKey(currency string) string {
	return fmt.Sprintf("handshake_exchange.currency_rates.%s", currency)
}

func GetCryptoRateItemPath(currency string, exchange string) string {
	return fmt.Sprintf("crypto_rates/%s/exchanges/%s", currency, exchange)
}

func GetCryptoRateItemCacheKey(currency string) string {
	return fmt.Sprintf("handshake_exchange.crypto_rates.%s", currency)
}

func GetSystemFeePath() string {
	return "system_fees"
}

func GetSystemFeeCacheKey(fee string) string {
	return fmt.Sprintf("handshake_exchange.system_fees.%s", fee)
}

func GetCCLimitPath() string {
	return "cc_limits"
}

func GetCCLimitItemPath(level string) string {
	return fmt.Sprintf("cc_limits/%s", level)
}

func GetCCLimitCacheKey(level string) string {
	return fmt.Sprintf("handshake_exchange.cc_limits.%s", level)
}
