package appointments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAppointmentsHandler
// @Summary Get All Appointments
// @Description Retrieve a list of all appointments
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Appointment
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /appointments [get]
func GetAllAppointmentsHandler(c *gin.Context) {
	appointments, err := GetAllAppointmentsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// GetAppointmentHandler
// @Summary Get Appointment by ID
// @Description Retrieve a specific appointment
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /appointments/{id} [get]
func GetAppointmentHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	appointment, err := GetAppointmentByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// CreateAppointmentHandler
// @Summary Create a new Appointment
// @Description Create a new appointment
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAppointmentRequest true "Appointment Details"
// @Success 201 {object} Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /appointments [post]
func CreateAppointmentHandler(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := CreateAppointmentService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, appointment)
}

// UpdateAppointmentHandler
// @Summary Update Appointment
// @Description Update appointment details
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param request body UpdateAppointmentRequest true "Update Details"
// @Success 200 {object} Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /appointments/{id} [put]
func UpdateAppointmentHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := UpdateAppointmentService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// DeleteAppointmentHandler
// @Summary Delete Appointment
// @Description Soft delete an appointment
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /appointments/{id} [delete]
func DeleteAppointmentHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteAppointmentService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
