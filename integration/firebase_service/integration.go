package firebase_service

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var firebaseService = NewFirebase()
var AuthClient = NewAuthClient()
var FirestoreClient = NewFirestoreClient()

func NewFirebase() *firebase.App {
	opt := option.WithCredentialsFile("./credentials/cred.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	return app
}

func NewAuthClient() *auth.Client {
	authClient, err := firebaseService.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	return authClient
}

func NewFirestoreClient() *firestore.Client {
	client, err := firebaseService.Firestore(context.Background())
	if err != nil {
		panic(err)
	}

	return client
}
