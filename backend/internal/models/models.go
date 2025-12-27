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
	ID                 int64               `gorm:"primaryKey" json:"id"`
	Name               string              `gorm:"not null;type:varchar(255)" json:"name"`
	Address            string              `gorm:"type:text" json:"address"`
	Latitude           float64             `gorm:"not null;type:double precision" json:"latitude"`
	Longitude          float64             `gorm:"not null;type:double precision" json:"longitude"`
	Capacity           float64             `gorm:"type:double precision;default:0" json:"capacity"`
	CurrentStock       float64             `gorm:"column:current_stock;type:double precision;default:0" json:"current_stock"`
	HoldingCost        float64             `gorm:"column:holding_cost;type:double precision;default:0" json:"holding_cost"`
	ReplenishmentQty   float64             `gorm:"column:replenishment_qty;type:double precision;default:0" json:"replenishment_qty"`
	CreatedAt          time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	Vehicles           []Vehicle           `gorm:"foreignKey:WarehouseID" json:"vehicles,omitempty"`
	Plans              []Plan              `gorm:"foreignKey:WarehouseID" json:"plans,omitempty"`
	InventorySnapshots []InventorySnapshot `gorm:"foreignKey:EntityID" json:"inventory_snapshots,omitempty"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

// Customer represents a customer location
type Customer struct {
	ID                 int64                      `gorm:"primaryKey" json:"id"`
	Name               string                     `gorm:"not null;type:varchar(255)" json:"name"`
	Address            string                     `gorm:"type:text" json:"address"`
	Latitude           float64                    `gorm:"not null;type:double precision" json:"latitude"`
	Longitude          float64                    `gorm:"not null;type:double precision" json:"longitude"`
	DemandRate         float64                    `gorm:"column:demand_rate;type:double precision;default:0" json:"demand_rate"`
	MaxInventory       float64                    `gorm:"column:max_inventory;type:double precision;default:0" json:"max_inventory"`
	CurrentInventory   float64                    `gorm:"column:current_inventory;type:double precision;default:0" json:"current_inventory"`
	MinInventory       float64                    `gorm:"column:min_inventory;type:double precision;default:0" json:"min_inventory"`
	HoldingCost        float64                    `gorm:"column:holding_cost;type:double precision;default:0" json:"holding_cost"`
	Priority           int                        `gorm:"type:integer;default:1" json:"priority"`
	CreatedAt          time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
	Stops              []Stop                     `gorm:"foreignKey:CustomerID" json:"stops,omitempty"`
	InventorySnapshots []InventorySnapshot        `gorm:"foreignKey:EntityID" json:"inventory_snapshots,omitempty"`
	ProductInventory   []CustomerProductInventory `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"product_inventory,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

// Vehicle represents a delivery vehicle
type Vehicle struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null;type:varchar(255)" json:"name"`
	Capacity    float64    `gorm:"not null;type:double precision" json:"capacity"`
	CostPerKm   float64    `gorm:"column:cost_per_km;type:double precision;default:0" json:"cost_per_km"`
	FixedCost   float64    `gorm:"column:fixed_cost;type:double precision;default:0" json:"fixed_cost"`
	MaxDistance float64    `gorm:"column:max_distance;type:double precision;default:0" json:"max_distance"`
	Available   bool       `gorm:"type:boolean;default:true" json:"available"`
	WarehouseID *int64     `gorm:"index;type:integer" json:"warehouse_id"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Warehouse   *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Routes      []Route    `gorm:"foreignKey:VehicleID" json:"routes,omitempty"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

// Plan represents a delivery plan
type Plan struct {
	ID                 int64               `gorm:"primaryKey" json:"id"`
	Name               string              `gorm:"not null;type:varchar(255)" json:"name"`
	StartDate          time.Time           `gorm:"column:start_date;type:date;not null" json:"start_date"`
	EndDate            time.Time           `gorm:"column:end_date;type:date;not null" json:"end_date"`
	Status             string              `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, optimizing, optimized, executed
	TotalCost          float64             `gorm:"column:total_cost;type:double precision;default:0" json:"total_cost"`
	TotalDistance      float64             `gorm:"column:total_distance;type:double precision;default:0" json:"total_distance"`
	WarehouseID        *int64              `gorm:"index;type:integer" json:"warehouse_id"`
	CreatedBy          *int64              `gorm:"index;type:integer" json:"created_by"`
	CreatedAt          time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	Warehouse          *Warehouse          `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	User               *User               `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	Routes             []Route             `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"routes,omitempty"`
	Executions         []RouteExecution    `gorm:"foreignKey:RouteID" json:"executions,omitempty"`
	InventorySnapshots []InventorySnapshot `gorm:"foreignKey:PlanID" json:"inventory_snapshots,omitempty"`
}

func (Plan) TableName() string {
	return "plans"
}

// Route represents a delivery route for a specific day
type Route struct {
	ID            int64            `gorm:"primaryKey" json:"id"`
	PlanID        int64            `gorm:"index;not null;type:integer" json:"plan_id"`
	VehicleID     *int64           `gorm:"index;type:integer" json:"vehicle_id"`
	Day           int              `gorm:"not null;type:integer" json:"day"`
	Date          time.Time        `gorm:"type:date;not null" json:"date"`
	TotalDistance float64          `gorm:"column:total_distance;type:double precision;default:0" json:"total_distance"`
	TotalCost     float64          `gorm:"column:total_cost;type:double precision;default:0" json:"total_cost"`
	TotalLoad     float64          `gorm:"column:total_load;type:double precision;default:0" json:"total_load"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"created_at"`
	Plan          *Plan            `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Vehicle       *Vehicle         `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Stops         []Stop           `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE" json:"stops,omitempty"`
	Executions    []RouteExecution `gorm:"foreignKey:RouteID" json:"executions,omitempty"`
}

