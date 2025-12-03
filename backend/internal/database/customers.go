package database

import (
	"database/sql"

	"LogiTrackPro/backend/internal/models"
)

func ListCustomers(db *sql.DB) ([]models.Customer, error) {
	query := `SELECT id, name, address, latitude, longitude, demand_rate, 
			  max_inventory, current_inventory, min_inventory, holding_cost, priority,
			  created_at, updated_at FROM customers ORDER BY name`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude,
			&c.DemandRate, &c.MaxInventory, &c.CurrentInventory, &c.MinInventory,
			&c.HoldingCost, &c.Priority, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func GetCustomer(db *sql.DB, id int64) (*models.Customer, error) {
	c := &models.Customer{}
	query := `SELECT id, name, address, latitude, longitude, demand_rate, 
			  max_inventory, current_inventory, min_inventory, holding_cost, priority,
			  created_at, updated_at FROM customers WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude,
		&c.DemandRate, &c.MaxInventory, &c.CurrentInventory, &c.MinInventory,
		&c.HoldingCost, &c.Priority, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func CreateCustomer(db *sql.DB, c *models.Customer) error {
	query := `INSERT INTO customers (name, address, latitude, longitude, demand_rate, 
			  max_inventory, current_inventory, min_inventory, holding_cost, priority) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
			  RETURNING id, created_at, updated_at`
	
	return db.QueryRow(query, c.Name, c.Address, c.Latitude, c.Longitude,
		c.DemandRate, c.MaxInventory, c.CurrentInventory, c.MinInventory,
		c.HoldingCost, c.Priority).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func UpdateCustomer(db *sql.DB, c *models.Customer) error {
	query := `UPDATE customers SET name = $1, address = $2, latitude = $3, 
			  longitude = $4, demand_rate = $5, max_inventory = $6, 
			  current_inventory = $7, min_inventory = $8, holding_cost = $9, 
			  priority = $10, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $11`
	
	result, err := db.Exec(query, c.Name, c.Address, c.Latitude, c.Longitude,
		c.DemandRate, c.MaxInventory, c.CurrentInventory, c.MinInventory,
		c.HoldingCost, c.Priority, c.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteCustomer(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func CountCustomers(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	return count, err
}

