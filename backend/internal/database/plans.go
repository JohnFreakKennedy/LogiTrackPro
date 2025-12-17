package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

func ListPlans(db *gorm.DB) ([]models.Plan, error) {
	var plans []models.Plan
	err := db.Order("created_at DESC").Find(&plans).Error
	return plans, err
}

func GetPlan(db *gorm.DB, id int64) (*models.Plan, error) {
	p := &models.Plan{}
	err := db.First(p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func CreatePlan(db *gorm.DB, p *models.Plan) error {
	return db.Create(p).Error
}

func UpdatePlanStatus(db *gorm.DB, id int64, status string, totalCost, totalDistance float64) error {
	result := db.Model(&models.Plan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         status,
		"total_cost":     totalCost,
		"total_distance": totalDistance,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func UpdatePlanStatusTx(tx *gorm.DB, id int64, status string, totalCost, totalDistance float64) error {
	result := tx.Model(&models.Plan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         status,
		"total_cost":     totalCost,
		"total_distance": totalDistance,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func DeletePlan(db *gorm.DB, id int64) error {
	result := db.Delete(&models.Plan{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func CountActivePlans(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&models.Plan{}).
		Where("status IN ?", []string{"draft", "optimizing", "optimized"}).
		Count(&count).Error
	return int(count), err
}

func GetRecentPlans(db *gorm.DB, limit int) ([]models.Plan, error) {
	var plans []models.Plan
	err := db.Order("created_at DESC").Limit(limit).Find(&plans).Error
	return plans, err
}

