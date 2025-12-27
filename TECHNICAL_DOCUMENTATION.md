# LogiTrackPro Technical Documentation

## 5.1 Overall System Architecture and Logic

### System Overview

LogiTrackPro is a logistics planning platform that solves the Inventory Routing Problem (IRP) through a three-tier architecture: React frontend, Go backend API, and Python optimization microservice. The system manages warehouses, customers, vehicles, and generates optimized multi-day delivery plans.

### Architecture Components

The system consists of three main components:

1. Frontend: React application (Vite) running on port 3000, located in frontend/
2. Backend API: Go application (Gin framework) running on port 8080, located in backend/
3. Optimizer Service: Python application (FastAPI) running on port 8000, located in optimizer/

### System Flow

User requests originate from the React frontend (frontend/src/App.jsx, frontend/src/pages/*.jsx). The frontend communicates with the Go backend API via HTTP REST calls defined in frontend/src/api.js. The backend API handles authentication, data persistence, and business logic. When optimization is requested, the backend calls the Python optimizer service via HTTP. The optimizer service uses Google OR-Tools to solve the IRP problem and returns optimized routes. The backend persists these routes in PostgreSQL database using GORM ORM.

GRAPHIC RECOMMENDATION: System architecture diagram showing three tiers (Frontend, Backend, Optimizer) with arrows indicating request/response flow, database connection, and HTTP communication paths.

### Technology Stack

Frontend: React 18, Vite build tool, Tailwind CSS, Axios for HTTP, React Router for navigation. Main entry point is frontend/src/main.jsx which renders App.jsx.

Backend: Go 1.23, Gin web framework, GORM v1.25 for database operations, PostgreSQL driver, JWT for authentication. Main entry point is backend/cmd/api/main.go.

Optimizer: Python 3.11+, FastAPI web framework, OR-Tools 9.9+ for optimization, Pydantic for data validation. Main entry point is optimizer/main.py.

Database: PostgreSQL 15, accessed via GORM ORM. Connection configured in backend/internal/database/database.go.

### Request Flow

Authentication flow: User submits credentials via POST /api/v1/auth/login (handled by backend/internal/handlers/auth.go Login function). Backend validates credentials against users table, generates JWT token using JWT_SECRET from config, returns token to frontend. Frontend stores token in localStorage and includes it in Authorization header for subsequent requests.

Data management flow: CRUD operations follow pattern: Frontend calls API endpoint (e.g., GET /api/v1/customers), backend handler (backend/internal/handlers/customers.go) receives request, handler calls database function (backend/internal/database/customers.go), database function uses GORM to query PostgreSQL, results returned through handler to frontend as JSON.

Optimization flow: User creates plan via POST /api/v1/plans (backend/internal/handlers/plans.go CreatePlan), then triggers optimization via POST /api/v1/plans/:id/optimize (OptimizePlan function). Handler fetches warehouse, customers, vehicles from database, constructs OptimizeRequest, calls optimizer client (backend/internal/optimizer/client.go Optimize function), optimizer client sends HTTP POST to Python service /optimize endpoint, Python service (optimizer/main.py optimize function) receives request, instantiates IRPSolver (optimizer/solver.py), solver executes multi-day optimization, returns OptimizeResponse with routes, backend receives response, begins database transaction, deletes existing routes, creates new routes and stops, updates plan status, commits transaction, returns plan with routes to frontend.

GRAPHIC RECOMMENDATION: Sequence diagram showing optimization request flow from frontend through backend to optimizer service, including database operations and transaction boundaries.

### Database Architecture

Database schema consists of seven tables: users, warehouses, customers, vehicles, plans, routes, stops. Defined as GORM models in backend/internal/models/models.go. Relationships: User creates Plans (created_by foreign key), Warehouse has Vehicles (warehouse_id foreign key), Warehouse serves Plans (warehouse_id foreign key), Plan contains Routes (plan_id foreign key with CASCADE DELETE), Vehicle assigned to Routes (vehicle_id foreign key), Route has Stops (route_id foreign key with CASCADE DELETE), Customer visited at Stops (customer_id foreign key).

Migrations executed automatically on backend startup via GORM AutoMigrate in backend/internal/database/database.go RunMigrations function. AutoMigrate creates missing tables, adds missing columns, creates indexes, but does not delete unused columns.

GRAPHIC RECOMMENDATION: Entity-relationship diagram showing all seven tables with foreign key relationships, cardinalities, and cascade delete rules. See DATABASE_SCHEMA.md for existing ER diagram.

### Configuration Management

Configuration loaded from environment variables in backend/internal/config/config.go Load function. Required variables: DATABASE_URL (PostgreSQL connection string), JWT_SECRET (token signing key), OPTIMIZER_URL (Python service endpoint), PORT (backend server port, default 8080), JWT_EXPIRY_HOURS (token expiration, default 24). Configuration struct defined in config.go with validation for insecure defaults in production.

### Security Architecture

Authentication uses JWT tokens. Token generation in backend/internal/handlers/auth.go Login function uses golang-jwt/jwt/v5 library. Tokens signed with JWT_SECRET, contain user ID and role, expire after JWT_EXPIRY_HOURS. Middleware in backend/internal/handlers/auth.go AuthMiddleware validates tokens on protected routes, extracts user ID and role, stores in request context.

Password hashing uses bcrypt via golang.org/x/crypto. Passwords hashed during registration in Register function, verified during login in Login function.

CORS configured in backend/cmd/api/main.go corsMiddleware function. Whitelist includes localhost:3000 and localhost:8080, allows credentials, handles OPTIONS preflight requests.

### Error Handling

Backend uses standard HTTP status codes: 200 OK for success, 201 Created for resource creation, 400 Bad Request for validation errors, 401 Unauthorized for authentication failures, 404 Not Found for missing resources, 500 Internal Server Error for server errors. Error responses formatted as JSON with success boolean and error message in backend/internal/handlers/handler.go errorResponse function.

Database errors handled via custom errors: ErrNotFound and ErrDuplicate defined in backend/internal/database/users.go, propagated through handlers to HTTP responses.

Optimizer errors: If optimizer service unavailable, HealthCheck returns error, optimization request fails with 500 status. If optimization fails, OptimizePlan reverts plan status to draft, returns error response.

### Deployment Architecture

Docker Compose configuration in docker-compose.yml defines three services: backend (Go application), optimizer (Python application), frontend (Nginx serving React build), postgres (PostgreSQL database). Services communicate via Docker network, volumes used for database persistence.

Backend Dockerfile (backend/Dockerfile) uses multi-stage build: build stage compiles Go binary, final stage uses minimal base image. Optimizer Dockerfile (optimizer/Dockerfile) uses Python base image, installs dependencies from requirements.txt. Frontend Dockerfile (frontend/Dockerfile) builds React app with Vite, serves with Nginx.

## 5.2 Implementation of the Algorithmic Component

### IRP Problem Formulation

The Inventory Routing Problem combines inventory management with vehicle routing. Problem inputs: Warehouse location (latitude, longitude), Customer locations with inventory constraints (current_inventory, min_inventory, max_inventory, demand_rate), Vehicle fleet with capacity and cost parameters (capacity, cost_per_km, fixed_cost, max_distance), Planning horizon (start_date, end_date defining number of days).

Problem objective: Minimize total delivery cost while ensuring no customer stockouts. Cost components: Fixed cost per vehicle route, variable cost per kilometer traveled, inventory holding costs (not currently optimized).

Constraints: Vehicle capacity limits, maximum travel distance per vehicle, customer inventory bounds (min_inventory <= current_inventory <= max_inventory), temporal constraints (inventory depletes daily at demand_rate).

### Algorithm Architecture

The IRP solver implements a sequential rolling horizon heuristic. Located in optimizer/solver.py IRPSolver class. Algorithm decomposes multi-day IRP into series of single-day VRP subproblems.

Main solving loop in IRPSolver.solve method: For each day in planning_horizon, determine customers needing delivery via _get_customers_needing_delivery, solve VRP for that day via _solve_day_vrp using OR-Tools, apply deliveries to inventory, consume daily demand via _update_inventory, aggregate results.

GRAPHIC RECOMMENDATION: Flowchart showing day-by-day optimization loop with decision points for customer selection, VRP solving, and inventory updates.

### Customer Selection Logic

Customer selection implemented in _get_customers_needing_delivery method. For each customer, calculates days_until_stockout = (current_inventory - min_inventory) / demand_rate. If days_until_stockout <= 2 OR current_inventory <= min_inventory, customer added to delivery list. Customers sorted by priority (descending) then demand_rate (descending).

This implements (s,S) inventory policy variant with 2-day look-ahead window. Policy parameters: s (reorder point) = min_inventory + 2 * demand_rate, S (order-up-to level) = max_inventory.

Delivery quantity calculated in demand_callback function within _solve_day_vrp: delivery_qty = min(max_inventory - current_inventory, max_inventory). This implements fill-up policy maximizing inventory per visit.

### VRP Solving with OR-Tools

Daily VRP solved using Google OR-Tools in _solve_day_vrp method. OR-Tools components: RoutingIndexManager manages node indices (warehouse = 0, customers = 1+), RoutingModel defines VRP structure, callbacks provide distance and demand data.

Distance calculation: Haversine formula implemented in _haversine static method calculates great-circle distance between geographic coordinates. Distance matrix pre-computed in _compute_distance_matrix, stored as integer meters (OR-Tools requires integers).

Constraints: Capacity constraint via AddDimensionWithVehicleCapacity, uses demand_callback returning delivery quantities in grams (integer precision). Distance constraint via AddDimension with max_distance limit per vehicle, global span cost coefficient set to 100 to balance route lengths.

Search strategy: First solution uses PATH_CHEAPEST_ARC strategy building greedy solution by selecting cheapest arcs. Improvement uses GUIDED_LOCAL_SEARCH metaheuristic with adaptive penalty mechanism to escape local optima. Time limit set to 30 seconds per day's VRP problem.

GRAPHIC RECOMMENDATION: Diagram showing OR-Tools VRP model structure with RoutingIndexManager, RoutingModel, callbacks, dimensions, and search parameters.

### Solution Extraction

After OR-Tools solves VRP, solution extracted in _solve_day_vrp method. For each vehicle, iterates through route using routing.Start and routing.NextVar, extracts customer visits, calculates route distance by summing arc distances, computes route cost as fixed_cost + (distance * cost_per_km), generates arrival times assuming 50 km/h average speed and 15 minutes service time per stop.

If OR-Tools fails to find solution, fallback algorithm _create_fallback_routes uses nearest neighbor heuristic: starts from warehouse, repeatedly selects nearest unassigned customer that fits capacity, builds route until capacity exhausted or no customers remain.

### Inventory State Management

Inventory tracked per customer in IRPSolver.inventory dictionary. Initialized from customer.current_inventory in constructor. Updated after deliveries: inventory[customer_id] += delivery_quantity. Updated daily consumption: inventory[customer_id] = max(0, inventory[customer_id] - demand_rate).

State transitions: State(day) -> Deliveries -> State'(day) -> Consumption -> State(day+1). This deterministic state machine assumes constant demand rates, no uncertainty modeling.

### Algorithm Complexity

Time complexity: O(planning_horizon * (n^2 + VRP_time)) where n = number of customers, VRP_time = OR-Tools solving time (typically 30 seconds). Distance matrix computation O(n^2) done once, inventory projection O(n) per day, customer selection O(n log n) per day, VRP optimization O(30s) per day.

Space complexity: O(n^2) for distance matrix, O(n) for inventory state, O(v) for vehicles, O(r) for routes where v = vehicles, r = routes.

Solution quality: Typically 5-15% from optimal for IRP, 1-5% from optimal for individual VRP subproblems. Quality depends on problem size, time limit, and OR-Tools convergence.

## 5.3 Implementation of the Server-Side Logic

### Backend Application Structure

Backend organized as modular monolith. Entry point: backend/cmd/api/main.go initializes components in sequence: loads environment variables via godotenv, loads configuration via config.Load, connects to database via database.Connect, runs migrations via database.RunMigrations, initializes optimizer client via optimizer.NewClient, creates handlers via handlers.New, sets up router via setupRouter, starts HTTP server.

Package structure: backend/internal/config (configuration management), backend/internal/database (data access layer with GORM), backend/internal/handlers (HTTP request handlers), backend/internal/models (domain models), backend/internal/optimizer (optimizer service client).

### HTTP Routing

Router setup in backend/cmd/api/main.go setupRouter function. Uses Gin framework with route groups: /api/v1/auth (public authentication routes), /api/v1/* (protected routes requiring authentication middleware).

Route definitions: Auth routes (POST /register, POST /login, POST /refresh) in backend/internal/handlers/auth.go. Warehouse routes (GET, POST, GET /:id, PUT /:id, DELETE /:id) in backend/internal/handlers/warehouses.go. Customer routes in backend/internal/handlers/customers.go. Vehicle routes in backend/internal/handlers/vehicles.go. Plan routes (GET, POST, GET /:id, DELETE /:id, POST /:id/optimize, GET /:id/routes) in backend/internal/handlers/plans.go. Analytics routes (GET /dashboard, GET /summary) in backend/internal/handlers/analytics.go.

Middleware chain: CORS middleware (corsMiddleware) applied globally, authentication middleware (AuthMiddleware) applied to protected routes, request validation via Gin binding.

### Handler Implementation Pattern

Handlers follow consistent pattern: Extract request parameters (path params, query params, JSON body), validate input (Gin binding, manual validation), call database functions (passing gorm.DB instance), handle errors (convert database errors to HTTP responses), format response (successResponse, createdResponse, errorResponse helpers).

Example: CreateCustomer handler in backend/internal/handlers/customers.go: Binds JSON to CustomerRequest struct, validates required fields, creates Customer model, calls database.CreateCustomer, returns created customer or error.

Error handling: Database ErrNotFound converted to 404 Not Found, validation errors return 400 Bad Request with error message, server errors return 500 Internal Server Error.

### Database Layer Implementation

Database functions use GORM ORM. Connection established in backend/internal/database/database.go Connect function: Opens connection via gorm.Open with PostgreSQL driver, configures connection pool (25 max open, 5 idle), tests connection via Ping.

CRUD operations: List functions use db.Find with Order clauses, Get functions use db.First with ID, Create functions use db.Create, Update functions use db.Model().Updates, Delete functions use db.Delete with RowsAffected check.

Relationships: Preload used for eager loading. Example: GetRoutesByPlan in backend/internal/database/routes.go uses Preload("Vehicle") and Preload("Stops.Customer") to load related entities in single query, avoiding N+1 problem.

Transactions: GORM transactions used in OptimizePlan handler. Pattern: db.Transaction(func(tx *gorm.DB) error { ... }) automatically commits on nil return, rolls back on error return. Transaction-aware functions accept *gorm.DB parameter (regular) or *gorm.DB from transaction (transactional).

### Authentication Implementation

JWT token generation in backend/internal/handlers/auth.go Login function: Validates email/password against database, generates token with jwt.NewWithClaims, signs with JWT_SECRET, sets expiration from JWT_EXPIRY_HOURS, returns token to client.

Token validation in AuthMiddleware: Extracts token from Authorization header, parses with jwt.Parse, validates signature and expiration, extracts claims (user ID, role), stores in context via c.Set, calls c.Next to continue request.

Password hashing: Registration hashes password via bcrypt.GenerateFromPassword, stores hash in database. Login verifies password via bcrypt.CompareHashAndPassword.

### Optimization Request Handling

OptimizePlan handler in backend/internal/handlers/plans.go implements complex workflow: Fetches plan from database, validates plan exists, fetches warehouse, customers, vehicles, calculates planning_horizon from dates, constructs OptimizeRequest with all data, updates plan status to "optimizing", calls optimizer.Optimize with 5-minute timeout, handles optimizer errors (reverts status to "draft"), begins database transaction, deletes existing routes, creates routes and stops from optimizer response, updates plan status to "optimized" with totals, commits transaction, returns plan with routes.

Error recovery: If optimizer fails, plan status reverted to "draft". If transaction fails, all changes rolled back, status reverted. This ensures data consistency.

### Optimizer Client Implementation

Optimizer client in backend/internal/optimizer/client.go provides HTTP client for Python service. Client struct contains baseURL and http.Client with 5-minute timeout.

HealthCheck method: Sends GET request to /health endpoint, returns error if service unavailable or non-200 status.

Optimize method: Marshals OptimizeRequest to JSON, sends POST to /optimize endpoint, decodes OptimizeResponse from JSON, returns response or error.

Request/Response types: OptimizeRequest contains WarehouseData, CustomerData array, VehicleData array, planning_horizon integer, start_date string. OptimizeResponse contains success boolean, message string, total_cost, total_distance, RouteResult array. RouteResult contains day, date, vehicle_id, totals, StopResult array.

### Data Validation

Input validation: Gin binding validates JSON structure against struct tags. Required fields marked with binding:"required". Manual validation for business rules (e.g., end_date >= start_date in CreatePlan).

Database constraints: GORM tags enforce NOT NULL, UNIQUE, defaults. Foreign key constraints enforced at database level. Cascade deletes configured for plan->routes and route->stops relationships.

### Logging and Monitoring

Logging: Standard Go log package used throughout. Log levels: Info for normal operations, Error for failures, Fatal for startup failures. Optimizer service uses Python logging module with INFO level.

Health checks: Backend /health endpoint checks database connection (db.Ping) and optimizer service (optimizer.HealthCheck), returns status for each component.

## 5.4 Module Integration and Data Exchange

### Frontend-Backend Integration

Frontend API client in frontend/src/api.js defines Axios instance with baseURL pointing to backend, interceptors for adding Authorization header from localStorage token, error handling for 401 responses (redirects to login).

API functions: login, register, refreshToken for authentication. listWarehouses, createWarehouse, getWarehouse, updateWarehouse, deleteWarehouse for warehouse management. Similar functions for customers, vehicles, plans. getPlanRoutes for route retrieval. getDashboard, getSummary for analytics.

Data flow: Frontend components (frontend/src/pages/*.jsx) call API functions, API functions send HTTP requests, backend handlers process requests, return JSON responses, frontend updates React state, components re-render.

GRAPHIC RECOMMENDATION: Sequence diagram showing frontend component -> API call -> backend handler -> database -> response flow for typical CRUD operation.

### Backend-Database Integration

GORM ORM provides abstraction layer. Models defined in backend/internal/models/models.go with GORM tags: primaryKey, foreignKey, index, not null, default values, column names, table names.

Database functions in backend/internal/database/*.go use GORM query builder: Where clauses for filtering, Order for sorting, Preload for relationships, Find for multiple records, First for single record, Create for insertion, Updates for modification, Delete for removal.

Connection management: Single *gorm.DB instance created at startup, passed to handlers via Handler struct, handlers pass to database functions, GORM manages connection pooling internally.

Transaction boundaries: Transactions used for multi-step operations (e.g., creating routes and stops). Transaction scope defined by db.Transaction callback function, all operations within callback use same transaction context.

### Backend-Optimizer Integration

HTTP-based integration via REST API. Backend acts as client, optimizer acts as server. Communication protocol: JSON over HTTP POST.

Request format: OptimizeRequest JSON contains warehouse object (id, latitude, longitude, stock), customers array (id, latitude, longitude, demand_rate, inventory fields, priority), vehicles array (id, capacity, cost fields), planning_horizon integer, start_date string (YYYY-MM-DD format).

Response format: OptimizeResponse JSON contains success boolean, message string, total_cost, total_distance, routes array. Each route contains day, date, vehicle_id, totals, stops array. Each stop contains customer_id, sequence, quantity, arrival_time.

Error handling: Network errors (connection refused, timeout) handled by HTTP client, returns error to handler. HTTP errors (non-200 status) handled by checking resp.StatusCode, returns error. Business logic errors (optimization failure) returned in OptimizeResponse.success=false, handler checks success field.

Timeout configuration: HTTP client timeout set to 5 minutes in optimizer.NewClient, accommodates long-running optimizations for large problems.

GRAPHIC RECOMMENDATION: Data flow diagram showing backend -> optimizer request transformation (Go structs to JSON), optimizer processing, response transformation (JSON to Go structs), with data structure examples.

### Data Transformation

Backend to Optimizer: Go structs (models.Warehouse, models.Customer, models.Vehicle) converted to optimizer types (WarehouseData, CustomerData, VehicleData) in OptimizePlan handler. Field mapping: Direct mapping for most fields, date conversion (time.Time to string YYYY-MM-DD), planning_horizon calculated from dates.

Optimizer to Backend: JSON response decoded into OptimizeResponse struct, RouteResult and StopResult arrays converted to models.Route and models.Stop. Field mapping: Date strings parsed to time.Time, vehicle_id converted to *int64 pointer, customer_id converted to *int64 pointer, sequence and quantities directly mapped.

### State Synchronization

Database as source of truth: All persistent state stored in PostgreSQL. Frontend state derived from API responses, not maintained independently. Optimizer state ephemeral, only exists during optimization request.

Optimization state: Plan status field tracks optimization lifecycle: "draft" (created), "optimizing" (request sent), "optimized" (complete), "executed" (optional future state). Status updated atomically within transaction.

Inventory state: Customer inventory levels stored in customers.current_inventory. Updated manually by users or via optimization results. Optimizer uses current_inventory as input, does not directly update database (backend handles updates).

### Error Propagation

Error flow: Database errors (ErrNotFound, ErrDuplicate) propagated from database layer through handlers to HTTP responses. Optimizer errors propagated from HTTP client through handlers to HTTP responses. Network errors (timeout, connection refused) returned as 500 Internal Server Error with descriptive message.

Error recovery: Optimizer failures trigger plan status reversion to "draft". Transaction failures trigger automatic rollback. Partial failures prevented by transaction boundaries.

### Performance Considerations

Database queries: GORM Preload prevents N+1 queries for relationships. Indexes on foreign keys optimize joins. Connection pooling (25 max open, 5 idle) manages database connections efficiently.

Optimizer calls: Single HTTP request per optimization, no polling or streaming. Timeout prevents indefinite blocking. Large requests (many customers) may take minutes, handled asynchronously from user perspective (status updates indicate progress).

Frontend updates: React state updates trigger re-renders. API calls made on component mount and user actions. No automatic polling, manual refresh required.

### Testing Integration Points

Unit tests: Database functions testable with in-memory SQLite via GORM. Handler functions testable with mock database and optimizer client. Optimizer solver testable independently with test data.

Integration tests: End-to-end tests require all three services running. Docker Compose enables full stack testing. Test data seeded in database, optimizer called with known inputs, results verified.

GRAPHIC RECOMMENDATION: Component interaction diagram showing all three services (Frontend, Backend, Optimizer) with database, indicating data flow directions, protocols (HTTP, SQL), and integration points.

### Deployment Integration

Docker Compose orchestrates all services: Services defined in docker-compose.yml, network enables inter-service communication, volumes persist database data, environment variables configured per service, health checks verify service availability.

Service dependencies: Backend depends on postgres (waits for database ready), Frontend depends on backend (API calls), Optimizer independent (called by backend). Startup order: postgres -> optimizer -> backend -> frontend.

Configuration management: Environment variables passed via docker-compose.yml environment section, .env file supported for local development, production uses external configuration management.

This documentation covers the complete technical implementation of LogiTrackPro, from high-level architecture through detailed component interactions and data exchange patterns.
