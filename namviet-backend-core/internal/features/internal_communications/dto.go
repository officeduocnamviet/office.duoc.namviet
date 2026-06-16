package internal_communications

// Channel DTOs
type CreateInternalChannelRequest struct {
	Name string  `json:"name" binding:"required"`
	Type *string `json:"type"`
}

type UpdateInternalChannelRequest struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
}

// Message DTOs
type CreateInternalMessageRequest struct {
	ChannelID int64  `json:"channel_id" binding:"required"`
	SenderID  string `json:"sender_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
