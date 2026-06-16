package finance_transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllFinanceTransactionsHandler
// @Summary Get All Finance Transactions
// @Description Retrieve a list of all finance transactions
// @Tags Finance Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} FinanceTransaction
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /finance-transactions [get]
func GetAllFinanceTransactionsHandler(c *gin.Context) {
	fts, err := GetAllFinanceTransactionsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fts)
}

// GetFinanceTransactionHandler
// @Summary Get Finance Transaction by ID
// @Description Retrieve a specific finance transaction
// @Tags Finance Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Finance Transaction ID"
// @Success 200 {object} FinanceTransaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /finance-transactions/{id} [get]
func GetFinanceTransactionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	ft, err := GetFinanceTransactionByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ft)
}

// CreateFinanceTransactionHandler
// @Summary Create a new Finance Transaction
// @Description Create a new finance transaction
// @Tags Finance Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFinanceTransactionRequest true "Finance Transaction Details"
// @Success 201 {object} FinanceTransaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /finance-transactions [post]
func CreateFinanceTransactionHandler(c *gin.Context) {
	var req CreateFinanceTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ft, err := CreateFinanceTransactionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ft)
}

// UpdateFinanceTransactionHandler
// @Summary Update Finance Transaction
// @Description Update finance transaction details
// @Tags Finance Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Finance Transaction ID"
// @Param request body UpdateFinanceTransactionRequest true "Update Details"
// @Success 200 {object} FinanceTransaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /finance-transactions/{id} [put]
func UpdateFinanceTransactionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateFinanceTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ft, err := UpdateFinanceTransactionService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ft)
}

// DeleteFinanceTransactionHandler
// @Summary Delete Finance Transaction
// @Description Soft delete a finance transaction
// @Tags Finance Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Finance Transaction ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /finance-transactions/{id} [delete]
func DeleteFinanceTransactionHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteFinanceTransactionService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
