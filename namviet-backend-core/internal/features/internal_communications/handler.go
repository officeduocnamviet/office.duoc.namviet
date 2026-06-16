package internal_communications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Internal Channels API ---

// GetAllInternalChannelsHandler
// @Summary Get All Internal Channels
// @Description Retrieve a list of all internal channels
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} InternalChannel
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-channels [get]
func GetAllInternalChannelsHandler(c *gin.Context) {
	items, err := GetAllInternalChannelsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetInternalChannelHandler
// @Summary Get Internal Channel by ID
// @Description Retrieve a specific internal channel
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Channel ID"
// @Success 200 {object} InternalChannel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /internal-channels/{id} [get]
func GetInternalChannelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetInternalChannelByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateInternalChannelHandler
// @Summary Create a new Internal Channel
// @Description Create a new internal channel
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateInternalChannelRequest true "Channel Details"
// @Success 201 {object} InternalChannel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-channels [post]
func CreateInternalChannelHandler(c *gin.Context) {
	var req CreateInternalChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateInternalChannelService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateInternalChannelHandler
// @Summary Update Internal Channel
// @Description Update internal channel details
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Channel ID"
// @Param request body UpdateInternalChannelRequest true "Update Details"
// @Success 200 {object} InternalChannel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-channels/{id} [put]
func UpdateInternalChannelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateInternalChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateInternalChannelService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteInternalChannelHandler
// @Summary Delete Internal Channel
// @Description Delete an internal channel
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Channel ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-channels/{id} [delete]
func DeleteInternalChannelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteInternalChannelService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Internal Messages API ---

// GetMessagesByChannelIDHandler
// @Summary Get Messages By Channel ID
// @Description Retrieve a list of messages for a channel
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Channel ID"
// @Success 200 {array} InternalMessage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-channels/{id}/messages [get]
func GetMessagesByChannelIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	items, err := GetMessagesByChannelIDService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateInternalMessageHandler
// @Summary Create a new Internal Message
// @Description Create a new internal message
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateInternalMessageRequest true "Message Details"
// @Success 201 {object} InternalMessage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-messages [post]
func CreateInternalMessageHandler(c *gin.Context) {
	var req CreateInternalMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateInternalMessageService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// DeleteInternalMessageHandler
// @Summary Delete Internal Message
// @Description Delete an internal message
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Message ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal-messages/{id} [delete]
func DeleteInternalMessageHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteInternalMessageService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
