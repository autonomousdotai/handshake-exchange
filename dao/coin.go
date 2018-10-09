package dao

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
)

type CoinDao struct {
}

func (dao CoinDao) ListCoinCenter(country string) (t TransferObject) {
	ListObjects(GetCoinCenterCountryCurrenyPath(country), &t, nil, snapshotToCoinCenter)
	return
}

func GetCoinCenterCountryCurrenyPath(country string) string {
	return fmt.Sprintf("coin_centers/%s/currency", country)
}

func snapshotToCoinCenter(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CoinCenter
	snapshot.DataTo(&obj)
	return obj
}
