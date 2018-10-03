package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/common"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"github.com/shopspring/decimal"
	"google.golang.org/api/iterator"
)

type CreditCardDao struct {
}

func (dao CreditCardDao) AddCCTransaction(ccTran bean.CCTransaction) (bean.CCTransaction, error) {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	userDocRef := dbClient.Collection(GetUserCCTransactionPath(ccTran.UID)).NewDoc()
	ccTran.Id = userDocRef.ID
	docRef := dbClient.Doc(GetCCTransactionItemPath(fmt.Sprintf("%s_%s", ccTran.UID, ccTran.Id)))

	batch.Set(userDocRef, ccTran.GetAddCCTransaction())
	batch.Set(docRef, ccTran.GetAddCCTransaction())

	_, err := batch.Commit(context.Background())

	return ccTran, err
}

func (dao CreditCardDao) UpdateCCTransaction(ccTran bean.CCTransaction) (bean.CCTransaction, error) {
	dbClient := firebase_service.FirestoreClient
	batch := dbClient.Batch()

	userDocRef := dbClient.Doc(GetUserCCTransactionItemPath(ccTran.UID, ccTran.Id))
	docRef := dbClient.Doc(GetCCTransactionItemPath(fmt.Sprintf("%s_%s", ccTran.UID, ccTran.Id)))

	batch.Set(userDocRef, ccTran.GetUpdateCCTransaction(), firestore.MergeAll)
	batch.Set(docRef, ccTran.GetUpdateCCTransaction(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return ccTran, err
}

func (dao CreditCardDao) UpdateCCTransactionStatus(ccTran bean.CCTransaction) (bean.CCTransaction, error) {
	dbClient := firebase_service.FirestoreClient

	batch := dbClient.Batch()
	userDocRef := dbClient.Doc(GetUserCCTransactionItemPath(ccTran.UID, ccTran.Id))
	docRef := dbClient.Doc(GetCCTransactionItemPath(fmt.Sprintf("%s_%s", ccTran.UID, ccTran.Id)))

	batch.Set(userDocRef, ccTran.GetUpdateStatus(), firestore.MergeAll)
	batch.Set(docRef, ccTran.GetUpdateStatus(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return ccTran, err
}

func (dao CreditCardDao) ListCCTransactions(userId string, limit int, startAt interface{}) (t TransferObject) {
	ListPagingObjects(GetUserCCTransactionPath(userId), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.OrderBy("created_at", firestore.Desc)
		return query
	}, snapshotToCCTransaction)

	return
}

func (dao CreditCardDao) GetCCTransaction(userId string, ccTranId string) TransferObject {
	return dao.GetCCTransactionByPath(GetUserCCTransactionItemPath(userId, ccTranId))
}

func (dao CreditCardDao) GetCCTransactionByPath(path string) (t TransferObject) {
	// users/{uid}/cc_transactions/{id}
	GetObject(path, &t, snapshotToCCTransaction)
	return
}

func (dao CreditCardDao) AddInstantOffer(offer bean.InstantOffer, transaction bean.Transaction, providerId string) (bean.InstantOffer, error) {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Collection(GetInstantOfferPath(offer.UID)).NewDoc()
	offer.Id = docRef.ID

	pendingOffer := bean.PendingInstantOffer{
		UID:             offer.UID,
		InstantOffer:    offer.Id,
		InstantOfferRef: GetInstantOfferItemPath(offer.UID, offer.Id),
		Duration:        offer.Duration,
		Provider:        offer.Provider,
		ProviderId:      providerId,
		CCMode:          offer.CCMode,
	}
	pendingOfferId := fmt.Sprintf("%s-%s", offer.UID, offer.Id)
	docPendingRef := dbClient.Doc(GetPendingInstantOfferItemPath(pendingOfferId))
	pendingOffer.Id = pendingOfferId

	docTransactionRef := dbClient.Collection(GetTransactionPath(offer.UID)).NewDoc()
	offer.TransactionRef = GetTransactionItemPath(offer.UID, docTransactionRef.ID)

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetAddInstantOffer())
	batch.Set(docPendingRef, pendingOffer.GetAddInstantOffer())
	batch.Set(docTransactionRef, transaction.GetAddTransaction())
	_, err := batch.Commit(context.Background())

	return offer, err
}

func (dao CreditCardDao) UpdateInstantOffer(offer bean.InstantOffer, transaction bean.Transaction) (bean.InstantOffer, error) {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetInstantOfferItemPath(offer.UID, offer.Id))

	pendingOfferId := fmt.Sprintf("%s-%s", offer.UID, offer.Id)
	pendingOfferDocRef := dbClient.Doc(GetPendingInstantOfferItemPath(pendingOfferId))

	docTransactionRef := dbClient.Doc(offer.TransactionRef)

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdate(), firestore.MergeAll)
	batch.Delete(pendingOfferDocRef)
	batch.Set(docTransactionRef, transaction.GetUpdateStatus(), firestore.MergeAll)
	_, err := batch.Commit(context.Background())

	return offer, err
}

