package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetInventoryHandler
// @Summary Get Inventory
// @Description Retrieve inventory data with optional filters
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param warehouse_id query int false "Warehouse ID"
// @Param product_id query int false "Product ID"
// @Success 200 {array} InventoryBatch
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory [get]
func GetInventoryHandler(c *gin.Context) {
	whIDStr := c.Query("warehouse_id")
	prodIDStr := c.Query("product_id")

	var warehouseID, productID int64
	if whIDStr != "" {
		warehouseID, _ = strconv.ParseInt(whIDStr, 10, 64)
	}
	if prodIDStr != "" {
		productID, _ = strconv.ParseInt(prodIDStr, 10, 64)
	}

	data, err := GetInventoryService(warehouseID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// CreateTransactionHandler
// @Summary Create Inventory Transaction
// @Description Create an IN or OUT inventory transaction
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTransactionRequest true "Transaction Details"
// @Success 201 {object} InventoryTransaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory/transactions [post]
func CreateTransactionHandler(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := CreateTransactionService(req)
	if err != nil {
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tx)
}
