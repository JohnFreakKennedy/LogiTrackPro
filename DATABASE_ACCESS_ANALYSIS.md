# Database Access Implementation Analysis

## Executive Summary

**SQLAlchemy is NOT used** in this project. The LogiTrackPro backend is written in **Go** and uses a **raw SQL approach** with Go's standard `database/sql` package and the PostgreSQL driver `github.com/lib/pq`. This document provides a comprehensive analysis of how database operations are implemented.

---

## Table of Contents

1. [Technology Stack](#technology-stack)
2. [Database Connection Management](#database-connection-management)
3. [Raw SQL Query Execution](#raw-sql-query-execution)
4. [Manual Result Scanning](#manual-result-scanning)
5. [NULL Value Handling](#null-value-handling)
6. [Transaction Management](#transaction-management)
7. [Migration System](#migration-system)
8. [Error Handling](#error-handling)
9. [Comparison: Raw SQL vs ORM](#comparison-raw-sql-vs-orm)
10. [Code Examples](#code-examples)

---

## 1. Technology Stack

### 1.1 Core Components

| Component | Package | Purpose |
|-----------|---------|---------|
| **Database Interface** | `database/sql` (stdlib) | Go's standard database interface |
| **PostgreSQL Driver** | `github.com/lib/pq` | PostgreSQL-specific driver implementation |
| **Connection Pooling** | Built into `database/sql` | Automatic connection pool management |
| **Query Builder** | None | Raw SQL strings |
| **ORM** | None | Manual struct mapping |

### 1.2 Dependencies

From `go.mod`:
```go
require (
    github.com/lib/pq v1.10.9  // PostgreSQL driver
)
```

**Note**: `database/sql` is part of Go's standard library, so no explicit import needed in `go.mod`.

---

## 2. Database Connection Management

### 2.1 Connection Setup

**Location**: `backend/internal/database/database.go`

**Implementation**:
```go
func Connect(databaseURL string) (*sql.DB, error) {
    // Open connection (doesn't actually connect yet)
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)  // Maximum open connections
    db.SetMaxIdleConns(5)    // Maximum idle connections

    return db, nil
}
```

### 2.2 Connection Pool Configuration

**Settings**:
- **MaxOpenConns**: 25 (maximum concurrent connections)
- **MaxIdleConns**: 5 (connections kept in pool when idle)
- **Connection Lifetime**: Default (no explicit setting)

**How It Works**:
1. `sql.Open()` creates a connection pool but doesn't establish connections
2. `db.Ping()` tests the first connection
3. Connections are created lazily as needed
4. Pool manages connection reuse automatically
5. Connections are closed when pool is closed (`db.Close()`)

### 2.3 Connection String Format

**PostgreSQL Connection URL**:
```
postgres://username:password@host:port/database?sslmode=disable
```

**Example**:
```
postgres://logitrack:password@localhost:5432/logitrackpro?sslmode=disable
```

---

## 3. Raw SQL Query Execution

### 3.1 Query Types

The project uses three main query execution methods:

#### 3.1.1 Single Row Queries (`QueryRow`)

**Use Case**: Fetching a single record (e.g., by ID)

**Example**: `GetCustomer(db *sql.DB, id int64)`
```go
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
```

**Characteristics**:
- Returns single row or error
- Uses `$1, $2, ...` for parameter placeholders (PostgreSQL style)
- Must call `.Scan()` immediately
- Returns `sql.ErrNoRows` if no record found

#### 3.1.2 Multiple Row Queries (`Query`)

**Use Case**: Fetching multiple records (e.g., list operations)

**Example**: `ListCustomers(db *sql.DB)`
```go
func ListCustomers(db *sql.DB) ([]models.Customer, error) {
    query := `SELECT id, name, address, latitude, longitude, demand_rate, 
              max_inventory, current_inventory, min_inventory, holding_cost, priority,
              created_at, updated_at FROM customers ORDER BY name`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // CRITICAL: Must close rows

    var customers []models.Customer
    for rows.Next() {  // Iterate through rows
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
```

**Characteristics**:
- Returns `*sql.Rows` iterator
- Must iterate with `rows.Next()`
- Must call `rows.Close()` (use `defer`)
- Check `rows.Err()` after iteration for errors

#### 3.1.3 Exec Queries (`Exec`)

**Use Case**: INSERT, UPDATE, DELETE operations

**Example**: `UpdateCustomer(db *sql.DB, c *models.Customer)`
```go
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
        return ErrNotFound  // No rows updated
    }
    return nil
}
```

**Characteristics**:
- Returns `sql.Result` with `RowsAffected()` and `LastInsertId()`
- Use `RowsAffected()` to check if update succeeded
- PostgreSQL doesn't support `LastInsertId()` (use `RETURNING` instead)

### 3.2 Parameter Binding

**PostgreSQL Placeholder Style**: `$1, $2, $3, ...`

**Example**:
```go
query := `SELECT * FROM customers WHERE id = $1 AND name = $2`
rows, err := db.Query(query, customerID, customerName)
```

**Benefits**:
- SQL injection prevention (automatic escaping)
- Type safety (driver handles type conversion)
- Performance (prepared statements)

**Note**: Unlike some ORMs, parameter order matters!

---

## 4. Manual Result Scanning

### 4.1 Scanning Process

**Manual Mapping**: Each query result must be manually scanned into struct fields.

**Example**:
```go
var c models.Customer
err := rows.Scan(
    &c.ID,              // int64
    &c.Name,            // string
    &c.Address,         // string
    &c.Latitude,        // float64
    &c.Longitude,       // float64
    &c.DemandRate,      // float64
    &c.MaxInventory,    // float64
    &c.CurrentInventory, // float64
    &c.MinInventory,    // float64
    &c.HoldingCost,     // float64
    &c.Priority,         // int
    &c.CreatedAt,       // time.Time
    &c.UpdatedAt,        // time.Time
)
```

**Requirements**:
1. **Field Order**: Must match SELECT column order exactly
2. **Field Count**: Must match number of columns
3. **Field Types**: Must be compatible with database types
4. **Pointers**: Pass pointers to `Scan()` (e.g., `&c.ID`)

### 4.2 Type Mapping

| Go Type | PostgreSQL Type | Notes |
|---------|----------------|-------|
| `int64` | `INTEGER`, `BIGINT`, `SERIAL` | Standard integer |
| `string` | `VARCHAR`, `TEXT` | String types |
| `float64` | `DOUBLE PRECISION`, `REAL` | Floating point |
| `bool` | `BOOLEAN` | Boolean |
| `time.Time` | `TIMESTAMP`, `DATE` | Time types |
| `sql.NullInt64` | `INTEGER` (nullable) | NULL handling |
| `sql.NullString` | `VARCHAR` (nullable) | NULL handling |
| `sql.NullFloat64` | `DOUBLE PRECISION` (nullable) | NULL handling |
| `sql.NullBool` | `BOOLEAN` (nullable) | NULL handling |
| `sql.NullTime` | `TIMESTAMP` (nullable) | NULL handling |

### 4.3 Struct Tags

**Purpose**: Struct tags (`db:"column_name"`) are used for documentation/reference only. They are **NOT automatically used** by the database layer.

**Example**:
```go
type Customer struct {
    ID               int64     `json:"id" db:"id"`           // Tag exists but not used
    Name             string    `json:"name" db:"name"`       // Manual mapping required
    CurrentInventory float64   `json:"current_inventory" db:"current_inventory"`
    // ...
}
```

**Note**: These tags might be used by other libraries (e.g., JSON marshaling), but the database layer ignores them.

---

## 5. NULL Value Handling

### 5.1 The Problem

PostgreSQL allows `NULL` values, but Go's primitive types (int, string, bool) cannot be `NULL`. This requires special handling.

### 5.2 Solution: `sql.Null*` Types

**Available Types**:
- `sql.NullInt64`
- `sql.NullString`
- `sql.NullFloat64`
- `sql.NullBool`
- `sql.NullTime`

**Example**: Handling nullable foreign keys
```go
func GetRoutesByPlan(db *sql.DB, planID int64) ([]models.Route, error) {
    // ... query setup ...
    
    for rows.Next() {
        var r models.Route
        var vehicleID sql.NullInt64  // Use NullInt64 for nullable column
        var vID sql.NullInt64
        var vName sql.NullString
        
        err := rows.Scan(
            &r.ID, &r.PlanID, &vehicleID, &r.Day, &r.Date,
            &r.TotalDistance, &r.TotalCost, &r.TotalLoad, &r.CreatedAt,
            &vID, &vName, /* ... more fields ... */,
        )
        
        // Check if value is valid (not NULL)
        if vehicleID.Valid {
            vIDVal := vehicleID.Int64
            r.VehicleID = &vIDVal  // Convert to pointer
        }
        
        // Build nested struct if all fields valid
        if vID.Valid && vName.Valid {
            r.Vehicle = &models.Vehicle{
                ID:   vID.Int64,
                Name: vName.String,
                // ... other fields ...
            }
        }
    }
}
```

### 5.3 NULL Handling Patterns

#### Pattern 1: Using Pointers
```go
type Route struct {
    VehicleID *int64    // Pointer allows nil (NULL)
    Vehicle   *Vehicle  // Nested struct pointer
}
```

#### Pattern 2: Using COALESCE
```go
query := `SELECT COALESCE(warehouse_id, 0), COALESCE(created_by, 0) 
          FROM plans WHERE id = $1`
// Returns 0 instead of NULL, can use regular int64
```

#### Pattern 3: Using sql.Null* Types
```go
var warehouseID sql.NullInt64
rows.Scan(&warehouseID)
if warehouseID.Valid {
    // Use warehouseID.Int64
}
```

### 5.4 Inserting NULL Values

**Example**: Inserting with optional foreign key
```go
func CreatePlan(db *sql.DB, p *models.Plan) error {
    var warehouseID, createdBy interface{} = nil, nil
    if p.WarehouseID > 0 {
        warehouseID = p.WarehouseID  // Convert to interface{}
    }
    if p.CreatedBy > 0 {
        createdBy = p.CreatedBy
    }
    
    query := `INSERT INTO plans (name, warehouse_id, created_by) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at, updated_at`
    
    return db.QueryRow(query, p.Name, warehouseID, createdBy).Scan(...)
}
```

**Key**: Use `interface{}` type and set to `nil` for NULL values.

---

## 6. Transaction Management

### 6.1 Manual Transaction Control

**No Automatic Transactions**: Unlike ORMs, transactions must be explicitly managed.

**Example**: Atomic route creation
```go
func (h *Handler) OptimizePlan(c *gin.Context) {
    // ... validation and optimization ...
    
    // Begin transaction
    tx, err := h.db.Begin()
    if err != nil {
        // Handle error
        return
    }
    
    // Delete existing routes (within transaction)
    if err := database.DeleteRoutesByPlanTx(tx, id); err != nil {
        tx.Rollback()  // Rollback on error
        return
    }
    
    // Create routes (within transaction)
    for _, routeResult := range optResp.Routes {
        if err := database.CreateRouteTx(tx, route); err != nil {
            tx.Rollback()  // Rollback on error
            return
        }
        
        // Create stops (within transaction)
        for _, stopResult := range routeResult.Stops {
            if err := database.CreateStopTx(tx, stop); err != nil {
                tx.Rollback()  // Rollback on error
                return
            }
        }
    }
    
    // Update plan status (within transaction)
    if err := database.UpdatePlanStatusTx(tx, id, "optimized", ...); err != nil {
        tx.Rollback()
        return
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        tx.Rollback()  // Rollback if commit fails
        return
    }
}
```

### 6.2 Transaction-Aware Functions

**Pattern**: Functions accept either `*sql.DB` or `*sql.Tx`

**Example**:
```go
// Regular function (uses *sql.DB)
func CreateRoute(db *sql.DB, r *models.Route) error {
    query := `INSERT INTO routes (...) VALUES (...) RETURNING id, created_at`
    return db.QueryRow(query, ...).Scan(&r.ID, &r.CreatedAt)
}

// Transaction-aware function (uses *sql.Tx)
func CreateRouteTx(tx *sql.Tx, r *models.Route) error {
    query := `INSERT INTO routes (...) VALUES (...) RETURNING id, created_at`
    return tx.QueryRow(query, ...).Scan(&r.ID, &r.CreatedAt)
}
```

**Benefits**:
- Same logic, different execution context
- Can be used in or out of transactions
- Explicit transaction control

### 6.3 Transaction Best Practices

1. **Always Rollback on Error**: `defer tx.Rollback()` or explicit rollback
2. **Commit Only on Success**: Commit after all operations succeed
3. **Use Defer for Safety**: `defer tx.Rollback()` ensures cleanup
4. **Check Commit Errors**: Commit can fail, check return value

**Improved Pattern**:
```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()  // Safe: no-op if already committed

// ... operations ...

if err := tx.Commit(); err != nil {
    return err  // defer will rollback
}
// Success: defer won't rollback committed transaction
```

---

## 7. Migration System

### 7.1 In-Code Migrations

**Location**: `backend/internal/database/database.go`

**Implementation**:
```go
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
        `CREATE TABLE IF NOT EXISTS warehouses (...)`,
        `CREATE TABLE IF NOT EXISTS customers (...)`,
        // ... more migrations ...
    }

    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return fmt.Errorf("migration failed: %w", err)
        }
    }
    return nil
}
```

### 7.2 Migration Characteristics

**Properties**:
- **Idempotent**: Uses `CREATE TABLE IF NOT EXISTS`
- **Sequential**: Executed in order
- **No Versioning**: No migration tracking table
- **Simple**: Basic DDL statements only

**Limitations**:
- No rollback support
- No version tracking
- No migration history
- Cannot detect already-run migrations
- All migrations run every time

### 7.3 Migration Execution

**Called on Startup**: `backend/cmd/api/main.go`
```go
func main() {
    // ... connection setup ...
    
    // Run migrations
    if err := database.RunMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // ... start server ...
}
```

**Note**: This is a simple approach. Production systems typically use migration tools like:
- `golang-migrate/migrate`
- `pressly/goose`
- `rubenv/sql-migrate`

---

## 8. Error Handling

### 8.1 Error Types

**Standard Errors**:
- `sql.ErrNoRows`: No rows returned (for `QueryRow`)
- `sql.ErrConnDone`: Connection closed
- Database-specific errors: Constraint violations, etc.

### 8.2 Custom Error Wrapping

**Pattern**: Define custom errors for business logic
```go
var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("record already exists")

func GetCustomer(db *sql.DB, id int64) (*models.Customer, error) {
    // ... query ...
    err := db.QueryRow(query, id).Scan(...)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound  // Convert to custom error
    }
    if err != nil {
        return nil, err  // Pass through other errors
    }
    return c, nil
}
```

### 8.3 Error Checking Patterns

**Pattern 1**: Check for specific errors
```go
if err == sql.ErrNoRows {
    return nil, ErrNotFound
}
```

**Pattern 2**: Check RowsAffected
```go
result, err := db.Exec(query, ...)
if err != nil {
    return err
}
rows, _ := result.RowsAffected()
if rows == 0 {
    return ErrNotFound
}
```

**Pattern 3**: Check constraint violations
```go
func isUniqueViolation(err error) bool {
    return err != nil && (
        contains(err.Error(), "unique") || 
        contains(err.Error(), "duplicate")
    )
}

func CreateUser(db *sql.DB, user *models.User) error {
    err := db.QueryRow(query, ...).Scan(...)
    if err != nil {
        if isUniqueViolation(err) {
            return ErrDuplicate
        }
        return err
    }
    return nil
}
```

---

## 9. Comparison: Raw SQL vs ORM

### 9.1 Current Approach (Raw SQL)

**Advantages**:
- ✅ **Full Control**: Complete control over SQL queries
- ✅ **Performance**: No ORM overhead, direct SQL execution
- ✅ **Explicit**: Clear what SQL is executed
- ✅ **Flexibility**: Can use any PostgreSQL feature
- ✅ **Learning**: Developers learn SQL directly
- ✅ **Debugging**: Easy to debug (see exact SQL)

**Disadvantages**:
- ❌ **Boilerplate**: Lots of repetitive scanning code
- ❌ **Error-Prone**: Manual mapping is error-prone
- ❌ **Maintenance**: Schema changes require code changes
- ❌ **No Relationships**: Manual JOIN handling
- ❌ **No Validation**: No automatic type validation
- ❌ **Verbose**: More code for simple operations

### 9.2 Alternative: ORM Approach (GORM example)

**What It Would Look Like**:
```go
// With GORM (hypothetical)
func GetCustomer(db *gorm.DB, id int64) (*models.Customer, error) {
    var c models.Customer
    err := db.First(&c, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrNotFound
    }
    return &c, err
}

func ListCustomers(db *gorm.DB) ([]models.Customer, error) {
    var customers []models.Customer
    err := db.Order("name").Find(&customers).Error
    return customers, err
}
```

**Advantages**:
- ✅ **Less Code**: Concise operations
- ✅ **Type Safety**: Compile-time checks
- ✅ **Relationships**: Automatic relationship loading
- ✅ **Migrations**: Built-in migration support
- ✅ **Validation**: Automatic validation

**Disadvantages**:
- ❌ **Learning Curve**: ORM-specific syntax
- ❌ **Performance**: Potential N+1 queries
- ❌ **Less Control**: Limited SQL control
- ❌ **Complexity**: ORM abstraction layer

### 9.3 When to Use Each Approach

**Use Raw SQL When**:
- Performance is critical
- Complex queries (window functions, CTEs)
- Full control needed
- Small team familiar with SQL
- Simple data models

**Use ORM When**:
- Rapid development needed
- Complex relationships
- Team prefers abstraction
- Standard CRUD operations
- Built-in features needed (migrations, validation)

---

## 10. Code Examples

### 10.1 Complete CRUD Example: Customers

**Create**:
```go
func CreateCustomer(db *sql.DB, c *models.Customer) error {
    query := `INSERT INTO customers (name, address, latitude, longitude, demand_rate, 
              max_inventory, current_inventory, min_inventory, holding_cost, priority) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
              RETURNING id, created_at, updated_at`
    
    return db.QueryRow(query, c.Name, c.Address, c.Latitude, c.Longitude,
        c.DemandRate, c.MaxInventory, c.CurrentInventory, c.MinInventory,
        c.HoldingCost, c.Priority).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}
```

**Read (Single)**:
```go
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
```

**Read (Multiple)**:
```go
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
```

**Update**:
```go
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
```

**Delete**:
```go
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
```

### 10.2 Complex Query Example: JOINs

**Example**: Fetching routes with vehicle and customer information
```go
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
        // Use sql.Null* types for nullable columns
        var vehicleID sql.NullInt64
        var vID sql.NullInt64
        var vName sql.NullString
        var vCapacity sql.NullFloat64
        // ... more nullable fields ...

        err := rows.Scan(
            &r.ID, &r.PlanID, &vehicleID, &r.Day, &r.Date,
            &r.TotalDistance, &r.TotalCost, &r.TotalLoad, &r.CreatedAt,
            &vID, &vName, &vCapacity, /* ... */,
        )
        if err != nil {
            return nil, err
        }

        // Handle nullable vehicle_id
        if vehicleID.Valid {
            vIDVal := vehicleID.Int64
            r.VehicleID = &vIDVal
        }

        // Build nested Vehicle struct if valid
        if vID.Valid && vName.Valid {
            r.Vehicle = &models.Vehicle{
                ID:          vID.Int64,
                Name:        vName.String,
                Capacity:    vCapacity.Float64,
                // ... other fields ...
            }
        }

        routes = append(routes, r)
    }

    // Load stops for each route (N+1 query pattern)
    for i := range routes {
        stops, err := GetStopsByRoute(db, routes[i].ID)
        if err != nil {
            return nil, err
        }
        routes[i].Stops = stops
    }

    return routes, nil
}
```

**Note**: This example shows the N+1 query problem. A more efficient approach would use a single query with multiple JOINs or a batch query.

### 10.3 Transaction Example: Atomic Operations

**Example**: Creating a plan with routes and stops atomically
```go
func (h *Handler) OptimizePlan(c *gin.Context) {
    // ... validation and optimization ...
    
    // Begin transaction
    tx, err := h.db.Begin()
    if err != nil {
        errorResponse(c, http.StatusInternalServerError, "Failed to begin transaction")
        return
    }
    defer tx.Rollback()  // Safe: no-op if committed

    // Delete existing routes
    if err := database.DeleteRoutesByPlanTx(tx, id); err != nil {
        errorResponse(c, http.StatusInternalServerError, "Failed to clear routes")
        return
    }

    // Create routes and stops
    for _, routeResult := range optResp.Routes {
        route := &models.Route{...}
        if err := database.CreateRouteTx(tx, route); err != nil {
            errorResponse(c, http.StatusInternalServerError, "Failed to save route")
            return
        }

        for _, stopResult := range routeResult.Stops {
            stop := &models.Stop{...}
            if err := database.CreateStopTx(tx, stop); err != nil {
                errorResponse(c, http.StatusInternalServerError, "Failed to save stop")
                return
            }
        }
    }

    // Update plan status
    if err := database.UpdatePlanStatusTx(tx, id, "optimized", ...); err != nil {
        errorResponse(c, http.StatusInternalServerError, "Failed to update plan")
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        errorResponse(c, http.StatusInternalServerError, "Failed to commit")
        return
    }
    // Success: defer won't rollback committed transaction
}
```

---

## 11. Summary

### 11.1 Key Characteristics

1. **No ORM**: Uses raw SQL with `database/sql`
2. **Manual Mapping**: All struct scanning is manual
3. **Explicit Transactions**: Transactions must be explicitly managed
4. **NULL Handling**: Uses `sql.Null*` types or pointers
5. **PostgreSQL-Specific**: Uses `$1, $2, ...` placeholders
6. **Connection Pooling**: Built into `database/sql`
7. **Simple Migrations**: In-code migrations without versioning

### 11.2 File Structure

```
backend/
├── internal/
│   ├── database/
│   │   ├── database.go      # Connection & migrations
│   │   ├── users.go         # User CRUD operations
│   │   ├── customers.go     # Customer CRUD operations
│   │   ├── vehicles.go      # Vehicle CRUD operations
│   │   ├── warehouses.go    # Warehouse CRUD operations
│   │   ├── plans.go         # Plan CRUD operations
│   │   └── routes.go        # Route & Stop operations
│   └── models/
│       └── models.go        # Struct definitions
└── cmd/api/
    └── main.go              # Application entry point
```

### 11.3 Advantages of Current Approach

- ✅ **Performance**: Direct SQL execution, no ORM overhead
- ✅ **Control**: Full control over queries
- ✅ **Clarity**: Explicit SQL, easy to understand
- ✅ **Flexibility**: Can use any PostgreSQL feature
- ✅ **Learning**: Developers learn SQL directly

### 11.4 Areas for Improvement

- ⚠️ **N+1 Queries**: Some queries could be optimized with JOINs
- ⚠️ **Migration System**: Could use proper migration tooling
- ⚠️ **Error Handling**: Could be more consistent
- ⚠️ **Code Duplication**: Scanning code is repetitive
- ⚠️ **Type Safety**: Manual mapping is error-prone

---

**Document Version**: 1.0  
**Last Updated**: 2024  
**Author**: Database Access Analysis for LogiTrackPro
