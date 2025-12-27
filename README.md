# LogiTrackPro

A comprehensive logistics planning platform designed to solve the **Inventory Routing Problem (IRP)**. The system manages warehouses, customers, vehicles, inventory levels, and generates optimized multi-day delivery plans.

## Architecture

LogiTrackPro is built as a **modular monolith backend with one external microservice**:

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  React Frontend │────▶│  Go Backend API  │────▶│   PostgreSQL    │
│     (Vite)      │     │    (Gin)         │     │    Database     │
└─────────────────┘     └────────┬─────────┘     └─────────────────┘
                                 │
                                 ▼
                        ┌──────────────────┐
                        │ Python Optimizer │
                        │   (FastAPI)      │
                        └──────────────────┘
```

## Features

- **User Authentication**: JWT-based auth with registration and login
- **Warehouse Management**: CRUD operations for distribution centers
- **Customer Management**: Track customer locations, demand rates, and inventory levels
- **Vehicle Fleet Management**: Manage vehicles with capacity and cost parameters
- **Delivery Planning**: Create multi-day delivery plans
- **Route Optimization**: IRP solver using Google OR-Tools with Guided Local Search metaheuristic
- **Analytics Dashboard**: Overview of logistics operations

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend API | Go 1.23+ with Gin framework |
| Database ORM | GORM v1.25+ (PostgreSQL driver) |
| Optimizer Service | Python 3.11+ with FastAPI |
| Frontend | React 18 with Vite, Tailwind CSS |
| Database | PostgreSQL 15 |
| Auth | JWT tokens |
| Containerization | Docker & Docker Compose |

## Project Structure

```
LogiTrackPro/
├── backend/                  # Go backend API
│   ├── cmd/api/             # Application entry point
│   └── internal/            # Internal packages
│       ├── config/          # Configuration management
│       ├── database/        # Database layer (GORM) & migrations
│       ├── handlers/        # HTTP request handlers
│       ├── models/          # Domain models (GORM models with relationships)
│       └── optimizer/       # Optimizer client
├── optimizer/               # Python optimization service
│   ├── main.py             # FastAPI application
│   ├── solver.py           # IRP solver implementation
│   └── requirements.txt    # Python dependencies
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # Reusable components
│   │   ├── pages/          # Page components
│   │   ├── api.js          # API client
│   │   └── App.jsx         # Main application
│   └── package.json
├── docker-compose.yml      # Container orchestration
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23+
- Python 3.11+
- Node.js 20+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Option 1: Docker Compose (Recommended)

```bash
# Clone and navigate to project
cd LogiTrackPro

# Start all services
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Optimizer: http://localhost:8000
```

### Option 2: Manual Setup

#### 1. Database Setup

```bash
# Create PostgreSQL database
createdb logitrackpro

# Or using psql
psql -c "CREATE DATABASE logitrackpro;"
```

#### 2. Backend Setup

```bash
cd backend

# Install dependencies (includes GORM and PostgreSQL driver)
go mod download

# Set environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/logitrackpro?sslmode=disable"
export OPTIMIZER_URL="http://localhost:8000"
export JWT_SECRET="your-secret-key"
export PORT="8080"

# Run the backend
# GORM AutoMigrate will automatically create/update database schema on startup
go run cmd/api/main.go
```

**Note**: The backend uses GORM for database operations. On first run, GORM will automatically create all required tables. Subsequent runs will only update the schema if models have changed.

#### 3. Optimizer Setup

```bash
cd optimizer

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Run the optimizer
uvicorn main:app --host 0.0.0.0 --port 8000
```

#### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Warehouses
- `GET /api/v1/warehouses` - List all warehouses
- `POST /api/v1/warehouses` - Create warehouse
- `GET /api/v1/warehouses/:id` - Get warehouse by ID
- `PUT /api/v1/warehouses/:id` - Update warehouse
- `DELETE /api/v1/warehouses/:id` - Delete warehouse

