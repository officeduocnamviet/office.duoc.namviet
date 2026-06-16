package chats

import "github.com/namviet/backend-core/internal/features/roles"

// Chat Session DTOs
type CreateChatSessionRequest struct {
	CustomerID *string `json:"customer_id"`
	AgentID    *string `json:"agent_id"`
}

type UpdateChatSessionRequest struct {
	Status *string `json:"status"`
}

// Chat Message DTOs
type CreateChatMessageRequest struct {
	SessionID   *string      `json:"session_id"`
	SenderType  string       `json:"sender_type" binding:"required"`
	SenderID    *string      `json:"sender_id"`
	Content     string       `json:"content" binding:"required"`
	MessageType *string      `json:"message_type"`
	Metadata    *roles.JSONB `json:"metadata"`
}
