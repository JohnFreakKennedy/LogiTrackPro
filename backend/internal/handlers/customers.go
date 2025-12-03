package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CustomerRequest struct {
	Name             string  `json:"name" binding:"required"`
	Address          string  `json:"address"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	DemandRate       float64 `json:"demand_rate"`
	MaxInventory     float64 `json:"max_inventory"`
	CurrentInventory float64 `json:"current_inventory"`
	MinInventory     float64 `json:"min_inventory"`
	HoldingCost      float64 `json:"holding_cost"`
	Priority         int     `json:"priority"`
}

// ListCustomers handles GET /api/v1/customers
func (h *Handler) ListCustomers(c *gin.Context) {
	customers, err := database.ListCustomers(h.db)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}
	if customers == nil {
		customers = []models.Customer{}
	}
	successResponse(c, customers)
}

// GetCustomer handles GET /api/v1/customers/:id
func (h *Handler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := database.GetCustomer(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Customer not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch customer")
		return
	}
	successResponse(c, customer)
}

// CreateCustomer handles POST /api/v1/customers
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	customer := &models.Customer{
		Name:             req.Name,
		Address:          req.Address,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		DemandRate:       req.DemandRate,
		MaxInventory:     req.MaxInventory,
		CurrentInventory: req.CurrentInventory,
		MinInventory:     req.MinInventory,
		HoldingCost:      req.HoldingCost,
		Priority:         req.Priority,
	}

	if err := database.CreateCustomer(h.db, customer); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create customer")
		return
	}
	createdResponse(c, customer)
}

// UpdateCustomer handles PUT /api/v1/customers/:id
func (h *Handler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	customer := &models.Customer{
		ID:               id,
		Name:             req.Name,
		Address:          req.Address,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		DemandRate:       req.DemandRate,
		MaxInventory:     req.MaxInventory,
		CurrentInventory: req.CurrentInventory,
		MinInventory:     req.MinInventory,
		HoldingCost:      req.HoldingCost,
		Priority:         req.Priority,
	}

	if err := database.UpdateCustomer(h.db, customer); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Customer not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to update customer")
		return
	}
	successResponse(c, customer)
}

// DeleteCustomer handles DELETE /api/v1/customers/:id
func (h *Handler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := database.DeleteCustomer(h.db, id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Customer not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to delete customer")
		return
	}
	successResponse(c, gin.H{"message": "Customer deleted successfully"})
}

