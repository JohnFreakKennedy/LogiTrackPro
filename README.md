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
- **Route Optimization**: IRP solver using heuristic algorithms (nearest neighbor + 2-opt)
- **Analytics Dashboard**: Overview of logistics operations

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend API | Go 1.21+ with Gin framework |
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
│       ├── database/        # Database layer & migrations
│       ├── handlers/        # HTTP request handlers
│       ├── models/          # Domain models
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

- Go 1.21+
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

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/logitrackpro?sslmode=disable"
export OPTIMIZER_URL="http://localhost:8000"
export JWT_SECRET="your-secret-key"
export PORT="8080"

# Run the backend
go run cmd/api/main.go
```

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

The IRP solver uses a heuristic approach:

1. **Inventory Projection**: Determine which customers need delivery based on projected inventory levels
2. **Customer Assignment**: Assign customers to vehicles using nearest neighbor heuristic
3. **Route Optimization**: Apply 2-opt improvement to reduce route distances
4. **Multi-day Planning**: Generate routes for each day in the planning horizon

The algorithm considers:
- Vehicle capacity constraints
- Maximum distance constraints
- Customer priority levels
- Demand rates and inventory levels
- Delivery costs (fixed + per-km)

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

Migrations run automatically on backend startup. The schema includes:
- `users` - User authentication
- `warehouses` - Distribution centers
- `customers` - Customer locations
- `vehicles` - Delivery vehicles
- `plans` - Delivery plans
- `routes` - Daily routes per plan
- `stops` - Route stops with delivery quantities

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

