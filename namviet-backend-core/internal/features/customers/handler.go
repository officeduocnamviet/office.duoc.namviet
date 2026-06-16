package customers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllCustomersHandler
// @Summary Get All Customers
// @Description Retrieve a list of all customers
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Customer
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers [get]
func GetAllCustomersHandler(c *gin.Context) {
	customers, err := GetAllCustomersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// GetCustomerHandler
// @Summary Get Customer by ID
// @Description Retrieve a specific customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /customers/{id} [get]
func GetCustomerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	customer, err := GetCustomerByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// CreateCustomerHandler
// @Summary Create a new Customer
// @Description Create a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCustomerRequest true "Customer Details"
// @Success 201 {object} Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers [post]
func CreateCustomerHandler(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := CreateCustomerService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)
}

// UpdateCustomerHandler
// @Summary Update Customer
// @Description Update customer details
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param request body UpdateCustomerRequest true "Update Details"
// @Success 200 {object} Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [put]
func UpdateCustomerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := UpdateCustomerService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// DeleteCustomerHandler
// @Summary Delete Customer
// @Description Soft delete a customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [delete]
func DeleteCustomerHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteCustomerService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
