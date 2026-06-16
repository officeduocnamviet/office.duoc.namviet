package user_notifications

import "time"

// FCM Tokens
func GetAllFCMTokensService() ([]UserFCMToken, error) {
	return GetAllFCMTokens()
}

func GetFCMTokenByIDService(id int64) (*UserFCMToken, error) {
	return GetFCMTokenByID(id)
}

func CreateFCMTokenService(req CreateFCMTokenRequest) (*UserFCMToken, error) {
	token := &UserFCMToken{
		UserID:     req.UserID,
		FCMToken:   req.FCMToken,
		DeviceID:   req.DeviceID,
		DeviceType: req.DeviceType,
	}

	if err := CreateFCMToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func UpdateFCMTokenService(id int64, req UpdateFCMTokenRequest) (*UserFCMToken, error) {
	token, err := GetFCMTokenByID(id)
	if err != nil {
		return nil, err
	}

	if req.FCMToken != nil {
		token.FCMToken = *req.FCMToken
	}
	if req.DeviceID != nil {
		token.DeviceID = req.DeviceID
	}
	if req.DeviceType != nil {
		token.DeviceType = req.DeviceType
	}

	now := time.Now()
	token.UpdatedAt = &now

	if err := UpdateFCMToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func DeleteFCMTokenService(id int64) error {
	return DeleteFCMToken(id)
}

// Social Mappings
func GetAllSocialMappingsService() ([]UserSocialMapping, error) {
	return GetAllSocialMappings()
}

func GetSocialMappingByIDService(id int64) (*UserSocialMapping, error) {
	return GetSocialMappingByID(id)
}

func CreateSocialMappingService(req CreateSocialMappingRequest) (*UserSocialMapping, error) {
	mapping := &UserSocialMapping{
		UserID:         req.UserID,
		SocialProvider: req.SocialProvider,
		SocialID:       req.SocialID,
		SocialName:     req.SocialName,
		SocialAvatar:   req.SocialAvatar,
	}

	if err := CreateSocialMapping(mapping); err != nil {
		return nil, err
	}
	return mapping, nil
}

func DeleteSocialMappingService(id int64) error {
	return DeleteSocialMapping(id)
}
