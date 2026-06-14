package users

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/namviet/backend-core/internal/platform/supabase"
	"golang.org/x/crypto/bcrypt"
)

// LoginService verifies credentials and generates a JWT token
func LoginService(req LoginRequest) (*LoginResponse, error) {
	user, err := FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	mockHash, _ := bcrypt.GenerateFromPassword([]byte("namviet123"), bcrypt.DefaultCost)
	err = bcrypt.CompareHashAndPassword(mockHash, []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("namviet-secret-key-1234"))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

// GetAllUsersService returns all employees
func GetAllUsersService() ([]User, error) {
	return FindAll()
}

// CreateUserService creates a Supabase Auth user and then a DB record
func CreateUserService(req CreateUserRequest) (*User, error) {
	// 1. Create in Supabase Auth via Admin API
	authUser, err := supabase.CreateAuthUser(req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}

	// 2. Insert into public.users
	user := &User{
		ID:        authUser.ID,
		Email:     req.Email,
		FullName:  req.FullName,
		Phone:     req.Phone,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := Create(user); err != nil {
		// Rollback Auth user if DB insert fails
		_ = supabase.DeleteAuthUser(authUser.ID)
		return nil, err
	}

	return user, nil
}

// UpdateUserService updates a user profile
func UpdateUserService(id string, req UpdateUserRequest) (*User, error) {
	user, err := FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.RoleID != "" {
		user.RoleID = req.RoleID
	}

	if err := Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUserService deletes a user
func DeleteUserService(id string) error {
	// 1. Delete from DB (Soft Delete)
	if err := Delete(id); err != nil {
		return err
	}
	// 2. Delete from Supabase Auth so they can't login anymore
	_ = supabase.DeleteAuthUser(id)
	return nil
}

// RegisterFCMTokenService saves an FCM token for the user
func RegisterFCMTokenService(userID string, req RegisterFCMTokenRequest) error {
	fcmToken := &UserFCMToken{
		UserID:     userID,
		Token:      req.Token,
		DeviceInfo: req.DeviceInfo,
	}
	return SaveFCMToken(fcmToken)
}
