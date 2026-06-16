package attendance_logs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAttendanceLogsHandler
// @Summary Get All Attendance Logs
// @Description Retrieve a list of all attendance logs
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} AttendanceLog
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance-logs [get]
func GetAllAttendanceLogsHandler(c *gin.Context) {
	items, err := GetAllAttendanceLogsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAttendanceLogHandler
// @Summary Get Attendance Log by ID
// @Description Retrieve a specific attendance log
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Success 200 {object} AttendanceLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /attendance-logs/{id} [get]
func GetAttendanceLogHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetAttendanceLogByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateAttendanceLogHandler
// @Summary Create a new Attendance Log (Check-in)
// @Description Create a new attendance log
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAttendanceLogRequest true "Log Details"
// @Success 201 {object} AttendanceLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance-logs [post]
func CreateAttendanceLogHandler(c *gin.Context) {
	var req CreateAttendanceLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateAttendanceLogService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateAttendanceLogHandler
// @Summary Update Attendance Log (Check-out)
// @Description Update attendance log details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Param request body UpdateAttendanceLogRequest true "Update Details"
// @Success 200 {object} AttendanceLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance-logs/{id} [put]
func UpdateAttendanceLogHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateAttendanceLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateAttendanceLogService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteAttendanceLogHandler
// @Summary Delete Attendance Log
// @Description Delete an attendance log
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Log ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance-logs/{id} [delete]
func DeleteAttendanceLogHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteAttendanceLogService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
