package chats

import "time"

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
	return message, nil
}
