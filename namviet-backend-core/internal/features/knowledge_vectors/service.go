package knowledge_vectors

// Medical Knowledge Vectors
func GetAllMedicalKnowledgeVectorsService() ([]MedicalKnowledgeVector, error) {
	return GetAllMedicalKnowledgeVectors()
}

func GetMedicalKnowledgeVectorByIDService(id string) (*MedicalKnowledgeVector, error) {
	return GetMedicalKnowledgeVectorByID(id)
}

func CreateMedicalKnowledgeVectorService(req CreateMedicalKnowledgeVectorRequest) (*MedicalKnowledgeVector, error) {
	vector := &MedicalKnowledgeVector{
		Title:     req.Title,
		Content:   req.Content,
		Embedding: req.Embedding,
		Metadata:  req.Metadata,
	}

	if err := CreateMedicalKnowledgeVector(vector); err != nil {
		return nil, err
	}
	return vector, nil
}

func UpdateMedicalKnowledgeVectorService(id string, req UpdateMedicalKnowledgeVectorRequest) (*MedicalKnowledgeVector, error) {
	vector, err := GetMedicalKnowledgeVectorByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		vector.Title = *req.Title
	}
	if req.Content != nil {
		vector.Content = *req.Content
	}
	if req.Embedding != nil {
		vector.Embedding = *req.Embedding
	}
	if req.Metadata != nil {
		vector.Metadata = req.Metadata
	}

	if err := UpdateMedicalKnowledgeVector(vector); err != nil {
		return nil, err
	}
	return vector, nil
}

func DeleteMedicalKnowledgeVectorService(id string) error {
	return DeleteMedicalKnowledgeVector(id)
}

// Product Vectors
func GetAllProductVectorsService() ([]ProductVector, error) {
	return GetAllProductVectors()
}

func GetProductVectorByIDService(id string) (*ProductVector, error) {
	return GetProductVectorByID(id)
}

func CreateProductVectorService(req CreateProductVectorRequest) (*ProductVector, error) {
	vector := &ProductVector{
		ProductID: req.ProductID,
		Content:   req.Content,
		Embedding: req.Embedding,
		Metadata:  req.Metadata,
	}

	if err := CreateProductVector(vector); err != nil {
		return nil, err
	}
	return vector, nil
}

func UpdateProductVectorService(id string, req UpdateProductVectorRequest) (*ProductVector, error) {
	vector, err := GetProductVectorByID(id)
	if err != nil {
		return nil, err
	}

	if req.ProductID != nil {
		vector.ProductID = req.ProductID
	}
	if req.Content != nil {
		vector.Content = *req.Content
	}
	if req.Embedding != nil {
		vector.Embedding = *req.Embedding
	}
	if req.Metadata != nil {
		vector.Metadata = req.Metadata
	}

	if err := UpdateProductVector(vector); err != nil {
		return nil, err
	}
	return vector, nil
}

func DeleteProductVectorService(id string) error {
	return DeleteProductVector(id)
}
