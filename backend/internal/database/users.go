package database

import (
	"database/sql"
	"errors"

	"LogiTrackPro/backend/internal/models"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("record already exists")

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at 
			  FROM users WHERE email = $1`
	
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(db *sql.DB, id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at 
			  FROM users WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (email, password_hash, name, role) 
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	
	err := db.QueryRow(query, user.Email, user.Password, user.Name, user.Role).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "unique") || contains(err.Error(), "duplicate"))
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

