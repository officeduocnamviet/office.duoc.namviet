package categories

// CreateCategoryRequest represents the payload for creating a category
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
	ParentID *int64 `json:"parent_id"`
	Status   string `json:"status"`
}

// UpdateCategoryRequest represents the payload for updating a category
type UpdateCategoryRequest struct {
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
	ParentID *int64  `json:"parent_id"`
	Status   *string `json:"status"`
}
