package user_notifications

// FCM Token DTOs
type CreateFCMTokenRequest struct {
	TargetID   string  `json:"target_id" binding:"required"`
	TargetType string  `json:"target_type" binding:"required"` // employee, retail_customer, wholesale_customer
	FCMToken   string  `json:"fcm_token" binding:"required"`
	DeviceID   *string `json:"device_id"`
	DeviceType *string `json:"device_type"`
}

type UpdateFCMTokenRequest struct {
	FCMToken   *string `json:"fcm_token"`
	DeviceID   *string `json:"device_id"`
	DeviceType *string `json:"device_type"`
}

// Social Mapping DTOs
type CreateSocialMappingRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	SocialProvider string  `json:"social_provider" binding:"required"`
	SocialID       string  `json:"social_id" binding:"required"`
	SocialName     *string `json:"social_name"`
	SocialAvatar   *string `json:"social_avatar"`
}