### Customers
- `GET /api/v1/customers` - List all customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/:id` - Get customer by ID
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

### Vehicles
- `GET /api/v1/vehicles` - List all vehicles
- `POST /api/v1/vehicles` - Create vehicle
- `GET /api/v1/vehicles/:id` - Get vehicle by ID
- `PUT /api/v1/vehicles/:id` - Update vehicle
- `DELETE /api/v1/vehicles/:id` - Delete vehicle

### Plans
- `GET /api/v1/plans` - List all plans
- `POST /api/v1/plans` - Create plan
- `GET /api/v1/plans/:id` - Get plan by ID
- `DELETE /api/v1/plans/:id` - Delete plan
- `POST /api/v1/plans/:id/optimize` - Run optimization
- `GET /api/v1/plans/:id/routes` - Get plan routes

### Analytics
- `GET /api/v1/analytics/dashboard` - Get dashboard data
- `GET /api/v1/analytics/summary` - Get summary statistics

## Optimization Algorithm

### IRP vs VRP: Key Differences

**Vehicle Routing Problem (VRP)** is a classic optimization problem where:
- A fleet of vehicles must visit a set of customers
- Each customer has a known demand that must be satisfied
- The goal is to minimize total travel distance/cost
- All customers must be visited exactly once

**Inventory Routing Problem (IRP)** extends VRP with inventory management:
- Customers have **inventory levels** that change over time
- Customers have **demand rates** (consumption per day)
- Customers have **min/max inventory constraints**
- The problem spans **multiple days** (planning horizon)
- We must decide **when** to visit customers (not just route optimization)
- Goal: Minimize total cost while preventing stockouts

**Key IRP Challenges:**
- **Temporal dimension**: Inventory depletes over time
- **Delivery timing**: Must deliver before stockout, but not too early (wasteful)
- **Multi-day optimization**: Routes on day N affect inventory on day N+1
- **Inventory constraints**: Must respect min/max inventory levels

### OR-Tools Implementation

LogiTrackPro uses **Google OR-Tools** (version 9.9+) to solve IRP by decomposing it into a series of VRP subproblems, one per day.

#### Architecture

```
For each day in planning horizon:
  1. Project inventory levels forward
  2. Identify customers needing delivery (stockout risk)
  3. Solve VRP using OR-Tools for that day
  4. Update inventory levels after deliveries
  5. Consume daily demand
```

#### OR-Tools Configuration

**Routing Model Setup:**
- **RoutingIndexManager**: Manages node indices (warehouse = 0, customers = 1+)
- **RoutingModel**: Core VRP solver with constraint programming
- **Distance Matrix**: Pre-computed haversine distances (in meters, as integers)

**Constraints Implemented:**

1. **Capacity Constraint** (`AddDimensionWithVehicleCapacity`):
   - Each vehicle has a maximum capacity
   - Total deliveries on a route cannot exceed vehicle capacity
   - Delivery quantities calculated based on inventory gaps

2. **Distance Constraint** (`AddDimension`):
   - Maximum travel distance per vehicle (if specified)
   - Prevents routes exceeding vehicle range
   - Uses cumulative distance tracking

**Search Strategies:**

1. **First Solution Strategy**: `PATH_CHEAPEST_ARC`
   - Builds initial solution by iteratively adding cheapest arcs
   - Fast greedy construction heuristic
   - Provides good starting point for improvement

2. **Local Search Metaheuristic**: `GUIDED_LOCAL_SEARCH`
   - Advanced metaheuristic that guides local search
   - Uses penalty mechanism to escape local optima
   - Balances exploration vs exploitation
   - Typically finds solutions within 1-5% of optimal

**Time Limit:**
- 30 seconds per day's VRP problem
- Balances solution quality vs computation time
- Allows real-time optimization for practical use cases

#### Algorithm Flow

```python
# Day-by-day IRP solving
for day in planning_horizon:
    # Step 1: Inventory projection
    customers_needing_delivery = []
    for customer in customers:
        days_until_stockout = (inventory - min_inventory) / demand_rate
        if days_until_stockout <= 2 or inventory <= min_inventory:
            customers_needing_delivery.append(customer)
    
    # Step 2: Solve VRP using OR-Tools
    solution = routing.SolveWithParameters(search_parameters)
    
    # Step 3: Extract routes and update inventory
    for route in solution.routes:
        for stop in route.stops:
            inventory[stop.customer] += stop.delivery_qty
    
    # Step 4: Consume daily demand
    for customer in customers:
        inventory[customer] -= customer.demand_rate
