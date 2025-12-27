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

type CreateRouteExecutionRequest struct {
	RouteID int64 `json:"route_id" binding:"required"`
}

type UpdateRouteExecutionRequest struct {
	Status          string     `json:"status"`
	ActualDistance  float64    `json:"actual_distance"`
	ActualCost      float64    `json:"actual_cost"`
	ActualLoad      float64    `json:"actual_load"`
	ActualStartTime *time.Time `json:"actual_start_time"`
	ActualEndTime   *time.Time `json:"actual_end_time"`
	DriverNotes     string     `json:"driver_notes"`
	DeviationReason string     `json:"deviation_reason"`
}

type StartRouteExecutionRequest struct {
	ActualStartTime *time.Time `json:"actual_start_time"`
}

type CompleteRouteExecutionRequest struct {
	ActualDistance  float64    `json:"actual_distance"`
	ActualCost      float64    `json:"actual_cost"`
	ActualLoad      float64    `json:"actual_load"`
	ActualEndTime   *time.Time `json:"actual_end_time"`
	DriverNotes     string     `json:"driver_notes"`
	DeviationReason string     `json:"deviation_reason"`
}

// CreateRouteExecution handles POST /api/v1/routes/:id/executions
func (h *Handler) CreateRouteExecution(c *gin.Context) {
	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid route ID")
		return
	}

	// Verify route exists
	route, err := database.GetRouteByID(h.db, routeID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Route not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch route")
		return
	}

	// Create execution with planned values
	execution := &models.RouteExecution{
		RouteID:         routeID,
		Status:          "pending",
		PlannedDistance: route.TotalDistance,
		PlannedCost:     route.TotalCost,
		PlannedLoad:     route.TotalLoad,
	}

	if err := database.CreateRouteExecution(h.db, execution); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create route execution")
		return
	}

	createdResponse(c, execution)
}

// GetRouteExecution handles GET /api/v1/executions/:id
func (h *Handler) GetRouteExecution(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid execution ID")
		return
	}

	execution, err := database.GetRouteExecution(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Route execution not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch route execution")
		return
	}

	successResponse(c, execution)
}

// GetRouteExecutions handles GET /api/v1/routes/:id/executions
func (h *Handler) GetRouteExecutions(c *gin.Context) {
	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid route ID")
		return
	}

	executions, err := database.GetRouteExecutionsByRoute(h.db, routeID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch route executions")
		return
	}

	if executions == nil {
		executions = []models.RouteExecution{}
	}

	successResponse(c, executions)
}

// StartRouteExecution handles POST /api/v1/executions/:id/start
func (h *Handler) StartRouteExecution(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid execution ID")
		return
	}

	var req StartRouteExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	execution := &models.RouteExecution{
		ID:              id,
		Status:          "in_progress",
		ActualStartTime: req.ActualStartTime,
	}

	if execution.ActualStartTime == nil {
		now := time.Now()
		execution.ActualStartTime = &now
	}

	if err := database.UpdateRouteExecution(h.db, execution); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Route execution not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to start route execution")
		return
	}

	successResponse(c, execution)
}

// CompleteRouteExecution handles POST /api/v1/executions/:id/complete
func (h *Handler) CompleteRouteExecution(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid execution ID")
		return
	}

	var req CompleteRouteExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.ActualEndTime == nil {
		now := time.Now()
		req.ActualEndTime = &now
	}

	err = database.CompleteRouteExecution(h.db, id, req.ActualDistance, req.ActualCost, req.ActualLoad)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Route execution not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to complete route execution")
		return
	}

	// Update notes and deviation reason if provided
	if req.DriverNotes != "" || req.DeviationReason != "" {
		execution := &models.RouteExecution{
			ID:              id,
			DriverNotes:     req.DriverNotes,
			DeviationReason: req.DeviationReason,
			ActualEndTime:   req.ActualEndTime,
		}
		database.UpdateRouteExecution(h.db, execution)
	}

	execution, _ := database.GetRouteExecution(h.db, id)
	successResponse(c, execution)
}

// UpdateRouteExecution handles PUT /api/v1/executions/:id
func (h *Handler) UpdateRouteExecution(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid execution ID")
		return
	}

	var req UpdateRouteExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	execution := &models.RouteExecution{
		ID:              id,
		Status:          req.Status,
		ActualDistance:  req.ActualDistance,
		ActualCost:      req.ActualCost,
		ActualLoad:      req.ActualLoad,
		ActualStartTime: req.ActualStartTime,
		ActualEndTime:   req.ActualEndTime,
		DriverNotes:     req.DriverNotes,
		DeviationReason: req.DeviationReason,
	}

	if err := database.UpdateRouteExecution(h.db, execution); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Route execution not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to update route execution")
		return
	}

	successResponse(c, execution)
}

// GetPlanExecutionStats handles GET /api/v1/plans/:id/execution-stats
func (h *Handler) GetPlanExecutionStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	stats, err := database.GetExecutionStats(h.db, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch execution statistics")
		return
	}

	successResponse(c, stats)
}
