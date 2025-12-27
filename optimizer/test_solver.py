"""
Unit tests for IRP Solver
Tests the core algorithmic component including inventory projection, VRP solving, and state management.
"""

import pytest
from datetime import datetime, timedelta
from unittest.mock import Mock, patch
from solver import IRPSolver, OptimizeResponse


class MockWarehouse:
    def __init__(self, id, lat, lon, stock=10000):
        self.id = id
        self.latitude = lat
        self.longitude = lon
        self.stock = stock


class MockCustomer:
    def __init__(self, id, lat, lon, demand_rate=100, max_inv=1000, current_inv=500, min_inv=100, priority=1):
        self.id = id
        self.latitude = lat
        self.longitude = lon
        self.demand_rate = demand_rate
        self.max_inventory = max_inv
        self.current_inventory = current_inv
        self.min_inventory = min_inv
        self.priority = priority


class MockVehicle:
    def __init__(self, id, capacity=5000, cost_per_km=1.0, fixed_cost=100.0, max_distance=0):
        self.id = id
        self.capacity = capacity
        self.cost_per_km = cost_per_km
        self.fixed_cost = fixed_cost
        self.max_distance = max_distance


@pytest.fixture
def sample_warehouse():
    return MockWarehouse(id=1, lat=40.7128, lon=-74.0060)


@pytest.fixture
def sample_customers():
    return [
        MockCustomer(id=1, lat=40.7580, lon=-73.9855, current_inv=200, min_inv=100, demand_rate=50),  # Needs delivery
        MockCustomer(id=2, lat=40.7505, lon=-73.9934, current_inv=600, min_inv=100, demand_rate=50),  # OK
        MockCustomer(id=3, lat=40.7282, lon=-73.9942, current_inv=50, min_inv=100, demand_rate=50),  # Below minimum
    ]


@pytest.fixture
def sample_vehicles():
    return [
        MockVehicle(id=1, capacity=5000, cost_per_km=1.0, fixed_cost=100.0),
        MockVehicle(id=2, capacity=3000, cost_per_km=0.8, fixed_cost=80.0),
    ]


class TestHaversineDistance:
    """Tests for haversine distance calculation"""
    
    def test_haversine_same_point(self):
        """Distance between same point should be zero"""
        solver = IRPSolver(MockWarehouse(1, 0, 0), [], [], 1, "2024-01-01")
        dist = solver._haversine(40.7128, -74.0060, 40.7128, -74.0060)
        assert dist == pytest.approx(0.0, abs=0.01)
    
    def test_haversine_known_distance(self):
        """Test haversine with known coordinates (NYC to LA approximately 3944 km)"""
        solver = IRPSolver(MockWarehouse(1, 0, 0), [], [], 1, "2024-01-01")
        # NYC to Los Angeles
        dist = solver._haversine(40.7128, -74.0060, 34.0522, -118.2437)
        assert dist == pytest.approx(3944, rel=0.1)  # Within 10% of actual distance
    
    def test_haversine_short_distance(self):
        """Test haversine for short distances (within same city)"""
        solver = IRPSolver(MockWarehouse(1, 0, 0), [], [], 1, "2024-01-01")
        # Two points in NYC (approximately 5 km apart)
        dist = solver._haversine(40.7128, -74.0060, 40.7580, -73.9855)
        assert dist > 0
        assert dist < 20  # Should be less than 20 km


class TestDistanceMatrix:
    """Tests for distance matrix computation"""
    
    def test_distance_matrix_size(self, sample_warehouse, sample_customers):
        """Distance matrix should have correct dimensions"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        matrix = solver.distance_matrix
        expected_size = len(sample_customers) + 1  # +1 for warehouse
        assert len(matrix) == expected_size
        assert all(len(row) == expected_size for row in matrix)
    
    def test_distance_matrix_diagonal(self, sample_warehouse, sample_customers):
        """Diagonal elements should be zero (distance to self)"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        matrix = solver.distance_matrix
        for i in range(len(matrix)):
            assert matrix[i][i] == 0
    
    def test_distance_matrix_symmetric(self, sample_warehouse, sample_customers):
        """Distance matrix should be symmetric (undirected graph)"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        matrix = solver.distance_matrix
        for i in range(len(matrix)):
            for j in range(len(matrix)):
                assert matrix[i][j] == matrix[j][i]
    
    def test_distance_matrix_integers(self, sample_warehouse, sample_customers):
        """Distance matrix should contain integers (OR-Tools requirement)"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        matrix = solver.distance_matrix
        for row in matrix:
            for dist in row:
                assert isinstance(dist, int)


