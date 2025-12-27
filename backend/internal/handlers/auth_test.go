package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// setupTestHandler creates a test handler with in-memory database
func setupTestHandler(t *testing.T) *Handler {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing-only",
		JWTExpiry: 24,
	}

	optimizerClient := optimizer.NewClient("http://localhost:8000")

	return New(db, optimizerClient, cfg)
}

// TestRegister tests user registration
func TestRegister(t *testing.T) {
	h := setupTestHandler(t)

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder) bool
	}{
		{
			name: "valid registration",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Data    AuthResponse
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return response.Success && response.Data.Token != "" && response.Data.User.Email == "test@example.com"
			},
		},
		{
			name: "duplicate email",
			requestBody: RegisterRequest{
				Email:    "duplicate@example.com",
				Password: "password123",
				Name:     "First User",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				// Register first user
				return true
			},
		},
		{
			name: "invalid email format",
			requestBody: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
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
			name: "password too short",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "12345",
				Name:     "Test User",
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
			name: "missing required fields",
			requestBody: RegisterRequest{
				Email: "test@example.com",
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
			// Handle duplicate email case
			if tt.name == "duplicate email" {
				// First registration
				body, _ := json.Marshal(tt.requestBody)
				req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router := gin.New()
				router.POST("/api/v1/auth/register", h.Register)
				router.ServeHTTP(w, req)

				// Second registration with same email
				w2 := httptest.NewRecorder()
				router.ServeHTTP(w2, req)
				if w2.Code != http.StatusConflict {
					t.Errorf("Register() status = %d, want %d", w2.Code, http.StatusConflict)
				}
				return
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := gin.New()
			router.POST("/api/v1/auth/register", h.Register)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Register() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkResponse != nil && !tt.checkResponse(w) {
				t.Error("Register() response validation failed")
			}
		})
	}
}

// TestLogin tests user login
func TestLogin(t *testing.T) {
	h := setupTestHandler(t)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    "login@example.com",
		Password: string(hashedPassword),
		Name:     "Login User",
		Role:     "user",
	}
	database.CreateUser(h.db, user)

	tests := []struct {
		name           string
		requestBody    LoginRequest
		expectedStatus int
		checkResponse  func(*httptest.ResponseRecorder) bool
	}{
		{
			name: "valid login",
			requestBody: LoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Data    AuthResponse
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return response.Success && response.Data.Token != "" && response.Data.User.Email == "login@example.com"
			},
		},
		{
			name: "invalid email",
			requestBody: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Error   string
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return !response.Success && response.Error == "Invalid credentials"
			},
		},
		{
			name: "invalid password",
			requestBody: LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(w *httptest.ResponseRecorder) bool {
				var response struct {
					Success bool
					Error   string
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				return !response.Success && response.Error == "Invalid credentials"
			},
		},
		{
			name: "missing credentials",
			requestBody: LoginRequest{
				Email: "login@example.com",
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
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := gin.New()
			router.POST("/api/v1/auth/login", h.Login)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Login() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkResponse != nil && !tt.checkResponse(w) {
				t.Error("Login() response validation failed")
			}
		})
	}
}

// TestAuthMiddleware tests JWT authentication middleware
func TestAuthMiddleware(t *testing.T) {
	h := setupTestHandler(t)

	// Create user and get token
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    "middleware@example.com",
		Password: string(hashedPassword),
		Name:     "Middleware User",
		Role:     "user",
	}
	database.CreateUser(h.db, user)

	// Login to get token
	loginBody, _ := json.Marshal(LoginRequest{
		Email:    "middleware@example.com",
		Password: "password123",
	})
	loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router := gin.New()
	router.POST("/api/v1/auth/login", h.Login)
	router.ServeHTTP(loginW, loginW)

	var loginResponse struct {
		Success bool
		Data    AuthResponse
	}
	json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	token := loginResponse.Data.Token

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
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
		{
			name:           "malformed header",
			authHeader:     "InvalidFormat token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/me", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router := gin.New()
			router.GET("/api/v1/me", h.AuthMiddleware(), h.GetCurrentUser)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("AuthMiddleware() status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

// TestGenerateToken tests JWT token generation
func TestGenerateToken(t *testing.T) {
	h := setupTestHandler(t)

	user := &models.User{
		ID:    1,
		Email: "token@example.com",
		Name:  "Token User",
		Role:  "user",
	}

	token, expiresAt, err := h.generateToken(user)
	if err != nil {
		t.Fatalf("generateToken() error = %v", err)
	}

	if token == "" {
		t.Error("generateToken() returned empty token")
	}

	if expiresAt.IsZero() {
		t.Error("generateToken() returned zero expiration time")
	}

	// Verify token can be parsed
	claims, err := h.parseToken(token)
	if err != nil {
		t.Fatalf("parseToken() error = %v", err)
	}

	if claims.Subject != "1" {
		t.Errorf("parseToken() Subject = %v, want 1", claims.Subject)
	}
}

// TestParseToken tests JWT token parsing
func TestParseToken(t *testing.T) {
	h := setupTestHandler(t)

	user := &models.User{
		ID:    1,
		Email: "parse@example.com",
		Name:  "Parse User",
		Role:  "user",
	}

	token, _, _ := h.generateToken(user)

	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{
			name:      "valid token",
			token:     token,
			wantError: false,
		},
		{
			name:      "invalid token",
			token:     "invalid.token.here",
			wantError: true,
		},
		{
			name:      "empty token",
			token:     "",
			wantError: true,
		},
		{
			name:      "wrong secret",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.parseToken(tt.token)
			if (err != nil) != tt.wantError {
				t.Errorf("parseToken() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
