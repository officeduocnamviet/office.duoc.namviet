package fund_accounts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllFundAccountsHandler
// @Summary Get All Fund Accounts
// @Description Retrieve a list of all fund accounts
// @Tags Fund Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} FundAccount
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fund-accounts [get]
func GetAllFundAccountsHandler(c *gin.Context) {
	fas, err := GetAllFundAccountsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fas)
}

// GetFundAccountHandler
// @Summary Get Fund Account by ID
// @Description Retrieve a specific fund account
// @Tags Fund Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Fund Account ID"
// @Success 200 {object} FundAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /fund-accounts/{id} [get]
func GetFundAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	fa, err := GetFundAccountByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fa)
}

// CreateFundAccountHandler
// @Summary Create a new Fund Account
// @Description Create a new fund account
// @Tags Fund Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFundAccountRequest true "Fund Account Details"
// @Success 201 {object} FundAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fund-accounts [post]
func CreateFundAccountHandler(c *gin.Context) {
	var req CreateFundAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fa, err := CreateFundAccountService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fa)
}

// UpdateFundAccountHandler
// @Summary Update Fund Account
// @Description Update fund account details
// @Tags Fund Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Fund Account ID"
// @Param request body UpdateFundAccountRequest true "Update Details"
// @Success 200 {object} FundAccount
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fund-accounts/{id} [put]
func UpdateFundAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateFundAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fa, err := UpdateFundAccountService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fa)
}

// DeleteFundAccountHandler
// @Summary Delete Fund Account
// @Description Soft delete a fund account
// @Tags Fund Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Fund Account ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fund-accounts/{id} [delete]
func DeleteFundAccountHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteFundAccountService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
