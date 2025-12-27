package database

import (
	"fmt"

	"LogiTrackPro/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	// AutoMigrate will create tables, missing columns, missing indexes, etc.
	// It will NOT delete unused columns to protect your data.
	err := db.AutoMigrate(
		&models.User{},
		&models.Warehouse{},
		&models.Customer{},
		&models.Vehicle{},
		&models.Plan{},
		&models.Route{},
		&models.Stop{},
		&models.RouteExecution{},
		&models.StopExecution{},
		&models.InventorySnapshot{},
		&models.Product{},
		&models.CustomerProductInventory{},
		&models.StopProductQuantity{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
