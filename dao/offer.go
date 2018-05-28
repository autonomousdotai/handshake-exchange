package dao

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/autonomousdotai/handshake-exchange/bean"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
)

type OfferDao struct {
}

func (dao OfferDao) AddOffer(offer bean.Offer, profile bean.Profile) (bean.Offer, error) {
	dbClient := firebase_service.FirestoreClient

	profileDocRef := dbClient.Doc(GetUserPath(offer.UID))
	docRef := dbClient.Collection(GetOfferPath()).NewDoc()
	offer.Id = docRef.ID

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetAddOffer())
	batch.Set(profileDocRef, profile.GetUpdateOfferProfile(), firestore.MergeAll)

	if offer.SystemAddress != "" {
		mapping := bean.OfferAddressMap{
			Address:  offer.SystemAddress,
			Offer:    offer.Id,
			OfferRef: GetOfferItemPath(offer.Id),
			UID:      offer.UID,
		}
		mappingDocRef := dbClient.Doc(GetOfferAddressMapItemPath(offer.Id))
		batch.Set(mappingDocRef, mapping.GetAddOfferAddressMap())
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

func (dao OfferDao) GetOffer(offerId string) (t TransferObject) {
	// offers/{id}
	GetObject(GetOfferItemPath(offerId), &t, snapshotToOffer)

	return
}

func (dao OfferDao) UpdateOffer(offer bean.Offer, updateData map[string]interface{}) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))
	_, err := docRef.Set(context.Background(), updateData, firestore.MergeAll)

	return err
}

func (dao OfferDao) UpdateOfferCompleting(offer bean.Offer, externalId string) error {
	dbClient := firebase_service.FirestoreClient

	docRef := dbClient.Doc(GetOfferItemPath(offer.Id))
	transferDocRef := dbClient.Collection(GetOfferTransferMapItemPath(offer.Id)).NewDoc()

	batch := dbClient.Batch()
	batch.Set(docRef, offer.GetUpdateOfferCompleting())
	batch.Set(transferDocRef, bean.OfferTransferMap{
		UID:        offer.UID,
		Address:    offer.UserAddress,
		Offer:      offer.Id,
		OfferRef:   GetOfferItemPath(offer.Id),
		ExternalId: externalId,
	})

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

func GetOfferAddressMapPath() string {
	return "offer_addresses"
}

func GetOfferAddressMapItemPath(id string) string {
	return fmt.Sprintf("offer_addresses/%s", id)
}

func GetOfferTransferMapPath() string {
	return "offer_transfers"
}

func GetOfferTransferMapItemPath(id string) string {
	return fmt.Sprintf("offer_transfers/%s", id)
}

func snapshotToOffer(snapshot *firestore.DocumentSnapshot) interface{} {
	var obj bean.Offer
	snapshot.DataTo(&obj)
	obj.Id = snapshot.Ref.ID
	return obj
}