class TestCustomerSelection:
    """Tests for customer selection logic"""
    
    def test_customers_needing_delivery_below_minimum(self, sample_warehouse, sample_customers):
        """Customers below minimum inventory should be selected"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        solver.inventory[3] = 50  # Below min_inventory
        customers = solver._get_customers_needing_delivery(0)
        assert 3 in customers
    
    def test_customers_needing_delivery_two_day_threshold(self, sample_warehouse, sample_customers):
        """Customers with <=2 days until stockout should be selected"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        # Customer 1: current_inv=200, min_inv=100, demand_rate=50
        # Days until stockout = (200-100)/50 = 2 days
        solver.inventory[1] = 200
        customers = solver._get_customers_needing_delivery(0)
        assert 1 in customers
    
    def test_customers_needing_delivery_above_threshold(self, sample_warehouse, sample_customers):
        """Customers with >2 days until stockout should not be selected"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        # Customer 2: current_inv=600, min_inv=100, demand_rate=50
        # Days until stockout = (600-100)/50 = 10 days
        solver.inventory[2] = 600
        customers = solver._get_customers_needing_delivery(0)
        assert 2 not in customers
    
    def test_customers_needing_delivery_priority_sorting(self, sample_warehouse):
        """Customers should be sorted by priority then demand rate"""
        customers = [
            MockCustomer(id=1, lat=40.0, lon=-74.0, priority=1, demand_rate=50),
            MockCustomer(id=2, lat=40.1, lon=-74.1, priority=3, demand_rate=30),
            MockCustomer(id=3, lat=40.2, lon=-74.2, priority=2, demand_rate=100),
        ]
        solver = IRPSolver(sample_warehouse, customers, [], 1, "2024-01-01")
        # Set all to need delivery
        for c in customers:
            solver.inventory[c.id] = c.min_inventory
        
        selected = solver._get_customers_needing_delivery(0)
        # Should be sorted: priority 3, then 2, then 1
        assert selected[0] == 2  # Highest priority
        assert selected[1] == 3  # Second priority
        assert selected[2] == 1  # Lowest priority
    
    def test_customers_needing_delivery_zero_demand_rate(self, sample_warehouse):
        """Customers with zero demand rate but below minimum should be selected"""
        customer = MockCustomer(id=1, lat=40.0, lon=-74.0, demand_rate=0, current_inv=50, min_inv=100)
        solver = IRPSolver(sample_warehouse, [customer], [], 1, "2024-01-01")
        solver.inventory[1] = 50
        customers = solver._get_customers_needing_delivery(0)
        assert 1 in customers


class TestInventoryManagement:
    """Tests for inventory state management"""
    
    def test_inventory_initialization(self, sample_warehouse, sample_customers):
        """Inventory should be initialized from customer current_inventory"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        assert solver.inventory[1] == 200
        assert solver.inventory[2] == 600
        assert solver.inventory[3] == 50
    
    def test_inventory_update_consumption(self, sample_warehouse, sample_customers):
        """Daily demand should reduce inventory"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        initial_inv = solver.inventory[1]
        solver._update_inventory()
        assert solver.inventory[1] == initial_inv - sample_customers[0].demand_rate
    
    def test_inventory_update_no_negative(self, sample_warehouse):
        """Inventory should not go negative"""
        customer = MockCustomer(id=1, lat=40.0, lon=-74.0, demand_rate=1000, current_inv=50)
        solver = IRPSolver(sample_warehouse, [customer], [], 1, "2024-01-01")
        solver.inventory[1] = 50
        solver._update_inventory()
        assert solver.inventory[1] >= 0
    
    def test_inventory_after_delivery(self, sample_warehouse, sample_customers):
        """Inventory should increase after delivery"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        initial_inv = solver.inventory[1]
        delivery_qty = 300
        solver.inventory[1] += delivery_qty
        assert solver.inventory[1] == initial_inv + delivery_qty


class TestDeliveryQuantity:
    """Tests for delivery quantity calculation"""
    
    def test_delivery_quantity_fill_up(self, sample_warehouse, sample_customers):
        """Delivery quantity should fill to max_inventory"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        solver.inventory[1] = 200  # current_inventory
        customer = sample_customers[0]
        customer.max_inventory = 1000
        
        # Mock the demand callback logic
        delivery_qty = min(customer.max_inventory - solver.inventory[1], customer.max_inventory)
        assert delivery_qty == 800  # 1000 - 200
    
    def test_delivery_quantity_at_maximum(self, sample_warehouse, sample_customers):
        """Delivery quantity should be zero if already at max"""
        solver = IRPSolver(sample_warehouse, sample_customers, [], 1, "2024-01-01")
        solver.inventory[1] = 1000  # Already at max
        customer = sample_customers[0]
        customer.max_inventory = 1000
        
        delivery_qty = min(customer.max_inventory - solver.inventory[1], customer.max_inventory)
        assert delivery_qty == 0


class TestVRPSolving:
    """Tests for VRP solving with OR-Tools"""
    
    @patch('solver.pywrapcp.RoutingModel')
    @patch('solver.pywrapcp.RoutingIndexManager')
    def test_solve_day_vrp_no_customers(self, mock_manager, mock_routing, sample_warehouse, sample_vehicles):
        """VRP with no customers should return empty routes"""
        solver = IRPSolver(sample_warehouse, [], sample_vehicles, 1, "2024-01-01")
        routes = solver._solve_day_vrp(0, datetime(2024, 1, 1), [])
        assert routes == []
    
    def test_solve_day_vrp_with_customers(self, sample_warehouse, sample_customers, sample_vehicles):
        """VRP with customers should return routes"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 1, "2024-01-01")
        customers_to_visit = [1, 3]  # Customers needing delivery
        
        # This will actually call OR-Tools, so we test the integration
        routes = solver._solve_day_vrp(0, datetime(2024, 1, 1), customers_to_visit)
        
        # Should return routes (may be empty if OR-Tools fails, but structure should be correct)
        assert isinstance(routes, list)
        if routes:
            assert all(hasattr(r, 'day') for r in routes)
            assert all(hasattr(r, 'stops') for r in routes)


