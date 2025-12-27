# Testing Guide

## Quick Start

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Optimizer Tests
```bash
cd optimizer
pip install -r test_requirements.txt
pytest -v
```

### Frontend Tests
```bash
cd frontend
npm install
npm test
```

## Test Structure

### Backend Tests (Go)

**Location**: `backend/internal/*/**_test.go`

**Framework**: Go testing package with GORM (SQLite in-memory)

**Key Test Files**:

1. **database/customers_test.go** - Database layer unit tests
   - Validates: CRUD operations, error handling, data integrity
   - Uses: In-memory SQLite database
   - Coverage: All customer database functions

2. **database/users_test.go** - User database operations
   - Validates: User creation, retrieval, duplicate handling
   - Tests: CreateUser, GetUserByEmail, GetUserByID

3. **handlers/auth_test.go** - Authentication handler tests
   - Validates: Registration, login, JWT generation, middleware
   - Tests: Register, Login, AuthMiddleware, GenerateToken, ParseToken
   - Mocks: None (uses real database)

4. **handlers/integration_test.go** - API integration tests
   - Validates: Complete workflows, authentication integration
   - Tests: Customer CRUD flow, plan creation, protected routes
   - Uses: Full database schema, authenticated requests

5. **handlers/plans_test.go** - Plan handler tests
   - Validates: Plan creation, validation, route retrieval
   - Tests: CreatePlan, GetPlan, DeletePlan, GetPlanRoutes

6. **optimizer/client_test.go** - Optimizer client tests
   - Validates: HTTP client behavior, timeout handling
   - Tests: HealthCheck, Optimize, request marshaling
   - Mocks: HTTP server (httptest)

### Optimizer Tests (Python)

**Location**: `optimizer/test_*.py`

**Framework**: pytest with FastAPI TestClient

**Key Test Files**:

1. **test_solver.py** - Solver algorithm unit tests
   - Test Classes: 9 test classes covering all solver functionality
   - Validates: Distance calculation, customer selection, inventory management, VRP solving
   - Uses: Mock objects (MockWarehouse, MockCustomer, MockVehicle)
   - Coverage: Core algorithm logic, edge cases, error scenarios

2. **test_api.py** - API integration tests
   - Test Classes: 5 test classes covering API endpoints
   - Validates: Endpoint behavior, request validation, error handling, performance
   - Uses: FastAPI TestClient, real solver (may be slow)
   - Coverage: All API endpoints, various scenarios

### Frontend Tests (JavaScript)

**Location**: `frontend/src/**/*.test.js`

**Framework**: Vitest with React Testing Library

**Key Test Files**:

1. **api.test.js** - API client tests
   - Validates: HTTP requests, error handling, token management
   - Mocks: Axios (vi.mock)
   - Coverage: All API methods, authentication, error scenarios

2. **components/__tests__/Modal.test.jsx** - Component tests
   - Validates: Component rendering, user interactions
   - Uses: React Testing Library, mocked framer-motion
   - Coverage: Modal visibility, close behavior, children rendering

## What Each Test Group Validates

### Backend Database Tests

**customers_test.go**:
- CreateCustomer: Validates customer creation with all fields, zero values, ID and timestamp setting
- GetCustomer: Tests retrieval of existing and non-existent customers, error handling
- ListCustomers: Validates listing with ordering, empty list handling
- UpdateCustomer: Tests field updates, UpdatedAt timestamp changes, non-existent customer handling
- DeleteCustomer: Validates deletion, verification of deletion, non-existent customer handling
- CountCustomers: Tests accurate counting after operations

**users_test.go**:
- CreateUser: Validates user creation, duplicate email detection (ErrDuplicate)
- GetUserByEmail: Tests email-based retrieval, non-existent user handling
- GetUserByID: Tests ID-based retrieval, error handling

### Backend Handler Tests

**auth_test.go**:
- Register: Validates successful registration, duplicate email (409), invalid input (400), password validation
- Login: Tests successful login, invalid credentials (401), missing fields (400)
- AuthMiddleware: Validates token validation, missing token (401), invalid token (401), malformed header (401)
- GenerateToken: Tests JWT generation, token structure, expiration setting
- ParseToken: Validates token parsing, invalid token handling, wrong secret handling

**integration_test.go**:
- CustomerCRUDIntegration: Complete workflow from create to delete, verifies each step
- PlanCreationFlow: Tests plan creation with warehouse dependency
- ProtectedRouteAccess: Validates authentication requirement on protected routes
- ErrorHandling: Tests error response structure and status codes

**plans_test.go**:
- CreatePlan: Validates plan creation, date format validation, end date before start date check
- GetPlan: Tests plan retrieval, invalid ID handling, non-existent plan (404)
- DeletePlan: Validates plan deletion, cascade delete verification
- GetPlanRoutes: Tests route retrieval with relationships

**optimizer/client_test.go**:
- HealthCheck: Tests health check endpoint, service unavailable handling
- Optimize: Validates optimization request, response handling, error scenarios
- OptimizeTimeout: Tests 5-minute timeout behavior
- OptimizeRequestMarshaling: Validates JSON serialization/deserialization

### Optimizer Solver Tests

**test_solver.py**:

**TestHaversineDistance**:
- Same point returns zero distance
- Known distance calculation (NYC to LA)
- Short distance accuracy

