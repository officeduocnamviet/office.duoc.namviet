package chats

import (
	"time"

	"github.com/namviet/backend-core/internal/features/user_notifications"
	"github.com/namviet/backend-core/internal/platform/firebase"
)

// Chat Sessions
func GetAllChatSessionsService() ([]ChatSession, error) {
	return GetAllChatSessions()
}

func GetChatSessionByIDService(id string) (*ChatSession, error) {
	return GetChatSessionByID(id)
}

func CreateChatSessionService(req CreateChatSessionRequest) (*ChatSession, error) {
	session := &ChatSession{
		CustomerID: req.CustomerID,
		AgentID:    req.AgentID,
		Status:     "active",
	}

	if err := CreateChatSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func UpdateChatSessionService(id string, req UpdateChatSessionRequest) (*ChatSession, error) {
	session, err := GetChatSessionByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		session.Status = *req.Status
		if *req.Status == "closed" || *req.Status == "ended" {
			now := time.Now()
			session.EndedAt = &now
		}
	}

	if err := UpdateChatSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func DeleteChatSessionService(id string) error {
	return DeleteChatSession(id)
}

// Chat Messages
func GetMessagesBySessionIDService(sessionID string) ([]ChatMessage, error) {
	return GetMessagesBySessionID(sessionID)
}

func CreateChatMessageService(req CreateChatMessageRequest) (*ChatMessage, error) {
	msgType := "text"
	if req.MessageType != nil {
		msgType = *req.MessageType
	}

	message := &ChatMessage{
		SessionID:   req.SessionID,
		SenderType:  req.SenderType,
		SenderID:    req.SenderID,
		Content:     req.Content,
		MessageType: msgType,
		Metadata:    req.Metadata,
	}

	if err := CreateChatMessage(message); err != nil {
		return nil, err
	}

	if req.SessionID != nil {
		session, err := GetChatSessionByID(*req.SessionID)
		if err == nil && session != nil {
			go sendPushNotificationForChat(message, session)
		}
	}

	return message, nil
}

func sendPushNotificationForChat(message *ChatMessage, session *ChatSession) {
	var targetID string
	var targetType string

	if message.SenderType == "customer" {
		if session.AgentID == nil {
			return
		}
		targetID = *session.AgentID
		targetType = "employee"
	} else if message.SenderType == "agent" || message.SenderType == "employee" {
		if session.CustomerID == nil {
			return
		}
		targetID = *session.CustomerID
		targetType = "retail_customer" // Can be enhanced later to check if wholesale
	} else {
		return // Do not send for system/bot messages
	}

	tokens, err := user_notifications.GetFCMTokensByTargetService(targetID, targetType)
	if err != nil || len(tokens) == 0 {
		return
	}

	var tokenStrings []string
	for _, t := range tokens {
		if t.FCMToken != "" {
			tokenStrings = append(tokenStrings, t.FCMToken)
		}
	}

	if len(tokenStrings) > 0 {
		firebase.SendMulticastNotification(tokenStrings, "Tin nhắn mới", message.Content, map[string]string{
			"session_id": *message.SessionID,
			"type":       "chat_message",
		})
	}
}
