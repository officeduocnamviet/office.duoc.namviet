package payrolls

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllPayrollsHandler
// @Summary Get All Payrolls
// @Description Retrieve a list of all payrolls
// @Tags Payrolls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Payroll
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payrolls [get]
func GetAllPayrollsHandler(c *gin.Context) {
	payrolls, err := GetAllPayrollsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payrolls)
}

// GetPayrollHandler
// @Summary Get Payroll by ID
// @Description Retrieve a specific payroll
// @Tags Payrolls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payroll ID"
// @Success 200 {object} Payroll
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /payrolls/{id} [get]
func GetPayrollHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	payroll, err := GetPayrollByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payroll)
}

// CreatePayrollHandler
// @Summary Create a new Payroll
// @Description Create a new payroll record
// @Tags Payrolls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePayrollRequest true "Payroll Details"
// @Success 201 {object} Payroll
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payrolls [post]
func CreatePayrollHandler(c *gin.Context) {
	var req CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payroll, err := CreatePayrollService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payroll)
}

// UpdatePayrollHandler
// @Summary Update Payroll
// @Description Update payroll details
// @Tags Payrolls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payroll ID"
// @Param request body UpdatePayrollRequest true "Update Details"
// @Success 200 {object} Payroll
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payrolls/{id} [put]
func UpdatePayrollHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payroll, err := UpdatePayrollService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payroll)
}

// DeletePayrollHandler
// @Summary Delete Payroll
// @Description Delete a payroll record
// @Tags Payrolls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payroll ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payrolls/{id} [delete]
func DeletePayrollHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeletePayrollService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
