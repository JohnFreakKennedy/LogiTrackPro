# LogiTrackPro - Complete Project Context

**Last Updated:** December 2024  
**Status:** Active Development

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Tech Stack](#tech-stack)
4. [Project Structure](#project-structure)
5. [Key Features](#key-features)
6. [Recent Bug Fixes & Improvements](#recent-bug-fixes--improvements)
7. [Code Style & Standards](#code-style--standards)
8. [Database Schema](#database-schema)
9. [API Endpoints](#api-endpoints)
10. [Development Setup](#development-setup)
11. [Design System](#design-system)
12. [MCP Configuration](#mcp-configuration)
13. [Docker Configuration](#docker-configuration)
14. [Testing](#testing)
15. [Security Considerations](#security-considerations)

---

## Project Overview

**LogiTrackPro** is a comprehensive logistics planning platform designed to solve the **Inventory Routing Problem (IRP)**. The system manages warehouses, customers, vehicles, inventory levels, and generates optimized multi-day delivery plans.

### Core Purpose
- Solve complex inventory routing optimization problems
- Manage logistics operations through a modern web interface
- Provide analytics and insights for logistics planning
- Support multi-day delivery planning with route optimization

### Target Users
- Logistics managers
- Distribution center operators
- Supply chain planners
- Academic researchers (IRP problem domain)

---

## Architecture

LogiTrackPro follows a **modular monolith backend with one external microservice** architecture:

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

### Architecture Principles
- **Modular Monolith**: Backend organized into internal packages
- **Microservice**: Separate Python service for optimization algorithms
- **RESTful API**: Standard HTTP endpoints for all operations
- **JWT Authentication**: Stateless authentication mechanism
- **Database Migrations**: Automatic migrations on startup

---

## Tech Stack

| Component | Technology | Version |
|-----------|------------|---------|
| **Backend API** | Go | 1.21+ |
| **Backend Framework** | Gin | Latest |
| **Optimizer Service** | Python | 3.11+ |
| **Optimizer Framework** | FastAPI | Latest |
| **Frontend** | React | 18+ |
| **Frontend Build Tool** | Vite | Latest |
| **Styling** | Tailwind CSS | 3.4+ |
| **Animations** | Framer Motion | Latest |
| **Database** | PostgreSQL | 15+ |
| **Authentication** | JWT | Latest |
| **Containerization** | Docker & Docker Compose | Latest |

---

## Project Structure

```
LogiTrackPro/
├── backend/                  # Go backend API
│   ├── cmd/api/             # Application entry point
│   │   └── main.go          # Main server file
│   ├── Dockerfile           # Backend container definition
│   ├── go.mod               # Go dependencies
│   ├── go.sum               # Go dependency checksums
│   └── internal/            # Internal packages
│       ├── config/          # Configuration management
│       │   └── config.go    # Config loading & validation
│       ├── database/        # Database layer & migrations
│       │   ├── database.go  # Connection & migrations
│       │   ├── users.go     # User operations
│       │   ├── warehouses.go # Warehouse operations
│       │   ├── customers.go  # Customer operations
│       │   ├── vehicles.go   # Vehicle operations
│       │   ├── plans.go      # Plan operations
│       │   └── routes.go    # Route & stop operations
│       ├── handlers/        # HTTP request handlers
│       │   ├── handler.go   # Handler setup
│       │   ├── auth.go      # Authentication handlers
│       │   ├── warehouses.go
│       │   ├── customers.go
│       │   ├── vehicles.go
│       │   ├── plans.go     # Plan handlers (includes optimization)
│       │   └── analytics.go
│       ├── models/          # Domain models
│       │   └── models.go    # All data models
│       └── optimizer/       # Optimizer client
│           └── client.go    # HTTP client for optimizer service
│
├── optimizer/               # Python optimization service
│   ├── Dockerfile          # Optimizer container definition
│   ├── main.py             # FastAPI application
│   ├── solver.py           # IRP solver implementation
│   └── requirements.txt    # Python dependencies
│
├── frontend/               # React frontend
│   ├── Dockerfile          # Frontend container definition
│   ├── nginx.conf         # Nginx configuration for production
│   ├── package.json        # Node.js dependencies
│   ├── postcss.config.js   # PostCSS configuration
│   ├── tailwind.config.js  # Tailwind CSS configuration
│   ├── vite.config.js      # Vite configuration
│   ├── index.html          # HTML entry point
│   ├── public/            # Static assets
│   │   └── favicon.svg
│   └── src/
│       ├── main.jsx       # React entry point
│       ├── App.jsx        # Main application component
│       ├── api.js         # API client (Axios)
│       ├── index.css      # Global styles
│       ├── components/    # Reusable components
│       │   ├── Layout.jsx  # Main layout with sidebar
│       │   └── Modal.jsx  # Modal component
│       └── pages/         # Page components
│           ├── Login.jsx
│           ├── Dashboard.jsx
│           ├── Warehouses.jsx
│           ├── Customers.jsx
│           ├── Vehicles.jsx
│           ├── Plans.jsx
│           └── PlanDetail.jsx
│
├── scripts/                # Helper scripts
│   ├── db-helper.sh       # Database operations
│   ├── docker-helper.sh    # Docker Compose management
│   ├── mcp-setup.sh       # MCP configuration setup
│   ├── setup-mcp.sh       # Alternative MCP setup
│   └── test-mcp.sh        # MCP connectivity test
│
├── thesis/                # Thesis-related documentation
│   └── docs/              # Thesis documents
│
├── .cursorrules           # Cursor IDE project rules
├── docker-compose.yml     # Container orchestration
├── DESIGN_SPEC.md        # Design system specification
├── MCP_SETUP.md          # MCP server setup guide
├── mcp-config.example.json # MCP configuration example
├── README.md              # Main project README
└── PROJECT_CONTEXT.md    # This file
```

---

## Key Features

### 1. User Authentication
- JWT-based authentication system
- User registration and login
- Token refresh mechanism
- Secure password handling

### 2. Warehouse Management
- CRUD operations for distribution centers
- Location tracking (latitude/longitude)
- Capacity and inventory management
- Holding cost configuration

### 3. Customer Management
- Customer location tracking
- Demand rate configuration
- Inventory level monitoring
- Min/max inventory thresholds
- Priority levels

### 4. Vehicle Fleet Management
- Vehicle capacity configuration
- Cost parameters (fixed + per-km)
- Maximum distance constraints
- Warehouse assignment
- Availability status

### 5. Delivery Planning
- Multi-day delivery plan creation
- Plan status tracking (draft, optimizing, completed, failed)
- Date range configuration
- Warehouse association

### 6. Route Optimization
- IRP solver using heuristic algorithms
- Nearest neighbor + 2-opt improvement
- Multi-day route generation
- Vehicle capacity constraints
- Distance optimization

### 7. Analytics Dashboard
- Overview of logistics operations
- Statistics and metrics
- Performance indicators
- Recent plans summary

---

## Recent Bug Fixes & Improvements

### 1. Database NULL Handling (routes.go)
**Issue**: `GetRoutesByPlan` and `GetStopsByRoute` functions used LEFT JOINs but scanned directly into non-pointer Go types. When vehicles or customers were deleted (ON DELETE SET NULL), scanning NULL values caused scan errors.

**Fix**: 
- Modified to use `sql.NullInt64`, `sql.NullString`, `sql.NullFloat64`, etc.
- Only populate structs when IDs are valid (not NULL)
- Prevents crashes when associated entities are deleted

**Files**: `backend/internal/database/routes.go`

### 2. Python Optimizer RuntimeError Fix
**Issue**: In `solver.py`, modifying a set (`available.discard(cid)`) during iteration caused `RuntimeError: Set changed size during iteration`.

**Fix**: 
- Collect items to remove in a separate list
- Perform discard operations after iteration completes

**Files**: `optimizer/solver.py`

### 3. JWT Secret Security Enhancement
**Issue**: Hardcoded insecure default JWT secret values in `config.go` and `docker-compose.yml`.

**Fix**:
- Added security warnings for insecure defaults
- Enforced secure JWT secret in production (fatal error if insecure)
- Removed hardcoded secrets from docker-compose.yml
- Added list of known insecure defaults to check against

**Files**: 
- `backend/internal/config/config.go`
- `docker-compose.yml`

### 4. Database Transaction Atomicity
**Issue**: `OptimizePlan` function deleted routes and created new ones without transactions. Partial failures left database in inconsistent state.

**Fix**:
- Wrapped entire optimization process in database transaction
- Added transaction-aware versions: `CreateRouteTx`, `DeleteRoutesByPlanTx`, `CreateStopTx`, `UpdatePlanStatusTx`
- Proper rollback on any failure
- Ensures atomicity of route creation

**Files**: 
- `backend/internal/handlers/plans.go`
- `backend/internal/database/routes.go`

### 5. IDE Configuration Cleanup
**Issue**: `.idea/` directory (JetBrains IDE config) was committed to repository despite `.gitignore`.

**Fix**: Removed `.idea/` directory from Git tracking using `git rm -r --cached .idea/`

---

## Code Style & Standards

### Go Backend
- Use `gofmt` for formatting
- Follow Go naming conventions:
  - Exported: PascalCase
  - Private: camelCase
- Error handling: Always check and return errors explicitly
- Use context for request cancellation
- Database: Use prepared statements, handle NULL values with pointers (`sql.Null*` types)

### Python Optimizer
- Follow PEP 8 style guide
- Use type hints (Pydantic models)
- Async/await for I/O operations
- Docstrings for public functions

### React Frontend
- Functional components with hooks
- Use Tailwind CSS for styling
- Axios for API calls
- Framer Motion for animations
- Component structure: `pages/` and `components/`

---

## Database Schema

### Tables

1. **users**
   - User authentication and profiles

2. **warehouses**
   - Distribution centers
   - Fields: name, address, lat, lng, capacity, stock, holding_cost, replenishment_qty

3. **customers**
   - Customer locations and inventory
   - Fields: name, address, lat, lng, demand_rate, max_inventory, min_inventory, current_inventory, holding_cost, priority

4. **vehicles**
   - Delivery vehicles
   - Fields: name, capacity, cost_per_km, fixed_cost, max_distance, warehouse_id, available

5. **plans**
   - Delivery plans
   - Fields: name, start_date, end_date, warehouse_id, status, total_cost, total_distance

6. **routes**
   - Daily routes per plan
   - Fields: plan_id, day, vehicle_id, distance, cost, load
   - Foreign keys: `vehicle_id` (ON DELETE SET NULL)

7. **stops**
   - Route stops with delivery quantities
   - Fields: route_id, sequence, customer_id, delivery_qty
   - Foreign keys: `customer_id` (ON DELETE SET NULL)

### Key Relationships
- Routes → Plans (many-to-one)
- Routes → Vehicles (many-to-one, nullable)
- Stops → Routes (many-to-one)
- Stops → Customers (many-to-one, nullable)
- Plans → Warehouses (many-to-one)

### Migration Strategy
- Migrations run automatically on backend startup
- Defined in `backend/internal/database/database.go`
- Uses raw SQL for schema creation

---

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

---

## Development Setup

### Prerequisites
- Go 1.21+
- Python 3.11+
- Node.js 20+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Option 1: Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# Access services
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Optimizer: http://localhost:8000
# PostgreSQL: localhost:5432
```

### Option 2: Manual Setup

#### Database Setup
```bash
createdb logitrackpro
```

#### Backend Setup
```bash
cd backend
go mod download
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/logitrackpro?sslmode=disable"
export OPTIMIZER_URL="http://localhost:8000"
export JWT_SECRET="your-secure-secret-key"
export PORT="8080"
go run cmd/api/main.go
```

#### Optimizer Setup
```bash
cd optimizer
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000
```

#### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

---

## Design System

### Color Palette
- **Primary Green**: #22c55e (optimization, success, efficiency)
- **Accent Orange**: #f97316 (important actions, alerts)
- **Dark Theme**: #0f172a (main background), #1e293b (cards)
- **Text**: #ffffff (primary), #94a3b8 (secondary)

### Typography
- **Body**: DM Sans (400, 500, 600, 700)
- **Headings**: Outfit (400, 500, 600, 700, 800)
- **Data/Monospace**: JetBrains Mono (400, 500)

### Component Library
- Buttons (Primary, Secondary, Accent, Danger)
- Cards (Standard, Interactive)
- Forms (Inputs, Selects, Textareas)
- Tables
- Badges (Success, Warning, Info, Danger)
- Modals
- Navigation (Sidebar, Mobile Menu)

See `DESIGN_SPEC.md` for complete design system documentation.

---

## MCP Configuration

LogiTrackPro includes MCP (Model Context Protocol) server configuration for enhanced AI assistance.

### Configured MCP Servers
1. **PostgreSQL MCP** - Database operations and queries
2. **Filesystem MCP** - File operations and navigation
3. **Git MCP** - Version control operations
4. **Docker MCP** - Container management
5. **Brave Search MCP** - Web search (optional, requires API key)

### Setup
```bash
./scripts/setup-mcp.sh
./scripts/test-mcp.sh
```

See `MCP_SETUP.md` for detailed documentation.

---

## Docker Configuration

### Services
1. **backend** - Go API server (port 8080)
2. **optimizer** - Python optimization service (port 8000)
3. **frontend** - React application (port 3000)
4. **postgres** - PostgreSQL database (port 5432)

### Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
- `OPTIMIZER_URL` - Optimizer service URL
- `JWT_SECRET` - Secret key for JWT signing (MUST be set securely)
- `PORT` - Backend server port

### Health Checks
- All services include health check endpoints
- Docker Compose monitors service health

---

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Optimizer Tests
```bash
cd optimizer
pytest
```

### Frontend Tests
```bash
cd frontend
npm test
```

---

## Security Considerations

### Authentication
- JWT tokens with expiration (default: 24 hours)
- Secure secret key required (no insecure defaults in production)
- Password hashing (bcrypt)

### Database
- Parameterized queries (SQL injection prevention)
- NULL value handling for foreign keys
- Transaction atomicity for critical operations

### API
- Input validation on all endpoints
- CORS whitelist for production
- Error messages don't expose sensitive information

### Best Practices
- Never commit secrets (use .env files)
- Validate all user inputs
- Use prepared statements
- Handle NULL values properly
- Use transactions for multi-step operations

---

## Optimization Algorithm

The IRP solver uses a heuristic approach:

1. **Inventory Projection**: Determine which customers need delivery based on projected inventory levels
2. **Customer Assignment**: Assign customers to vehicles using nearest neighbor heuristic
3. **Route Optimization**: Apply 2-opt improvement to reduce route distances
4. **Multi-day Planning**: Generate routes for each day in the planning horizon

### Algorithm Considerations
- Vehicle capacity constraints
- Maximum distance constraints
- Customer priority levels
- Demand rates and inventory levels
- Delivery costs (fixed + per-km)

---

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Backend server port | `8080` | No |
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `OPTIMIZER_URL` | Optimizer service URL | `http://localhost:8000` | No |
| `JWT_SECRET` | Secret key for JWT signing | - | Yes |
| `JWT_EXPIRY_HOURS` | Token expiration time | `24` | No |

---

## Future Enhancements

### Potential Features
- Real-time route tracking
- Map visualization (Mapbox/Leaflet)
- Advanced analytics and reporting
- Export functionality (PDF, Excel)
- Multi-warehouse support
- Vehicle tracking integration
- Mobile app (React Native)

### Algorithm Improvements
- More sophisticated optimization algorithms
- Machine learning for demand prediction
- Dynamic route adjustment
- Multi-objective optimization

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow code style guidelines
4. Write tests for new features
5. Commit your changes
6. Push to the branch
7. Create a Pull Request

---

## License

This project is developed for academic purposes as part of logistics optimization research.

---

## Additional Resources

- **Design Specification**: See `DESIGN_SPEC.md`
- **MCP Setup**: See `MCP_SETUP.md`
- **Main README**: See `README.md`
- **Project Rules**: See `.cursorrules`

---

**Document Version**: 1.0  
**Last Updated**: December 2024  
**Maintained By**: LogiTrackPro Development Team



