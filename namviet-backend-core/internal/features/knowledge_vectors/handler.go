package knowledge_vectors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Medical Knowledge Vectors API ---

// GetAllMedicalKnowledgeVectorsHandler
// @Summary Get All Medical Knowledge Vectors
// @Description Retrieve a list of all medical knowledge vectors
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} MedicalKnowledgeVector
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-knowledge-vectors [get]
func GetAllMedicalKnowledgeVectorsHandler(c *gin.Context) {
	items, err := GetAllMedicalKnowledgeVectorsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetMedicalKnowledgeVectorHandler
// @Summary Get Medical Knowledge Vector by ID
// @Description Retrieve a specific medical knowledge vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Success 200 {object} MedicalKnowledgeVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /medical-knowledge-vectors/{id} [get]
func GetMedicalKnowledgeVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetMedicalKnowledgeVectorByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateMedicalKnowledgeVectorHandler
// @Summary Create a new Medical Knowledge Vector
// @Description Create a new medical knowledge vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMedicalKnowledgeVectorRequest true "Vector Details"
// @Success 201 {object} MedicalKnowledgeVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-knowledge-vectors [post]
func CreateMedicalKnowledgeVectorHandler(c *gin.Context) {
	var req CreateMedicalKnowledgeVectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateMedicalKnowledgeVectorService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateMedicalKnowledgeVectorHandler
// @Summary Update Medical Knowledge Vector
// @Description Update medical knowledge vector details
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Param request body UpdateMedicalKnowledgeVectorRequest true "Update Details"
// @Success 200 {object} MedicalKnowledgeVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-knowledge-vectors/{id} [put]
func UpdateMedicalKnowledgeVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateMedicalKnowledgeVectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateMedicalKnowledgeVectorService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteMedicalKnowledgeVectorHandler
// @Summary Delete Medical Knowledge Vector
// @Description Delete a medical knowledge vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-knowledge-vectors/{id} [delete]
func DeleteMedicalKnowledgeVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteMedicalKnowledgeVectorService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Product Vectors API ---

// GetAllProductVectorsHandler
// @Summary Get All Product Vectors
// @Description Retrieve a list of all product vectors
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ProductVector
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /product-vectors [get]
func GetAllProductVectorsHandler(c *gin.Context) {
	items, err := GetAllProductVectorsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetProductVectorHandler
// @Summary Get Product Vector by ID
// @Description Retrieve a specific product vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Success 200 {object} ProductVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /product-vectors/{id} [get]
func GetProductVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetProductVectorByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateProductVectorHandler
// @Summary Create a new Product Vector
// @Description Create a new product vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductVectorRequest true "Vector Details"
// @Success 201 {object} ProductVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /product-vectors [post]
func CreateProductVectorHandler(c *gin.Context) {
	var req CreateProductVectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateProductVectorService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateProductVectorHandler
// @Summary Update Product Vector
// @Description Update product vector details
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Param request body UpdateProductVectorRequest true "Update Details"
// @Success 200 {object} ProductVector
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /product-vectors/{id} [put]
func UpdateProductVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateProductVectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateProductVectorService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteProductVectorHandler
// @Summary Delete Product Vector
// @Description Delete a product vector
// @Tags AI Ecosystem
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vector ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /product-vectors/{id} [delete]
func DeleteProductVectorHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteProductVectorService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
