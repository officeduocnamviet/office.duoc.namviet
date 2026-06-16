package ai_agent_memories

func GetAllAIAgentMemoriesService() ([]AIAgentMemory, error) {
	return GetAllAIAgentMemories()
}

func GetAIAgentMemoryByIDService(id string) (*AIAgentMemory, error) {
	return GetAIAgentMemoryByID(id)
}

func CreateAIAgentMemoryService(req CreateAIAgentMemoryRequest) (*AIAgentMemory, error) {
	memory := &AIAgentMemory{
		SessionID: req.SessionID,
		UserID:    req.UserID,
		Key:       req.Key,
		Value:     req.Value,
		ExpiresAt: req.ExpiresAt,
	}

	if err := CreateAIAgentMemory(memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func UpdateAIAgentMemoryService(id string, req UpdateAIAgentMemoryRequest) (*AIAgentMemory, error) {
	memory, err := GetAIAgentMemoryByID(id)
	if err != nil {
		return nil, err
	}

	if req.SessionID != nil {
		memory.SessionID = req.SessionID
	}
	if req.UserID != nil {
		memory.UserID = req.UserID
	}
	if req.Key != nil {
		memory.Key = *req.Key
	}
	if req.Value != nil {
		memory.Value = *req.Value
	}
	if req.ExpiresAt != nil {
		memory.ExpiresAt = req.ExpiresAt
	}

	if err := UpdateAIAgentMemory(memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func DeleteAIAgentMemoryService(id string) error {
	return DeleteAIAgentMemory(id)
}
