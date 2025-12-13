package handlers

import (
	"database/sql"
	"net/http"

	"LogiTrackPro/backend/internal/config"
	"LogiTrackPro/backend/internal/optimizer"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db        *sql.DB
	optimizer *optimizer.Client
	config    *config.Config
}

func New(db *sql.DB, optimizerClient *optimizer.Client, cfg *config.Config) *Handler {
	return &Handler{
		db:        db,
		optimizer: optimizerClient,
		config:    cfg,
	}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := "connected"
	if err := h.db.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	// Check optimizer service
	optimizerStatus := "connected"
	if err := h.optimizer.HealthCheck(); err != nil {
		optimizerStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "LogiTrackPro API",
		"database":  dbStatus,
		"optimizer": optimizerStatus,
	})
}

// Response helpers
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func createdResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func errorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}

