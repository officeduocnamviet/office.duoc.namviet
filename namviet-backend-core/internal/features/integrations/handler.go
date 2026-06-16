package integrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Connections API ---

// GetAllConnectionsHandler
// @Summary Get All Connections
// @Description Retrieve a list of all connections
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ThirdPartyConnection
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /connections [get]
func GetAllConnectionsHandler(c *gin.Context) {
	items, err := GetAllConnectionsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetConnectionHandler
// @Summary Get Connection by ID
// @Description Retrieve a specific connection
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Connection ID"
// @Success 200 {object} ThirdPartyConnection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /connections/{id} [get]
func GetConnectionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetConnectionByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateConnectionHandler
// @Summary Create a new Connection
// @Description Create a new connection
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateConnectionRequest true "Connection Details"
// @Success 201 {object} ThirdPartyConnection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /connections [post]
func CreateConnectionHandler(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateConnectionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateConnectionHandler
// @Summary Update Connection
// @Description Update connection details
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Connection ID"
// @Param request body UpdateConnectionRequest true "Update Details"
// @Success 200 {object} ThirdPartyConnection
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /connections/{id} [put]
func UpdateConnectionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateConnectionService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteConnectionHandler
// @Summary Delete Connection
// @Description Delete a connection
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Connection ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /connections/{id} [delete]
func DeleteConnectionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteConnectionService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Webhook Logs API ---

// GetAllWebhookLogsHandler
// @Summary Get All Webhook Logs
// @Description Retrieve a list of webhook logs
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} WebhookLog
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /webhook-logs [get]
func GetAllWebhookLogsHandler(c *gin.Context) {
	items, err := GetAllWebhookLogsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetWebhookLogHandler
// @Summary Get Webhook Log by ID
// @Description Retrieve a specific webhook log
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Success 200 {object} WebhookLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /webhook-logs/{id} [get]
func GetWebhookLogHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetWebhookLogByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateWebhookLogHandler
// @Summary Create a new Webhook Log
// @Description Create a new webhook log
// @Tags Integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWebhookLogRequest true "Log Details"
// @Success 201 {object} WebhookLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /webhook-logs [post]
func CreateWebhookLogHandler(c *gin.Context) {
	var req CreateWebhookLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateWebhookLogService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}
