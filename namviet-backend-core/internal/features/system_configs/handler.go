package system_configs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllSystemConfigsHandler
// @Summary Get All System Configs
// @Description Retrieve a list of all system configs
// @Tags System Configs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} SystemConfig
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /system-configs [get]
func GetAllSystemConfigsHandler(c *gin.Context) {
	items, err := GetAllSystemConfigsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetSystemConfigHandler
// @Summary Get System Config by Key
// @Description Retrieve a specific system config by config_key
// @Tags System Configs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Config Key"
// @Success 200 {object} SystemConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /system-configs/{key} [get]
func GetSystemConfigHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Key format"})
		return
	}
	item, err := GetSystemConfigByKeyService(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateSystemConfigHandler
// @Summary Create a new System Config
// @Description Create a new system config
// @Tags System Configs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSystemConfigRequest true "Config Details"
// @Success 201 {object} SystemConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /system-configs [post]
func CreateSystemConfigHandler(c *gin.Context) {
	var req CreateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateSystemConfigService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateSystemConfigHandler
// @Summary Update System Config
// @Description Update system config details
// @Tags System Configs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Config Key"
// @Param request body UpdateSystemConfigRequest true "Update Details"
// @Success 200 {object} SystemConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /system-configs/{key} [put]
func UpdateSystemConfigHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Key format"})
		return
	}
	var req UpdateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateSystemConfigService(key, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteSystemConfigHandler
// @Summary Delete System Config
// @Description Delete a system config
// @Tags System Configs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Config Key"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /system-configs/{key} [delete]
func DeleteSystemConfigHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Key format"})
		return
	}
	
	if err := DeleteSystemConfigService(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
