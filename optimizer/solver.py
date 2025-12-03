"""
IRP (Inventory Routing Problem) Solver
Implements a heuristic algorithm for solving the inventory routing problem.
"""

import math
from datetime import datetime, timedelta
from typing import List, Dict, Tuple
from dataclasses import dataclass

import numpy as np


@dataclass
class StopResult:
    customer_id: int
    sequence: int
    quantity: float
    arrival_time: str


@dataclass
class RouteResult:
    day: int
    date: str
    vehicle_id: int
    total_distance: float
    total_cost: float
    total_load: float
    stops: List[StopResult]


@dataclass
class OptimizeResponse:
    success: bool
    message: str
    total_cost: float
    total_distance: float
    routes: List[RouteResult]


class IRPSolver:
    """
    Inventory Routing Problem Solver using a heuristic approach.
    
    The algorithm:
    1. For each day in the planning horizon:
       a. Determine which customers need delivery (inventory projection)
       b. Assign customers to vehicles using nearest neighbor insertion
       c. Optimize route for each vehicle using 2-opt improvement
    """
    
    def __init__(self, warehouse, customers, vehicles, planning_horizon, start_date):
        self.warehouse = warehouse
        self.customers = {c.id: c for c in customers}
        self.vehicles = {v.id: v for v in vehicles}
        self.planning_horizon = planning_horizon
        self.start_date = datetime.strptime(start_date, "%Y-%m-%d")
        
        # Build distance matrix
        self.locations = self._build_locations()
        self.distance_matrix = self._compute_distance_matrix()
        
        # Track customer inventory levels
        self.inventory = {c.id: c.current_inventory for c in customers}
    
    def _build_locations(self) -> Dict[int, Tuple[float, float]]:
        """Build location dictionary with warehouse as ID 0"""
        locations = {0: (self.warehouse.latitude, self.warehouse.longitude)}
        for cid, customer in self.customers.items():
            locations[cid] = (customer.latitude, customer.longitude)
        return locations
    
    def _compute_distance_matrix(self) -> Dict[Tuple[int, int], float]:
        """Compute haversine distances between all locations"""
        matrix = {}
        ids = list(self.locations.keys())
        
        for i in ids:
            for j in ids:
                if i == j:
                    matrix[(i, j)] = 0
                else:
                    matrix[(i, j)] = self._haversine(
                        self.locations[i][0], self.locations[i][1],
                        self.locations[j][0], self.locations[j][1]
                    )
        return matrix
    
    @staticmethod
    def _haversine(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
        """Calculate haversine distance in kilometers"""
        R = 6371  # Earth radius in km
        
        lat1, lon1, lat2, lon2 = map(math.radians, [lat1, lon1, lat2, lon2])
        dlat = lat2 - lat1
        dlon = lon2 - lon1
        
        a = math.sin(dlat/2)**2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon/2)**2
        c = 2 * math.asin(math.sqrt(a))
        
        return R * c
    
    def solve(self) -> OptimizeResponse:
        """Main solving method"""
        all_routes = []
        total_cost = 0
        total_distance = 0
        
        for day in range(self.planning_horizon):
            current_date = self.start_date + timedelta(days=day)
            
            # Determine customers needing delivery
            customers_to_visit = self._get_customers_needing_delivery(day)
            
            if not customers_to_visit:
                # Update inventory for next day (consume demand)
                self._update_inventory()
                continue
            
            # Create routes for this day
            day_routes = self._create_day_routes(day, current_date, customers_to_visit)
            
            for route in day_routes:
                all_routes.append(route)
                total_cost += route.total_cost
                total_distance += route.total_distance
            
            # Update inventory levels
            self._update_inventory()
        
        return OptimizeResponse(
            success=True,
            message=f"Optimization complete: {len(all_routes)} routes generated",
            total_cost=total_cost,
            total_distance=total_distance,
            routes=all_routes
        )
    
    def _get_customers_needing_delivery(self, day: int) -> List[int]:
        """
        Determine which customers need delivery based on inventory projections.
        A customer needs delivery if their projected inventory will drop below minimum.
        """
        customers_needing_delivery = []
        
        for cid, customer in self.customers.items():
            current_inv = self.inventory[cid]
            # Project inventory: days until stockout
            if customer.demand_rate > 0:
                days_until_stockout = (current_inv - customer.min_inventory) / customer.demand_rate
                
                # Deliver if we'll run out within 2 days or if inventory is below minimum
                if days_until_stockout <= 2 or current_inv <= customer.min_inventory:
                    customers_needing_delivery.append(cid)
            elif current_inv <= customer.min_inventory:
                customers_needing_delivery.append(cid)
        
        # Sort by priority (higher priority first)
        customers_needing_delivery.sort(
            key=lambda cid: (-self.customers[cid].priority, -self.customers[cid].demand_rate)
        )
        
        return customers_needing_delivery
    
    def _create_day_routes(self, day: int, date: datetime, 
                           customers_to_visit: List[int]) -> List[RouteResult]:
        """Create routes for a single day"""
        routes = []
        unassigned = set(customers_to_visit)
        vehicle_ids = list(self.vehicles.keys())
        vehicle_index = 0
        
        while unassigned and vehicle_index < len(vehicle_ids):
            vehicle_id = vehicle_ids[vehicle_index]
            vehicle = self.vehicles[vehicle_id]
            
            # Build route for this vehicle
            route_customers, deliveries = self._build_vehicle_route(
                vehicle, list(unassigned)
            )
            
            if route_customers:
                # Remove assigned customers
                for cid in route_customers:
                    unassigned.discard(cid)
                
                # Calculate route metrics
                route_distance = self._calculate_route_distance(route_customers)
                route_cost = vehicle.fixed_cost + (route_distance * vehicle.cost_per_km)
                total_load = sum(deliveries.values())
                
                # Apply 2-opt improvement
                route_customers = self._improve_route_2opt(route_customers)
                route_distance = self._calculate_route_distance(route_customers)
                route_cost = vehicle.fixed_cost + (route_distance * vehicle.cost_per_km)
                
                # Create stops
                stops = []
                current_time = datetime.combine(date.date(), datetime.min.time().replace(hour=8))
                avg_speed = 50  # km/h
                
                prev_loc = 0  # warehouse
                for seq, cid in enumerate(route_customers, 1):
                    # Calculate travel time
                    dist = self.distance_matrix[(prev_loc, cid)]
                    travel_time = timedelta(hours=dist / avg_speed)
                    current_time += travel_time
                    
                    stops.append(StopResult(
                        customer_id=cid,
                        sequence=seq,
                        quantity=deliveries[cid],
                        arrival_time=current_time.strftime("%H:%M")
                    ))
                    
                    # Update inventory for delivered customer
                    self.inventory[cid] += deliveries[cid]
                    
                    # Add service time (15 min per stop)
                    current_time += timedelta(minutes=15)
                    prev_loc = cid
                
                routes.append(RouteResult(
                    day=day + 1,
                    date=date.strftime("%Y-%m-%d"),
                    vehicle_id=vehicle_id,
                    total_distance=round(route_distance, 2),
                    total_cost=round(route_cost, 2),
                    total_load=round(total_load, 2),
                    stops=stops
                ))
            
            vehicle_index += 1
        
        return routes
    
    def _build_vehicle_route(self, vehicle, candidates: List[int]) -> Tuple[List[int], Dict[int, float]]:
        """
        Build a route for a vehicle using nearest neighbor heuristic.
        Returns list of customer IDs and delivery quantities.
        """
        route = []
        deliveries = {}
        remaining_capacity = vehicle.capacity
        remaining_distance = vehicle.max_distance if vehicle.max_distance > 0 else float('inf')
        current_location = 0  # warehouse
        
        available = set(candidates)
        
        while available and remaining_capacity > 0:
            # Find nearest customer we can serve
            best_customer = None
            best_distance = float('inf')
            
            for cid in available:
                customer = self.customers[cid]
                
                # Calculate delivery quantity
                delivery_qty = min(
                    customer.max_inventory - self.inventory[cid],
                    remaining_capacity,
                    customer.max_inventory  # Don't over-deliver
                )
                
                if delivery_qty <= 0:
                    available.discard(cid)
                    continue
                
                # Check if we can reach customer and return to warehouse
                dist_to_customer = self.distance_matrix[(current_location, cid)]
                dist_to_warehouse = self.distance_matrix[(cid, 0)]
                
                if dist_to_customer + dist_to_warehouse <= remaining_distance:
                    if dist_to_customer < best_distance:
                        best_distance = dist_to_customer
                        best_customer = cid
            
            if best_customer is None:
                break
            
            # Add customer to route
            customer = self.customers[best_customer]
            delivery_qty = min(
                customer.max_inventory - self.inventory[best_customer],
                remaining_capacity
            )
            
            route.append(best_customer)
            deliveries[best_customer] = delivery_qty
            remaining_capacity -= delivery_qty
            remaining_distance -= best_distance
            current_location = best_customer
            available.discard(best_customer)
        
        return route, deliveries
    
    def _calculate_route_distance(self, route: List[int]) -> float:
        """Calculate total distance for a route (warehouse -> customers -> warehouse)"""
        if not route:
            return 0
        
        distance = self.distance_matrix[(0, route[0])]  # warehouse to first
        
        for i in range(len(route) - 1):
            distance += self.distance_matrix[(route[i], route[i+1])]
        
        distance += self.distance_matrix[(route[-1], 0)]  # last to warehouse
        
        return distance
    
    def _improve_route_2opt(self, route: List[int]) -> List[int]:
        """Apply 2-opt improvement to reduce route distance"""
        if len(route) <= 2:
            return route
        
        improved = True
        best_route = route.copy()
        best_distance = self._calculate_route_distance(best_route)
        
        while improved:
            improved = False
            for i in range(len(best_route) - 1):
                for j in range(i + 2, len(best_route)):
                    # Create new route by reversing segment between i and j
                    new_route = (
                        best_route[:i+1] +
                        best_route[i+1:j+1][::-1] +
                        best_route[j+1:]
                    )
                    new_distance = self._calculate_route_distance(new_route)
                    
                    if new_distance < best_distance:
                        best_route = new_route
                        best_distance = new_distance
                        improved = True
                        break
                if improved:
                    break
        
        return best_route
    
    def _update_inventory(self):
        """Update inventory levels by consuming daily demand"""
        for cid, customer in self.customers.items():
            self.inventory[cid] = max(0, self.inventory[cid] - customer.demand_rate)

