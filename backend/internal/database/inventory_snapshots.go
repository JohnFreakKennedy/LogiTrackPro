package database

import (
	"errors"
	"time"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/gorm"
)

// CreateInventorySnapshot creates a new inventory snapshot
func CreateInventorySnapshot(db *gorm.DB, snapshot *models.InventorySnapshot) error {
	return db.Create(snapshot).Error
}

// GetInventorySnapshots retrieves inventory snapshots with filters
func GetInventorySnapshots(db *gorm.DB, entityType string, entityID int64, startDate, endDate *time.Time) ([]models.InventorySnapshot, error) {
	query := db.Where("entity_type = ? AND entity_id = ?", entityType, entityID)

	if startDate != nil {
		query = query.Where("snapshot_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("snapshot_date <= ?", endDate)
	}

	var snapshots []models.InventorySnapshot
	err := query.Order("snapshot_time DESC").Find(&snapshots).Error
	return snapshots, err
}

// GetInventorySnapshotByDate retrieves snapshot for specific date
func GetInventorySnapshotByDate(db *gorm.DB, entityType string, entityID int64, date time.Time) (*models.InventorySnapshot, error) {
	snapshot := &models.InventorySnapshot{}
	err := db.Where("entity_type = ? AND entity_id = ? AND snapshot_date = ?",
		entityType, entityID, date).
		Order("snapshot_time DESC").
		First(snapshot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return snapshot, nil
}

// GetLatestInventorySnapshot retrieves the most recent snapshot
func GetLatestInventorySnapshot(db *gorm.DB, entityType string, entityID int64) (*models.InventorySnapshot, error) {
	snapshot := &models.InventorySnapshot{}
	err := db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("snapshot_time DESC").
		First(snapshot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return snapshot, nil
}

// CreateDailyInventorySnapshots creates snapshots for all customers/warehouses for a date
func CreateDailyInventorySnapshots(db *gorm.DB, snapshotDate time.Time, reason string) error {
	// Create snapshots for all customers
	var customers []models.Customer
	if err := db.Find(&customers).Error; err != nil {
		return err
	}

	now := time.Now()
	for _, customer := range customers {
		snapshot := &models.InventorySnapshot{
			EntityType:     "customer",
			EntityID:       customer.ID,
			SnapshotDate:   snapshotDate,
			SnapshotTime:   now,
			InventoryLevel: customer.CurrentInventory,
			DemandRate:     customer.DemandRate,
			MinInventory:   customer.MinInventory,
			MaxInventory:   customer.MaxInventory,
			SnapshotReason: reason,
		}
		if err := db.Create(snapshot).Error; err != nil {
			return err
		}
	}

	// Create snapshots for all warehouses
	var warehouses []models.Warehouse
	if err := db.Find(&warehouses).Error; err != nil {
		return err
	}

	for _, warehouse := range warehouses {
		snapshot := &models.InventorySnapshot{
			EntityType:     "warehouse",
			EntityID:       warehouse.ID,
			SnapshotDate:   snapshotDate,
			SnapshotTime:   now,
			InventoryLevel: warehouse.CurrentStock,
			SnapshotReason: reason,
		}
		if err := db.Create(snapshot).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetInventoryHistory retrieves inventory history for analytics
func GetInventoryHistory(db *gorm.DB, entityType string, entityID int64, days int) ([]models.InventorySnapshot, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	var snapshots []models.InventorySnapshot
	err := db.Where("entity_type = ? AND entity_id = ? AND snapshot_date >= ?",
		entityType, entityID, startDate).
		Order("snapshot_time ASC").
		Find(&snapshots).Error
	return snapshots, err
}

// GetInventorySnapshotsByPlan retrieves snapshots associated with a plan
func GetInventorySnapshotsByPlan(db *gorm.DB, planID int64) ([]models.InventorySnapshot, error) {
	var snapshots []models.InventorySnapshot
	err := db.Where("plan_id = ?", planID).
		Order("snapshot_time ASC").
		Find(&snapshots).Error
	return snapshots, err
}
