# Test Suite Documentation

## Overview

This document describes the complete testing suite for LogiTrackPro, covering unit tests, integration tests, and end-to-end tests for all components.

## Test Structure

```
LogiTrackPro/
├── backend/
│   └── internal/
│       ├── database/
│       │   └── customers_test.go          # Database unit tests
│       └── handlers/
│           ├── auth_test.go               # Authentication unit tests
│           └── integration_test.go        # API integration tests
├── optimizer/
│   ├── test_solver.py                    # Solver unit tests
│   └── test_api.py                        # API integration tests
└── frontend/
    └── src/
        ├── api.test.js                    # API client unit tests
        └── components/
            └── __tests__/
                └── Modal.test.jsx         # Component unit tests
```

## Backend Tests (Go)

### Unit Tests

#### database/customers_test.go

**Purpose**: Tests database layer CRUD operations for customers using GORM.

**Test Groups**:

1. **TestCreateCustomer**: Validates customer creation
   - Valid customer with all fields
   - Customer with zero/default values
   - Verifies ID and timestamps are set

2. **TestGetCustomer**: Tests customer retrieval
   - Existing customer retrieval
   - Non-existent customer (ErrNotFound)
   - Verifies correct data returned

3. **TestListCustomers**: Tests customer listing
   - Multiple customers returned
   - Correct ordering by name
   - Empty list handling

4. **TestUpdateCustomer**: Tests customer updates
   - Field updates persist
   - UpdatedAt timestamp changes
   - Non-existent customer (ErrNotFound)

5. **TestDeleteCustomer**: Tests customer deletion
   - Successful deletion
   - Non-existent customer (ErrNotFound)
   - Verifies deletion with GetCustomer

6. **TestCountCustomers**: Tests customer counting
   - Initial count is zero
   - Count increases with creations
   - Accurate count after operations

**Test Database**: Uses in-memory SQLite via GORM for isolation and speed.

**Coverage**: All CRUD operations, error cases, edge cases (zero values, empty lists).

#### handlers/auth_test.go

**Purpose**: Tests authentication handlers including registration, login, token generation, and middleware.

**Test Groups**:

1. **TestRegister**: User registration tests
   - Valid registration creates user and returns token
   - Duplicate email returns 409 Conflict
   - Invalid email format returns 400 Bad Request
   - Password too short returns 400 Bad Request
   - Missing required fields returns 400 Bad Request

2. **TestLogin**: User login tests
   - Valid credentials return token
   - Invalid email returns 401 Unauthorized
   - Invalid password returns 401 Unauthorized
   - Missing credentials return 400 Bad Request

3. **TestAuthMiddleware**: JWT middleware tests
   - Valid token allows access
   - No token returns 401 Unauthorized
   - Invalid token returns 401 Unauthorized
   - Malformed header returns 401 Unauthorized

4. **TestGenerateToken**: Token generation tests
   - Token generated successfully
   - Token contains user ID in subject
   - Expiration time set correctly
   - Token can be parsed and validated

5. **TestParseToken**: Token parsing tests
   - Valid token parsed successfully
   - Invalid token returns error
   - Empty token returns error
   - Wrong secret returns error

**Test Setup**: Uses in-memory database, test JWT secret, Gin test mode.

**Coverage**: Authentication flow, token lifecycle, error scenarios, security validation.

### Integration Tests

#### handlers/integration_test.go

**Purpose**: Tests complete API workflows and component interactions.

**Test Groups**:

1. **TestCustomerCRUDIntegration**: Complete CRUD workflow
   - Create customer via API
   - Retrieve created customer
   - Update customer via API
   - List customers includes created customer
   - Delete customer via API
   - Verify deletion with GET request

2. **TestPlanCreationFlow**: Plan creation with dependencies
   - Create warehouse first
   - Create plan referencing warehouse
   - Verify plan status is "draft"
   - Verify warehouse association

3. **TestProtectedRouteAccess**: Authentication requirements
   - Protected routes require valid token
   - No token returns 401
   - Invalid token returns 401

4. **TestErrorHandling**: Error response validation
   - Invalid ID format returns 400
   - Non-existent resource returns 404
   - Error response structure validation

**Test Setup**: Full database schema, authenticated requests, complete request/response cycle.

**Coverage**: End-to-end API workflows, authentication integration, error propagation.

## Optimizer Tests (Python)

### Unit Tests

#### test_solver.py

**Purpose**: Tests core IRP solver algorithm including distance calculation, customer selection, inventory management, and VRP solving.

**Test Classes**:

1. **TestHaversineDistance**: Distance calculation tests
   - Same point returns zero distance
   - Known distance (NYC to LA) approximately correct
   - Short distances calculated correctly

2. **TestDistanceMatrix**: Distance matrix tests
   - Correct matrix dimensions (customers + warehouse)
   - Diagonal elements are zero
   - Matrix is symmetric
   - All values are integers (OR-Tools requirement)

3. **TestCustomerSelection**: Customer selection logic
   - Customers below minimum inventory selected
   - Customers with <=2 days until stockout selected
   - Customers with >2 days not selected
   - Priority sorting (priority then demand_rate)
   - Zero demand rate handling

4. **TestInventoryManagement**: Inventory state management
   - Inventory initialized from customer data
   - Daily consumption reduces inventory
   - Inventory never goes negative
   - Inventory increases after delivery

5. **TestDeliveryQuantity**: Delivery quantity calculation
   - Fill-up policy fills to max_inventory
   - Zero quantity if already at max
   - Respects inventory constraints

6. **TestVRPSolving**: VRP solving integration
   - Empty customer list returns empty routes
   - Routes generated for customers needing delivery
   - Route structure validation

