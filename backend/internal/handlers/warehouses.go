package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type WarehouseRequest struct {
	Name            string  `json:"name" binding:"required"`
	Address         string  `json:"address"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
	Capacity        float64 `json:"capacity"`
	CurrentStock    float64 `json:"current_stock"`
	HoldingCost     float64 `json:"holding_cost"`
	ReplenishmentQty float64 `json:"replenishment_qty"`
}

// ListWarehouses handles GET /api/v1/warehouses
func (h *Handler) ListWarehouses(c *gin.Context) {
	warehouses, err := database.ListWarehouses(h.db)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch warehouses")
		return
	}
	if warehouses == nil {
		warehouses = []models.Warehouse{}
	}
	successResponse(c, warehouses)
}

// GetWarehouse handles GET /api/v1/warehouses/:id
func (h *Handler) GetWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	warehouse, err := database.GetWarehouse(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Warehouse not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch warehouse")
		return
	}
	successResponse(c, warehouse)
}

// CreateWarehouse handles POST /api/v1/warehouses
func (h *Handler) CreateWarehouse(c *gin.Context) {
	var req WarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	warehouse := &models.Warehouse{
		Name:            req.Name,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Capacity:        req.Capacity,
		CurrentStock:    req.CurrentStock,
		HoldingCost:     req.HoldingCost,
		ReplenishmentQty: req.ReplenishmentQty,
	}

	if err := database.CreateWarehouse(h.db, warehouse); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create warehouse")
		return
	}
	createdResponse(c, warehouse)
}

// UpdateWarehouse handles PUT /api/v1/warehouses/:id
func (h *Handler) UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	var req WarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	warehouse := &models.Warehouse{
		ID:              id,
		Name:            req.Name,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Capacity:        req.Capacity,
		CurrentStock:    req.CurrentStock,
		HoldingCost:     req.HoldingCost,
		ReplenishmentQty: req.ReplenishmentQty,
	}

	if err := database.UpdateWarehouse(h.db, warehouse); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Warehouse not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to update warehouse")
		return
	}
	successResponse(c, warehouse)
}

// DeleteWarehouse handles DELETE /api/v1/warehouses/:id
func (h *Handler) DeleteWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	if err := database.DeleteWarehouse(h.db, id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Warehouse not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to delete warehouse")
		return
	}
	successResponse(c, gin.H{"message": "Warehouse deleted successfully"})
}