func (Route) TableName() string {
	return "routes"
}

// Stop represents a stop on a route
type Stop struct {
	ID                int64                 `gorm:"primaryKey" json:"id"`
	RouteID           int64                 `gorm:"index;not null;type:integer" json:"route_id"`
	CustomerID        *int64                `gorm:"index;type:integer" json:"customer_id"`
	Sequence          int                   `gorm:"not null;type:integer" json:"sequence"`
	Quantity          float64               `gorm:"type:double precision;default:0" json:"quantity"`
	ArrivalTime       string                `gorm:"type:varchar(10)" json:"arrival_time"`
	CreatedAt         time.Time             `gorm:"autoCreateTime" json:"created_at"`
	Route             *Route                `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	Customer          *Customer             `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	StopExecutions    []StopExecution       `gorm:"foreignKey:StopID" json:"stop_executions,omitempty"`
	ProductQuantities []StopProductQuantity `gorm:"foreignKey:StopID;constraint:OnDelete:CASCADE" json:"product_quantities,omitempty"`
}

func (Stop) TableName() string {
	return "stops"
}

// RouteExecution represents the actual execution of a planned route
type RouteExecution struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	RouteID          int64           `gorm:"index;not null;type:integer" json:"route_id"`
	Status           string          `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, in_progress, completed, cancelled
	PlannedDistance  float64         `gorm:"column:planned_distance;type:double precision;default:0" json:"planned_distance"`
	ActualDistance   float64         `gorm:"column:actual_distance;type:double precision;default:0" json:"actual_distance"`
	PlannedCost      float64         `gorm:"column:planned_cost;type:double precision;default:0" json:"planned_cost"`
	ActualCost       float64         `gorm:"column:actual_cost;type:double precision;default:0" json:"actual_cost"`
	PlannedLoad      float64         `gorm:"column:planned_load;type:double precision;default:0" json:"planned_load"`
	ActualLoad       float64         `gorm:"column:actual_load;type:double precision;default:0" json:"actual_load"`
	PlannedStartTime *time.Time      `gorm:"type:timestamp" json:"planned_start_time"`
	ActualStartTime  *time.Time      `gorm:"type:timestamp" json:"actual_start_time"`
	PlannedEndTime   *time.Time      `gorm:"type:timestamp" json:"planned_end_time"`
	ActualEndTime    *time.Time      `gorm:"type:timestamp" json:"actual_end_time"`
	DriverNotes      string          `gorm:"type:text" json:"driver_notes"`
	DeviationReason  string          `gorm:"type:text" json:"deviation_reason"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	Route            *Route          `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	StopExecutions   []StopExecution `gorm:"foreignKey:RouteExecutionID;constraint:OnDelete:CASCADE" json:"stop_executions,omitempty"`
}

func (RouteExecution) TableName() string {
	return "route_executions"
}

// StopExecution represents the actual execution of a planned stop
type StopExecution struct {
	ID                   int64           `gorm:"primaryKey" json:"id"`
	RouteExecutionID     int64           `gorm:"index;not null;type:integer" json:"route_execution_id"`
	StopID               int64           `gorm:"index;not null;type:integer" json:"stop_id"`
	Status               string          `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, in_progress, completed, skipped, failed
	PlannedQuantity      float64         `gorm:"column:planned_quantity;type:double precision;default:0" json:"planned_quantity"`
	ActualQuantity       float64         `gorm:"column:actual_quantity;type:double precision;default:0" json:"actual_quantity"`
	PlannedArrivalTime   *time.Time      `gorm:"type:timestamp" json:"planned_arrival_time"`
	ActualArrivalTime    *time.Time      `gorm:"type:timestamp" json:"actual_arrival_time"`
	PlannedDepartureTime *time.Time      `gorm:"type:timestamp" json:"planned_departure_time"`
	ActualDepartureTime  *time.Time      `gorm:"type:timestamp" json:"actual_departure_time"`
	ServiceDuration      int             `gorm:"type:integer;default:0" json:"service_duration"` // minutes
	Notes                string          `gorm:"type:text" json:"notes"`
	CreatedAt            time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	RouteExecution       *RouteExecution `gorm:"foreignKey:RouteExecutionID" json:"route_execution,omitempty"`
	Stop                 *Stop           `gorm:"foreignKey:StopID" json:"stop,omitempty"`
}

