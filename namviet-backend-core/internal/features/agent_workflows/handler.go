package agent_workflows

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAgentWorkflowsHandler
// @Summary Get All Agent Workflows
// @Description Retrieve a list of all agent workflows
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} AgentWorkflow
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agent-workflows [get]
func GetAllAgentWorkflowsHandler(c *gin.Context) {
	items, err := GetAllAgentWorkflowsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAgentWorkflowHandler
// @Summary Get Agent Workflow by ID
// @Description Retrieve a specific agent workflow
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workflow ID"
// @Success 200 {object} AgentWorkflow
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agent-workflows/{id} [get]
func GetAgentWorkflowHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetAgentWorkflowByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateAgentWorkflowHandler
// @Summary Create a new Agent Workflow
// @Description Create a new agent workflow
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAgentWorkflowRequest true "Workflow Details"
// @Success 201 {object} AgentWorkflow
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agent-workflows [post]
func CreateAgentWorkflowHandler(c *gin.Context) {
	var req CreateAgentWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateAgentWorkflowService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateAgentWorkflowHandler
// @Summary Update Agent Workflow
// @Description Update agent workflow details
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workflow ID"
// @Param request body UpdateAgentWorkflowRequest true "Update Details"
// @Success 200 {object} AgentWorkflow
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agent-workflows/{id} [put]
func UpdateAgentWorkflowHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateAgentWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateAgentWorkflowService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteAgentWorkflowHandler
// @Summary Delete Agent Workflow
// @Description Delete an agent workflow
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workflow ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agent-workflows/{id} [delete]
func DeleteAgentWorkflowHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteAgentWorkflowService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
