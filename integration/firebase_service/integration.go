package firebase_service

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"
)

var firebaseService *firebase.App
var notificationFirebaseService *firebase.App

// var AuthClient *auth.Client
var FirestoreClient *firestore.Client
var NotificationFirebaseClient *db.Client

func NewFirestore(credFile string) *firebase.App {
	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	return app
}

func NewFirebase(credFile string) *firebase.App {
	opt := option.WithCredentialsFile(credFile)
	conf := &firebase.Config{
		DatabaseURL: os.Getenv("NOTIFICATION_FBDB"),
	}
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		panic(err)
	}

	return app
}

//func NewAuthClient() *auth.Client {
//	authClient, err := firebaseService.Auth(context.Background())
//	if err != nil {
//		panic(err)
//	}
//
//	return authClient
//}

func NewFirestoreClient() *firestore.Client {
	client, err := firebaseService.Firestore(context.Background())
	if err != nil {
		panic(err)
	}

	return client
}

func NewFirebaseClient() *db.Client {
	client, _ := notificationFirebaseService.Database(context.Background())

	return client
}

func Intialize() {
	firebaseService = NewFirestore("./credentials/cred.json")
	notificationFirebaseService = NewFirebase("./credentials/cred.json")

	// AuthClient = NewAuthClient()
	FirestoreClient = NewFirestoreClient()
	NotificationFirebaseClient = NewFirebaseClient()
}
