package users

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@namviet.com"`
	Password string `json:"password" binding:"required" example:"namviet123"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FullName  string `json:"full_name" binding:"required"`
	Phone     string `json:"phone"`
	RoleID      string `json:"role_id" binding:"required"`
	CompanyID   string `json:"company_id" binding:"required"`
	WarehouseID *int64 `json:"warehouse_id"`
}

type UpdateUserRequest struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Status      string `json:"status"`
	RoleID      string `json:"role_id"`
	WarehouseID *int64 `json:"warehouse_id"`
}

type RegisterFCMTokenRequest struct {
	Token      string `json:"token" binding:"required"`
	DeviceInfo string `json:"device_info"`
}
