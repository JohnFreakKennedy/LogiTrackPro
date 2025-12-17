package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("record already exists")

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	user := &models.User{}
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func GetUserByID(db *gorm.DB, id int64) (*models.User, error) {
	user := &models.User{}
	err := db.First(user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func CreateUser(db *gorm.DB, user *models.User) error {
	err := db.Create(user).Error
	if err != nil {
		// Check for unique constraint violation
		if isUniqueViolation(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func isUniqueViolation(err error) bool {
	// GORM wraps PostgreSQL errors, check for unique constraint violations
	return err != nil && (
		errors.Is(err, gorm.ErrDuplicatedKey) ||
		contains(err.Error(), "unique") ||
		contains(err.Error(), "duplicate") ||
		contains(err.Error(), "violates unique constraint"),
	)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

