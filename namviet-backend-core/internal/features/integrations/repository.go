package integrations

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Third Party Connections
func GetAllConnections() ([]ThirdPartyConnection, error) {
	var results []ThirdPartyConnection
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetConnectionByID(id string) (*ThirdPartyConnection, error) {
	var result ThirdPartyConnection
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("connection not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateConnection(data *ThirdPartyConnection) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateConnection(data *ThirdPartyConnection) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteConnection(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ThirdPartyConnection{}).Error
}

// Webhook Logs
func GetAllWebhookLogs() ([]WebhookLog, error) {
	var results []WebhookLog
	db := supabase.DB
	err := db.Order("created_at DESC").Limit(100).Find(&results).Error
	return results, err
}

func GetWebhookLogByID(id string) (*WebhookLog, error) {
	var result WebhookLog
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("webhook log not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateWebhookLog(data *WebhookLog) error {
	db := supabase.DB
	return db.Create(data).Error
}
