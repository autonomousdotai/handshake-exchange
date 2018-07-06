package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/integration/firebase_service"
	"google.golang.org/api/iterator"
	"strings"
)

type OfferDao struct {
}

func (dao OfferDao) AddOffer(offer bean.Offer, profile bean.Profile) (bean.Offer, error) {
	dbClient := firebase_service.FirestoreClient

	// profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Collection(GetOfferPath()).NewDoc()
	offer.Id = docRef.ID

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetAddOffer())
	// batch.Set(profileDocRef, profile.GetUpdateOfferProfile(), firestore.MergeAll)

	offerPath := GetOfferItemPath(offer.Id)
	if offer.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  offer.SystemAddress,
			Offer:    offer.Id,
			OfferRef: offerPath,
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	if offer.Currency == bean.ETH.Code && (offer.Status == bean.OFFER_STATUS_CREATED) {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offer.Status,
			Currency: offer.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return offer, err
}

func (dao OfferDao) ListOffers(userId string, offerType string, currency string, status string, limit int, startAt interface{}) (t TransferObject) {
	ListPagingObjects(GetOfferPath(), &t, limit, startAt, func(collRef *firestore.CollectionRef) firestore.Query {
		query := collRef.Where("uid", "==", userId)
		if offerType != "" {
			query = query.Where("type", "==", offerType)
		}
		if status != "" {
			query = query.Where("status", "==", status)
		}
		if currency != "" {
			query = query.Where("currency", "==", currency)
		}
		query = query.OrderBy("created_at", firestore.Desc)
		return query
	}, snapshotToOffer)

	return
}

func (dao OfferDao) ListTransferMaps() ([]bean.OfferTransferMap, error) {
	dbClient := firebase_service.FirestoreClient

	// pending_instant_offers
	iter := dbClient.Collection(GetOfferTransferMapPath()).Documents(context.Background())
	offers := make([]bean.OfferTransferMap, 0)

	for {
		var offer bean.OfferTransferMap
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

func (dao OfferDao) GetOffer(offerId string) (t TransferObject) {
	// offers/{id}
	GetObject(GetOfferItemPath(offerId), &t, snapshotToOffer)

	return
}

func (dao OfferDao) GetOfferByPath(path string) (t TransferObject) {
	// offers/{id}
	GetObject(path, &t, snapshotToOffer)

	return
}

func (dao OfferDao) UpdateOffer(offer bean.Offer, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferItemPath(offer.Id)
	docRef := dbClient.Doc(offerPath)

	batch := dbClient.Batch()
	batch.Set(docRef, updateData, firestore.MergeAll)

	if offer.SystemAddress != "" &&
		(offer.Status == bean.OFFER_STATUS_CREATE_FAILED ||
			offer.Status == bean.OFFER_STATUS_PRE_SHAKE_FAILED) {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	if offer.Currency == bean.ETH.Code &&
		(offer.Status == bean.OFFER_STATUS_CREATED ||
			offer.Status == bean.OFFER_STATUS_PRE_SHAKING ||
			offer.Status == bean.OFFER_STATUS_SHAKING ||
			offer.Status == bean.OFFER_STATUS_CANCELLING ||
			offer.Status == bean.OFFER_STATUS_REJECTING ||
			offer.Status == bean.OFFER_STATUS_CLOSING ||
			offer.Status == bean.OFFER_STATUS_COMPLETING) {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offer.Status,
			Currency: offer.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferActive(offer bean.Offer) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferActive(), firestore.MergeAll)
	if offer.SystemAddress != "" {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferShaking(offer bean.Offer) error {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferItemPath(offer.Id)
	docRef := dbClient.Doc(offerPath)

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferShake(), firestore.MergeAll)
	if offer.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  offer.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferItemPath(offer.Id),
			UID:      offer.UID,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
	}

	if offer.Currency == bean.ETH.Code && (offer.Status == bean.OFFER_STATUS_PRE_SHAKING || offer.Status == bean.OFFER_STATUS_SHAKING) {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offer.Status,
			Currency: offer.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferShake(offer bean.Offer) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))
	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferShake(), firestore.MergeAll)
	if offer.SystemAddress != "" {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferClose(offer bean.Offer, profile bean.Profile) error {
	dbClient := firebase_service.FirestoreClient

	// profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	offerPath := GetOfferItemPath(offer.Id)
	docRef := dbClient.Doc(offerPath)

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferClose(), firestore.MergeAll)
	// batch.Set(profileDocRef, profile.GetUpdateOfferProfile(), firestore.MergeAll)

	if offer.SystemAddress != "" {
		addressMapDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.SystemAddress))
		batch.Delete(addressMapDocRef)
	}

	if offer.Currency == bean.ETH.Code && (offer.Status == bean.OFFER_STATUS_CREATED) {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offer.Status,
			Currency: offer.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferReject(offer bean.Offer, profile bean.Profile, transactionCount bean.TransactionCount) error {
	dbClient := firebase_service.FirestoreClient

	offerPath := GetOfferItemPath(offer.Id)
	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))
	transCountDocRef := dbClient.Doc(GetTransactionCountItemPath(offer.UID, offer.Currency))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferReject(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferProfile(), firestore.MergeAll)
	batch.Set(transCountDocRef, transactionCount.GetUpdateFailed(), firestore.MergeAll)

	if offer.Currency == bean.ETH.Code && (offer.Status == bean.OFFER_STATUS_REJECTING) {
		// Store a record to check onchain
		docId := strings.Replace(offerPath, "/", "-", -1)
		onChainTrackingRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, docId))
		batch.Set(onChainTrackingRef, bean.OfferOnChainActionTracking{
			Id:       docId,
			Action:   offer.Status,
			Currency: offer.Currency,
			Offer:    offer.Id,
			OfferRef: offerPath,
			Type:     bean.OFFER_ADDRESS_MAP_OFFER,
			UID:      offer.UID,
		}.GetAddOfferOnChainActionTracking())
	}

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferCompleted(offer bean.Offer, profile bean.Profile, transactionCount bean.TransactionCount) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))
	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	// transferDocRef := dbClient.Doc(GetOfferTransferMapItemPath(offer.Id))
	transCountDocRef := dbClient.Doc(GetTransactionCountItemPath(offer.UID, offer.Currency))

	trans1, trans2 := bean.NewTransactionFromOfferHandshake(offer)
	trans1DocRef := dbClient.Collection(GetTransactionPath(offer.UID)).NewDoc()
	trans2DocRef := dbClient.Collection(GetTransactionPath(offer.ToUID)).NewDoc()

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferCompleted(), firestore.MergeAll)
	batch.Set(profileDocRef, profile.GetUpdateOfferProfile(), firestore.MergeAll)
	batch.Set(transCountDocRef, transactionCount.GetUpdateSuccess(), firestore.MergeAll)
	batch.Set(trans1DocRef, trans1.GetAddTransaction(), firestore.MergeAll)
	batch.Set(trans2DocRef, trans2.GetAddTransaction(), firestore.MergeAll)
	// batch.Delete(transferDocRef)

	_, err := batch.Commit(context.Background())

	return err
}

