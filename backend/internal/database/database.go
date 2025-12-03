package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS warehouses (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			address TEXT,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			capacity DOUBLE PRECISION DEFAULT 0,
			current_stock DOUBLE PRECISION DEFAULT 0,
			holding_cost DOUBLE PRECISION DEFAULT 0,
			replenishment_qty DOUBLE PRECISION DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			address TEXT,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			demand_rate DOUBLE PRECISION DEFAULT 0,
			max_inventory DOUBLE PRECISION DEFAULT 0,
			current_inventory DOUBLE PRECISION DEFAULT 0,
			min_inventory DOUBLE PRECISION DEFAULT 0,
			holding_cost DOUBLE PRECISION DEFAULT 0,
			priority INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS vehicles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			capacity DOUBLE PRECISION NOT NULL,
			cost_per_km DOUBLE PRECISION DEFAULT 0,
			fixed_cost DOUBLE PRECISION DEFAULT 0,
			max_distance DOUBLE PRECISION DEFAULT 0,
			available BOOLEAN DEFAULT true,
			warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS plans (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			status VARCHAR(50) DEFAULT 'draft',
			total_cost DOUBLE PRECISION DEFAULT 0,
			total_distance DOUBLE PRECISION DEFAULT 0,
			warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE SET NULL,
			created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS routes (
			id SERIAL PRIMARY KEY,
			plan_id INTEGER REFERENCES plans(id) ON DELETE CASCADE,
			vehicle_id INTEGER REFERENCES vehicles(id) ON DELETE SET NULL,
			day INTEGER NOT NULL,
			date DATE NOT NULL,
			total_distance DOUBLE PRECISION DEFAULT 0,
			total_cost DOUBLE PRECISION DEFAULT 0,
			total_load DOUBLE PRECISION DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS stops (
			id SERIAL PRIMARY KEY,
			route_id INTEGER REFERENCES routes(id) ON DELETE CASCADE,
			customer_id INTEGER REFERENCES customers(id) ON DELETE SET NULL,
			sequence INTEGER NOT NULL,
			quantity DOUBLE PRECISION DEFAULT 0,
			arrival_time VARCHAR(10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_routes_plan_id ON routes(plan_id)`,
		`CREATE INDEX IF NOT EXISTS idx_stops_route_id ON stops(route_id)`,
		`CREATE INDEX IF NOT EXISTS idx_vehicles_warehouse_id ON vehicles(warehouse_id)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_warehouse_id ON plans(warehouse_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