7. **TestFallbackAlgorithm**: Fallback nearest neighbor
   - Creates routes when OR-Tools fails
   - Respects vehicle capacity
   - Generates valid route structure

8. **TestEndToEndSolver**: Complete optimization
   - Multi-day optimization succeeds
   - Empty customers handled
   - Inventory tracked across days
   - Response structure validation

9. **TestEdgeCases**: Edge cases and error scenarios
   - Zero planning horizon
   - Single customer
   - Single vehicle
   - Very large capacity
   - Customers far apart

**Test Fixtures**: MockWarehouse, MockCustomer, MockVehicle classes for test data.

**Coverage**: Algorithm logic, state management, edge cases, error handling.

### Integration Tests

#### test_api.py

**Purpose**: Tests FastAPI endpoints, request/response handling, validation, and error scenarios.

**Test Classes**:

1. **TestHealthEndpoint**: Health check tests
   - Returns 200 OK
   - Contains status, service, timestamp

2. **TestOptimizeEndpoint**: Optimization endpoint tests
   - Valid request returns success
   - No customers returns error
   - No vehicles returns error
   - Invalid warehouse data returns 422
   - Missing required fields returns 422
   - Invalid date format handled

3. **TestOptimizeScenarios**: Various optimization scenarios
   - Single day optimization
   - Multi-day optimization
   - Large planning horizon (30 days)
   - Multiple vehicles utilized
   - Capacity constraints respected

4. **TestErrorHandling**: Error handling tests
   - Malformed JSON returns 422
   - Wrong HTTP method returns 405
   - Negative values handled
   - Zero capacity vehicle handled

5. **TestPerformance**: Performance tests
   - Many customers (50) handled
   - Concurrent requests (5) handled
   - Response time acceptable

**Test Client**: Uses FastAPI TestClient for HTTP testing.

**Coverage**: API contract, validation, error handling, performance.

## Frontend Tests (JavaScript/React)

### Unit Tests

#### api.test.js

**Purpose**: Tests API client functions, HTTP requests, error handling, and token management.

**Test Groups**:

1. **Authentication**: Login and registration
   - Successful login stores token
   - Login errors handled
   - Registration creates user
   - Registration errors handled

2. **Customer Management**: CRUD operations
   - List customers
   - Create customer
   - Update customer
   - Delete customer
   - API errors handled

3. **Plan Management**: Plan operations
   - Create plan
   - Optimize plan
   - Get plan routes

4. **Token Management**: Authentication token handling
   - Token added to requests
   - Missing token handled
   - 401 responses clear token

5. **Error Handling**: Various error scenarios
   - Network errors
   - Timeout errors
   - Invalid JSON responses

**Test Framework**: Vitest with axios mocking.

**Coverage**: API client functionality, error handling, token management.

#### components/__tests__/Modal.test.jsx

**Purpose**: Tests Modal component rendering and interactions.

**Test Groups**:

1. **Rendering**: Component visibility
   - Renders when isOpen is true
   - Does not render when isOpen is false
   - Renders children content

2. **Interactions**: User interactions
   - Close button calls onClose
   - Backdrop click calls onClose
   - Missing onClose handled gracefully

**Test Framework**: Vitest with React Testing Library.

**Coverage**: Component behavior, props handling, user interactions.

## Running Tests

### Backend Tests

```bash
cd backend
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # Coverage report
go test ./internal/database/...  # Specific package
go test -run TestCreateCustomer  # Specific test
```

### Optimizer Tests

```bash
cd optimizer
pytest                           # Run all tests
pytest -v                        # Verbose output
pytest --cov=solver              # Coverage for solver
pytest test_solver.py           # Specific test file
pytest -k "test_haversine"      # Specific test pattern
```

### Frontend Tests

```bash
cd frontend
npm test                         # Run all tests
npm test -- --coverage          # Coverage report
npm test -- api.test.js         # Specific test file
npm test -- --watch             # Watch mode
```

## Test Coverage Goals

- **Backend**: 80%+ coverage for handlers and database functions
- **Optimizer**: 85%+ coverage for solver logic, 70%+ for API
- **Frontend**: 70%+ coverage for API client, 60%+ for components

## Test Data Management

**Backend**: Uses in-memory SQLite databases, reset between tests, no external dependencies.

**Optimizer**: Uses mock objects (MockWarehouse, MockCustomer, MockVehicle), isolated test data.

**Frontend**: Mocks axios, localStorage, no real HTTP calls, isolated component tests.

## Continuous Integration

Tests should run on:
- Pull request creation
- Code push to main branch
- Pre-commit hooks (optional)

## Assumptions and Limitations

1. **Backend**: Assumes GORM AutoMigrate works correctly, does not test migration edge cases.

2. **Optimizer**: Assumes OR-Tools is installed and functional, some tests may be slow (30s timeout).

3. **Frontend**: Assumes API contract matches backend, does not test actual network calls.

4. **Integration**: Requires all services running for full E2E tests, uses Docker Compose in CI.

5. **Performance**: Load tests are basic, not comprehensive stress testing.

## Future Test Additions

1. **Backend**: 
   - Optimizer client error scenarios
   - Transaction rollback tests
   - Concurrent request handling
   - Database connection pool tests

2. **Optimizer**:
   - OR-Tools failure scenarios
   - Large problem size tests (200+ customers)
   - Memory usage tests
   - Optimization quality benchmarks

3. **Frontend**:
   - Component integration tests
   - Form validation tests
   - Routing tests
   - E2E tests with Playwright/Cypress

4. **E2E**:
   - Complete user workflows
   - Cross-browser testing
   - Mobile responsiveness
   - Performance monitoring

## Test Maintenance

- Update tests when API contracts change
- Add tests for new features
- Review coverage reports regularly
- Refactor tests for clarity and maintainability
- Document test assumptions and limitations
