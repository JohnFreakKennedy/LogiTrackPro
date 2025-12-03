package database

import (
	"database/sql"

	"LogiTrackPro/backend/internal/models"
)

func ListVehicles(db *sql.DB) ([]models.Vehicle, error) {
	query := `SELECT id, name, capacity, cost_per_km, fixed_cost, max_distance, 
			  available, COALESCE(warehouse_id, 0), created_at, updated_at 
			  FROM vehicles ORDER BY name`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.ID, &v.Name, &v.Capacity, &v.CostPerKm, &v.FixedCost,
			&v.MaxDistance, &v.Available, &v.WarehouseID,
			&v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func ListAvailableVehiclesByWarehouse(db *sql.DB, warehouseID int64) ([]models.Vehicle, error) {
	query := `SELECT id, name, capacity, cost_per_km, fixed_cost, max_distance, 
			  available, warehouse_id, created_at, updated_at 
			  FROM vehicles WHERE warehouse_id = $1 AND available = true ORDER BY name`
	
	rows, err := db.Query(query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.ID, &v.Name, &v.Capacity, &v.CostPerKm, &v.FixedCost,
			&v.MaxDistance, &v.Available, &v.WarehouseID,
			&v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func GetVehicle(db *sql.DB, id int64) (*models.Vehicle, error) {
	v := &models.Vehicle{}
	query := `SELECT id, name, capacity, cost_per_km, fixed_cost, max_distance, 
			  available, COALESCE(warehouse_id, 0), created_at, updated_at 
			  FROM vehicles WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&v.ID, &v.Name, &v.Capacity, &v.CostPerKm, &v.FixedCost,
		&v.MaxDistance, &v.Available, &v.WarehouseID,
		&v.CreatedAt, &v.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

func CreateVehicle(db *sql.DB, v *models.Vehicle) error {
	var warehouseID interface{} = nil
	if v.WarehouseID > 0 {
		warehouseID = v.WarehouseID
	}
	
	query := `INSERT INTO vehicles (name, capacity, cost_per_km, fixed_cost, 
			  max_distance, available, warehouse_id) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING id, created_at, updated_at`
	
	return db.QueryRow(query, v.Name, v.Capacity, v.CostPerKm, v.FixedCost,
		v.MaxDistance, v.Available, warehouseID).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func UpdateVehicle(db *sql.DB, v *models.Vehicle) error {
	var warehouseID interface{} = nil
	if v.WarehouseID > 0 {
		warehouseID = v.WarehouseID
	}
	
	query := `UPDATE vehicles SET name = $1, capacity = $2, cost_per_km = $3, 
			  fixed_cost = $4, max_distance = $5, available = $6, warehouse_id = $7, 
			  updated_at = CURRENT_TIMESTAMP WHERE id = $8`
	
	result, err := db.Exec(query, v.Name, v.Capacity, v.CostPerKm, v.FixedCost,
		v.MaxDistance, v.Available, warehouseID, v.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteVehicle(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM vehicles WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func CountVehicles(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM vehicles").Scan(&count)
	return count, err
}

