package agent_workflows

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAgentWorkflows() ([]AgentWorkflow, error) {
	var results []AgentWorkflow
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetAgentWorkflowByID(id string) (*AgentWorkflow, error) {
	var result AgentWorkflow
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("agent workflow not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateAgentWorkflow(data *AgentWorkflow) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateAgentWorkflow(data *AgentWorkflow) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteAgentWorkflow(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&AgentWorkflow{}).Error
}
