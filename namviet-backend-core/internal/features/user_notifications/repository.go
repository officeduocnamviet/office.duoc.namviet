package user_notifications

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// FCM Tokens
func GetAllFCMTokens() ([]UserFCMToken, error) {
	var results []UserFCMToken
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetFCMTokenByID(id string) (*UserFCMToken, error) {
	var result UserFCMToken
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fcm token not found")
		}
		return nil, err
	}
	return &result, nil
}

func GetFCMTokensByTarget(targetID string, targetType string) ([]UserFCMToken, error) {
	var results []UserFCMToken
	db := supabase.DB
	err := db.Where("target_id = ? AND target_type = ?", targetID, targetType).Find(&results).Error
	return results, err
}

func CreateFCMToken(data *UserFCMToken) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateFCMToken(data *UserFCMToken) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteFCMToken(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&UserFCMToken{}).Error
}

// Social Mappings
func GetAllSocialMappings() ([]UserSocialMapping, error) {
	var results []UserSocialMapping
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetSocialMappingByID(id int64) (*UserSocialMapping, error) {
	var result UserSocialMapping
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("social mapping not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateSocialMapping(data *UserSocialMapping) error {
	db := supabase.DB
	return db.Create(data).Error
}

func DeleteSocialMapping(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&UserSocialMapping{}).Error
}
