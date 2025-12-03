package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Warehouse represents a warehouse/distribution center
type Warehouse struct {
	ID              int64   `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	Address         string  `json:"address" db:"address"`
	Latitude        float64 `json:"latitude" db:"latitude"`
	Longitude       float64 `json:"longitude" db:"longitude"`
	Capacity        float64 `json:"capacity" db:"capacity"`
	CurrentStock    float64 `json:"current_stock" db:"current_stock"`
	HoldingCost     float64 `json:"holding_cost" db:"holding_cost"`
	ReplenishmentQty float64 `json:"replenishment_qty" db:"replenishment_qty"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Customer represents a customer location
type Customer struct {
	ID               int64   `json:"id" db:"id"`
	Name             string  `json:"name" db:"name"`
	Address          string  `json:"address" db:"address"`
	Latitude         float64 `json:"latitude" db:"latitude"`
	Longitude        float64 `json:"longitude" db:"longitude"`
	DemandRate       float64 `json:"demand_rate" db:"demand_rate"`
	MaxInventory     float64 `json:"max_inventory" db:"max_inventory"`
	CurrentInventory float64 `json:"current_inventory" db:"current_inventory"`
	MinInventory     float64 `json:"min_inventory" db:"min_inventory"`
	HoldingCost      float64 `json:"holding_cost" db:"holding_cost"`
	Priority         int     `json:"priority" db:"priority"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Vehicle represents a delivery vehicle
type Vehicle struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Capacity    float64 `json:"capacity" db:"capacity"`
	CostPerKm   float64 `json:"cost_per_km" db:"cost_per_km"`
	FixedCost   float64 `json:"fixed_cost" db:"fixed_cost"`
	MaxDistance float64 `json:"max_distance" db:"max_distance"`
	Available   bool    `json:"available" db:"available"`
	WarehouseID int64   `json:"warehouse_id" db:"warehouse_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Plan represents a delivery plan
type Plan struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Status       string    `json:"status" db:"status"` // draft, optimizing, optimized, executed
	TotalCost    float64   `json:"total_cost" db:"total_cost"`
	TotalDistance float64  `json:"total_distance" db:"total_distance"`
	WarehouseID  int64     `json:"warehouse_id" db:"warehouse_id"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Routes       []Route   `json:"routes,omitempty"`
}

// Route represents a delivery route for a specific day
type Route struct {
	ID           int64     `json:"id" db:"id"`
	PlanID       int64     `json:"plan_id" db:"plan_id"`
	VehicleID    *int64    `json:"vehicle_id" db:"vehicle_id"`
	Day          int       `json:"day" db:"day"`
	Date         time.Time `json:"date" db:"date"`
	TotalDistance float64  `json:"total_distance" db:"total_distance"`
	TotalCost    float64   `json:"total_cost" db:"total_cost"`
	TotalLoad    float64   `json:"total_load" db:"total_load"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Stops        []Stop    `json:"stops,omitempty"`
	Vehicle      *Vehicle  `json:"vehicle,omitempty"`
}

// Stop represents a stop on a route
type Stop struct {
	ID          int64   `json:"id" db:"id"`
	RouteID     int64   `json:"route_id" db:"route_id"`
	CustomerID  int64   `json:"customer_id" db:"customer_id"`
	Sequence    int     `json:"sequence" db:"sequence"`
	Quantity    float64 `json:"quantity" db:"quantity"`
	ArrivalTime string  `json:"arrival_time" db:"arrival_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Customer    *Customer `json:"customer,omitempty"`
}

// Dashboard represents analytics dashboard data
type Dashboard struct {
	TotalWarehouses    int     `json:"total_warehouses"`
	TotalCustomers     int     `json:"total_customers"`
	TotalVehicles      int     `json:"total_vehicles"`
	ActivePlans        int     `json:"active_plans"`
	TotalDeliveries    int     `json:"total_deliveries"`
	TotalDistanceKm    float64 `json:"total_distance_km"`
	TotalCost          float64 `json:"total_cost"`
	AvgUtilization     float64 `json:"avg_utilization"`
	RecentPlans        []Plan  `json:"recent_plans"`
}

