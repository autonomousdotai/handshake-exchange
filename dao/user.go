package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/shopspring/decimal"
)

type UserDaoInterface interface {
	GetProfile(userId string) (t TransferObject)
	AddProfile(profile bean.Profile) error
	UpdateProfileCreditCard(userId string, creditCard bean.UserCreditCard, userCCLimit bean.UserCreditCardLimit) error
	UpdateProfileOfferRejectLock(profile bean.Profile) error
	UpdateUserCCLimitAmount(userId string, token string, amount decimal.Decimal) error
	UpdateUserCCLimitTracks() (userIds []string, t TransferObject)
	GetCCLimit(userId string, token string) (t TransferObject)
	GetUserCCLimitEndTracks() (t TransferObject)
	UpgradeCCLimitLevel(userId string, token string, limit bean.UserCreditCardLimit) error
}

type UserDao struct {
}

func (dao UserDao) GetProfile(userId string) (t TransferObject) {
	// users/{uid}
	GetObject(GetUserPath(userId), &t, func(snapshot *firestore.DocumentSnapshot) interface{} {
		var obj bean.Profile
		snapshot.DataTo(&obj)
		return obj
	})

	return
}

func (dao UserDao) AddProfile(profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	// users/{uid}
	profileRef := dbClient.Collection("users").Doc(profile.UserId)
	batch.Set(profileRef, profile.GetAddProfile())
	_, err := batch.Commit(context.Background())

	return err
}

func (dao UserDao) UpdateProfileCreditCard(userId string, creditCard bean.UserCreditCard, userCCLimit bean.UserCreditCardLimit) error {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	profileRef := dbClient.Collection("users").Doc(userId)
	userCCLimitRef := dbClient.Doc(GetUserCCLimitItemPath(userId, creditCard.Token))
	userCCLimitTrackRef := dbClient.Doc(GetUserCCLimitTrackItemPath(userId))

	batch.Set(profileRef, creditCard.GetUpdateProfileCreditCard(), firestore.MergeAll)
	batch.Set(userCCLimitRef, userCCLimit.GetAddUserCreditCardLimit(), firestore.MergeAll)
	batch.Set(userCCLimitTrackRef, bean.UserCreditCardLimitTrack{
		UID:      userId,
		Level:    userCCLimit.Level,
		Duration: userCCLimit.Duration,
		Left:     userCCLimit.Duration,
	}.GetAddUserCreditCardLimitTrack())

	_, err := batch.Commit(context.Background())

	return err
}

func (dao UserDao) UpdateUserCCLimitAmount(userId string, token string, amount decimal.Decimal) error {
	dbClient := firebase_service.FirestoreClient
	// TODO: Until stripe save CC
	userCCLimitRef := dbClient.Doc(GetUserCCLimitItemPath(userId, userId))
	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(userCCLimitRef)
		if err != nil {
			return err
		}
		currentAmount, err := common.ConvertToDecimal(doc, "amount")
		if err != nil {
			return err
		}
		currentAmount = currentAmount.Add(amount)
		if currentAmount.LessThan(common.Zero) {
			currentAmount = common.Zero
		}
		return tx.Set(userCCLimitRef, bean.UserCreditCardLimit{Amount: currentAmount.String()}.GetUpdateAmount(), firestore.MergeAll)
	})

	return err
}

func (dao UserDao) UpdateProfileOfferRejectLock(profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()

	profileRef := dbClient.Collection("users").Doc(profile.UserId)
	profileRef.Set(context.Background(), profile.GetUpdateOfferRejectLock(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao UserDao) GetCCLimit(userId string, token string) (t TransferObject) {
	// users/{uid}/cc_limit/{token}
	GetObject(GetUserCCLimitItemPath(userId, token), &t, snapshotUserCCLimit)
	return
}

func (dao UserDao) UpgradeCCLimitLevel(userId string, token string, limit bean.UserCreditCardLimit) error {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	_, err := batch.Commit(context.Background())
	docRef := dbClient.Doc(GetUserCCLimitItemPath(userId, token))
	trackDocRef := dbClient.Doc(GetUserCCLimitTrackItemPath(userId))

	batch.Set(docRef, limit.GetUpdateLevel(), firestore.MergeAll)
	batch.Set(trackDocRef, bean.UserCreditCardLimitTrack{
		UID:      userId,
		Level:    limit.Level,
		Duration: limit.Duration,
		Left:     limit.Duration,
	}.GetAddUserCreditCardLimitTrack())
	batch.Commit(context.Background())

	return err
}

func (dao UserDao) GetUserCCLimitEndTracks() (t TransferObject) {
	ListObjects(GetUserCCLimitTracksPath(), &t, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.Where("left", "==", 1)
		return query
	}, snapshotToUserCCLimitTrack)

	return
}

func (dao UserDao) UpdateUserCCLimitTracks() (userIds []string, t TransferObject) {
	ListObjects(GetUserCCLimitTracksPath(), &t, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.Where("left", ">", 1).OrderBy("left", firestore.Asc)
		return query
	}, snapshotToUserCCLimitTrack)

	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	for _, obj := range t.Objects {
		track := obj.(bean.UserCreditCardLimitTrack)
		trackDocRef := dbClient.Doc(GetUserCCLimitTrackItemPath(track.UID))
		track.Left -= 1

		batch.Set(trackDocRef, track.GetUpdateLeft(), firestore.MergeAll)

		userIds = append(userIds, track.UID)
	}
	batch.Commit(context.Background())

	return
}

func GetUserPath(userId string) string {
	return fmt.Sprintf("users/%s", userId)
}

func GetUserCCLimitItemPath(userId string, token string) string {
	return fmt.Sprintf("users/%s/cc_limit/%s", userId, token)
}

func GetUserCCLimitTracksPath() string {
	return fmt.Sprintf("user_cc_limit_tracks")
}

func GetUserCCLimitTrackItemPath(userId string) string {
	return fmt.Sprintf("user_cc_limit_tracks/%s", userId)
}

func snapshotUserCCLimit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.UserCreditCardLimit
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToUserCCLimitTrack(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.UserCreditCardLimitTrack
	snapshot.DataTo(&obj)
	return obj
}
