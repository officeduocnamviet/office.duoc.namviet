package medical_visits

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllMedicalVisitsHandler
// @Summary Get All Medical Visits
// @Description Retrieve a list of all medical visits
// @Tags Medical Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} MedicalVisit
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-visits [get]
func GetAllMedicalVisitsHandler(c *gin.Context) {
	visits, err := GetAllMedicalVisitsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, visits)
}

// GetMedicalVisitHandler
// @Summary Get Medical Visit by ID
// @Description Retrieve a specific medical visit
// @Tags Medical Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Medical Visit ID"
// @Success 200 {object} MedicalVisit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /medical-visits/{id} [get]
func GetMedicalVisitHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	visit, err := GetMedicalVisitByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, visit)
}

// CreateMedicalVisitHandler
// @Summary Create a new Medical Visit
// @Description Create a new medical visit
// @Tags Medical Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMedicalVisitRequest true "Medical Visit Details"
// @Success 201 {object} MedicalVisit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-visits [post]
func CreateMedicalVisitHandler(c *gin.Context) {
	var req CreateMedicalVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit, err := CreateMedicalVisitService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, visit)
}

// UpdateMedicalVisitHandler
// @Summary Update Medical Visit
// @Description Update medical visit details
// @Tags Medical Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Medical Visit ID"
// @Param request body UpdateMedicalVisitRequest true "Update Details"
// @Success 200 {object} MedicalVisit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-visits/{id} [put]
func UpdateMedicalVisitHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateMedicalVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit, err := UpdateMedicalVisitService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, visit)
}

// DeleteMedicalVisitHandler
// @Summary Delete Medical Visit
// @Description Soft delete a medical visit
// @Tags Medical Visits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Medical Visit ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-visits/{id} [delete]
func DeleteMedicalVisitHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteMedicalVisitService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
