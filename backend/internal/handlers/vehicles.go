package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type VehicleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Capacity    float64 `json:"capacity" binding:"required"`
	CostPerKm   float64 `json:"cost_per_km"`
	FixedCost   float64 `json:"fixed_cost"`
	MaxDistance float64 `json:"max_distance"`
	Available   bool    `json:"available"`
	WarehouseID int64   `json:"warehouse_id"`
}

// ListVehicles handles GET /api/v1/vehicles
func (h *Handler) ListVehicles(c *gin.Context) {
	vehicles, err := database.ListVehicles(h.db)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch vehicles")
		return
	}
	if vehicles == nil {
		vehicles = []models.Vehicle{}
	}
	successResponse(c, vehicles)
}

// GetVehicle handles GET /api/v1/vehicles/:id
func (h *Handler) GetVehicle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	vehicle, err := database.GetVehicle(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Vehicle not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch vehicle")
		return
	}
	successResponse(c, vehicle)
}

// CreateVehicle handles POST /api/v1/vehicles
func (h *Handler) CreateVehicle(c *gin.Context) {
	var req VehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	vehicle := &models.Vehicle{
		Name:        req.Name,
		Capacity:    req.Capacity,
		CostPerKm:   req.CostPerKm,
		FixedCost:   req.FixedCost,
		MaxDistance: req.MaxDistance,
		Available:   req.Available,
		WarehouseID: req.WarehouseID,
	}

	if err := database.CreateVehicle(h.db, vehicle); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create vehicle")
		return
	}
	createdResponse(c, vehicle)
}

// UpdateVehicle handles PUT /api/v1/vehicles/:id
func (h *Handler) UpdateVehicle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	var req VehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	vehicle := &models.Vehicle{
		ID:          id,
		Name:        req.Name,
		Capacity:    req.Capacity,
		CostPerKm:   req.CostPerKm,
		FixedCost:   req.FixedCost,
		MaxDistance: req.MaxDistance,
		Available:   req.Available,
		WarehouseID: req.WarehouseID,
	}

	if err := database.UpdateVehicle(h.db, vehicle); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Vehicle not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to update vehicle")
		return
	}
	successResponse(c, vehicle)
}

// DeleteVehicle handles DELETE /api/v1/vehicles/:id
func (h *Handler) DeleteVehicle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	if err := database.DeleteVehicle(h.db, id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Vehicle not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to delete vehicle")
		return
	}
	successResponse(c, gin.H{"message": "Vehicle deleted successfully"})
}

