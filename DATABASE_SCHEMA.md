# Database Schema Documentation

## Entity-Relationship Diagram

```mermaid
erDiagram
    USER ||--o{ PLAN : "creates"
    WAREHOUSE ||--o{ VEHICLE : "has"
    WAREHOUSE ||--o{ PLAN : "serves"
    PLAN ||--o{ ROUTE : "contains"
    VEHICLE ||--o{ ROUTE : "assigned_to"
    ROUTE ||--o{ STOP : "has"
    CUSTOMER ||--o{ STOP : "visited_at"

    USER {
        int64 id PK
        string email UK "unique"
        string password_hash
        string name
        string role "default: 'user'"
        timestamp created_at
        timestamp updated_at
    }

    WAREHOUSE {
        int64 id PK
        string name
        text address
        double precision latitude
        double precision longitude
        double precision capacity "default: 0"
        double precision current_stock "default: 0"
        double precision holding_cost "default: 0"
        double precision replenishment_qty "default: 0"
        timestamp created_at
        timestamp updated_at
    }

    CUSTOMER {
        int64 id PK
        string name
        text address
        double precision latitude
        double precision longitude
        double precision demand_rate "default: 0"
        double precision max_inventory "default: 0"
        double precision current_inventory "default: 0"
        double precision min_inventory "default: 0"
        double precision holding_cost "default: 0"
        integer priority "default: 1"
        timestamp created_at
        timestamp updated_at
    }

    VEHICLE {
        int64 id PK
        string name
        double precision capacity
        double precision cost_per_km "default: 0"
        double precision fixed_cost "default: 0"
        double precision max_distance "default: 0"
        boolean available "default: true"
        int64 warehouse_id FK "nullable"
        timestamp created_at
        timestamp updated_at
    }

    PLAN {
        int64 id PK
        string name
        date start_date
        date end_date
        string status "default: 'draft'"
        double precision total_cost "default: 0"
        double precision total_distance "default: 0"
        int64 warehouse_id FK "nullable"
        int64 created_by FK "nullable"
        timestamp created_at
        timestamp updated_at
    }

    ROUTE {
        int64 id PK
        int64 plan_id FK "not null"
        int64 vehicle_id FK "nullable"
        integer day "not null"
        date date "not null"
        double precision total_distance "default: 0"
        double precision total_cost "default: 0"
        double precision total_load "default: 0"
        timestamp created_at
    }

    STOP {
        int64 id PK
        int64 route_id FK "not null"
        int64 customer_id FK "nullable"
        integer sequence "not null"
        double precision quantity "default: 0"
        varchar arrival_time "varchar(10)"
        timestamp created_at
    }
```

## Detailed Entity Descriptions

### 1. USER

**Purpose**: Stores system user accounts for authentication and authorization.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `email` (UK): Unique email address for login, indexed
- `password_hash`: Bcrypt-hashed password (never stored in plain text)
- `name`: User's full name
- `role`: User role (default: 'user'), can be 'user' or 'admin'
- `created_at`: Timestamp of account creation
- `updated_at`: Timestamp of last update

**Relationships**:
- **One-to-Many** with `PLAN` (via `created_by`): A user can create multiple plans

**Constraints**:
- Email must be unique
- Email and password_hash cannot be NULL
- Role defaults to 'user'

**Indexes**:
- Primary key on `id`
- Unique index on `email`

---

### 2. WAREHOUSE

**Purpose**: Represents distribution centers or warehouses that store inventory and dispatch vehicles.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `name`: Warehouse name/identifier
- `address`: Physical address (text field)
- `latitude`: Geographic latitude coordinate (required)
- `longitude`: Geographic longitude coordinate (required)
- `capacity`: Maximum storage capacity (default: 0)
- `current_stock`: Current inventory level (default: 0)
- `holding_cost`: Cost per unit per time period for holding inventory (default: 0)
- `replenishment_qty`: Standard replenishment quantity (default: 0)
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

**Relationships**:
- **One-to-Many** with `VEHICLE` (via `warehouse_id`): A warehouse can have multiple vehicles
- **One-to-Many** with `PLAN` (via `warehouse_id`): A warehouse can be associated with multiple delivery plans

**Constraints**:
- Latitude and longitude are required (NOT NULL)
- All numeric fields default to 0

**Indexes**:
- Primary key on `id`

**Business Rules**:
- Warehouses serve as the origin point for all delivery routes
- Vehicles are typically assigned to a specific warehouse
- Plans are created for a specific warehouse's operations

---

### 3. CUSTOMER

