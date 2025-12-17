package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

func ListVehicles(db *gorm.DB) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := db.Order("name").Find(&vehicles).Error
	return vehicles, err
}

func ListAvailableVehiclesByWarehouse(db *gorm.DB, warehouseID int64) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := db.Where("warehouse_id = ? AND available = ?", warehouseID, true).
		Order("name").
		Find(&vehicles).Error
	return vehicles, err
}

func GetVehicle(db *gorm.DB, id int64) (*models.Vehicle, error) {
	v := &models.Vehicle{}
	err := db.First(v, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return v, nil
}

func CreateVehicle(db *gorm.DB, v *models.Vehicle) error {
	return db.Create(v).Error
}

func UpdateVehicle(db *gorm.DB, v *models.Vehicle) error {
	result := db.Model(v).Updates(models.Vehicle{
		Name:        v.Name,
		Capacity:    v.Capacity,
		CostPerKm:   v.CostPerKm,
		FixedCost:   v.FixedCost,
		MaxDistance: v.MaxDistance,
		Available:   v.Available,
		WarehouseID: v.WarehouseID,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteVehicle(db *gorm.DB, id int64) error {
	result := db.Delete(&models.Vehicle{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func CountVehicles(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&models.Vehicle{}).Count(&count).Error
	return int(count), err
}

