package categories

import "time"

func GetAllCategoriesService() ([]Category, error) {
	return GetAllCategories()
}

func GetCategoryByIDService(id int64) (*Category, error) {
	return GetCategoryByID(id)
}

func CreateCategoryService(req CreateCategoryRequest) (*Category, error) {
	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	category := &Category{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		Status:   status,
	}

	if err := CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func UpdateCategoryService(id int64, req UpdateCategoryRequest) (*Category, error) {
	category, err := GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Status != nil {
		category.Status = *req.Status
	}
	
	now := time.Now()
	category.UpdatedAt = &now

	if err := UpdateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}
