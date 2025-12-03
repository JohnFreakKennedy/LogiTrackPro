package handlers

import (
	"net/http"

	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDashboard handles GET /api/v1/analytics/dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	dashboard := &models.Dashboard{}

	// Get counts
	warehouseCount, _ := database.CountWarehouses(h.db)
	customerCount, _ := database.CountCustomers(h.db)
	vehicleCount, _ := database.CountVehicles(h.db)
	activePlans, _ := database.CountActivePlans(h.db)
	deliveries, _ := database.CountTotalDeliveries(h.db)
	distance, cost, _ := database.GetTotalDistanceAndCost(h.db)
	recentPlans, _ := database.GetRecentPlans(h.db, 5)

	dashboard.TotalWarehouses = warehouseCount
	dashboard.TotalCustomers = customerCount
	dashboard.TotalVehicles = vehicleCount
	dashboard.ActivePlans = activePlans
	dashboard.TotalDeliveries = deliveries
	dashboard.TotalDistanceKm = distance
	dashboard.TotalCost = cost
	dashboard.RecentPlans = recentPlans

	if dashboard.RecentPlans == nil {
		dashboard.RecentPlans = []models.Plan{}
	}

	successResponse(c, dashboard)
}

// GetSummary handles GET /api/v1/analytics/summary
func (h *Handler) GetSummary(c *gin.Context) {
	warehouseCount, _ := database.CountWarehouses(h.db)
	customerCount, _ := database.CountCustomers(h.db)
	vehicleCount, _ := database.CountVehicles(h.db)
	activePlans, _ := database.CountActivePlans(h.db)

	successResponse(c, gin.H{
		"warehouses":   warehouseCount,
		"customers":    customerCount,
		"vehicles":     vehicleCount,
		"active_plans": activePlans,
	})
}

