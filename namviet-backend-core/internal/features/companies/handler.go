package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Companies API ---

// GetAllCompaniesHandler
// @Summary Get All Companies
// @Description Retrieve a list of all companies
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Company
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies [get]
func GetAllCompaniesHandler(c *gin.Context) {
	items, err := GetAllCompaniesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetCompanyHandler
// @Summary Get Company by ID
// @Description Retrieve a specific company
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /companies/{id} [get]
func GetCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetCompanyByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateCompanyHandler
// @Summary Create a new Company
// @Description Create a new company
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCompanyRequest true "Company Details"
// @Success 201 {object} Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies [post]
func CreateCompanyHandler(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateCompanyService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateCompanyHandler
// @Summary Update Company
// @Description Update company details
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param request body UpdateCompanyRequest true "Update Details"
// @Success 200 {object} Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/{id} [put]
func UpdateCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateCompanyService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteCompanyHandler
// @Summary Delete Company
// @Description Soft delete a company
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/{id} [delete]
func DeleteCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteCompanyService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}


// --- Branches API ---

// GetAllBranchesHandler
// @Summary Get All Branches
// @Description Retrieve a list of all branches
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Branch
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /branches [get]
func GetAllBranchesHandler(c *gin.Context) {
	items, err := GetAllBranchesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetBranchHandler
// @Summary Get Branch by ID
// @Description Retrieve a specific branch
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Success 200 {object} Branch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /branches/{id} [get]
func GetBranchHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetBranchByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateBranchHandler
// @Summary Create a new Branch
// @Description Create a new branch
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBranchRequest true "Branch Details"
// @Success 201 {object} Branch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /branches [post]
func CreateBranchHandler(c *gin.Context) {
	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateBranchService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateBranchHandler
// @Summary Update Branch
// @Description Update branch details
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Param request body UpdateBranchRequest true "Update Details"
// @Success 200 {object} Branch
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /branches/{id} [put]
func UpdateBranchHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateBranchService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteBranchHandler
// @Summary Delete Branch
// @Description Soft delete a branch
// @Tags Branches
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /branches/{id} [delete]
func DeleteBranchHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteBranchService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