func (dao CreditCardDao) ListInstantOffers(userId string, currency string, limit int, startAt interface{}) (t TransferObject) {
	//ListPagingObjects(GetInstantOfferPath(userId), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
	//	query := collRef.Where("currency", "==", currency).OrderBy("created_at", firestore.Desc)
	//	return query
	//}, snapshotToInstantOffer)

	ListPagingObjects(GetInstantOfferPath(userId), &t, limit, startAt, nil, snapshotToInstantOffer)

	return
}

func (dao CreditCardDao) GetInstantOffer(userId string, instantOfferId string) TransferObject {
	return dao.GetInstantOfferByPath(GetInstantOfferItemPath(userId, instantOfferId))
}

func (dao CreditCardDao) GetInstantOfferByPath(path string) (t TransferObject) {
	// users/{uid}/instant_offers/{id}
	GetObject(path, &t, snapshotToInstantOffer)
	return
}

func (dao CreditCardDao) ListPendingInstantOffer() ([]bean.PendingInstantOffer, error) {
	dbClient := firebase_service.FirestoreClient

	// pending_instant_offers
	iter := dbClient.Collection(GetPendingInstantOfferPath()).Documents(context.Background())
	offers := make([]bean.PendingInstantOffer, 0)

	for {
		var offer bean.PendingInstantOffer
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return offers, err
		}
		doc.DataTo(&offer)
		offers = append(offers, offer)
	}

	return offers, nil
}

func (dao CreditCardDao) UpdateNotificationInstantOffer(offer bean.InstantOffer) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationInstantOfferItemPath(offer.UID, offer.Id))
	err := ref.Set(context.Background(), offer.GetNotificationUpdate())

	return err
}

func (dao CreditCardDao) GetCCGlobalLimit() TransferObject {
	return dao.GetCCGlobalLimitByPath(GetGlobalCCLimitPath())
}

func (dao CreditCardDao) GetCCGlobalLimitByPath(path string) (t TransferObject) {
	GetObject(path, &t, snapshotToGlobalCCLimit)
	return
}

func (dao CreditCardDao) UpdateCCGlobalLimitAmount(amount decimal.Decimal) error {
	dbClient := firebase_service.FirestoreClient
	limitRef := dbClient.Doc(GetGlobalCCLimitPath())
	err := dbClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(limitRef)
		if err != nil {
			return err
		}
		usage, err := common.ConvertToDecimal(doc, "usage")
		if err != nil {
			return err
		}
		usage = usage.Add(amount)
		if usage.LessThan(common.Zero) {
			usage = common.Zero
		}
		return tx.Set(limitRef, bean.GlobalCCLimit{Usage: usage.String()}.GetUpdateUsage(), firestore.MergeAll)
	})

	return err
}

func (dao CreditCardDao) GetPendingInstantOfferTransfer(id string) (t TransferObject) {
	GetObject(GetPendingInstantOfferTransferItemPath(id), &t, snapshotToPendingInstantOfferTransfer)
	return
}

func (dao CreditCardDao) ListPendingInstantOfferTransfer() (t TransferObject) {
	ListObjects(GetPendingInstantOfferTransferPath(), &t, nil, snapshotToPendingInstantOfferTransfer)
	return
}

func (dao CreditCardDao) AddPendingInstantOfferTransfer(pendingTransfer *bean.PendingInstantOfferTransfer) (err error) {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Collection(GetPendingInstantOfferTransferPath()).NewDoc()
	pendingTransfer.Id = docRef.ID

	_, err = docRef.Set(context.Background(), pendingTransfer.GetAdd())
	return err
}