**Purpose**: Represents customer locations that receive deliveries and maintain inventory.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `name`: Customer name/identifier
- `address`: Physical address (text field)
- `latitude`: Geographic latitude coordinate (required)
- `longitude`: Geographic longitude coordinate (required)
- `demand_rate`: Daily consumption/demand rate (default: 0)
- `max_inventory`: Maximum inventory capacity (default: 0)
- `current_inventory`: Current inventory level (default: 0)
- `min_inventory`: Minimum inventory threshold (default: 0)
- `holding_cost`: Cost per unit per time period for holding inventory (default: 0)
- `priority`: Customer priority level (default: 1, higher = more important)
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

**Relationships**:
- **One-to-Many** with `STOP` (via `customer_id`): A customer can be visited multiple times across different routes

**Constraints**:
- Latitude and longitude are required (NOT NULL)
- All numeric fields default to 0
- Priority defaults to 1

**Indexes**:
- Primary key on `id`

**Business Rules**:
- Customers consume inventory at `demand_rate` per day
- When `current_inventory` drops to `min_inventory` or below, delivery is triggered
- Delivery quantity cannot exceed `max_inventory - current_inventory`
- Higher priority customers are served first during optimization

---

### 4. VEHICLE

**Purpose**: Represents delivery vehicles in the fleet with capacity and cost parameters.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `name`: Vehicle identifier/name
- `capacity`: Maximum load capacity (required)
- `cost_per_km`: Variable cost per kilometer traveled (default: 0)
- `fixed_cost`: Fixed cost per route/assignment (default: 0)
- `max_distance`: Maximum travel distance per route (default: 0, 0 = unlimited)
- `available`: Whether vehicle is available for assignment (default: true)
- `warehouse_id` (FK): Reference to assigned warehouse (nullable)
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

**Relationships**:
- **Many-to-One** with `WAREHOUSE` (via `warehouse_id`): Vehicle belongs to a warehouse (nullable)
- **One-to-Many** with `ROUTE` (via `vehicle_id`): Vehicle can be assigned to multiple routes

**Constraints**:
- Capacity is required (NOT NULL)
- Warehouse_id is nullable (vehicle can be unassigned)
- Available defaults to true

**Indexes**:
- Primary key on `id`
- Index on `warehouse_id` for efficient warehouse queries

**Business Rules**:
- Vehicle capacity limits total load on a route
- Total route distance cannot exceed `max_distance` (if specified)
- Route cost = `fixed_cost + (distance × cost_per_km)`
- Only available vehicles are considered during optimization

---

### 5. PLAN

**Purpose**: Represents a multi-day delivery plan spanning a planning horizon.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `name`: Plan name/identifier
- `start_date`: Start date of planning horizon (required)
- `end_date`: End date of planning horizon (required)
- `status`: Plan status (default: 'draft')
  - `draft`: Plan created but not optimized
  - `optimizing`: Optimization in progress
  - `optimized`: Optimization complete
  - `executed`: Plan has been executed
- `total_cost`: Total cost across all routes (default: 0)
- `total_distance`: Total distance across all routes (default: 0)
- `warehouse_id` (FK): Reference to associated warehouse (nullable)
- `created_by` (FK): Reference to user who created the plan (nullable)
- `created_at`: Timestamp of creation
- `updated_at`: Timestamp of last update

**Relationships**:
- **Many-to-One** with `WAREHOUSE` (via `warehouse_id`): Plan is associated with a warehouse
- **Many-to-One** with `USER` (via `created_by`): Plan creator
- **One-to-Many** with `ROUTE` (via `plan_id`): Plan contains multiple routes (CASCADE DELETE)

**Constraints**:
- Start_date and end_date are required
- End_date must be >= start_date
- Status defaults to 'draft'
- Warehouse_id and created_by are nullable

**Indexes**:
- Primary key on `id`
- Index on `warehouse_id` for efficient warehouse queries
- Index on `created_by` for user queries

**Business Rules**:
- Planning horizon = `end_date - start_date + 1` days
- When plan is deleted, all associated routes are automatically deleted (CASCADE)
- Total cost and distance are calculated after optimization
- Status tracks plan lifecycle from creation to execution

---

### 6. ROUTE

**Purpose**: Represents a single day's delivery route within a plan, assigned to a vehicle.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `plan_id` (FK): Reference to parent plan (required)
- `vehicle_id` (FK): Reference to assigned vehicle (nullable)
- `day`: Day number within plan (1-indexed, required)
- `date`: Actual date of route execution (required)
- `total_distance`: Total route distance in kilometers (default: 0)
- `total_cost`: Total route cost (default: 0)
- `total_load`: Total load delivered on route (default: 0)
- `created_at`: Timestamp of creation

