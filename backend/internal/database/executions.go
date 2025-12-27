package database

import (
	"errors"
	"time"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

// CreateRouteExecution creates a new route execution record
func CreateRouteExecution(db *gorm.DB, execution *models.RouteExecution) error {
	return db.Create(execution).Error
}

// GetRouteExecution retrieves a route execution by ID
func GetRouteExecution(db *gorm.DB, id int64) (*models.RouteExecution, error) {
	execution := &models.RouteExecution{}
	err := db.Preload("Route").Preload("StopExecutions.Stop").
		First(execution, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return execution, nil
}

// GetRouteExecutionsByRoute retrieves all executions for a route
func GetRouteExecutionsByRoute(db *gorm.DB, routeID int64) ([]models.RouteExecution, error) {
	var executions []models.RouteExecution
	err := db.Where("route_id = ?", routeID).
		Preload("StopExecutions").
		Order("created_at DESC").
		Find(&executions).Error
	return executions, err
}

// UpdateRouteExecution updates a route execution
func UpdateRouteExecution(db *gorm.DB, execution *models.RouteExecution) error {
	result := db.Model(execution).Updates(models.RouteExecution{
		Status:          execution.Status,
		ActualDistance:  execution.ActualDistance,
		ActualCost:      execution.ActualCost,
		ActualLoad:      execution.ActualLoad,
		ActualStartTime: execution.ActualStartTime,
		ActualEndTime:   execution.ActualEndTime,
		DriverNotes:     execution.DriverNotes,
		DeviationReason: execution.DeviationReason,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// StartRouteExecution marks a route execution as in progress
func StartRouteExecution(db *gorm.DB, executionID int64) error {
	now := time.Now()
	result := db.Model(&models.RouteExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]interface{}{
			"status":            "in_progress",
			"actual_start_time": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// CompleteRouteExecution marks a route execution as completed
func CompleteRouteExecution(db *gorm.DB, executionID int64, actualDistance, actualCost, actualLoad float64) error {
	now := time.Now()
	result := db.Model(&models.RouteExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]interface{}{
			"status":          "completed",
			"actual_distance": actualDistance,
			"actual_cost":     actualCost,
			"actual_load":     actualLoad,
			"actual_end_time": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateStopExecution creates a new stop execution record
func CreateStopExecution(db *gorm.DB, execution *models.StopExecution) error {
	return db.Create(execution).Error
}

// GetStopExecutionsByRouteExecution retrieves all stop executions for a route execution
func GetStopExecutionsByRouteExecution(db *gorm.DB, routeExecutionID int64) ([]models.StopExecution, error) {
	var executions []models.StopExecution
	err := db.Where("route_execution_id = ?", routeExecutionID).
		Preload("Stop").
		Order("sequence").
		Find(&executions).Error
	return executions, err
}

// UpdateStopExecution updates a stop execution
func UpdateStopExecution(db *gorm.DB, execution *models.StopExecution) error {
	result := db.Model(execution).Updates(models.StopExecution{
		Status:              execution.Status,
		ActualQuantity:      execution.ActualQuantity,
		ActualArrivalTime:   execution.ActualArrivalTime,
		ActualDepartureTime: execution.ActualDepartureTime,
		ServiceDuration:     execution.ServiceDuration,
		Notes:               execution.Notes,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetExecutionStats calculates execution statistics for a plan
func GetExecutionStats(db *gorm.DB, planID int64) (map[string]interface{}, error) {
	var stats struct {
		TotalExecutions      int64
		CompletedExecutions  int64
		TotalPlannedCost     float64
		TotalActualCost      float64
		TotalPlannedDistance float64
		TotalActualDistance  float64
	}

	err := db.Table("route_executions").
		Select(`
			COUNT(*) as total_executions,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_executions,
			COALESCE(SUM(planned_cost), 0) as total_planned_cost,
			COALESCE(SUM(actual_cost), 0) as total_actual_cost,
			COALESCE(SUM(planned_distance), 0) as total_planned_distance,
			COALESCE(SUM(actual_distance), 0) as total_actual_distance
		`).
		Joins("JOIN routes ON route_executions.route_id = routes.id").
		Where("routes.plan_id = ?", planID).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_executions":       stats.TotalExecutions,
		"completed_executions":   stats.CompletedExecutions,
		"total_planned_cost":     stats.TotalPlannedCost,
		"total_actual_cost":      stats.TotalActualCost,
		"total_planned_distance": stats.TotalPlannedDistance,
		"total_actual_distance":  stats.TotalActualDistance,
	}

	if stats.TotalPlannedCost > 0 {
		result["cost_deviation_percent"] = ((stats.TotalActualCost - stats.TotalPlannedCost) / stats.TotalPlannedCost) * 100
	}
	if stats.TotalPlannedDistance > 0 {
		result["distance_deviation_percent"] = ((stats.TotalActualDistance - stats.TotalPlannedDistance) / stats.TotalPlannedDistance) * 100
	}

	return result, nil
}
