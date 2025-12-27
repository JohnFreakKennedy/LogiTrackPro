package database

import (
	"testing"
	"time"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.Customer{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TestCreateCustomer tests customer creation
func TestCreateCustomer(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name     string
		customer *models.Customer
		wantErr  bool
	}{
		{
			name: "valid customer",
			customer: &models.Customer{
				Name:             "Test Customer",
				Address:          "123 Test St",
				Latitude:         40.7128,
				Longitude:        -74.0060,
				DemandRate:       100.0,
				MaxInventory:     1000.0,
				CurrentInventory: 500.0,
				MinInventory:     100.0,
				Priority:         1,
			},
			wantErr: false,
		},
		{
			name: "customer with zero values",
			customer: &models.Customer{
				Name:      "Minimal Customer",
				Latitude:  0.0,
				Longitude: 0.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateCustomer(db, tt.customer)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.customer.ID == 0 {
				t.Error("CreateCustomer() did not set ID")
			}
			if !tt.wantErr && tt.customer.CreatedAt.IsZero() {
				t.Error("CreateCustomer() did not set CreatedAt")
			}
		})
	}
}

// TestGetCustomer tests customer retrieval
func TestGetCustomer(t *testing.T) {
	db := setupTestDB(t)

	// Create test customer
	customer := &models.Customer{
		Name:      "Test Customer",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	CreateCustomer(db, customer)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
		errType error
	}{
		{
			name:    "existing customer",
			id:      customer.ID,
			wantErr: false,
		},
		{
			name:    "non-existent customer",
			id:      99999,
			wantErr: true,
			errType: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCustomer(db, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetCustomer() returned nil")
					return
				}
				if got.ID != tt.id {
					t.Errorf("GetCustomer() ID = %v, want %v", got.ID, tt.id)
				}
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("GetCustomer() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

// TestListCustomers tests customer listing
func TestListCustomers(t *testing.T) {
	db := setupTestDB(t)

	// Create multiple customers
	customers := []*models.Customer{
		{Name: "Customer A", Latitude: 40.7128, Longitude: -74.0060},
		{Name: "Customer B", Latitude: 34.0522, Longitude: -118.2437},
		{Name: "Customer C", Latitude: 41.8781, Longitude: -87.6298},
	}

	for _, c := range customers {
		CreateCustomer(db, c)
	}

	got, err := ListCustomers(db)
	if err != nil {
		t.Fatalf("ListCustomers() error = %v", err)
	}

	if len(got) != len(customers) {
		t.Errorf("ListCustomers() returned %d customers, want %d", len(got), len(customers))
	}

	// Verify ordering (should be ordered by name)
	if len(got) > 1 && got[0].Name > got[1].Name {
		t.Error("ListCustomers() not ordered by name")
	}
}

// TestUpdateCustomer tests customer updates
func TestUpdateCustomer(t *testing.T) {
	db := setupTestDB(t)

	// Create test customer
	customer := &models.Customer{
		Name:      "Original Name",
		Latitude:  40.7128,
		Longitude: -74.0060,
		Priority:  1,
	}
	CreateCustomer(db, customer)

	originalUpdatedAt := customer.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt changes

	// Update customer
	customer.Name = "Updated Name"
	customer.Priority = 2
	err := UpdateCustomer(db, customer)
	if err != nil {
		t.Fatalf("UpdateCustomer() error = %v", err)
	}

	// Verify update
	updated, err := GetCustomer(db, customer.ID)
	if err != nil {
		t.Fatalf("GetCustomer() error = %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("UpdateCustomer() Name = %v, want Updated Name", updated.Name)
	}
	if updated.Priority != 2 {
		t.Errorf("UpdateCustomer() Priority = %v, want 2", updated.Priority)
	}
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdateCustomer() did not update UpdatedAt")
	}
}

// TestUpdateCustomerNotFound tests updating non-existent customer
func TestUpdateCustomerNotFound(t *testing.T) {
	db := setupTestDB(t)

	customer := &models.Customer{
		ID:        99999,
		Name:      "Non-existent",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	err := UpdateCustomer(db, customer)
	if err != ErrNotFound {
		t.Errorf("UpdateCustomer() error = %v, want ErrNotFound", err)
	}
}

// TestDeleteCustomer tests customer deletion
func TestDeleteCustomer(t *testing.T) {
	db := setupTestDB(t)

	// Create test customer
	customer := &models.Customer{
		Name:      "To Delete",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	CreateCustomer(db, customer)

	// Delete customer
	err := DeleteCustomer(db, customer.ID)
	if err != nil {
		t.Fatalf("DeleteCustomer() error = %v", err)
	}

	// Verify deletion
	_, err = GetCustomer(db, customer.ID)
	if err != ErrNotFound {
		t.Errorf("DeleteCustomer() customer still exists, error = %v", err)
	}
}

// TestDeleteCustomerNotFound tests deleting non-existent customer
func TestDeleteCustomerNotFound(t *testing.T) {
	db := setupTestDB(t)

	err := DeleteCustomer(db, 99999)
	if err != ErrNotFound {
		t.Errorf("DeleteCustomer() error = %v, want ErrNotFound", err)
	}
}

// TestCountCustomers tests customer counting
func TestCountCustomers(t *testing.T) {
	db := setupTestDB(t)

	// Initially should be 0
	count, err := CountCustomers(db)
	if err != nil {
		t.Fatalf("CountCustomers() error = %v", err)
	}
	if count != 0 {
		t.Errorf("CountCustomers() = %d, want 0", count)
	}

	// Create customers
	for i := 0; i < 5; i++ {
		customer := &models.Customer{
			Name:      "Customer",
			Latitude:  float64(i),
			Longitude: float64(i),
		}
		CreateCustomer(db, customer)
	}

	count, err = CountCustomers(db)
	if err != nil {
		t.Fatalf("CountCustomers() error = %v", err)
	}
	if count != 5 {
		t.Errorf("CountCustomers() = %d, want 5", count)
	}
}