**Relationships**:
- **Many-to-One** with `PLAN` (via `plan_id`): Route belongs to a plan (CASCADE DELETE)
- **Many-to-One** with `VEHICLE` (via `vehicle_id`): Route assigned to a vehicle (nullable)
- **One-to-Many** with `STOP` (via `route_id`): Route contains multiple stops (CASCADE DELETE)

**Constraints**:
- Plan_id, day, and date are required (NOT NULL)
- Vehicle_id is nullable (route can be unassigned)
- Day must be >= 1

**Indexes**:
- Primary key on `id`
- Index on `plan_id` for efficient plan queries
- Index on `vehicle_id` for vehicle queries

**Business Rules**:
- Each route represents one day's deliveries
- Routes are ordered by `day` within a plan
- When plan is deleted, all routes are automatically deleted (CASCADE)
- When route is deleted, all stops are automatically deleted (CASCADE)
- Total cost = vehicle.fixed_cost + (total_distance × vehicle.cost_per_km)
- Total load must not exceed vehicle capacity

---

### 7. STOP

**Purpose**: Represents a single delivery stop on a route, visiting a customer.

**Attributes**:
- `id` (PK): Primary key, auto-incrementing integer
- `route_id` (FK): Reference to parent route (required)
- `customer_id` (FK): Reference to customer being visited (nullable)
- `sequence`: Stop sequence number on route (1-indexed, required)
- `quantity`: Delivery quantity at this stop (default: 0)
- `arrival_time`: Estimated arrival time (format: "HH:MM", varchar(10))
- `created_at`: Timestamp of creation

**Relationships**:
- **Many-to-One** with `ROUTE` (via `route_id`): Stop belongs to a route (CASCADE DELETE)
- **Many-to-One** with `CUSTOMER` (via `customer_id`): Stop visits a customer (nullable)

**Constraints**:
- Route_id and sequence are required (NOT NULL)
- Customer_id is nullable (stop can be unassigned)
- Sequence must be >= 1
- Sequence determines stop order on route

**Indexes**:
- Primary key on `id`
- Index on `route_id` for efficient route queries
- Index on `customer_id` for customer queries

**Business Rules**:
- Stops are ordered by `sequence` within a route
- Sequence 1 is the first stop after leaving warehouse
- When route is deleted, all stops are automatically deleted (CASCADE)
- Quantity delivered increases customer's current_inventory
- Arrival time is calculated based on distance and average speed

---

## Relationship Summary

### Cardinality Overview

| Relationship | Type | Cardinality | Foreign Key | Cascade Delete |
|--------------|------|-------------|-------------|----------------|
| USER → PLAN | One-to-Many | 1:N | `created_by` | No |
| WAREHOUSE → VEHICLE | One-to-Many | 1:N | `warehouse_id` | No |
| WAREHOUSE → PLAN | One-to-Many | 1:N | `warehouse_id` | No |
| PLAN → ROUTE | One-to-Many | 1:N | `plan_id` | **Yes** |
| VEHICLE → ROUTE | One-to-Many | 1:N | `vehicle_id` | No |
| ROUTE → STOP | One-to-Many | 1:N | `route_id` | **Yes** |
| CUSTOMER → STOP | One-to-Many | 1:N | `customer_id` | No |

### Cascade Delete Rules

**Critical**: The following deletions trigger cascade deletes:

1. **PLAN deletion** → Deletes all associated ROUTEs
2. **ROUTE deletion** → Deletes all associated STOPs

**Rationale**:
- Routes are meaningless without their parent plan
- Stops are meaningless without their parent route
- This ensures data integrity and prevents orphaned records

### Nullable Foreign Keys

The following foreign keys are nullable, allowing for flexible data modeling:

- `VEHICLE.warehouse_id`: Vehicle can be unassigned
- `PLAN.warehouse_id`: Plan can be created without warehouse assignment
- `PLAN.created_by`: Plan can exist without creator tracking
- `ROUTE.vehicle_id`: Route can be created before vehicle assignment
- `STOP.customer_id`: Stop can exist without customer assignment

---

## Indexes

### Primary Keys
All tables have an auto-incrementing `id` as primary key.

### Foreign Key Indexes
All foreign keys are indexed for efficient joins:
- `vehicles.warehouse_id`
- `plans.warehouse_id`
- `plans.created_by`
- `routes.plan_id`
- `routes.vehicle_id`
- `stops.route_id`
- `stops.customer_id`

