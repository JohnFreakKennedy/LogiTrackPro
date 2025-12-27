package database

import (
	"errors"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

// ListProducts retrieves all products
func ListProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	err := db.Order("name").Find(&products).Error
	return products, err
}

// GetProduct retrieves a product by ID
func GetProduct(db *gorm.DB, id int64) (*models.Product, error) {
	product := &models.Product{}
	err := db.First(product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return product, nil
}

// CreateProduct creates a new product
func CreateProduct(db *gorm.DB, product *models.Product) error {
	return db.Create(product).Error
}

// UpdateProduct updates a product
func UpdateProduct(db *gorm.DB, product *models.Product) error {
	result := db.Model(product).Updates(models.Product{
		Name:        product.Name,
		SKU:         product.SKU,
		Description: product.Description,
		Unit:        product.Unit,
		Weight:      product.Weight,
		Volume:      product.Volume,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteProduct deletes a product
func DeleteProduct(db *gorm.DB, id int64) error {
	result := db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetCustomerProductInventory retrieves product inventory for a customer
func GetCustomerProductInventory(db *gorm.DB, customerID int64) ([]models.CustomerProductInventory, error) {
	var inventory []models.CustomerProductInventory
	err := db.Where("customer_id = ?", customerID).
		Preload("Product").
		Find(&inventory).Error
	return inventory, err
}

// UpdateCustomerProductInventory updates product inventory for a customer
func UpdateCustomerProductInventory(db *gorm.DB, inventory *models.CustomerProductInventory) error {
	result := db.Model(inventory).
		Where("customer_id = ? AND product_id = ?", inventory.CustomerID, inventory.ProductID).
		Updates(models.CustomerProductInventory{
			CurrentInventory: inventory.CurrentInventory,
			MaxInventory:     inventory.MaxInventory,
			MinInventory:     inventory.MinInventory,
			DemandRate:       inventory.DemandRate,
			HoldingCost:      inventory.HoldingCost,
			Priority:         inventory.Priority,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Create if doesn't exist
		return db.Create(inventory).Error
	}
	return nil
}
