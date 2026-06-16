package user_notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- FCM Tokens API ---

// GetAllFCMTokensHandler
// @Summary Get All FCM Tokens
// @Description Retrieve a list of all user FCM tokens
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} UserFCMToken
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fcm-tokens [get]
func GetAllFCMTokensHandler(c *gin.Context) {
	items, err := GetAllFCMTokensService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetFCMTokenHandler
// @Summary Get FCM Token by ID
// @Description Retrieve a specific FCM token
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Token ID"
// @Success 200 {object} UserFCMToken
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fcm-tokens/{id} [get]
func GetFCMTokenHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetFCMTokenByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateFCMTokenHandler
// @Summary Create a new FCM Token
// @Description Create a new FCM token
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFCMTokenRequest true "Token Details"
// @Success 201 {object} UserFCMToken
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fcm-tokens [post]
func CreateFCMTokenHandler(c *gin.Context) {
	var req CreateFCMTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateFCMTokenService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateFCMTokenHandler
// @Summary Update FCM Token
// @Description Update FCM token details
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Token ID"
// @Param request body UpdateFCMTokenRequest true "Update Details"
// @Success 200 {object} UserFCMToken
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fcm-tokens/{id} [put]
func UpdateFCMTokenHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateFCMTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateFCMTokenService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteFCMTokenHandler
// @Summary Delete FCM Token
// @Description Delete an FCM token
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Token ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fcm-tokens/{id} [delete]
func DeleteFCMTokenHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteFCMTokenService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Social Mappings API ---

// GetAllSocialMappingsHandler
// @Summary Get All Social Mappings
// @Description Retrieve a list of all social mappings
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} UserSocialMapping
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /social-mappings [get]
func GetAllSocialMappingsHandler(c *gin.Context) {
	items, err := GetAllSocialMappingsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetSocialMappingHandler
// @Summary Get Social Mapping by ID
// @Description Retrieve a specific social mapping
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Mapping ID"
// @Success 200 {object} UserSocialMapping
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /social-mappings/{id} [get]
func GetSocialMappingHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetSocialMappingByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateSocialMappingHandler
// @Summary Create a new Social Mapping
// @Description Create a new social mapping
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSocialMappingRequest true "Mapping Details"
// @Success 201 {object} UserSocialMapping
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /social-mappings [post]
func CreateSocialMappingHandler(c *gin.Context) {
	var req CreateSocialMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateSocialMappingService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// DeleteSocialMappingHandler
// @Summary Delete Social Mapping
// @Description Delete a social mapping
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Mapping ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /social-mappings/{id} [delete]
func DeleteSocialMappingHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteSocialMappingService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
