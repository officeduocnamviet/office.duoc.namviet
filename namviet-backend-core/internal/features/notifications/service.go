package notifications

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/messaging"
	"github.com/namviet/backend-core/internal/platform/firebase"
)

// SendPushNotification sends an FCM message to a specific device token
func SendPushNotification(deviceToken, title, body string) error {
	if firebase.MessagingClient == nil {
		return errors.New("firebase messaging client is not initialized")
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: deviceToken,
	}

	response, err := firebase.MessagingClient.Send(context.Background(), message)
	if err != nil {
		return err
	}

	_ = response // Log or return response if needed
	return nil
}
