package chats

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Chat Sessions
func GetAllChatSessions() ([]ChatSession, error) {
	var results []ChatSession
	db := supabase.DB
	err := db.Order("created_at DESC").Find(&results).Error
	return results, err
}

func GetChatSessionByID(id string) (*ChatSession, error) {
	var result ChatSession
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat session not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateChatSession(data *ChatSession) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateChatSession(data *ChatSession) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteChatSession(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ChatSession{}).Error
}

// Chat Messages
func GetMessagesBySessionID(sessionID string) ([]ChatMessage, error) {
	var results []ChatMessage
	db := supabase.DB
	err := db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&results).Error
	return results, err
}

func CreateChatMessage(data *ChatMessage) error {
	db := supabase.DB
	return db.Create(data).Error
}
