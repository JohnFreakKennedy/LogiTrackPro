package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"LogiTrackPro/backend/internal/config"
	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"
	"LogiTrackPro/backend/internal/optimizer"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupIntegrationTestDB creates a test database with all models
func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Warehouse{},
		&models.Customer{},
		&models.Vehicle{},
		&models.Plan{},
		&models.Route{},
		&models.Stop{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// setupIntegrationHandler creates a handler with test database
func setupIntegrationHandler(t *testing.T) (*Handler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := setupIntegrationTestDB(t)

	cfg := &config.Config{
		JWTSecret:    "test-secret-key-for-testing-only",
		JWTExpiry:    24,
		OptimizerURL: "http://localhost:8000",
	}

	optimizerClient := optimizer.NewClient(cfg.OptimizerURL)

	return New(db, optimizerClient, cfg), db
}

// getAuthToken helper function to get authentication token
func getAuthToken(t *testing.T, h *Handler) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
		Role:     "user",
	}
	database.CreateUser(h.db, user)

	loginBody, _ := json.Marshal(LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/auth/login", h.Login)
	router.ServeHTTP(w, req)

	var response struct {
		Success bool
		Data    AuthResponse
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Data.Token
}

// TestCustomerCRUDIntegration tests complete CRUD flow for customers
func TestCustomerCRUDIntegration(t *testing.T) {
	h, db := setupIntegrationHandler(t)
	token := getAuthToken(t, h)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.GET("/api/v1/customers", h.ListCustomers)
	router.POST("/api/v1/customers", h.CreateCustomer)
	router.GET("/api/v1/customers/:id", h.GetCustomer)
	router.PUT("/api/v1/customers/:id", h.UpdateCustomer)
	router.DELETE("/api/v1/customers/:id", h.DeleteCustomer)

	// Create customer
	createBody, _ := json.Marshal(CustomerRequest{
		Name:             "Integration Test Customer",
		Address:          "123 Test St",
		Latitude:         40.7128,
		Longitude:        -74.0060,
		DemandRate:       100.0,
		MaxInventory:     1000.0,
		CurrentInventory: 500.0,
		MinInventory:     100.0,
		Priority:         1,
	})

	createReq := httptest.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("CreateCustomer() status = %d, want %d", createW.Code, http.StatusCreated)
	}

	var createResponse struct {
		Success bool
		Data    models.Customer
	}
	json.Unmarshal(createW.Body.Bytes(), &createResponse)
	customerID := createResponse.Data.ID

	// Get customer
	getReq := httptest.NewRequest("GET", "/api/v1/customers/"+strconv.FormatInt(customerID, 10), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("GetCustomer() status = %d, want %d", getW.Code, http.StatusOK)
	}

	// Update customer
	updateBody, _ := json.Marshal(CustomerRequest{
		Name:             "Updated Customer",
		Latitude:         40.7128,
		Longitude:        -74.0060,
		CurrentInventory: 600.0,
		Priority:         2,
	})

	updateReq := httptest.NewRequest("PUT", "/api/v1/customers/"+strconv.FormatInt(customerID, 10), bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Fatalf("UpdateCustomer() status = %d, want %d", updateW.Code, http.StatusOK)
	}

	// List customers
	listReq := httptest.NewRequest("GET", "/api/v1/customers", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)

	if listW.Code != http.StatusOK {
		t.Fatalf("ListCustomers() status = %d, want %d", listW.Code, http.StatusOK)
	}

	var listResponse struct {
		Success bool
		Data    []models.Customer
	}
	json.Unmarshal(listW.Body.Bytes(), &listResponse)
	if len(listResponse.Data) == 0 {
		t.Error("ListCustomers() returned empty list")
	}

	// Delete customer
	deleteReq := httptest.NewRequest("DELETE", "/api/v1/customers/"+strconv.FormatInt(customerID, 10), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusOK {
		t.Fatalf("DeleteCustomer() status = %d, want %d", deleteW.Code, http.StatusOK)
	}

	// Verify deletion
	getReq2 := httptest.NewRequest("GET", "/api/v1/customers/"+strconv.FormatInt(customerID, 10), nil)
	getReq2.Header.Set("Authorization", "Bearer "+token)
	getW2 := httptest.NewRecorder()
	router.ServeHTTP(getW2, getReq2)

	if getW2.Code != http.StatusNotFound {
		t.Fatalf("GetCustomer() after delete status = %d, want %d", getW2.Code, http.StatusNotFound)
	}
}

// TestPlanCreationFlow tests plan creation with warehouse
func TestPlanCreationFlow(t *testing.T) {
	h, db := setupIntegrationHandler(t)
	token := getAuthToken(t, h)

	// Create warehouse first
	warehouse := &models.Warehouse{
		Name:      "Test Warehouse",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	database.CreateWarehouse(h.db, warehouse)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.POST("/api/v1/plans", h.CreatePlan)

	planBody, _ := json.Marshal(PlanRequest{
		Name:        "Test Plan",
		StartDate:   "2024-01-01",
		EndDate:     "2024-01-07",
		WarehouseID: warehouse.ID,
	})

	req := httptest.NewRequest("POST", "/api/v1/plans", bytes.NewBuffer(planBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreatePlan() status = %d, want %d", w.Code, http.StatusCreated)
	}

	var response struct {
		Success bool
		Data    models.Plan
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Data.Status != "draft" {
		t.Errorf("CreatePlan() Status = %v, want draft", response.Data.Status)
	}
	if response.Data.WarehouseID == nil || *response.Data.WarehouseID != warehouse.ID {
		t.Error("CreatePlan() WarehouseID not set correctly")
	}
}

// TestProtectedRouteAccess tests that protected routes require authentication
func TestProtectedRouteAccess(t *testing.T) {
	h, _ := setupIntegrationHandler(t)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.GET("/api/v1/customers", h.ListCustomers)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "no token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/customers", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Protected route status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

// TestErrorHandling tests error responses
func TestErrorHandling(t *testing.T) {
	h, _ := setupIntegrationHandler(t)
	token := getAuthToken(t, h)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.GET("/api/v1/customers/:id", h.GetCustomer)

	// Test invalid ID format
	req := httptest.NewRequest("GET", "/api/v1/customers/invalid", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GetCustomer() with invalid ID status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	// Test non-existent resource
	req2 := httptest.NewRequest("GET", "/api/v1/customers/99999", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("GetCustomer() with non-existent ID status = %d, want %d", w2.Code, http.StatusNotFound)
	}
}