func (dao CreditCardDao) UpdatePendingInstantOfferTransfer(pendingTransfer *bean.PendingInstantOfferTransfer) (err error) {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetPendingInstantOfferTransferItemPath(pendingTransfer.Id))

	_, err = docRef.Set(context.Background(), pendingTransfer.GetUpdate(), firestore.MergeAll)
	return err
}

func (dao CreditCardDao) RemovePendingInstantOfferTransfer(pendingTransfer *bean.PendingInstantOfferTransfer, offer *bean.InstantOffer) (err error) {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetPendingInstantOfferTransferItemPath(pendingTransfer.Id))
	offerDocRef := dbClient.Doc(GetInstantOfferItemPath(offer.UID, offer.Id))

	batch := dbClient.Batch()

	batch.Set(offerDocRef, offer.GetUpdate(), firestore.MergeAll)
	batch.Delete(docRef)

	_, err = batch.Commit(context.Background())

	return err
}

func (dao CreditCardDao) UpdateNotificationInitInstantOffer(providerId string, data map[string]interface{}) error {
	dbClient := firebase_service.NotificationFirebaseClient
	ref := dbClient.NewRef(GetNotificationInitInstantOfferItemPath(providerId))
	err := ref.Set(context.Background(), data)
	fmt.Println(err)

	return err
}

func (dao CreditCardDao) GetInitInstantOffer(providerId string) (t TransferObject) {
	GetObject(GetInitInstantOfferItemPath(providerId), &t, snapshotToInitInstantOffer)
	return
}

func (dao CreditCardDao) AddInitInstantOffer(providerId string, data map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	path := GetInitInstantOfferItemPath(providerId)
	docRef := dbClient.Doc(path)
	_, err := docRef.Set(context.Background(), data)

	return err
}

func GetUserCCTransactionPath(userId string) string {
	return fmt.Sprintf("users/%s/cc_transactions", userId)
}

func GetUserCCTransactionItemPath(userId string, id string) string {
	return fmt.Sprintf("%s/%s", GetUserCCTransactionPath(userId), id)
}

func GetCCTransactionPath() string {
	return fmt.Sprintf("cc_transactions")
}

func GetCCTransactionItemPath(id string) string {
	return fmt.Sprintf("%s/%s", GetCCTransactionPath(), id)
}

func GetInstantOfferPath(userId string) string {
	return fmt.Sprintf("users/%s/instant_offers", userId)
}

func GetInstantOfferItemPath(userId string, id string) string {
	return fmt.Sprintf("%s/%s", GetInstantOfferPath(userId), id)
}

func GetPendingInstantOfferPath() string {
	return fmt.Sprintf("pending_instant_offers")
}

func GetPendingInstantOfferTransferPath() string {
	return fmt.Sprintf("pending_instant_offer_transfers")
}

func GetPendingInstantOfferTransferItemPath(id string) string {
	return fmt.Sprintf("pending_instant_offer_transfers/%s", id)
}

func GetPendingInstantOfferItemPath(pendingOfferId string) string {
	return fmt.Sprintf("%s/%s", GetPendingInstantOfferPath(), pendingOfferId)
}

func GetGlobalCCLimitPath() string {
	return "cc_global_limit/1"
}

func GetInitInstantOfferItemPath(providerId string) string {
	return fmt.Sprintf("init_instant_offers/%s", providerId)
}

// Firebase
func GetNotificationInstantOfferItemPath(userId string, offerId string) string {
	return fmt.Sprintf("users/%s/offers/instant_%s", userId, offerId)
}

func GetNotificationInitInstantOfferItemPath(providerId string) string {
	return fmt.Sprintf("users/000/init_offers/%s", providerId)
}

func snapshotToCCTransaction(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.CCTransaction
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID

	return obj
}

func snapshotToInstantOffer(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.InstantOffer
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID

	return obj
}

func snapshotToGlobalCCLimit(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.GlobalCCLimit
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToPendingInstantOfferTransfer(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.PendingInstantOfferTransfer
	snapshot.DataTo(&obj)

	return obj
}

func snapshotToInitInstantOffer(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj map[string]string
	snapshot.DataTo(&obj)

	return obj
}
