package database

import (
	"database/sql"

	"LogiTrackPro/backend/internal/models"
)

func ListWarehouses(db *sql.DB) ([]models.Warehouse, error) {
	query := `SELECT id, name, address, latitude, longitude, capacity, 
			  current_stock, holding_cost, replenishment_qty, created_at, updated_at 
			  FROM warehouses ORDER BY name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		err := rows.Scan(
			&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
			&w.Capacity, &w.CurrentStock, &w.HoldingCost, &w.ReplenishmentQty,
			&w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}
	return warehouses, nil
}

func GetWarehouse(db *sql.DB, id int64) (*models.Warehouse, error) {
	w := &models.Warehouse{}
	query := `SELECT id, name, address, latitude, longitude, capacity, 
			  current_stock, holding_cost, replenishment_qty, created_at, updated_at 
			  FROM warehouses WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&w.ID, &w.Name, &w.Address, &w.Latitude, &w.Longitude,
		&w.Capacity, &w.CurrentStock, &w.HoldingCost, &w.ReplenishmentQty,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return w, nil
}

func CreateWarehouse(db *sql.DB, w *models.Warehouse) error {
	query := `INSERT INTO warehouses (name, address, latitude, longitude, capacity, 
			  current_stock, holding_cost, replenishment_qty) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id, created_at, updated_at`

	return db.QueryRow(query, w.Name, w.Address, w.Latitude, w.Longitude,
		w.Capacity, w.CurrentStock, w.HoldingCost, w.ReplenishmentQty).Scan(
		&w.ID, &w.CreatedAt, &w.UpdatedAt,
	)
}

func UpdateWarehouse(db *sql.DB, w *models.Warehouse) error {
	query := `UPDATE warehouses SET name = $1, address = $2, latitude = $3, 
			  longitude = $4, capacity = $5, current_stock = $6, holding_cost = $7, 
			  replenishment_qty = $8, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $9 RETURNING updated_at`

	result, err := db.Exec(query, w.Name, w.Address, w.Latitude, w.Longitude,
		w.Capacity, w.CurrentStock, w.HoldingCost, w.ReplenishmentQty, w.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteWarehouse(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM warehouses WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func CountWarehouses(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM warehouses").Scan(&count)
	return count, err
}
