package roles

import (
	"errors"
	"strings"
)

// GetAllRolesService handles fetching all roles
func GetAllRolesService() ([]Role, error) {
	return GetAllRoles()
}

// GetRoleByIDService handles fetching a single role
func GetRoleByIDService(id string) (*Role, error) {
	return GetRoleByID(id)
}

// CreateRoleService handles creating a new role
func CreateRoleService(req CreateRoleRequest) (*Role, error) {
	role := &Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	}
	if role.Permissions == nil {
		role.Permissions = JSONB{}
	}

	err := CreateRole(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRoleService handles updating an existing role
func UpdateRoleService(id string, req UpdateRoleRequest) (*Role, error) {
	role, err := GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	// Basic protection against editing the core Admin role
	if strings.ToLower(role.Name) == "admin" && req.Name != nil && *req.Name != "admin" {
		return nil, errors.New("cannot rename the admin role")
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = req.Description
	}
	if req.Permissions != nil {
		role.Permissions = req.Permissions
	}

	err = UpdateRole(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRoleService handles deleting a role
func DeleteRoleService(id string) error {
	role, err := GetRoleByID(id)
	if err != nil {
		return err
	}

	// Protect Admin role
	if strings.ToLower(role.Name) == "admin" {
		return errors.New("cannot delete the admin role")
	}

	return DeleteRole(id)
}
