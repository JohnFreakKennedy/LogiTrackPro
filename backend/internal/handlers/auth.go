package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      *models.User `json:"user"`
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := database.CreateUser(h.db, user); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			errorResponse(c, http.StatusConflict, "Email already registered")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token
	token, expiresAt, err := h.generateToken(user)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	createdResponse(c, AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	user, err := database.GetUserByEmail(h.db, req.Email)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		errorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, expiresAt, err := h.generateToken(user)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	successResponse(c, AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errorResponse(c, http.StatusUnauthorized, "No token provided")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.parseToken(tokenString)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	user, err := database.GetUserByID(h.db, userID)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	newToken, expiresAt, err := h.generateToken(user)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	successResponse(c, AuthResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// GetCurrentUser handles GET /api/v1/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt64("userID")
	user, err := database.GetUserByID(h.db, userID)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	successResponse(c, user)
}

// AuthMiddleware verifies JWT token
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorResponse(c, http.StatusUnauthorized, "No token provided")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := h.parseToken(tokenString)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (h *Handler) generateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(h.config.JWTExpiry) * time.Hour)
	
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "LogiTrackPro",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (h *Handler) parseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