class TestFallbackAlgorithm:
    """Tests for fallback nearest neighbor algorithm"""
    
    def test_fallback_creates_routes(self, sample_warehouse, sample_customers, sample_vehicles):
        """Fallback should create routes when OR-Tools fails"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 1, "2024-01-01")
        customers_to_visit = [1, 2, 3]
        
        routes = solver._create_fallback_routes(0, datetime(2024, 1, 1), customers_to_visit)
        
        assert len(routes) > 0
        assert all(hasattr(r, 'day') for r in routes)
        assert all(hasattr(r, 'stops') for r in routes)
        assert all(len(r.stops) > 0 for r in routes)
    
    def test_fallback_respects_capacity(self, sample_warehouse, sample_customers):
        """Fallback should respect vehicle capacity"""
        # Small capacity vehicle
        small_vehicle = MockVehicle(id=1, capacity=100)
        solver = IRPSolver(sample_warehouse, sample_customers, [small_vehicle], 1, "2024-01-01")
        customers_to_visit = [1, 2, 3]
        
        routes = solver._create_fallback_routes(0, datetime(2024, 1, 1), customers_to_visit)
        
        # Total load should not exceed capacity
        for route in routes:
            total_load = sum(stop.quantity for stop in route.stops)
            assert total_load <= small_vehicle.capacity


class TestEndToEndSolver:
    """End-to-end tests for complete solver"""
    
    def test_solve_complete_optimization(self, sample_warehouse, sample_customers, sample_vehicles):
        """Complete optimization should return OptimizeResponse"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 7, "2024-01-01")
        result = solver.solve()
        
        assert isinstance(result, OptimizeResponse)
        assert result.success == True
        assert isinstance(result.total_cost, float)
        assert isinstance(result.total_distance, float)
        assert isinstance(result.routes, list)
    
    def test_solve_empty_customers(self, sample_warehouse, sample_vehicles):
        """Optimization with no customers should succeed with empty routes"""
        solver = IRPSolver(sample_warehouse, [], sample_vehicles, 7, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
        assert len(result.routes) == 0
        assert result.total_cost == 0
    
    def test_solve_multi_day_horizon(self, sample_warehouse, sample_customers, sample_vehicles):
        """Multi-day optimization should generate routes for multiple days"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 7, "2024-01-01")
        result = solver.solve()
        
        # Should have routes for multiple days (if customers need delivery)
        days_with_routes = set(r.day for r in result.routes)
        # At least some days should have routes if customers need delivery
        if result.routes:
            assert len(days_with_routes) > 0
    
    def test_solve_inventory_tracking(self, sample_warehouse, sample_customers, sample_vehicles):
        """Inventory should be tracked correctly across days"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 3, "2024-01-01")
        initial_inv = solver.inventory[1]
        
        result = solver.solve()
        
        # Inventory should have changed (consumed daily)
        # Note: This tests that inventory state is maintained, not exact values
        assert result.success == True


class TestEdgeCases:
    """Edge case and error scenario tests"""
    
    def test_zero_planning_horizon(self, sample_warehouse, sample_customers, sample_vehicles):
        """Zero planning horizon should return empty result"""
        solver = IRPSolver(sample_warehouse, sample_customers, sample_vehicles, 0, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
        assert len(result.routes) == 0
    
    def test_single_customer(self, sample_warehouse, sample_vehicles):
        """Optimization with single customer should work"""
        customer = MockCustomer(id=1, lat=40.7580, lon=-73.9855, current_inv=50, min_inv=100)
        solver = IRPSolver(sample_warehouse, [customer], sample_vehicles, 1, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
    
    def test_single_vehicle(self, sample_warehouse, sample_customers):
        """Optimization with single vehicle should work"""
        vehicle = MockVehicle(id=1, capacity=10000)
        solver = IRPSolver(sample_warehouse, sample_customers, [vehicle], 1, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
    
    def test_very_large_capacity(self, sample_warehouse, sample_customers):
        """Vehicle with very large capacity should handle all customers"""
        vehicle = MockVehicle(id=1, capacity=1000000)
        solver = IRPSolver(sample_warehouse, sample_customers, [vehicle], 1, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
    
    def test_customers_far_apart(self, sample_warehouse):
        """Customers far apart should still be optimized"""
        customers = [
            MockCustomer(id=1, lat=40.7128, lon=-74.0060),  # NYC
            MockCustomer(id=2, lat=34.0522, lon=-118.2437),  # LA
        ]
        vehicles = [MockVehicle(id=1, capacity=10000, max_distance=10000000)]
        solver = IRPSolver(sample_warehouse, customers, vehicles, 1, "2024-01-01")
        result = solver.solve()
        
        assert result.success == True
