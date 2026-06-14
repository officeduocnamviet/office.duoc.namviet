package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var MessagingClient *messaging.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("configs/firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase App: %v\n", err)
		return
	}

	MessagingClient, err = app.Messaging(context.Background())
	if err != nil {
		log.Printf("Error getting Messaging client: %v\n", err)
		return
	}

	log.Println("Firebase Admin SDK initialized successfully")
}