func (dao OfferDao) UpdateOfferWithdraw(offer bean.Offer) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferWithdraw(), firestore.MergeAll)

	_, err := batch.Commit(context.Background())

	return err
}

//Firebase
func (dao OfferDao) UpdateNotificationOffer(offer bean.Offer) error {
	dbClient := firebase_service.NotificationFirebaseClient

	ref := dbClient.NewRef(GetNotificationOfferItemPath(offer.UID, offer.Id))
	err := ref.Set(context.Background(), offer.GetNotificationUpdate())
	if offer.ToUID != "" {
		ref = dbClient.NewRef(GetNotificationOfferItemPath(offer.ToUID, offer.Id))
		err = ref.Set(context.Background(), offer.GetNotificationUpdate())
	}

	return err
}

func (dao OfferDao) GetOfferAddress(address string) (t TransferObject) {
	// offer_addresses/{id}
	GetObject(GetOfferAddressMapItemPath(address), &t, snapshotToOfferAddressMap)

	return
}

func (dao OfferDao) UpdateTickTransferMap(transferMap bean.OfferTransferMap) {
	dbClient := firebase_service.FirestoreClient

	transferDocRef := dbClient.Doc(GetOfferTransferMapItemPath(transferMap.Offer))
	transferDocRef.Set(context.Background(), transferMap.GetUpdateTick(), firestore.MergeAll)
}

func (dao OfferDao) ListOfferConfirmingAddressMap() ([]bean.OfferConfirmingAddressMap, error) {
	dbClient := firebase_service.FirestoreClient

	// pending_instant_offers
	iter := dbClient.Collection(GetOfferConfirmingAddressMapPath()).Documents(context.Background())
	offers := make([]bean.OfferConfirmingAddressMap, 0)

	for {
		var offer bean.OfferConfirmingAddressMap
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

func (dao OfferDao) AddOfferConfirmingAddressMap(offerMap bean.OfferConfirmingAddressMap) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetOfferConfirmingAddressMapItemPath(offerMap.TxHash))

	_, err := docRef.Set(context.Background(), offerMap.GetAddOfferConfirmingAddressMap(), firestore.MergeAll)

	return err
}

