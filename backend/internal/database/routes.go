package database

import (
	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

func GetRoutesByPlan(db *gorm.DB, planID int64) ([]models.Route, error) {
	var routes []models.Route
	err := db.Where("plan_id = ?", planID).
		Preload("Vehicle").
		Preload("Stops.Customer").
		Order("day, id").
		Find(&routes).Error
	return routes, err
}

func CreateRoute(db *gorm.DB, r *models.Route) error {
	return db.Create(r).Error
}

func CreateRouteTx(tx *gorm.DB, r *models.Route) error {
	return tx.Create(r).Error
}

func DeleteRoutesByPlan(db *gorm.DB, planID int64) error {
	return db.Where("plan_id = ?", planID).Delete(&models.Route{}).Error
}

func DeleteRoutesByPlanTx(tx *gorm.DB, planID int64) error {
	return tx.Where("plan_id = ?", planID).Delete(&models.Route{}).Error
}

func GetStopsByRoute(db *gorm.DB, routeID int64) ([]models.Stop, error) {
	var stops []models.Stop
	err := db.Where("route_id = ?", routeID).
		Preload("Customer").
		Order("sequence").
		Find(&stops).Error
	return stops, err
}

func CreateStop(db *gorm.DB, s *models.Stop) error {
	return db.Create(s).Error
}

func CreateStopTx(tx *gorm.DB, s *models.Stop) error {
	return tx.Create(s).Error
}

func CountTotalDeliveries(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&models.Stop{}).Count(&count).Error
	return int(count), err
}

func GetTotalDistanceAndCost(db *gorm.DB) (float64, float64, error) {
	var result struct {
		TotalDistance float64
		TotalCost     float64
	}
	err := db.Model(&models.Route{}).
		Select("COALESCE(SUM(total_distance), 0) as total_distance, COALESCE(SUM(total_cost), 0) as total_cost").
		Scan(&result).Error
	return result.TotalDistance, result.TotalCost, err
}
