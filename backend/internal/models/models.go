package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
	Password  string    `gorm:"column:password_hash;not null;type:varchar(255)" json:"-"`
	Name      string    `gorm:"not null;type:varchar(255)" json:"name"`
	Role      string    `gorm:"type:varchar(50);default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Warehouse represents a warehouse/distribution center
type Warehouse struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"not null;type:varchar(255)" json:"name"`
	Address         string    `gorm:"type:text" json:"address"`
	Latitude        float64   `gorm:"not null;type:double precision" json:"latitude"`
	Longitude       float64   `gorm:"not null;type:double precision" json:"longitude"`
	Capacity        float64   `gorm:"type:double precision;default:0" json:"capacity"`
	CurrentStock    float64   `gorm:"column:current_stock;type:double precision;default:0" json:"current_stock"`
	HoldingCost     float64   `gorm:"column:holding_cost;type:double precision;default:0" json:"holding_cost"`
	ReplenishmentQty float64  `gorm:"column:replenishment_qty;type:double precision;default:0" json:"replenishment_qty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Vehicles        []Vehicle `gorm:"foreignKey:WarehouseID" json:"vehicles,omitempty"`
	Plans           []Plan    `gorm:"foreignKey:WarehouseID" json:"plans,omitempty"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

// Customer represents a customer location
type Customer struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"not null;type:varchar(255)" json:"name"`
	Address          string    `gorm:"type:text" json:"address"`
	Latitude         float64   `gorm:"not null;type:double precision" json:"latitude"`
	Longitude        float64   `gorm:"not null;type:double precision" json:"longitude"`
	DemandRate       float64   `gorm:"column:demand_rate;type:double precision;default:0" json:"demand_rate"`
	MaxInventory     float64   `gorm:"column:max_inventory;type:double precision;default:0" json:"max_inventory"`
	CurrentInventory float64   `gorm:"column:current_inventory;type:double precision;default:0" json:"current_inventory"`
	MinInventory     float64   `gorm:"column:min_inventory;type:double precision;default:0" json:"min_inventory"`
	HoldingCost      float64   `gorm:"column:holding_cost;type:double precision;default:0" json:"holding_cost"`
	Priority         int       `gorm:"type:integer;default:1" json:"priority"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Stops            []Stop    `gorm:"foreignKey:CustomerID" json:"stops,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

// Vehicle represents a delivery vehicle
type Vehicle struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;type:varchar(255)" json:"name"`
	Capacity    float64   `gorm:"not null;type:double precision" json:"capacity"`
	CostPerKm   float64   `gorm:"column:cost_per_km;type:double precision;default:0" json:"cost_per_km"`
	FixedCost   float64   `gorm:"column:fixed_cost;type:double precision;default:0" json:"fixed_cost"`
	MaxDistance float64   `gorm:"column:max_distance;type:double precision;default:0" json:"max_distance"`
	Available   bool      `gorm:"type:boolean;default:true" json:"available"`
	WarehouseID *int64    `gorm:"index;type:integer" json:"warehouse_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Warehouse   *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Routes      []Route    `gorm:"foreignKey:VehicleID" json:"routes,omitempty"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

// Plan represents a delivery plan
type Plan struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null;type:varchar(255)" json:"name"`
	StartDate     time.Time `gorm:"column:start_date;type:date;not null" json:"start_date"`
	EndDate       time.Time `gorm:"column:end_date;type:date;not null" json:"end_date"`
	Status        string    `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, optimizing, optimized, executed
	TotalCost     float64   `gorm:"column:total_cost;type:double precision;default:0" json:"total_cost"`
	TotalDistance float64   `gorm:"column:total_distance;type:double precision;default:0" json:"total_distance"`
	WarehouseID   *int64    `gorm:"index;type:integer" json:"warehouse_id"`
	CreatedBy     *int64    `gorm:"index;type:integer" json:"created_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Warehouse     *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	User          *User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	Routes        []Route     `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"routes,omitempty"`
}

func (Plan) TableName() string {
	return "plans"
}

// Route represents a delivery route for a specific day
type Route struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	PlanID        int64     `gorm:"index;not null;type:integer" json:"plan_id"`
	VehicleID     *int64    `gorm:"index;type:integer" json:"vehicle_id"`
	Day           int       `gorm:"not null;type:integer" json:"day"`
	Date          time.Time `gorm:"type:date;not null" json:"date"`
	TotalDistance float64   `gorm:"column:total_distance;type:double precision;default:0" json:"total_distance"`
	TotalCost     float64   `gorm:"column:total_cost;type:double precision;default:0" json:"total_cost"`
	TotalLoad     float64   `gorm:"column:total_load;type:double precision;default:0" json:"total_load"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	Plan          *Plan     `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Vehicle       *Vehicle  `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Stops         []Stop    `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE" json:"stops,omitempty"`
}

func (Route) TableName() string {
	return "routes"
}

// Stop represents a stop on a route
type Stop struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	RouteID     int64     `gorm:"index;not null;type:integer" json:"route_id"`
	CustomerID  *int64    `gorm:"index;type:integer" json:"customer_id"`
	Sequence    int       `gorm:"not null;type:integer" json:"sequence"`
	Quantity    float64   `gorm:"type:double precision;default:0" json:"quantity"`
	ArrivalTime string    `gorm:"type:varchar(10)" json:"arrival_time"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	Route       *Route    `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	Customer    *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Stop) TableName() string {
	return "stops"
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

