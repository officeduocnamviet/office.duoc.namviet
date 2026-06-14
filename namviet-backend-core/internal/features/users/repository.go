package users

import (
	"errors"
	"time"

	"github.com/namviet/backend-core/internal/platform/supabase"
)

// FindByEmail mock fetching user from DB
func FindByEmail(email string) (*User, error) {
	var user User
	result := supabase.DB.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// FindAll fetches all active users
func FindAll() ([]User, error) {
	var usersList []User
	result := supabase.DB.Where("deleted_at IS NULL").Find(&usersList)
	if result.Error != nil {
		return nil, result.Error
	}
	return usersList, nil
}

// Create inserts a new user
func Create(user *User) error {
	return supabase.DB.Create(user).Error
}

// Update updates user fields
func Update(user *User) error {
	user.UpdatedAt = time.Now()
	return supabase.DB.Save(user).Error
}

// FindByID gets user by ID
func FindByID(id string) (*User, error) {
	var user User
	result := supabase.DB.Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// Delete logically deletes a user
func Delete(id string) error {
	return supabase.DB.Model(&User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// SaveFCMToken upserts an FCM token
func SaveFCMToken(fcmToken *UserFCMToken) error {
	// Upsert based on the token string
	return supabase.DB.Where("token = ?", fcmToken.Token).
		Assign(UserFCMToken{
			UserID:     fcmToken.UserID,
			DeviceInfo: fcmToken.DeviceInfo,
			UpdatedAt:  time.Now(),
		}).
		FirstOrCreate(fcmToken).Error
}
