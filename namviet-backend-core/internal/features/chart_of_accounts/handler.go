package chart_of_accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllChartOfAccountsHandler
// @Summary Get All Chart Of Accounts
// @Description Retrieve a list of all chart of accounts
// @Tags Chart Of Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ChartOfAccount
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chart-of-accounts [get]
func GetAllChartOfAccountsHandler(c *gin.Context) {
	coas, err := GetAllChartOfAccountsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coas)
}

// GetChartOfAccountHandler
// @Summary Get Chart Of Account by ID
// @Description Retrieve a specific chart of account
// @Tags Chart Of Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chart Of Account ID"
// @Success 200 {object} ChartOfAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /chart-of-accounts/{id} [get]
func GetChartOfAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	coa, err := GetChartOfAccountByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coa)
}

// CreateChartOfAccountHandler
// @Summary Create a new Chart Of Account
// @Description Create a new chart of account
// @Tags Chart Of Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateChartOfAccountRequest true "Chart Of Account Details"
// @Success 201 {object} ChartOfAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chart-of-accounts [post]
func CreateChartOfAccountHandler(c *gin.Context) {
	var req CreateChartOfAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coa, err := CreateChartOfAccountService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, coa)
}

// UpdateChartOfAccountHandler
// @Summary Update Chart Of Account
// @Description Update chart of account details
// @Tags Chart Of Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chart Of Account ID"
// @Param request body UpdateChartOfAccountRequest true "Update Details"
// @Success 200 {object} ChartOfAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chart-of-accounts/{id} [put]
func UpdateChartOfAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateChartOfAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coa, err := UpdateChartOfAccountService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coa)
}

// DeleteChartOfAccountHandler
// @Summary Delete Chart Of Account
// @Description Soft delete a chart of account
// @Tags Chart Of Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chart Of Account ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chart-of-accounts/{id} [delete]
func DeleteChartOfAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteChartOfAccountService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
