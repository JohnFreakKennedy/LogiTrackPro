package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateInventorySnapshotRequest struct {
	EntityType     string `json:"entity_type" binding:"required,oneof=customer warehouse"`
	EntityID       int64  `json:"entity_id" binding:"required"`
	SnapshotDate   string `json:"snapshot_date" binding:"required"`
	SnapshotReason string `json:"snapshot_reason"`
	PlanID         *int64 `json:"plan_id"`
	RouteID        *int64 `json:"route_id"`
}

type GetInventoryHistoryRequest struct {
	EntityType string `form:"entity_type" binding:"required,oneof=customer warehouse"`
	EntityID   int64  `form:"entity_id" binding:"required"`
	Days       int    `form:"days" binding:"min=1,max=365"`
}

// CreateInventorySnapshot handles POST /api/v1/inventory-snapshots
func (h *Handler) CreateInventorySnapshot(c *gin.Context) {
	var req CreateInventorySnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	snapshotDate, err := time.Parse("2006-01-02", req.SnapshotDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD)")
		return
	}

	// Get current inventory level based on entity type
	var inventoryLevel float64
	if req.EntityType == "customer" {
		customer, err := database.GetCustomer(h.db, req.EntityID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				errorResponse(c, http.StatusNotFound, "Customer not found")
				return
			}
			errorResponse(c, http.StatusInternalServerError, "Failed to fetch customer")
			return
		}
		inventoryLevel = customer.CurrentInventory
	} else {
		warehouse, err := database.GetWarehouse(h.db, req.EntityID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				errorResponse(c, http.StatusNotFound, "Warehouse not found")
				return
			}
			errorResponse(c, http.StatusInternalServerError, "Failed to fetch warehouse")
			return
		}
		inventoryLevel = warehouse.CurrentStock
	}

	snapshot := &models.InventorySnapshot{
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		SnapshotDate:   snapshotDate,
		SnapshotTime:   time.Now(),
		InventoryLevel: inventoryLevel,
		SnapshotReason: req.SnapshotReason,
		PlanID:         req.PlanID,
		RouteID:        req.RouteID,
	}

	if err := database.CreateInventorySnapshot(h.db, snapshot); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create inventory snapshot")
		return
	}

	createdResponse(c, snapshot)
}

// GetInventoryHistory handles GET /api/v1/inventory-history
func (h *Handler) GetInventoryHistory(c *gin.Context) {
	var req GetInventoryHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	if req.Days == 0 {
		req.Days = 30 // Default to 30 days
	}

	snapshots, err := database.GetInventoryHistory(h.db, req.EntityType, req.EntityID, req.Days)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch inventory history")
		return
	}

	if snapshots == nil {
		snapshots = []models.InventorySnapshot{}
	}

	successResponse(c, snapshots)
}

// GetInventorySnapshots handles GET /api/v1/inventory-snapshots
func (h *Handler) GetInventorySnapshots(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if entityType == "" || entityIDStr == "" {
		errorResponse(c, http.StatusBadRequest, "entity_type and entity_id are required")
		return
	}

	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid entity_id")
		return
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "Invalid start_date format")
			return
		}
		startDate = &parsed
	}
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "Invalid end_date format")
			return
		}
		endDate = &parsed
	}

	snapshots, err := database.GetInventorySnapshots(h.db, entityType, entityID, startDate, endDate)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch inventory snapshots")
		return
	}

	if snapshots == nil {
		snapshots = []models.InventorySnapshot{}
	}

	successResponse(c, snapshots)
}
