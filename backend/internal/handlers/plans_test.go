package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"LogiTrackPro/backend/internal/config"
	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"
	"LogiTrackPro/backend/internal/optimizer"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPlanTestHandler(t *testing.T) (*Handler, *gorm.DB) {
	gin.SetMode(gin.TestMode)
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

	cfg := &config.Config{
		JWTSecret:    "test-secret-key",
		JWTExpiry:    24,
		OptimizerURL: "http://localhost:8000",
	}

	optimizerClient := optimizer.NewClient(cfg.OptimizerURL)

	return New(db, optimizerClient, cfg), db
}

func getAuthTokenForPlanTests(t *testing.T, h *Handler, db *gorm.DB) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    "planuser@example.com",
		Password: string(hashedPassword),
		Name:     "Plan User",
		Role:     "user",
	}
	database.CreateUser(db, user)

	loginBody, _ := json.Marshal(LoginRequest{
		Email:    "planuser@example.com",
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

// TestCreatePlan tests plan creation
func TestCreatePlan(t *testing.T) {
	h, db := setupPlanTestHandler(t)
	token := getAuthTokenForPlanTests(t, h, db)

	// Create warehouse
	warehouse := &models.Warehouse{
		Name:      "Test Warehouse",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	database.CreateWarehouse(db, warehouse)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.POST("/api/v1/plans", h.CreatePlan)

	tests := []struct {
		name           string
		requestBody    PlanRequest
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder) bool
	}{
		{
			name: "valid plan",
			requestBody: PlanRequest{
				Name:        "Test Plan",
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-07",
				WarehouseID: warehouse.ID,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Data    models.Plan
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return response.Success && response.Data.Status == "draft"
			},
		},
		{
			name: "invalid date format",
			requestBody: PlanRequest{
				Name:        "Invalid Plan",
				StartDate:   "invalid-date",
				EndDate:     "2024-01-07",
				WarehouseID: warehouse.ID,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Error   string
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return !response.Success
			},
		},
		{
			name: "end date before start date",
			requestBody: PlanRequest{
				Name:        "Invalid Plan",
				StartDate:   "2024-01-07",
				EndDate:     "2024-01-01",
				WarehouseID: warehouse.ID,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Error   string
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return !response.Success && response.Error != ""
			},
		},
		{
			name: "missing required fields",
			requestBody: PlanRequest{
				Name: "Incomplete Plan",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Error   string
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return !response.Success
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/plans", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreatePlan() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkResponse != nil && !tt.checkResponse(w) {
				t.Error("CreatePlan() response validation failed")
			}
		})
	}
}

// TestGetPlan tests plan retrieval
func TestGetPlan(t *testing.T) {
	h, db := setupPlanTestHandler(t)
	token := getAuthTokenForPlanTests(t, h, db)

	// Create plan
	plan := &models.Plan{
		Name:        "Test Plan",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
		Status:      "draft",
		WarehouseID: nil,
	}
	database.CreatePlan(db, plan)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.GET("/api/v1/plans/:id", h.GetPlan)

	tests := []struct {
		name           string
		planID         string
		expectedStatus int
	}{
		{
			name:           "existing plan",
			planID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			planID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non-existent plan",
			planID:         "99999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/plans/"+tt.planID, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetPlan() status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

// TestDeletePlan tests plan deletion
func TestDeletePlan(t *testing.T) {
	h, db := setupPlanTestHandler(t)
	token := getAuthTokenForPlanTests(t, h, db)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.DELETE("/api/v1/plans/:id", h.DeletePlan)

	// Create plan to delete
	plan := &models.Plan{
		Name:      "To Delete",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
		Status:    "draft",
	}
	database.CreatePlan(db, plan)

	// Delete plan
	req := httptest.NewRequest("DELETE", "/api/v1/plans/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("DeletePlan() status = %d, want %d", w.Code, http.StatusOK)
	}

	// Verify deletion
	_, err := database.GetPlan(db, plan.ID)
	if err != database.ErrNotFound {
		t.Errorf("DeletePlan() plan still exists, error = %v", err)
	}
}

// TestGetPlanRoutes tests route retrieval
func TestGetPlanRoutes(t *testing.T) {
	h, db := setupPlanTestHandler(t)
	token := getAuthTokenForPlanTests(t, h, db)

	// Create plan with routes
	plan := &models.Plan{
		Name:      "Plan With Routes",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
		Status:    "optimized",
	}
	database.CreatePlan(db, plan)

	route := &models.Route{
		PlanID:        plan.ID,
		Day:           1,
		Date:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		TotalDistance: 100.0,
		TotalCost:     200.0,
		TotalLoad:     500.0,
	}
	database.CreateRoute(db, route)

	router := gin.New()
	router.Use(h.AuthMiddleware())
	router.GET("/api/v1/plans/:id/routes", h.GetPlanRoutes)

	req := httptest.NewRequest("GET", "/api/v1/plans/1/routes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetPlanRoutes() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response struct {
		Success bool
		Data    []models.Route
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Error("GetPlanRoutes() returned success=false")
	}
	if len(response.Data) == 0 {
		t.Error("GetPlanRoutes() returned empty routes")
	}
}