func (StopExecution) TableName() string {
	return "stop_executions"
}

// InventorySnapshot represents a historical snapshot of inventory levels
type InventorySnapshot struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	EntityType     string    `gorm:"type:varchar(20);not null" json:"entity_type"` // 'customer' or 'warehouse'
	EntityID       int64     `gorm:"index;not null;type:integer" json:"entity_id"`
	SnapshotDate   time.Time `gorm:"column:snapshot_date;type:date;not null" json:"snapshot_date"`
	SnapshotTime   time.Time `gorm:"column:snapshot_time;type:timestamp;not null" json:"snapshot_time"`
	InventoryLevel float64   `gorm:"column:inventory_level;type:double precision;not null" json:"inventory_level"`
	DemandRate     float64   `gorm:"column:demand_rate;type:double precision;default:0" json:"demand_rate"`
	MinInventory   float64   `gorm:"column:min_inventory;type:double precision;default:0" json:"min_inventory"`
	MaxInventory   float64   `gorm:"column:max_inventory;type:double precision;default:0" json:"max_inventory"`
	SnapshotReason string    `gorm:"type:varchar(50)" json:"snapshot_reason"` // daily, delivery, manual, optimization
	PlanID         *int64    `gorm:"index;type:integer" json:"plan_id"`
	RouteID        *int64    `gorm:"index;type:integer" json:"route_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	Plan           *Plan     `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Route          *Route    `gorm:"foreignKey:RouteID" json:"route,omitempty"`
}

func (InventorySnapshot) TableName() string {
	return "inventory_snapshots"
}

// Product represents a product type (optional multi-product support)
// If not used, system assumes single product
type Product struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;type:varchar(255)" json:"name"`
	SKU         string    `gorm:"uniqueIndex;type:varchar(100)" json:"sku"`
	Description string    `gorm:"type:text" json:"description"`
	Unit        string    `gorm:"type:varchar(50);default:'kg'" json:"unit"`     // kg, liters, units, etc.
	Weight      float64   `gorm:"type:double precision;default:0" json:"weight"` // per unit
	Volume      float64   `gorm:"type:double precision;default:0" json:"volume"` // per unit
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

// CustomerProductInventory represents product-specific inventory for customers (optional)
type CustomerProductInventory struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	CustomerID       int64     `gorm:"index;not null;type:integer" json:"customer_id"`
	ProductID        int64     `gorm:"index;not null;type:integer" json:"product_id"`
	CurrentInventory float64   `gorm:"column:current_inventory;type:double precision;default:0" json:"current_inventory"`
	MaxInventory     float64   `gorm:"column:max_inventory;type:double precision;default:0" json:"max_inventory"`
	MinInventory     float64   `gorm:"column:min_inventory;type:double precision;default:0" json:"min_inventory"`
	DemandRate       float64   `gorm:"column:demand_rate;type:double precision;default:0" json:"demand_rate"`
	HoldingCost      float64   `gorm:"column:holding_cost;type:double precision;default:0" json:"holding_cost"`
	Priority         int       `gorm:"type:integer;default:1" json:"priority"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Customer         *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Product          *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (CustomerProductInventory) TableName() string {
	return "customer_product_inventory"
}

// StopProductQuantity represents product-specific quantities in stops (optional)
type StopProductQuantity struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	StopID    int64     `gorm:"index;not null;type:integer" json:"stop_id"`
	ProductID int64     `gorm:"index;not null;type:integer" json:"product_id"`
	Quantity  float64   `gorm:"type:double precision;default:0" json:"quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Stop      *Stop     `gorm:"foreignKey:StopID" json:"stop,omitempty"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (StopProductQuantity) TableName() string {
	return "stop_product_quantities"
}

// Dashboard represents analytics dashboard data
type Dashboard struct {
	TotalWarehouses int     `json:"total_warehouses"`
	TotalCustomers  int     `json:"total_customers"`
	TotalVehicles   int     `json:"total_vehicles"`
	ActivePlans     int     `json:"active_plans"`
	TotalDeliveries int     `json:"total_deliveries"`
	TotalDistanceKm float64 `json:"total_distance_km"`
	TotalCost       float64 `json:"total_cost"`
	AvgUtilization  float64 `json:"avg_utilization"`
	RecentPlans     []Plan  `json:"recent_plans"`
}
