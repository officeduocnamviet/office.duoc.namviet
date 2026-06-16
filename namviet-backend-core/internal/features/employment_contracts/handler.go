package employment_contracts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllEmploymentContractsHandler
// @Summary Get All Employment Contracts
// @Description Retrieve a list of all employment contracts
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} EmploymentContract
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employment-contracts [get]
func GetAllEmploymentContractsHandler(c *gin.Context) {
	items, err := GetAllEmploymentContractsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetEmploymentContractHandler
// @Summary Get Employment Contract by ID
// @Description Retrieve a specific employment contract
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Success 200 {object} EmploymentContract
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /employment-contracts/{id} [get]
func GetEmploymentContractHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetEmploymentContractByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateEmploymentContractHandler
// @Summary Create a new Employment Contract
// @Description Create a new employment contract
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEmploymentContractRequest true "Contract Details"
// @Success 201 {object} EmploymentContract
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employment-contracts [post]
func CreateEmploymentContractHandler(c *gin.Context) {
	var req CreateEmploymentContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateEmploymentContractService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateEmploymentContractHandler
// @Summary Update Employment Contract
// @Description Update employment contract details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Param request body UpdateEmploymentContractRequest true "Update Details"
// @Success 200 {object} EmploymentContract
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employment-contracts/{id} [put]
func UpdateEmploymentContractHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateEmploymentContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateEmploymentContractService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteEmploymentContractHandler
// @Summary Delete Employment Contract
// @Description Delete an employment contract
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employment-contracts/{id} [delete]
func DeleteEmploymentContractHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteEmploymentContractService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
