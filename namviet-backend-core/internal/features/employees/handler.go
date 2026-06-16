package employees

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllEmployeesHandler
// @Summary Get All Employees
// @Description Retrieve a list of all employees
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Employee
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [get]
func GetAllEmployeesHandler(c *gin.Context) {
	emps, err := GetAllEmployeesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emps)
}

// GetEmployeeHandler
// @Summary Get Employee by ID
// @Description Retrieve a specific employee
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Employee ID"
// @Success 200 {object} Employee
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /employees/{id} [get]
func GetEmployeeHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	emp, err := GetEmployeeByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// CreateEmployeeHandler
// @Summary Create a new Employee
// @Description Create a new employee profile
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEmployeeRequest true "Employee Details"
// @Success 201 {object} Employee
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [post]
func CreateEmployeeHandler(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp, err := CreateEmployeeService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, emp)
}

// UpdateEmployeeHandler
// @Summary Update Employee
// @Description Update employee details
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Employee ID"
// @Param request body UpdateEmployeeRequest true "Update Details"
// @Success 200 {object} Employee
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [put]
func UpdateEmployeeHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp, err := UpdateEmployeeService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// DeleteEmployeeHandler
// @Summary Delete Employee
// @Description Soft delete an employee
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Employee ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [delete]
func DeleteEmployeeHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteEmployeeService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