```

#### Why OR-Tools?

**Advantages:**
- **Proven algorithms**: Industry-standard optimization library
- **Better solutions**: Typically 10-30% better than simple heuristics
- **Robust constraint handling**: Handles complex constraints efficiently
- **Active maintenance**: Regularly updated by Google
- **Well-documented**: Extensive documentation and examples
- **Scalability**: Handles problems with 100+ customers efficiently

**Compared to Custom Heuristics:**
- Custom nearest neighbor + 2-opt: Simple but limited
- OR-Tools: More sophisticated, better solutions, production-ready

#### Fallback Mechanism

If OR-Tools fails to find a solution (rare), the system falls back to:
- Simple nearest neighbor heuristic
- Ensures system always returns valid routes
- Maintains system reliability

### Algorithm Parameters

The solver considers:
- **Vehicle capacity constraints**: Maximum load per vehicle
- **Maximum distance constraints**: Vehicle range limitations
- **Customer priority levels**: Higher priority customers served first
- **Demand rates**: Daily consumption rates
- **Inventory levels**: Current and projected inventory
- **Delivery costs**: Fixed cost + per-km cost per vehicle
- **Min/max inventory**: Prevents over/under-stocking

### Performance Characteristics

- **Time Complexity**: O(n² × m × d) where n=customers, m=vehicles, d=days
- **Solution Quality**: Typically within 1-5% of optimal for VRP subproblems
- **Scalability**: Handles 50-200 customers per day efficiently
- **Real-time**: 30-second limit ensures responsive API

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Backend server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `OPTIMIZER_URL` | Optimizer service URL | `http://localhost:8000` |
| `JWT_SECRET` | Secret key for JWT signing | Required |
| `JWT_EXPIRY_HOURS` | Token expiration time | `24` |

## Development

### Running Tests

```bash
# Backend tests
cd backend && go test ./...

# Optimizer tests
cd optimizer && pytest

# Frontend tests
cd frontend && npm test
```

### Database Migrations

Migrations run automatically on backend startup using **GORM AutoMigrate**. The schema includes:
- `users` - User authentication
- `warehouses` - Distribution centers
- `customers` - Customer locations
- `vehicles` - Delivery vehicles
- `plans` - Delivery plans
- `routes` - Daily routes per plan
- `stops` - Route stops with delivery quantities

**GORM AutoMigrate** automatically:
- Creates missing tables
- Adds missing columns
- Creates missing indexes
- **Does NOT** delete unused columns (protects your data)

For production deployments, consider using migration tools like `golang-migrate` or `goose` for more control over schema changes.

**📊 Database Schema Documentation**: See [DATABASE_SCHEMA.md](./DATABASE_SCHEMA.md) for detailed ER diagram, entity descriptions, relationships, and constraints.

## Database Layer

### GORM ORM

LogiTrackPro uses **GORM (Go Object-Relational Mapping)** v1.25+ for all database operations, replacing the previous raw SQL approach. This provides a more maintainable and type-safe database layer.

**Benefits:**
- **Type Safety**: Compile-time checks for database operations
- **Automatic Relationships**: Preloads related entities (e.g., routes with stops)
- **Simplified Queries**: Concise API for common operations
- **Transaction Support**: Built-in transaction management with `db.Transaction()`
- **Auto Migrations**: Automatic schema synchronization
- **Connection Pooling**: Built-in connection pool management (25 max open, 5 idle)
- **Reduced Boilerplate**: No manual SQL scanning or NULL handling

**Key Features Used:**
- **Preload**: Eager loading of relationships (e.g., `Preload("Vehicle")`, `Preload("Stops.Customer")`)
- **Transactions**: Atomic operations using `db.Transaction()`
- **AutoMigrate**: Automatic schema creation and updates
- **Query Builder**: Fluent API for building complex queries
- **Hooks**: Model lifecycle hooks (BeforeCreate, AfterUpdate, etc.)

### Code Examples

**Simple Query:**
```go
// Get a customer by ID
customer := &models.Customer{}
db.First(customer, id)
```

**Query with Conditions:**
```go
// List available vehicles for a warehouse
var vehicles []models.Vehicle
db.Where("warehouse_id = ? AND available = ?", warehouseID, true).
   Order("name").
   Find(&vehicles)
```

**Query with Relationships (Preload):**
```go
// Get routes with related vehicle and stops
var routes []models.Route
db.Where("plan_id = ?", planID).
   Preload("Vehicle").
   Preload("Stops.Customer").
   Order("day, id").
   Find(&routes)
```

**Create Operation:**
```go
// Create a new customer
customer := &models.Customer{
    Name:             "Acme Corp",
    Latitude:         40.7128,
    Longitude:        -74.0060,
    DemandRate:       100.0,
    MaxInventory:     1000.0,
    CurrentInventory: 500.0,
    MinInventory:     100.0,
}
db.Create(customer)
// customer.ID is automatically populated
```

**Update Operation:**
```go
// Update customer (only specified fields)
db.Model(customer).Updates(models.Customer{
    Name:             "Updated Name",
    CurrentInventory: 600.0,
})
```

