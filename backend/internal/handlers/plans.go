package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"
	"LogiTrackPro/backend/internal/optimizer"

	"github.com/gin-gonic/gin"
)

type PlanRequest struct {
	Name        string `json:"name" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	WarehouseID int64  `json:"warehouse_id" binding:"required"`
}

// ListPlans handles GET /api/v1/plans
func (h *Handler) ListPlans(c *gin.Context) {
	plans, err := database.ListPlans(h.db)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch plans")
		return
	}
	if plans == nil {
		plans = []models.Plan{}
	}
	successResponse(c, plans)
}

// GetPlan handles GET /api/v1/plans/:id
func (h *Handler) GetPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	plan, err := database.GetPlan(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Plan not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch plan")
		return
	}

	// Load routes
	routes, err := database.GetRoutesByPlan(h.db, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch plan routes")
		return
	}
	plan.Routes = routes

	successResponse(c, plan)
}

// CreatePlan handles POST /api/v1/plans
func (h *Handler) CreatePlan(c *gin.Context) {
	var req PlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid start date format (use YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid end date format (use YYYY-MM-DD)")
		return
	}

	if endDate.Before(startDate) {
		errorResponse(c, http.StatusBadRequest, "End date must be after start date")
		return
	}

	userID := c.GetInt64("userID")

	plan := &models.Plan{
		Name:        req.Name,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "draft",
		WarehouseID: req.WarehouseID,
		CreatedBy:   userID,
	}

	if err := database.CreatePlan(h.db, plan); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to create plan")
		return
	}
	createdResponse(c, plan)
}

// DeletePlan handles DELETE /api/v1/plans/:id
func (h *Handler) DeletePlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	if err := database.DeletePlan(h.db, id); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Plan not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to delete plan")
		return
	}
	successResponse(c, gin.H{"message": "Plan deleted successfully"})
}

// GetPlanRoutes handles GET /api/v1/plans/:id/routes
func (h *Handler) GetPlanRoutes(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	routes, err := database.GetRoutesByPlan(h.db, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch routes")
		return
	}
	if routes == nil {
		routes = []models.Route{}
	}
	successResponse(c, routes)
}

// OptimizePlan handles POST /api/v1/plans/:id/optimize
func (h *Handler) OptimizePlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	// Get plan
	plan, err := database.GetPlan(h.db, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "Plan not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch plan")
		return
	}

	// Get warehouse
	warehouse, err := database.GetWarehouse(h.db, plan.WarehouseID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch warehouse")
		return
	}

	// Get customers
	customers, err := database.ListCustomers(h.db)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}

	if len(customers) == 0 {
		errorResponse(c, http.StatusBadRequest, "No customers to optimize")
		return
	}

	// Get available vehicles for this warehouse
	vehicles, err := database.ListAvailableVehiclesByWarehouse(h.db, warehouse.ID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch vehicles")
		return
	}

	if len(vehicles) == 0 {
		errorResponse(c, http.StatusBadRequest, "No available vehicles for optimization")
		return
	}

	// Calculate planning horizon (days)
	planningHorizon := int(plan.EndDate.Sub(plan.StartDate).Hours()/24) + 1

	// Build optimizer request
	optReq := &optimizer.OptimizeRequest{
		Warehouse: optimizer.WarehouseData{
			ID:        warehouse.ID,
			Latitude:  warehouse.Latitude,
			Longitude: warehouse.Longitude,
			Stock:     warehouse.CurrentStock,
		},
		Customers:       make([]optimizer.CustomerData, len(customers)),
		Vehicles:        make([]optimizer.VehicleData, len(vehicles)),
		PlanningHorizon: planningHorizon,
		StartDate:       plan.StartDate.Format("2006-01-02"),
	}

	for i, c := range customers {
		optReq.Customers[i] = optimizer.CustomerData{
			ID:               c.ID,
			Latitude:         c.Latitude,
			Longitude:        c.Longitude,
			DemandRate:       c.DemandRate,
			MaxInventory:     c.MaxInventory,
			CurrentInventory: c.CurrentInventory,
			MinInventory:     c.MinInventory,
			Priority:         c.Priority,
		}
	}

	for i, v := range vehicles {
		optReq.Vehicles[i] = optimizer.VehicleData{
			ID:          v.ID,
			Capacity:    v.Capacity,
			CostPerKm:   v.CostPerKm,
			FixedCost:   v.FixedCost,
			MaxDistance: v.MaxDistance,
		}
	}

	// Update plan status
	if err := database.UpdatePlanStatus(h.db, id, "optimizing", 0, 0); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to update plan status: "+err.Error())
		return
	}

	// Call optimizer
	optResp, err := h.optimizer.Optimize(optReq)
	if err != nil {
		if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
			errorResponse(c, http.StatusInternalServerError, "Optimization failed: "+err.Error()+". Revert failed: "+revertErr.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "Optimization failed: "+err.Error())
		}
		return
	}

	if !optResp.Success {
		if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
			errorResponse(c, http.StatusInternalServerError, "Optimization failed: "+optResp.Message+". Revert failed: "+revertErr.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "Optimization failed: "+optResp.Message)
		}
		return
	}

	// Delete existing routes
	if err := database.DeleteRoutesByPlan(h.db, id); err != nil {
		if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
			errorResponse(c, http.StatusInternalServerError, "Failed to clear routes: "+err.Error()+". Revert failed: "+revertErr.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "Failed to clear existing routes: "+err.Error())
		}
		return
	}

	// Save new routes
	for _, routeResult := range optResp.Routes {
		routeDate, err := time.Parse("2006-01-02", routeResult.Date)
		if err != nil {
			if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
				errorResponse(c, http.StatusInternalServerError, "Invalid date: "+err.Error()+". Revert failed: "+revertErr.Error())
			} else {
				errorResponse(c, http.StatusInternalServerError, "Failed to parse route date: "+err.Error())
			}
			return
		}
		var vehicleID *int64
		if routeResult.VehicleID != 0 {
			vID := routeResult.VehicleID
			vehicleID = &vID
		}
		route := &models.Route{
			PlanID:        id,
			VehicleID:     vehicleID,
			Day:           routeResult.Day,
			Date:          routeDate,
			TotalDistance: routeResult.TotalDistance,
			TotalCost:     routeResult.TotalCost,
			TotalLoad:     routeResult.TotalLoad,
		}

		if err := database.CreateRoute(h.db, route); err != nil {
			if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
				errorResponse(c, http.StatusInternalServerError, "Failed to save route: "+err.Error()+". Revert failed: "+revertErr.Error())
			} else {
				errorResponse(c, http.StatusInternalServerError, "Failed to save route: "+err.Error())
			}
			return
		}

		// Save stops
		for _, stopResult := range routeResult.Stops {
			stop := &models.Stop{
				RouteID:     route.ID,
				CustomerID:  stopResult.CustomerID,
				Sequence:    stopResult.Sequence,
				Quantity:    stopResult.Quantity,
				ArrivalTime: stopResult.ArrivalTime,
			}
			if err := database.CreateStop(h.db, stop); err != nil {
				if revertErr := database.UpdatePlanStatus(h.db, id, "draft", 0, 0); revertErr != nil {
					errorResponse(c, http.StatusInternalServerError, "Failed to save stop: "+err.Error()+". Revert failed: "+revertErr.Error())
				} else {
					errorResponse(c, http.StatusInternalServerError, "Failed to save stop: "+err.Error())
				}
				return
			}
		}
	}

	// Update plan status
	if err := database.UpdatePlanStatus(h.db, id, "optimized", optResp.TotalCost, optResp.TotalDistance); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to update plan status: "+err.Error())
		return
	}

	// Get updated plan with routes
	plan, err = database.GetPlan(h.db, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch updated plan: "+err.Error())
		return
	}

	routes, err := database.GetRoutesByPlan(h.db, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch updated routes: "+err.Error())
		return
	}
	plan.Routes = routes

	successResponse(c, plan)
}
