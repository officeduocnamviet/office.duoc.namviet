package time_attendance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllTimeAttendancesHandler
// @Summary Get All Time Attendances
// @Description Retrieve a list of all time attendances
// @Tags Time Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} TimeAttendance
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /time-attendance [get]
func GetAllTimeAttendancesHandler(c *gin.Context) {
	tas, err := GetAllTimeAttendancesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tas)
}

// GetTimeAttendanceHandler
// @Summary Get Time Attendance by ID
// @Description Retrieve a specific time attendance record
// @Tags Time Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Time Attendance ID"
// @Success 200 {object} TimeAttendance
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /time-attendance/{id} [get]
func GetTimeAttendanceHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	ta, err := GetTimeAttendanceByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ta)
}

// CreateTimeAttendanceHandler
// @Summary Create a new Time Attendance
// @Description Create a new time attendance record
// @Tags Time Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTimeAttendanceRequest true "Time Attendance Details"
// @Success 201 {object} TimeAttendance
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /time-attendance [post]
func CreateTimeAttendanceHandler(c *gin.Context) {
	var req CreateTimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ta, err := CreateTimeAttendanceService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ta)
}

// UpdateTimeAttendanceHandler
// @Summary Update Time Attendance
// @Description Update time attendance details
// @Tags Time Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Time Attendance ID"
// @Param request body UpdateTimeAttendanceRequest true "Update Details"
// @Success 200 {object} TimeAttendance
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /time-attendance/{id} [put]
func UpdateTimeAttendanceHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateTimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ta, err := UpdateTimeAttendanceService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ta)
}

// DeleteTimeAttendanceHandler
// @Summary Delete Time Attendance
// @Description Delete a time attendance record
// @Tags Time Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Time Attendance ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /time-attendance/{id} [delete]
func DeleteTimeAttendanceHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteTimeAttendanceService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
