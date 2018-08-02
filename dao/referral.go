package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"strings"
	"time"
)

type ReferralDao struct {
}

func (d ReferralDao) ListReferralSummary(userId string) (t TransferObject) {
	ListObjects(GerReferralOfferStoreSummaryPath(userId), &t, nil, snapshotToReferralOfferStoreShake)

	return
}

func (d ReferralDao) AddReferral(userId string, checkExists bool) error {
	dbClient := firebase_service.FirestoreClient

	referralDocRef := dbClient.Doc(GetReferralItemPath(userId))
	_, errCheck := referralDocRef.Get(context.Background())
	if errCheck != nil && strings.Contains(errCheck.Error(), "not found") {
		// If there is no node, create one
		_, errAdd := referralDocRef.Set(context.Background(), bean.ReferralCount{
			UID:   userId,
			Count: 0,
		}.GetAddData())
		if errAdd != nil {
			return errAdd
		}
	} else {
		if !checkExists {
			return errors.New("referral exists")
		}
	}

	return nil
}

func (d ReferralDao) AddReferralRecord(record bean.ReferralRecord) error {
	dbClient := firebase_service.FirestoreClient

	referralDocRef := dbClient.Doc(GetReferralItemPath(record.UID))
	referralRecordDocRef := dbClient.Doc(GetReferralRecordItemPath(record.UID, record.ToUID))

	errAdd := d.AddReferral(record.UID, true)
	if errAdd != nil {
		return errAdd
	}

	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		referralDoc, err := tx.Get(referralDocRef)
		if err != nil {
			return err
		}

		countField := "count"
		valueTmp, err := referralDoc.DataAt(countField)
		value := valueTmp.(int64)
		value += 1
		err = tx.Update(referralDocRef, []firestore.Update{
			{Path: countField, Value: value},
		})

		err = tx.Set(referralRecordDocRef, record.GetAddData(), firestore.MergeAll)
		return err
	})

	return err
}

func (d ReferralDao) AddReferralOfferStoreShake(referralCreatedAt time.Time, offerShake bean.ReferralOfferStoreShakeRecord) error {
	dbClient := firebase_service.FirestoreClient

	referralDocRef := dbClient.Doc(GetReferralOfferStoreShakeItemPath(offerShake.UID, offerShake.Currency))
	referralRecordDocRef := dbClient.Doc(GetReferralOfferStoreShakeRecordItemPath(offerShake.UID, offerShake.Currency, offerShake.OfferShake))
	referralSummaryDocRef := dbClient.Doc(GerReferralOfferStoreSummaryItemPath(offerShake.UID,
		fmt.Sprintf("%s-%s", offerShake.ToUID, offerShake.Currency)))

	_, errCheck := referralDocRef.Get(context.Background())
	if errCheck != nil && strings.Contains(errCheck.Error(), "not found") {
		// If there is no node, create one
		_, errAdd := referralDocRef.Set(context.Background(), bean.ReferralOfferStoreShake{
			UID:               offerShake.UID,
			ToUID:             offerShake.ToUID,
			ToUsername:        offerShake.ToUsername,
			Currency:          offerShake.Currency,
			Reward:            common.Zero.String(),
			PendingReward:     common.Zero.String(),
			TotalReward:       common.Zero.String(),
			ReferralCreatedAt: referralCreatedAt,
		}.GetAddData())
		if errAdd != nil {
			return errAdd
		}
	}

	reward := common.StringToDecimal(offerShake.Reward)
	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		referralDoc, err := tx.Get(referralDocRef)
		if err != nil {
			return err
		}
		totalRewardField := "total_reward"
		totalReward, err := common.ConvertToDecimal(referralDoc, totalRewardField)
		if err != nil {
			return err
		}
		pendingRewardField := "pending_reward"
		pendingReward, err := common.ConvertToDecimal(referralDoc, pendingRewardField)
		if err != nil {
			return err
		}
		totalReward = totalReward.Add(reward)
		pendingReward = pendingReward.Add(reward)

		err = tx.Update(referralDocRef, []firestore.Update{
			{Path: totalRewardField, Value: totalReward.String()},
			{Path: pendingRewardField, Value: pendingReward.String()},
		})
		if err != nil {
			return err
		}

		err = tx.Set(referralSummaryDocRef, bean.ReferralOfferStoreShake{
			UID:               offerShake.UID,
			ToUID:             offerShake.ToUID,
			ToUsername:        offerShake.ToUsername,
			Currency:          offerShake.Currency,
			PendingReward:     pendingReward.String(),
			TotalReward:       totalReward.String(),
			ReferralCreatedAt: referralCreatedAt,
		}.GetOverridePendingReward(), firestore.MergeAll)
		if err != nil {
			return err
		}

		err = tx.Set(referralRecordDocRef, offerShake.GetAddData(), firestore.MergeAll)
		return err
	})

	return err
}

func GetReferralItemPath(userId string) string {
	return fmt.Sprintf("referrals/%s", userId)
}

func GetReferralRecordItemPath(userId string, toUserId string) string {
	return fmt.Sprintf("referrals/%s/records/%s", userId, toUserId)
}

func GetReferralOfferStoreShakeItemPath(userId string, currency string) string {
	return fmt.Sprintf("referrals/%s/currencies/%s", userId, currency)
}

func GetReferralOfferStoreShakeRecordPath(userId string, currency string) string {
	return fmt.Sprintf("referrals/%s/currencies/%s/records", userId, currency)
}

func GetReferralOfferStoreShakeRecordItemPath(userId string, currency string, offerShakeId string) string {
	return fmt.Sprintf("referrals/%s/currencies/%s/records/%s", userId, currency, offerShakeId)
}

func GerReferralOfferStoreSummaryPath(userId string) string {
	return fmt.Sprintf("referrals/%s/summary", userId)
}

func GerReferralOfferStoreSummaryItemPath(userId string, toUserId string) string {
	return fmt.Sprintf("referrals/%s/summary/%s", userId, toUserId)
}

func snapshotToReferralOfferStoreShake(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.ReferralOfferStoreShake
	snapshot.DataTo(&obj)

	return obj
}
