package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

func ListWarehouses(db *gorm.DB) ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	err := db.Order("name").Find(&warehouses).Error
	return warehouses, err
}

func GetWarehouse(db *gorm.DB, id int64) (*models.Warehouse, error) {
	w := &models.Warehouse{}
	err := db.First(w, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return w, nil
}

func CreateWarehouse(db *gorm.DB, w *models.Warehouse) error {
	return db.Create(w).Error
}

func UpdateWarehouse(db *gorm.DB, w *models.Warehouse) error {
	result := db.Model(w).Updates(models.Warehouse{
		Name:             w.Name,
		Address:          w.Address,
		Latitude:         w.Latitude,
		Longitude:        w.Longitude,
		Capacity:         w.Capacity,
		CurrentStock:     w.CurrentStock,
		HoldingCost:      w.HoldingCost,
		ReplenishmentQty: w.ReplenishmentQty,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	// Reload to get updated_at
	return db.First(w, w.ID).Error
}

func DeleteWarehouse(db *gorm.DB, id int64) error {
	result := db.Delete(&models.Warehouse{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func CountWarehouses(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&models.Warehouse{}).Count(&count).Error
	return int(count), err
}
