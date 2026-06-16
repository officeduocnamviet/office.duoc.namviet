package internal_communications

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Internal Channels
func GetAllInternalChannels() ([]InternalChannel, error) {
	var results []InternalChannel
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetInternalChannelByID(id int64) (*InternalChannel, error) {
	var result InternalChannel
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("internal channel not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateInternalChannel(data *InternalChannel) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateInternalChannel(data *InternalChannel) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteInternalChannel(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&InternalChannel{}).Error
}

// Internal Messages
func GetMessagesByChannelID(channelID int64) ([]InternalMessage, error) {
	var results []InternalMessage
	db := supabase.DB
	err := db.Where("channel_id = ?", channelID).Order("created_at ASC").Find(&results).Error
	return results, err
}

func CreateInternalMessage(data *InternalMessage) error {
	db := supabase.DB
	return db.Create(data).Error
}

func DeleteInternalMessage(id int64) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&InternalMessage{}).Error
}
