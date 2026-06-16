package ai_agent_memories

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAIAgentMemories() ([]AIAgentMemory, error) {
	var results []AIAgentMemory
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetAIAgentMemoryByID(id string) (*AIAgentMemory, error) {
	var result AIAgentMemory
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ai agent memory not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateAIAgentMemory(data *AIAgentMemory) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateAIAgentMemory(data *AIAgentMemory) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteAIAgentMemory(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&AIAgentMemory{}).Error
}