**Delete Operation:**
```go
// Delete a customer
db.Delete(&models.Customer{}, id)
```

**Transaction Example:**
```go
// Atomic operation: create route and stops together
err := db.Transaction(func(tx *gorm.DB) error {
    // Create route
    if err := tx.Create(route).Error; err != nil {
        return err // Rollback on error
    }
    
    // Create stops
    for _, stop := range stops {
        stop.RouteID = route.ID
        if err := tx.Create(&stop).Error; err != nil {
            return err // Rollback on error
        }
    }
    
    // Update plan status
    if err := tx.Model(&plan).Update("status", "optimized").Error; err != nil {
        return err // Rollback on error
    }
    
    return nil // Commit on success
})
```

### Model Relationships

GORM automatically handles relationships defined in models:

- **Warehouse** → **Vehicles** (one-to-many)
- **Warehouse** → **Plans** (one-to-many)
- **Plan** → **Routes** (one-to-many, cascade delete)
- **Route** → **Stops** (one-to-many, cascade delete)
- **Route** → **Vehicle** (many-to-one, nullable)
- **Stop** → **Customer** (many-to-one, nullable)
- **Plan** → **User** (many-to-one, nullable, creator)
- **Plan** → **Warehouse** (many-to-one, nullable)

These relationships enable efficient data loading with `Preload()` and automatic foreign key management.

### Database Connection

The database connection is configured in `backend/internal/database/database.go`:

```go
db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Configure connection pool
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
```

### Migration from Raw SQL

The project previously used raw SQL with `database/sql` and `github.com/lib/pq`. The migration to GORM provides:

**Before (Raw SQL):**
```go
query := `SELECT id, name, latitude, longitude FROM customers WHERE id = $1`
err := db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Latitude, &c.Longitude)
```

**After (GORM):**
```go
err := db.First(&customer, id).Error
```

**Benefits of Migration:**
- **Less Code**: ~70% reduction in database code
- **Type Safety**: Compile-time validation
- **Relationships**: Automatic handling of foreign keys
- **Maintainability**: Easier to read and modify
- **Performance**: Built-in query optimization and connection pooling

### Development Tips

**Enable SQL Logging:**
```go
db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Silent), // Silent, Error, Warn, Info
})
```

**Debug Queries:**
```go
// See the actual SQL being executed
db.Debug().Where("id = ?", id).First(&customer)
```

**Raw SQL (when needed):**
```go
// GORM still supports raw SQL for complex queries
db.Raw("SELECT * FROM customers WHERE name LIKE ?", "%Acme%").Scan(&customers)
```

**Batch Operations:**
```go
// Batch insert
db.CreateInBatches(customers, 100)

// Batch update
db.Model(&models.Customer{}).Where("id IN ?", ids).Update("priority", 1)
```

### Troubleshooting

**Common Issues:**

1. **Migration doesn't create tables**: Ensure `RunMigrations()` is called on startup
2. **Relationships not loading**: Use `Preload()` explicitly or check foreign key definitions
3. **NULL values**: Use pointers (`*int64`, `*string`) for nullable fields
4. **Connection pool exhausted**: Increase `SetMaxOpenConns()` if needed
5. **Slow queries**: Use `Preload()` instead of N+1 queries, add indexes

**Checking Database State:**
```go
// Get connection stats
sqlDB, _ := db.DB()
stats := sqlDB.Stats()
fmt.Printf("Open connections: %d\n", stats.OpenConnections)
fmt.Printf("In use: %d\n", stats.InUse)
```

## License

This project is developed for academic purposes as part of logistics optimization research.

## MCP (Model Context Protocol) Setup

LogiTrackPro includes MCP server configuration for enhanced AI assistance. MCPs provide structured access to project resources.

### Quick Setup

```bash
# Run the setup script
./scripts/setup-mcp.sh

# Test the configuration
./scripts/test-mcp.sh
```

### Available MCP Servers

- **PostgreSQL MCP**: Database operations and queries
- **Filesystem MCP**: File operations and navigation
- **Git MCP**: Version control operations
- **Docker MCP**: Container management
- **Brave Search MCP**: Web search (optional, requires API key)

### Helper Scripts

- `scripts/db-helper.sh`: Database operations (backup, restore, stats)
- `scripts/docker-helper.sh`: Docker Compose management
- `scripts/setup-mcp.sh`: MCP configuration setup
- `scripts/test-mcp.sh`: Verify MCP server connectivity

For detailed MCP documentation, see [MCP_SETUP.md](./MCP_SETUP.md).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

