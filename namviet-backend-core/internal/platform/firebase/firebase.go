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
		credPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
		if credPath == "" {
			credPath = "internal/services/firebase-adminsdk.json"
		}
		opt = option.WithCredentialsFile(credPath)
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

// SendMulticastNotification sends a push notification to multiple device tokens
func SendMulticastNotification(tokens []string, title string, body string, data map[string]string) error {
	if MessagingClient == nil {
		return nil // Avoid crash if Firebase is not initialized
	}
	
	if len(tokens) == 0 {
		return nil
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}

	response, err := MessagingClient.SendMulticast(context.Background(), message)
	if err != nil {
		log.Printf("Error sending FCM multicast message: %v\n", err)
		return err
	}

	log.Printf("Successfully sent %d FCM messages, %d failed\n", response.SuccessCount, response.FailureCount)
	return nil
}
