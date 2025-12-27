package database

import (
	"testing"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TestCreateUser tests user creation
func TestCreateUser(t *testing.T) {
	db := setupUserTestDB(t)

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
		errType error
	}{
		{
			name: "valid user",
			user: &models.User{
				Email:    "test@example.com",
				Password: "hashed_password",
				Name:     "Test User",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &models.User{
				Email:    "duplicate@example.com",
				Password: "hashed_password",
				Name:     "First User",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name: "admin user",
			user: &models.User{
				Email:    "admin@example.com",
				Password: "hashed_password",
				Name:     "Admin User",
				Role:     "admin",
			},
			wantErr: false,
		},
	}

	// First user should succeed
	err := CreateUser(db, tests[0].user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	// Test duplicate email
	err = CreateUser(db, &models.User{
		Email:    tests[0].user.Email,
		Password: "hashed_password",
		Name:     "Duplicate User",
		Role:     "user",
	})
	if err != ErrDuplicate {
		t.Errorf("CreateUser() duplicate email error = %v, want ErrDuplicate", err)
	}

	// Test other valid users
	for i := 1; i < len(tests); i++ {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(db, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.user.ID == 0 {
				t.Error("CreateUser() did not set ID")
			}
		})
	}
}

// TestGetUserByEmail tests user retrieval by email
func TestGetUserByEmail(t *testing.T) {
	db := setupUserTestDB(t)

	// Create test user
	user := &models.User{
		Email:    "getbyemail@example.com",
		Password: "hashed_password",
		Name:     "Get User",
		Role:     "user",
	}
	CreateUser(db, user)

	tests := []struct {
		name    string
		email   string
		wantErr bool
		errType error
	}{
		{
			name:    "existing user",
			email:   "getbyemail@example.com",
			wantErr: false,
		},
		{
			name:    "non-existent user",
			email:   "nonexistent@example.com",
			wantErr: true,
			errType: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByEmail(db, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetUserByEmail() returned nil")
					return
				}
				if got.Email != tt.email {
					t.Errorf("GetUserByEmail() Email = %v, want %v", got.Email, tt.email)
				}
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("GetUserByEmail() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

// TestGetUserByID tests user retrieval by ID
func TestGetUserByID(t *testing.T) {
	db := setupUserTestDB(t)

	// Create test user
	user := &models.User{
		Email:    "getbyid@example.com",
		Password: "hashed_password",
		Name:     "Get ID User",
		Role:     "user",
	}
	CreateUser(db, user)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
		errType error
	}{
		{
			name:    "existing user",
			id:      user.ID,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      99999,
			wantErr: true,
			errType: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByID(db, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetUserByID() returned nil")
					return
				}
				if got.ID != tt.id {
					t.Errorf("GetUserByID() ID = %v, want %v", got.ID, tt.id)
				}
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("GetUserByID() error = %v, want %v", err, tt.errType)
			}
		})
	}
}
