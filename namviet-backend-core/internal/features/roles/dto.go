package roles

// CreateRoleRequest represents the payload for creating a role
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest represents the payload for updating a role
type UpdateRoleRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}
