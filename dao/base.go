package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/autonomousdotai/handshake-exchange/api_error"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/common"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
	"github.com/autonomousdotai/handshake-exchange/service/cache"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"google.golang.org/api/iterator"
	"log"
)

type TransferObject struct {
	Object    interface{}
	Objects   []interface{}
	Page      interface{}
	CanMove   bool
	Found     bool
	StatusKey string
	Error     error
}

func (t *TransferObject) SetError(statusKey string, err error) bool {
	// Only set to error and status key if there is really error
	if err != nil {
		t.StatusKey = statusKey
		t.Error = err
		return true
	}

	return false
}

func (t *TransferObject) SetStatusKey(statusKey string) {
	t.SetError(statusKey, errors.New(statusKey))
}

func (t *TransferObject) HasError() bool {
	if !t.Found || t.StatusKey != "" || t.Error != nil {
		return true
	}
	return false
}

func (t *TransferObject) ContextValidate(context *gin.Context) (invalid bool) {
	if !t.Found {
		api_error.AbortNotFound(context)
		return true
	}
	if t.Error != nil {
		api_error.PropagateErrorAndAbort(context, t.StatusKey, t.Error)
		return true
	}
	if t.StatusKey != "" {
		api_error.AbortWithValidateErrorSimple(context, t.StatusKey)
		return true
	}

	return
}

func (t *TransferObject) FeedPaging(limit int) {
	total := len(t.Objects)
	offerCount := len(t.Objects)
	if offerCount > limit {
		t.Objects = t.Objects[:limit]
		offerCount -= 1
	}

	// Process next
	t.Page = interface{}(nil)
	if total > 0 {
		t.Page = t.Objects[offerCount-1].(bean.Paging).GetPageValue()
		if total > limit {
			t.CanMove = true
		}
	}
}

func (t *TransferObject) FeedComplexPaging(limit int, getPageValue func(interface{}) interface{}) {
	total := len(t.Objects)
	offerCount := len(t.Objects)
	if offerCount > limit {
		t.Objects = t.Objects[:limit]
		offerCount -= 1
	}

	// Process next
	t.Page = interface{}(nil)
	if total > 0 {
		t.Page = getPageValue(t.Objects[offerCount-1])
		if total > limit {
			t.CanMove = true
		}
	}
}

func GetObject(docPath string, t *TransferObject, f func(*firestore.DocumentSnapshot) interface{}) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(docPath)
	docSnapshot, err := docRef.Get(context.Background())

	if err == nil {
		t.Object = f(docSnapshot)
		t.Found = true
	} else {
		err := common.CheckNotFound(err)
		if err != nil {
			t.SetError(api_error.FirebaseError, errors.New(err))
		}
	}
}

func GetCacheObject(key string, t *TransferObject, f func(string) interface{}) {
	val, err := cache.RedisClient.Get(key).Result()
	if err == nil {
		t.Object = f(val)
		t.Found = true
	} else {
		if err != redis.Nil {
			t.SetError(api_error.GetDataFailed, err)
		}
	}
}

func ListPagingObjects(collectionPath string, t *TransferObject, limit int, startAt interface{}, q func(*firestore.CollectionRef) firestore.Query, f func(*firestore.DocumentSnapshot) interface{}) {
	dbClient := firebase_service.FirestoreClient

	collRef := dbClient.Collection(collectionPath)
	var query firestore.Query
	isPaging := limit != 0
	t.Found = true
	if q != nil {
		query = q(collRef)
		log.Println("Go Here 1")
		// Get extra 1 record to process paging
		if isPaging {
			query = query.Limit(limit + 1)
		}
		if startAt != nil {
			query = query.StartAfter(startAt)
		}
		docRef := query.Documents(context.Background())
		docSnapshots, err := docRef.GetAll()
		if err == nil {
			log.Println("Go Here 2")
			for _, docSnapshot := range docSnapshots {
				t.Objects = append(t.Objects, f(docSnapshot))
			}

			if isPaging {
				t.FeedPaging(limit)
			}
		} else {
			t.SetError(api_error.FirebaseError, errors.New(err))
		}
	} else {
		docSnapshots := collRef.Documents(context.Background())
		for {
			docSnapshot, err := docSnapshots.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				t.SetError(api_error.FirebaseError, errors.New(err))
				break
			} else {
				t.Objects = append(t.Objects, f(docSnapshot))
			}
		}
	}

	return
}

func ListComplexPagingObjects(collectionPath string, t *TransferObject, limit int, q func(*firestore.CollectionRef) firestore.Query,
	start func(firestore.Query) firestore.Query, getPageValue func(interface{}) interface{},
	f func(*firestore.DocumentSnapshot) interface{}) {

	dbClient := firebase_service.FirestoreClient

	collRef := dbClient.Collection(collectionPath)
	var query firestore.Query

	isPaging := limit != 0
	t.Found = true
	if q != nil {
		query = q(collRef)

		// Get extra 1 record to process paging
		if isPaging {
			query = query.Limit(limit + 1)
		}
		query = start(query)

		docRef := query.Documents(context.Background())
		docSnapshots, err := docRef.GetAll()
		if err == nil {
			for _, docSnapshot := range docSnapshots {
				t.Objects = append(t.Objects, f(docSnapshot))
			}

			if isPaging {
				t.FeedComplexPaging(limit, getPageValue)
			}
		} else {
			t.SetError(api_error.FirebaseError, errors.New(err))
		}
	} else {
		docSnapshots := collRef.Documents(context.Background())
		for {
			docSnapshot, err := docSnapshots.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				t.SetError(api_error.FirebaseError, errors.New(err))
				break
			} else {
				t.Objects = append(t.Objects, f(docSnapshot))
			}
		}
	}

	return
}

func ListObjects(collectionPath string, t *TransferObject, q func(*firestore.CollectionRef) firestore.Query, f func(*firestore.DocumentSnapshot) interface{}) {
	ListPagingObjects(collectionPath, t, 0, nil, q, f)
	return
}

func ConvertTo(amount decimal.Decimal, rate decimal.Decimal) decimal.Decimal {
	return amount.Mul(rate)
}

func ConvertFrom(amount decimal.Decimal, rate decimal.Decimal) decimal.Decimal {
	return amount.Div(rate)
}

func AddFeePercentage(amount decimal.Decimal, feePercentage decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	fee := amount.Mul(feePercentage)
	return amount.Add(fee), fee
}

func DeductFeePercentage(amount decimal.Decimal, feePercentage decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	fee := amount.Mul(feePercentage)
	return amount.Sub(fee), fee
}

func ClearCache(keyPattern string) {
	keys, err := cache.RedisClient.Keys(keyPattern).Result()
	if err == nil {
		for _, key := range keys {
			cache.RedisClient.Del(key)
		}
	}
}
