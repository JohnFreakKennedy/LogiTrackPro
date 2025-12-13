package database

import (
	"database/sql"

	"LogiTrackPro/backend/internal/models"
)

func ListPlans(db *sql.DB) ([]models.Plan, error) {
	query := `SELECT id, name, start_date, end_date, status, total_cost, 
			  total_distance, COALESCE(warehouse_id, 0), COALESCE(created_by, 0), 
			  created_at, updated_at FROM plans ORDER BY created_at DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var p models.Plan
		err := rows.Scan(
			&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.Status,
			&p.TotalCost, &p.TotalDistance, &p.WarehouseID, &p.CreatedBy,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func GetPlan(db *sql.DB, id int64) (*models.Plan, error) {
	p := &models.Plan{}
	query := `SELECT id, name, start_date, end_date, status, total_cost, 
			  total_distance, COALESCE(warehouse_id, 0), COALESCE(created_by, 0), 
			  created_at, updated_at FROM plans WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.Status,
		&p.TotalCost, &p.TotalDistance, &p.WarehouseID, &p.CreatedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func CreatePlan(db *sql.DB, p *models.Plan) error {
	var warehouseID, createdBy interface{} = nil, nil
	if p.WarehouseID > 0 {
		warehouseID = p.WarehouseID
	}
	if p.CreatedBy > 0 {
		createdBy = p.CreatedBy
	}
	
	query := `INSERT INTO plans (name, start_date, end_date, status, warehouse_id, created_by) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, created_at, updated_at`
	
	return db.QueryRow(query, p.Name, p.StartDate, p.EndDate, p.Status,
		warehouseID, createdBy).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func UpdatePlanStatus(db *sql.DB, id int64, status string, totalCost, totalDistance float64) error {
	query := `UPDATE plans SET status = $1, total_cost = $2, total_distance = $3, 
			  updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	
	result, err := db.Exec(query, status, totalCost, totalDistance, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func UpdatePlanStatusTx(tx *sql.Tx, id int64, status string, totalCost, totalDistance float64) error {
	query := `UPDATE plans SET status = $1, total_cost = $2, total_distance = $3, 
			  updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	
	result, err := tx.Exec(query, status, totalCost, totalDistance, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func DeletePlan(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM plans WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func CountActivePlans(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM plans WHERE status IN ('draft', 'optimizing', 'optimized')").Scan(&count)
	return count, err
}

func GetRecentPlans(db *sql.DB, limit int) ([]models.Plan, error) {
	query := `SELECT id, name, start_date, end_date, status, total_cost, 
			  total_distance, COALESCE(warehouse_id, 0), COALESCE(created_by, 0), 
			  created_at, updated_at FROM plans ORDER BY created_at DESC LIMIT $1`
	
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var p models.Plan
		err := rows.Scan(
			&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.Status,
			&p.TotalCost, &p.TotalDistance, &p.WarehouseID, &p.CreatedBy,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

