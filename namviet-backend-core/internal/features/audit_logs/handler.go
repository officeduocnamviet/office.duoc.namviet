package audit_logs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAuditLogsHandler
// @Summary Get All Audit Logs (Limit 100)
// @Description Retrieve a list of recent audit logs
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} SystemAuditLog
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-logs [get]
func GetAllAuditLogsHandler(c *gin.Context) {
	items, err := GetAllAuditLogsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAuditLogHandler
// @Summary Get Audit Log by ID
// @Description Retrieve a specific audit log
// @Tags Audit Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit Log ID"
// @Success 200 {object} SystemAuditLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /audit-logs/{id} [get]
func GetAuditLogHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetAuditLogByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}
