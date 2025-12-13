package database

import (
	"database/sql"

	"LogiTrackPro/backend/internal/models"
)

func GetRoutesByPlan(db *sql.DB, planID int64) ([]models.Route, error) {
	query := `SELECT r.id, r.plan_id, r.vehicle_id, r.day, r.date, 
			  r.total_distance, r.total_cost, r.total_load, r.created_at,
			  v.id, v.name, v.capacity, v.cost_per_km, v.fixed_cost, v.max_distance, 
			  v.available, COALESCE(v.warehouse_id, 0), v.created_at, v.updated_at
			  FROM routes r
			  LEFT JOIN vehicles v ON r.vehicle_id = v.id
			  WHERE r.plan_id = $1 ORDER BY r.day, r.id`

	rows, err := db.Query(query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var r models.Route
		var vehicleID sql.NullInt64
		var vID sql.NullInt64
		var vName sql.NullString
		var vCapacity sql.NullFloat64
		var vCostPerKm sql.NullFloat64
		var vFixedCost sql.NullFloat64
		var vMaxDistance sql.NullFloat64
		var vAvailable sql.NullBool
		var vWarehouseID sql.NullInt64
		var vCreatedAt sql.NullTime
		var vUpdatedAt sql.NullTime

		err := rows.Scan(
			&r.ID, &r.PlanID, &vehicleID, &r.Day, &r.Date,
			&r.TotalDistance, &r.TotalCost, &r.TotalLoad, &r.CreatedAt,
			&vID, &vName, &vCapacity, &vCostPerKm, &vFixedCost, &vMaxDistance,
			&vAvailable, &vWarehouseID, &vCreatedAt, &vUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set VehicleID pointer
		if vehicleID.Valid {
			vIDVal := vehicleID.Int64
			r.VehicleID = &vIDVal
		}

		// Set Vehicle if all vehicle fields are valid
		if vID.Valid && vName.Valid {
			r.Vehicle = &models.Vehicle{
				ID:          vID.Int64,
				Name:        vName.String,
				Capacity:    vCapacity.Float64,
				CostPerKm:   vCostPerKm.Float64,
				FixedCost:   vFixedCost.Float64,
				MaxDistance: vMaxDistance.Float64,
				Available:   vAvailable.Bool,
				WarehouseID: vWarehouseID.Int64,
				CreatedAt:   vCreatedAt.Time,
				UpdatedAt:   vUpdatedAt.Time,
			}
		}

		routes = append(routes, r)
	}

	// Load stops for each route
	for i := range routes {
		stops, err := GetStopsByRoute(db, routes[i].ID)
		if err != nil {
			return nil, err
		}
		routes[i].Stops = stops
	}

	return routes, nil
}

func CreateRoute(db *sql.DB, r *models.Route) error {
	query := `INSERT INTO routes (plan_id, vehicle_id, day, date, total_distance, total_cost, total_load) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING id, created_at`

	return db.QueryRow(query, r.PlanID, r.VehicleID, r.Day, r.Date,
		r.TotalDistance, r.TotalCost, r.TotalLoad).Scan(&r.ID, &r.CreatedAt)
}

func CreateRouteTx(tx *sql.Tx, r *models.Route) error {
	query := `INSERT INTO routes (plan_id, vehicle_id, day, date, total_distance, total_cost, total_load) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING id, created_at`

	return tx.QueryRow(query, r.PlanID, r.VehicleID, r.Day, r.Date,
		r.TotalDistance, r.TotalCost, r.TotalLoad).Scan(&r.ID, &r.CreatedAt)
}

func DeleteRoutesByPlan(db *sql.DB, planID int64) error {
	_, err := db.Exec("DELETE FROM routes WHERE plan_id = $1", planID)
	return err
}

func DeleteRoutesByPlanTx(tx *sql.Tx, planID int64) error {
	_, err := tx.Exec("DELETE FROM routes WHERE plan_id = $1", planID)
	return err
}

func GetStopsByRoute(db *sql.DB, routeID int64) ([]models.Stop, error) {
	query := `SELECT s.id, s.route_id, s.customer_id, s.sequence, s.quantity, 
			  COALESCE(s.arrival_time, ''), s.created_at,
			  c.id, c.name, c.address, c.latitude, c.longitude
			  FROM stops s
			  LEFT JOIN customers c ON s.customer_id = c.id
			  WHERE s.route_id = $1 ORDER BY s.sequence`

	rows, err := db.Query(query, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stops []models.Stop
	for rows.Next() {
		var s models.Stop
		var customerID sql.NullInt64
		var cID sql.NullInt64
		var cName sql.NullString
		var cAddress sql.NullString
		var cLatitude sql.NullFloat64
		var cLongitude sql.NullFloat64

		err := rows.Scan(
			&s.ID, &s.RouteID, &customerID, &s.Sequence, &s.Quantity,
			&s.ArrivalTime, &s.CreatedAt,
			&cID, &cName, &cAddress, &cLatitude, &cLongitude,
		)
		if err != nil {
			return nil, err
		}

		if customerID.Valid {
			s.CustomerID = customerID.Int64
		}

		// Set Customer if all customer fields are valid
		if cID.Valid && cName.Valid {
			s.Customer = &models.Customer{
				ID:        cID.Int64,
				Name:      cName.String,
				Address:   cAddress.String,
				Latitude:  cLatitude.Float64,
				Longitude: cLongitude.Float64,
			}
		}

		stops = append(stops, s)
	}
	return stops, nil
}

func CreateStop(db *sql.DB, s *models.Stop) error {
	var customerID interface{} = nil
	if s.CustomerID > 0 {
		customerID = s.CustomerID
	}

	query := `INSERT INTO stops (route_id, customer_id, sequence, quantity, arrival_time) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, created_at`

	return db.QueryRow(query, s.RouteID, customerID, s.Sequence, s.Quantity,
		s.ArrivalTime).Scan(&s.ID, &s.CreatedAt)
}

func CreateStopTx(tx *sql.Tx, s *models.Stop) error {
	var customerID interface{} = nil
	if s.CustomerID > 0 {
		customerID = s.CustomerID
	}

	query := `INSERT INTO stops (route_id, customer_id, sequence, quantity, arrival_time) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, created_at`

	return tx.QueryRow(query, s.RouteID, customerID, s.Sequence, s.Quantity,
		s.ArrivalTime).Scan(&s.ID, &s.CreatedAt)
}

func CountTotalDeliveries(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM stops").Scan(&count)
	return count, err
}

func GetTotalDistanceAndCost(db *sql.DB) (float64, float64, error) {
	var distance, cost float64
	err := db.QueryRow("SELECT COALESCE(SUM(total_distance), 0), COALESCE(SUM(total_cost), 0) FROM routes").Scan(&distance, &cost)
	return distance, cost, err
}
