"""
LogiTrackPro Optimizer Service
FastAPI-based IRP (Inventory Routing Problem) optimization service
"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime, timedelta
import logging

from solver import IRPSolver

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="LogiTrackPro Optimizer",
    description="Inventory Routing Problem optimization service",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# Request/Response Models
class WarehouseData(BaseModel):
    id: int
    latitude: float
    longitude: float
    stock: float


class CustomerData(BaseModel):
    id: int
    latitude: float
    longitude: float
    demand_rate: float
    max_inventory: float
    current_inventory: float
    min_inventory: float
    priority: int = 1


class VehicleData(BaseModel):
    id: int
    capacity: float
    cost_per_km: float
    fixed_cost: float
    max_distance: float


class OptimizeRequest(BaseModel):
    warehouse: WarehouseData
    customers: List[CustomerData]
    vehicles: List[VehicleData]
    planning_horizon: int
    start_date: str


class StopResult(BaseModel):
    customer_id: int
    sequence: int
    quantity: float
    arrival_time: str


class RouteResult(BaseModel):
    day: int
    date: str
    vehicle_id: int
    total_distance: float
    total_cost: float
    total_load: float
    stops: List[StopResult]


class OptimizeResponse(BaseModel):
    success: bool
    message: str
    total_cost: float
    total_distance: float
    routes: List[RouteResult]


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "LogiTrackPro Optimizer",
        "timestamp": datetime.now().isoformat()
    }


@app.post("/optimize", response_model=OptimizeResponse)
async def optimize(request: OptimizeRequest):
    """
    Optimize delivery routes for the given warehouse, customers, and vehicles
    over the specified planning horizon.
    """
    logger.info(f"Received optimization request: {len(request.customers)} customers, "
                f"{len(request.vehicles)} vehicles, {request.planning_horizon} days")
    
    try:
        # Validate input
        if not request.customers:
            return OptimizeResponse(
                success=False,
                message="No customers provided",
                total_cost=0,
                total_distance=0,
                routes=[]
            )
        
        if not request.vehicles:
            return OptimizeResponse(
                success=False,
                message="No vehicles provided",
                total_cost=0,
                total_distance=0,
                routes=[]
            )
        
        # Initialize solver
        solver = IRPSolver(
            warehouse=request.warehouse,
            customers=request.customers,
            vehicles=request.vehicles,
            planning_horizon=request.planning_horizon,
            start_date=request.start_date
        )
        
        # Run optimization
        result = solver.solve()
        
        logger.info(f"Optimization complete: {result.total_cost:.2f} cost, "
                    f"{result.total_distance:.2f} km, {len(result.routes)} routes")
        
        return result
        
    except Exception as e:
        logger.error(f"Optimization failed: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

