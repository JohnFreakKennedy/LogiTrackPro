package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

func ListCustomers(db *gorm.DB) ([]models.Customer, error) {
	var customers []models.Customer
	err := db.Order("name").Find(&customers).Error
	return customers, err
}

func GetCustomer(db *gorm.DB, id int64) (*models.Customer, error) {
	c := &models.Customer{}
	err := db.First(c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func CreateCustomer(db *gorm.DB, c *models.Customer) error {
	return db.Create(c).Error
}

func UpdateCustomer(db *gorm.DB, c *models.Customer) error {
	result := db.Model(c).Updates(models.Customer{
		Name:             c.Name,
		Address:          c.Address,
		Latitude:         c.Latitude,
		Longitude:        c.Longitude,
		DemandRate:       c.DemandRate,
		MaxInventory:     c.MaxInventory,
		CurrentInventory: c.CurrentInventory,
		MinInventory:     c.MinInventory,
		HoldingCost:      c.HoldingCost,
		Priority:         c.Priority,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteCustomer(db *gorm.DB, id int64) error {
	result := db.Delete(&models.Customer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func CountCustomers(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&models.Customer{}).Count(&count).Error
	return int(count), err
}

