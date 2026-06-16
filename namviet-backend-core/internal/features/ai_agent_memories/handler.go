package ai_agent_memories

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAIAgentMemoriesHandler
// @Summary Get All AI Agent Memories
// @Description Retrieve a list of all ai agent memories
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} AIAgentMemory
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ai-agent-memories [get]
func GetAllAIAgentMemoriesHandler(c *gin.Context) {
	items, err := GetAllAIAgentMemoriesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAIAgentMemoryHandler
// @Summary Get AI Agent Memory by ID
// @Description Retrieve a specific ai agent memory
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Memory ID"
// @Success 200 {object} AIAgentMemory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /ai-agent-memories/{id} [get]
func GetAIAgentMemoryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetAIAgentMemoryByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateAIAgentMemoryHandler
// @Summary Create a new AI Agent Memory
// @Description Create a new ai agent memory
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAIAgentMemoryRequest true "Memory Details"
// @Success 201 {object} AIAgentMemory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ai-agent-memories [post]
func CreateAIAgentMemoryHandler(c *gin.Context) {
	var req CreateAIAgentMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateAIAgentMemoryService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateAIAgentMemoryHandler
// @Summary Update AI Agent Memory
// @Description Update ai agent memory details
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Memory ID"
// @Param request body UpdateAIAgentMemoryRequest true "Update Details"
// @Success 200 {object} AIAgentMemory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ai-agent-memories/{id} [put]
func UpdateAIAgentMemoryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateAIAgentMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateAIAgentMemoryService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteAIAgentMemoryHandler
// @Summary Delete AI Agent Memory
// @Description Delete an ai agent memory
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Memory ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ai-agent-memories/{id} [delete]
func DeleteAIAgentMemoryHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteAIAgentMemoryService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