func (dao OfferDao) RemoveOfferConfirmingAddressMap(txHash string) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetOfferConfirmingAddressMapItemPath(txHash))

	_, err := docRef.Delete(context.Background())

	return err
}

func (dao OfferDao) ListCryptoPendingTransfer() ([]bean.CryptoPendingTransfer, error) {
	dbClient := firebase_service.FirestoreClient

	// pending_instant_offers
	iter := dbClient.Collection(GetCryptoPendingTransferPath()).Documents(context.Background())
	transfers := make([]bean.CryptoPendingTransfer, 0)

	for {
		var offer bean.CryptoPendingTransfer
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return transfers, err
		}
		doc.DataTo(&offer)
		transfers = append(transfers, offer)
	}

	return transfers, nil
}

func (dao OfferDao) RemoveCryptoPendingTransfer(id string) error {
	dbClient := firebase_service.FirestoreClient
	docRef := dbClient.Doc(GetCryptoPendingTransferItemPath(id))

	_, err := docRef.Delete(context.Background())

	return err
}

func (dao OfferDao) ListOfferOnChainActionTracking(isOriginal bool) ([]bean.OfferOnChainActionTracking, error) {
	dbClient := firebase_service.FirestoreClient

	iter := dbClient.Collection(GetOfferOnChainActionTrackingPath(isOriginal)).Documents(context.Background())
	offers := make([]bean.OfferOnChainActionTracking, 0)

	for {
		var offer bean.OfferOnChainActionTracking
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

func (dao OfferDao) AddOfferOnChainActionTracking(offerTracking bean.OfferOnChainActionTracking) error {
	dbClient := firebase_service.FirestoreClient
	id := strings.Replace(offerTracking.OfferRef, "/", "-", -1)
	docRef := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(false, id))
	offerTracking.Id = id

	_, err := docRef.Set(context.Background(), offerTracking.GetAddOfferOnChainActionTracking(), firestore.MergeAll)

	return err
}

func (dao OfferDao) RemoveOfferOnChainActionTracking(id string, all bool) error {
	dbClient := firebase_service.FirestoreClient
	docRef2 := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(true, id))

	batch := dbClient.Batch()
	if all {
		docRef1 := dbClient.Doc(GetOfferOnChainActionTrackingItemPath(false, id))
		batch.Delete(docRef1)
	}
	batch.Delete(docRef2)
	_, err := batch.Commit(context.Background())

	return err
}

// DB path
func GetOfferPath() string {
	return "offers"
}

func GetOfferItemPath(id string) string {
	return fmt.Sprintf("offers/%s", id)
}

func GetNotificationOfferItemPath(userId string, id string) string {
	return fmt.Sprintf("users/%s/offers/exchange_%s", userId, id)
}

//func GetOfferAddressMapPath() string {
//	return "offer_addresses"
//}

func GetOfferAddressMapItemPath(id string) string {
	return fmt.Sprintf("offer_addresses/%s", id)
}

func GetOfferConfirmingAddressMapPath() string {
	return "offer_confirming_addresses"
}

func GetOfferConfirmingAddressMapItemPath(id string) string {
	return fmt.Sprintf("offer_confirming_addresses/%s", id)
}

func GetOfferTransferMapPath() string {
	return "offer_transfers"
}

func GetOfferTransferMapItemPath(id string) string {
	return fmt.Sprintf("offer_transfers/%s", id)
}

func GetOfferOnChainActionTrackingPath(isOriginal bool) string {
	str := "original"
	if !isOriginal {
		str = "check"
	}
	return fmt.Sprintf("offer_onchain_action_trackings/%s/trackings", str)
}

func GetOfferOnChainActionTrackingItemPath(isOriginal bool, id string) string {
	return fmt.Sprintf("%s/%s", GetOfferOnChainActionTrackingPath(isOriginal), id)
}

func snapshotToOffer(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.Offer
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}

func snapshotToOfferAddressMap(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferAddressMap
	snapshot.DataTo(&obj)
	return obj
}

func snapshotToOfferOnChainActionTracking(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.OfferOnChainActionTracking
	snapshot.DataTo(&obj)
	return obj
}