**TestDistanceMatrix**:
- Correct matrix dimensions
- Diagonal elements are zero
- Matrix symmetry
- Integer values (OR-Tools requirement)

**TestCustomerSelection**:
- Customers below minimum inventory selected
- 2-day threshold logic
- Priority and demand rate sorting
- Zero demand rate handling

**TestInventoryManagement**:
- Inventory initialization from customer data
- Daily consumption reduces inventory
- Inventory never goes negative
- Inventory increases after delivery

**TestDeliveryQuantity**:
- Fill-up policy implementation
- Maximum inventory constraint
- Zero quantity when at max

**TestVRPSolving**:
- Empty customer list handling
- Route generation for customers
- Route structure validation

**TestFallbackAlgorithm**:
- Route creation when OR-Tools fails
- Capacity constraint respect
- Valid route structure

**TestEndToEndSolver**:
- Complete optimization flow
- Empty customers handling
- Multi-day horizon
- Inventory tracking across days

**TestEdgeCases**:
- Zero planning horizon
- Single customer/vehicle
- Very large capacity
- Customers far apart

### Optimizer API Tests

**test_api.py**:

**TestHealthEndpoint**:
- Returns 200 OK
- Contains status, service, timestamp

**TestOptimizeEndpoint**:
- Valid request returns success
- No customers/vehicles returns error
- Invalid data returns 422
- Missing fields returns 422

**TestOptimizeScenarios**:
- Single day optimization
- Multi-day optimization
- Large planning horizon
- Multiple vehicles
- Capacity constraints

**TestErrorHandling**:
- Malformed JSON (422)
- Wrong HTTP method (405)
- Negative values handling
- Zero capacity vehicle

**TestPerformance**:
- Many customers (50) handled
- Concurrent requests (5) handled
- Response time acceptable

### Frontend Tests

**api.test.js**:

**Authentication**:
- Successful login stores token
- Login errors handled
- Registration creates user
- Registration errors handled

**Customer Management**:
- List, create, update, delete operations
- API errors handled
- Request structure validation

**Plan Management**:
- Plan creation
- Optimization request
- Route retrieval

**Token Management**:
- Token added to requests
- Missing token handled
- 401 responses handled

**Error Handling**:
- Network errors
- Timeout errors
- Invalid JSON responses

**Modal.test.jsx**:
- Renders when isOpen is true
- Does not render when isOpen is false
- Close button calls onClose
- Backdrop click calls onClose
- Children content rendered

## Running Specific Tests

### Backend
```bash
# Single test function
go test -run TestCreateCustomer ./internal/database/

# Single package
go test ./internal/database/

# With coverage
go test -cover ./internal/handlers/
```

### Optimizer
```bash
# Single test class
pytest test_solver.py::TestHaversineDistance

# Single test function
pytest test_solver.py::TestHaversineDistance::test_haversine_same_point

# Pattern match
pytest -k "haversine"

# Marked tests
pytest -m unit
pytest -m integration
```

### Frontend
```bash
# Single test file
npm test -- api.test.js

# Pattern match
npm test -- --grep "Authentication"

# Watch mode
npm test -- --watch
```

## Test Data

### Backend
- Uses in-memory SQLite databases
- Data reset between tests
- No external dependencies
- Real GORM operations

### Optimizer
- Uses mock objects (MockWarehouse, MockCustomer, MockVehicle)
- Isolated test data per test
- No database dependencies
- Real OR-Tools calls (may be slow)

### Frontend
- Mocks axios HTTP client
- Mocks localStorage
- No real network calls
- Isolated component tests

## Coverage Goals

- **Backend**: 80%+ for handlers and database
- **Optimizer**: 85%+ for solver, 70%+ for API
- **Frontend**: 70%+ for API client, 60%+ for components

## Continuous Integration

Tests should run on:
- Pull request creation
- Code push to main branch
- Pre-commit hooks (optional)

## Troubleshooting

### Backend Tests Fail
- Ensure Go 1.23+ installed
- Run `go mod tidy` to update dependencies
- Check database migrations run correctly

### Optimizer Tests Fail
- Ensure Python 3.11+ installed
- Install dependencies: `pip install -r requirements.txt`
- Install test dependencies: `pip install -r test_requirements.txt`
- Ensure OR-Tools installed: `pip install ortools`

### Frontend Tests Fail
- Ensure Node.js 20+ installed
- Run `npm install` to install dependencies
- Check Vitest configuration in vitest.config.js

## Test Best Practices

1. **Isolation**: Each test is independent, no shared state
2. **Naming**: Tests named descriptively (TestFunctionName_Scenario_ExpectedResult)
3. **Arrange-Act-Assert**: Clear test structure
4. **Edge Cases**: Test boundaries, zero values, empty lists
5. **Error Cases**: Test all error paths
6. **No Trivial Tests**: Don't test simple getters/setters unless critical

## Assumptions

1. **Backend**: GORM AutoMigrate works, SQLite sufficient for testing
2. **Optimizer**: OR-Tools functional, 30s timeout acceptable
3. **Frontend**: Axios mocking sufficient, localStorage available
4. **Integration**: Services can start independently

## Future Enhancements

1. E2E tests with Playwright/Cypress
2. Performance benchmarks
3. Load testing
4. Mutation testing
5. Visual regression testing (frontend)
6. API contract testing
