package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var MessagingClient *messaging.Client

func InitFirebase() {
	var opt option.ClientOption

	// 1. Ưu tiên đọc từ Biến môi trường (Dành cho Cloud Run)
	firebaseJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	if firebaseJSON != "" {
		opt = option.WithCredentialsJSON([]byte(firebaseJSON))
	} else {
		// 2. Dự phòng đọc từ file (Dành cho Local chạy máy tính)
		opt = option.WithCredentialsFile("configs/firebase-adminsdk.json")
	}

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
