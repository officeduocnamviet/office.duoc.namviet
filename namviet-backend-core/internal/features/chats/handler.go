package chats

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Chat Sessions API ---

// GetAllChatSessionsHandler
// @Summary Get All Chat Sessions
// @Description Retrieve a list of all chat sessions
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ChatSession
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-sessions [get]
func GetAllChatSessionsHandler(c *gin.Context) {
	items, err := GetAllChatSessionsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetChatSessionHandler
// @Summary Get Chat Session by ID
// @Description Retrieve a specific chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} ChatSession
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /chat-sessions/{id} [get]
func GetChatSessionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetChatSessionByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateChatSessionHandler
// @Summary Create a new Chat Session
// @Description Create a new chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateChatSessionRequest true "Session Details"
// @Success 201 {object} ChatSession
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-sessions [post]
func CreateChatSessionHandler(c *gin.Context) {
	var req CreateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateChatSessionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateChatSessionHandler
// @Summary Update Chat Session
// @Description Update chat session details
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Param request body UpdateChatSessionRequest true "Update Details"
// @Success 200 {object} ChatSession
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-sessions/{id} [put]
func UpdateChatSessionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateChatSessionService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteChatSessionHandler
// @Summary Delete Chat Session
// @Description Delete a chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-sessions/{id} [delete]
func DeleteChatSessionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteChatSessionService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Chat Messages API ---

// GetMessagesBySessionIDHandler
// @Summary Get Messages by Session ID
// @Description Retrieve messages for a specific session
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {array} ChatMessage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-sessions/{id}/messages [get]
func GetMessagesBySessionIDHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Session ID format"})
		return
	}
	items, err := GetMessagesBySessionIDService(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateChatMessageHandler
// @Summary Create a new Chat Message
// @Description Send a message in a chat session
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateChatMessageRequest true "Message Details"
// @Success 201 {object} ChatMessage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat-messages [post]
func CreateChatMessageHandler(c *gin.Context) {
	var req CreateChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateChatMessageService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}
