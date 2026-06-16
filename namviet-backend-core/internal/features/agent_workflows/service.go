package agent_workflows

import "time"

func GetAllAgentWorkflowsService() ([]AgentWorkflow, error) {
	return GetAllAgentWorkflows()
}

func GetAgentWorkflowByIDService(id string) (*AgentWorkflow, error) {
	return GetAgentWorkflowByID(id)
}

func CreateAgentWorkflowService(req CreateAgentWorkflowRequest) (*AgentWorkflow, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	workflow := &AgentWorkflow{
		Name:        req.Name,
		Description: req.Description,
		TriggerType: req.TriggerType,
		Steps:       req.Steps,
		IsActive:    isActive,
	}

	if err := CreateAgentWorkflow(workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

func UpdateAgentWorkflowService(id string, req UpdateAgentWorkflowRequest) (*AgentWorkflow, error) {
	workflow, err := GetAgentWorkflowByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		workflow.Name = *req.Name
	}
	if req.Description != nil {
		workflow.Description = req.Description
	}
	if req.TriggerType != nil {
		workflow.TriggerType = *req.TriggerType
	}
	if req.Steps != nil {
		workflow.Steps = *req.Steps
	}
	if req.IsActive != nil {
		workflow.IsActive = *req.IsActive
	}

	now := time.Now()
	workflow.UpdatedAt = &now

	if err := UpdateAgentWorkflow(workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

func DeleteAgentWorkflowService(id string) error {
	return DeleteAgentWorkflow(id)
}
