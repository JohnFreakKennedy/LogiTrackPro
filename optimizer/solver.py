"""
IRP (Inventory Routing Problem) Solver using Google OR-Tools
Implements a VRP-based approach for solving the inventory routing problem.
"""

import math
from datetime import datetime, timedelta
from typing import List, Dict, Tuple, Optional
from dataclasses import dataclass

from ortools.constraint_solver import routing_enums_pb2
from ortools.constraint_solver import pywrapcp


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
    Inventory Routing Problem Solver using Google OR-Tools.
    
    The algorithm:
    1. For each day in the planning horizon:
       a. Determine which customers need delivery (inventory projection)
       b. Solve VRP using OR-Tools for that day
       c. Update inventory levels
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
    
    def _compute_distance_matrix(self) -> List[List[int]]:
        """
        Compute distance matrix as integers (OR-Tools requires integers).
        Returns matrix where [i][j] is distance from location i to j in meters.
        """
        ids = sorted(self.locations.keys())
        n = len(ids)
        matrix = [[0] * n for _ in range(n)]
        
        for i, id_i in enumerate(ids):
            for j, id_j in enumerate(ids):
                if i != j:
                    # Calculate haversine distance in meters
                    dist_km = self._haversine(
                        self.locations[id_i][0], self.locations[id_i][1],
                        self.locations[id_j][0], self.locations[id_j][1]
                    )
                    # Convert to meters and round to integer
                    matrix[i][j] = int(dist_km * 1000)
        
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
            
            # Solve VRP for this day using OR-Tools
            day_routes = self._solve_day_vrp(day, current_date, customers_to_visit)
            
            for route in day_routes:
                all_routes.append(route)
                total_cost += route.total_cost
                total_distance += route.total_distance
                
                # Apply deliveries for next day planning
                for stop in route.stops:
                    self.inventory[stop.customer_id] += stop.quantity
            
            # Update inventory levels
            self._update_inventory()
        
        return OptimizeResponse(
            success=True,
            message=f"Optimization complete: {len(all_routes)} routes generated using OR-Tools",
            total_cost=round(total_cost, 2),
            total_distance=round(total_distance, 2),
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
    
    def _solve_day_vrp(self, day: int, date: datetime, 
                       customers_to_visit: List[int]) -> List[RouteResult]:
        """
        Solve Vehicle Routing Problem for a single day using OR-Tools.
        """
        if not customers_to_visit:
            return []
        
        # Map customer IDs to matrix indices (0 = warehouse, 1+ = customers)
        customer_id_to_index = {0: 0}  # Warehouse is always index 0
        index_to_customer_id = {0: 0}
        
        for idx, cid in enumerate(customers_to_visit, start=1):
            customer_id_to_index[cid] = idx
            index_to_customer_id[idx] = cid
        
        num_locations = len(customers_to_visit) + 1  # +1 for warehouse
        num_vehicles = len(self.vehicles)
        
        # Create distance callback
        def distance_callback(from_index, to_index):
            """Returns the distance between the two nodes."""
            from_node = index_to_customer_id[from_index]
            to_node = index_to_customer_id[to_index]
            
            # Get original IDs for distance lookup
            from_id = from_node if from_node != 0 else 0
            to_id = to_node if to_node != 0 else 0
            
            # Find indices in original location list
            all_ids = sorted(self.locations.keys())
            from_idx = all_ids.index(from_id)
            to_idx = all_ids.index(to_id)
            
            return self.distance_matrix[from_idx][to_idx]
        
        # Create demand callback (delivery quantities)
        def demand_callback(from_index):
            """Returns the demand at a node."""
            if from_index == 0:  # Warehouse
                return 0
            
            cid = index_to_customer_id[from_index]
            customer = self.customers[cid]
            
            # Calculate delivery quantity needed
            delivery_qty = min(
                customer.max_inventory - self.inventory[cid],
                customer.max_inventory  # Don't exceed max
            )
            
            # Convert to integer (OR-Tools requires integers)
            # Use grams as unit to maintain precision
            return int(delivery_qty * 1000)
        
        # Create routing model
        manager = pywrapcp.RoutingIndexManager(
            num_locations, num_vehicles, 0  # depot_index = 0 (warehouse)
        )
        routing = pywrapcp.RoutingModel(manager)
        
        # Register callbacks
        transit_callback_index = routing.RegisterTransitCallback(distance_callback)
        routing.SetArcCostEvaluatorOfAllVehicles(transit_callback_index)
        
        # Add capacity constraint
        demand_callback_index = routing.RegisterUnaryCallback(demand_callback)
        
        # Vehicle capacities as a list (in grams for precision)
        vehicle_capacities = []
        vehicle_ids_list = list(self.vehicles.keys())
        for vehicle_index in range(num_vehicles):
            vehicle_id = vehicle_ids_list[vehicle_index]
            vehicle = self.vehicles[vehicle_id]
            # Convert to grams
            vehicle_capacities.append(int(vehicle.capacity * 1000))
        
        routing.AddDimensionWithVehicleCapacity(
            demand_callback_index,
            0,  # null capacity slack
            vehicle_capacities,  # vehicle capacities list
            True,  # start cumul to zero
            'Capacity'
        )
        
        # Add distance constraint (max distance per vehicle)
        dimension_name = 'Distance'
        routing.AddDimension(
            transit_callback_index,
            0,  # no slack
            300000000,  # vehicle maximum travel distance (300,000 km in meters)
            True,  # start cumul to zero
            dimension_name
        )
        distance_dimension = routing.GetDimensionOrDie(dimension_name)
        distance_dimension.SetGlobalSpanCostCoefficient(100)
        
        # Set max distance per vehicle if specified
        for vehicle_index in range(num_vehicles):
            vehicle_id = list(self.vehicles.keys())[vehicle_index]
            vehicle = self.vehicles[vehicle_id]
            if vehicle.max_distance > 0:
                max_dist_meters = int(vehicle.max_distance * 1000)
                distance_dimension.CumulVar(
                    manager.Start(vehicle_index)
                ).SetMax(max_dist_meters)
        
        # Set search parameters
        search_parameters = pywrapcp.DefaultRoutingSearchParameters()
        search_parameters.first_solution_strategy = (
            routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
        )
        search_parameters.local_search_metaheuristic = (
            routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH
        )
        search_parameters.time_limit.seconds = 30  # 30 second time limit per day
        search_parameters.log_search = False
        
        # Solve
        solution = routing.SolveWithParameters(search_parameters)
        
        if not solution:
            # Fallback: create simple routes if OR-Tools fails
            return self._create_fallback_routes(day, date, customers_to_visit)
        
        # Extract routes from solution
        routes = []
        vehicle_ids = list(self.vehicles.keys())
        
        for vehicle_index in range(num_vehicles):
            vehicle_id = vehicle_ids[vehicle_index]
            vehicle = self.vehicles[vehicle_id]
            
            route_customers = []
            route_deliveries = {}
            index = routing.Start(vehicle_index)
            route_distance = 0
            prev_index = index
            
            while not routing.IsEnd(index):
                node_index = manager.IndexToNode(index)
                if node_index != 0:  # Not warehouse
                    cid = index_to_customer_id[node_index]
                    route_customers.append(cid)
                    
                    # Get delivery quantity from demand callback
                    demand = demand_callback(node_index)
                    delivery_qty = demand / 1000.0  # Convert back from grams
                    route_deliveries[cid] = delivery_qty
                
                # Calculate distance
                next_index = solution.Value(routing.NextVar(index))
                route_distance += distance_callback(prev_index, next_index)
                prev_index = next_index
                index = next_index
            
            if route_customers:
                # Convert distance from meters to km
                route_distance_km = route_distance / 1000.0
                route_cost = vehicle.fixed_cost + (route_distance_km * vehicle.cost_per_km)
                total_load = sum(route_deliveries.values())
                
                # Create stops with arrival times
                stops = []
                current_time = datetime.combine(date.date(), datetime.min.time().replace(hour=8))
                avg_speed = 50  # km/h
                
                prev_loc = 0  # warehouse
                for seq, cid in enumerate(route_customers, 1):
                    # Calculate travel time
                    all_ids = sorted(self.locations.keys())
                    prev_idx = all_ids.index(prev_loc)
                    curr_idx = all_ids.index(cid)
                    dist_km = self.distance_matrix[prev_idx][curr_idx] / 1000.0
                    travel_time = timedelta(hours=dist_km / avg_speed)
                    current_time += travel_time
                    
                    stops.append(StopResult(
                        customer_id=cid,
                        sequence=seq,
                        quantity=round(route_deliveries[cid], 2),
                        arrival_time=current_time.strftime("%H:%M")
                    ))
                    
                    # Add service time (15 min per stop)
                    current_time += timedelta(minutes=15)
                    prev_loc = cid
                
                # Add return to warehouse distance
                if route_customers:
                    last_cid = route_customers[-1]
                    all_ids = sorted(self.locations.keys())
                    last_idx = all_ids.index(last_cid)
                    return_dist_km = self.distance_matrix[last_idx][0] / 1000.0
                    route_distance_km += return_dist_km
                    route_cost = vehicle.fixed_cost + (route_distance_km * vehicle.cost_per_km)
                
                routes.append(RouteResult(
                    day=day + 1,
                    date=date.strftime("%Y-%m-%d"),
                    vehicle_id=vehicle_id,
                    total_distance=round(route_distance_km, 2),
                    total_cost=round(route_cost, 2),
                    total_load=round(total_load, 2),
                    stops=stops
                ))
        
        return routes
    
    def _create_fallback_routes(self, day: int, date: datetime, 
                                customers_to_visit: List[int]) -> List[RouteResult]:
        """
        Fallback route creation if OR-Tools fails.
        Uses simple nearest neighbor approach.
        """
        routes = []
        unassigned = set(customers_to_visit)
        vehicle_ids = list(self.vehicles.keys())
        vehicle_index = 0
        
        while unassigned and vehicle_index < len(vehicle_ids):
            vehicle_id = vehicle_ids[vehicle_index]
            vehicle = self.vehicles[vehicle_id]
            
            # Simple nearest neighbor route
            route_customers = []
            route_deliveries = {}
            current_location = 0  # warehouse
            remaining_capacity = vehicle.capacity
            
            while unassigned and remaining_capacity > 0:
                # Find nearest unassigned customer
                best_customer = None
                best_distance = float('inf')
                
                all_ids = sorted(self.locations.keys())
                current_idx = all_ids.index(current_location)
                
                for cid in list(unassigned):
                    customer = self.customers[cid]
                    delivery_qty = min(
                        customer.max_inventory - self.inventory[cid],
                        remaining_capacity,
                        customer.max_inventory
                    )
                    
                    if delivery_qty <= 0:
                        continue
                    
                    cid_idx = all_ids.index(cid)
                    dist = self.distance_matrix[current_idx][cid_idx]
                    
                    if dist < best_distance:
                        best_distance = dist
                        best_customer = cid
                
                if best_customer is None:
                    break
                
                customer = self.customers[best_customer]
                delivery_qty = min(
                    customer.max_inventory - self.inventory[best_customer],
                    remaining_capacity
                )
                
                route_customers.append(best_customer)
                route_deliveries[best_customer] = delivery_qty
                remaining_capacity -= delivery_qty
                current_location = best_customer
                unassigned.discard(best_customer)
            
            if route_customers:
                # Calculate route metrics
                all_ids = sorted(self.locations.keys())
                route_distance = 0
                
                # Warehouse to first
                route_distance += self.distance_matrix[0][all_ids.index(route_customers[0])]
                
                # Between customers
                for i in range(len(route_customers) - 1):
                    route_distance += self.distance_matrix[
                        all_ids.index(route_customers[i])
                    ][all_ids.index(route_customers[i+1])]
                
                # Last to warehouse
                route_distance += self.distance_matrix[all_ids.index(route_customers[-1])][0]
                
                route_distance_km = route_distance / 1000.0
                route_cost = vehicle.fixed_cost + (route_distance_km * vehicle.cost_per_km)
                total_load = sum(route_deliveries.values())
                
                # Create stops
                stops = []
                current_time = datetime.combine(date.date(), datetime.min.time().replace(hour=8))
                avg_speed = 50  # km/h
                
                prev_loc = 0
                for seq, cid in enumerate(route_customers, 1):
                    prev_idx = all_ids.index(prev_loc)
                    curr_idx = all_ids.index(cid)
                    dist_km = self.distance_matrix[prev_idx][curr_idx] / 1000.0
                    travel_time = timedelta(hours=dist_km / avg_speed)
                    current_time += travel_time
                    
                    stops.append(StopResult(
                        customer_id=cid,
                        sequence=seq,
                        quantity=round(route_deliveries[cid], 2),
                        arrival_time=current_time.strftime("%H:%M")
                    ))
                    
                    current_time += timedelta(minutes=15)
                    prev_loc = cid
                
                routes.append(RouteResult(
                    day=day + 1,
                    date=date.strftime("%Y-%m-%d"),
                    vehicle_id=vehicle_id,
                    total_distance=round(route_distance_km, 2),
                    total_cost=round(route_cost, 2),
                    total_load=round(total_load, 2),
                    stops=stops
                ))
            
            vehicle_index += 1
        
        return routes
    
    def _update_inventory(self):
        """Update inventory levels by consuming daily demand"""
        for cid, customer in self.customers.items():
            self.inventory[cid] = max(0, self.inventory[cid] - customer.demand_rate)