### Unique Indexes
- `users.email` (unique constraint)

### Composite Indexes
None currently defined, but could be added for:
- `routes(plan_id, day)` - for efficient plan route queries
- `stops(route_id, sequence)` - for efficient route stop ordering

---

## Data Types

### Numeric Types
- **Integers**: `id` fields, `day`, `sequence`, `priority` use `INTEGER` or `BIGINT`
- **Floating Point**: All distance, cost, capacity, inventory values use `DOUBLE PRECISION`
- **Boolean**: `available` uses `BOOLEAN`

### String Types
- **VARCHAR(255)**: Short strings (names, emails, status)
- **VARCHAR(50)**: Short codes (role, status)
- **VARCHAR(10)**: Very short strings (arrival_time)
- **TEXT**: Long strings (addresses)

### Date/Time Types
- **DATE**: Date-only fields (`start_date`, `end_date`, `date`)
- **TIMESTAMP**: Full timestamp fields (`created_at`, `updated_at`)

---

## Constraints and Business Rules

### Referential Integrity
- All foreign keys maintain referential integrity
- Cascade deletes ensure no orphaned records
- Nullable foreign keys allow flexible data entry

### Data Validation
- Email uniqueness enforced at database level
- Required fields prevent NULL values where inappropriate
- Default values provide sensible defaults

### Temporal Constraints
- `plan.end_date >= plan.start_date` (enforced at application level)
- `route.day >= 1` (enforced at application level)
- `stop.sequence >= 1` (enforced at application level)

### Inventory Constraints
- `customer.current_inventory >= 0` (enforced at application level)
- `customer.current_inventory <= customer.max_inventory` (enforced at application level)
- `customer.current_inventory >= customer.min_inventory` (enforced at application level)

### Route Constraints
- `route.total_load <= vehicle.capacity` (enforced at application level)
- `route.total_distance <= vehicle.max_distance` (if max_distance > 0, enforced at application level)

---

## Query Patterns

### Common Queries

**1. Get plan with all routes and stops:**
```sql
SELECT p.*, r.*, s.*
FROM plans p
LEFT JOIN routes r ON r.plan_id = p.id
LEFT JOIN stops s ON s.route_id = r.id
WHERE p.id = ?
ORDER BY r.day, s.sequence;
```

**2. Get available vehicles for warehouse:**
```sql
SELECT * FROM vehicles
WHERE warehouse_id = ? AND available = true
ORDER BY name;
```

**3. Get customers needing delivery:**
```sql
SELECT * FROM customers
WHERE current_inventory <= min_inventory
   OR (current_inventory - min_inventory) / demand_rate <= 2
ORDER BY priority DESC, demand_rate DESC;
```

**4. Get route statistics:**
```sql
SELECT 
    COUNT(*) as total_routes,
    SUM(total_distance) as total_distance,
    SUM(total_cost) as total_cost,
    AVG(total_load) as avg_load
FROM routes
WHERE plan_id = ?;
```

---

## Performance Considerations

### Indexing Strategy
- All foreign keys are indexed for efficient joins
- Primary keys provide fast lookups
- Consider composite indexes for common query patterns

### Query Optimization
- Use `Preload()` in GORM to avoid N+1 queries
- Batch operations for bulk inserts/updates
- Connection pooling configured (25 max open, 5 idle)

### Scalability
- Table structure supports horizontal scaling
- No table-level locks required for common operations
- Cascade deletes are efficient (indexed foreign keys)

---

## Migration Notes

### GORM AutoMigrate Behavior
- Creates missing tables
- Adds missing columns
- Creates missing indexes
- **Does NOT** delete unused columns (protects data)
- **Does NOT** modify column types (manual migration required)

### Manual Migrations
For production deployments, consider:
- Using `golang-migrate` or `goose` for versioned migrations
- Testing migrations on staging environment first
- Backing up database before major schema changes

---

## Future Enhancements

### Potential Additions
1. **Soft Deletes**: Add `deleted_at` timestamp for soft delete pattern
2. **Audit Trail**: Track who modified records and when
3. **Versioning**: Track plan versions for historical analysis
4. **Constraints**: Add database-level check constraints for business rules
5. **Partitioning**: Partition routes/stops by date for large datasets
6. **Full-Text Search**: Add full-text indexes for name/address searches

---

**Document Version**: 1.0  
**Last Updated**: 2024  
**Database Engine**: PostgreSQL 15  
**ORM**: GORM v1.25+
