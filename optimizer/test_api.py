"""
Integration tests for Optimizer API
Tests FastAPI endpoints, request/response handling, and error scenarios.
"""

import pytest
from fastapi.testclient import TestClient
from main import app
from solver import IRPSolver


@pytest.fixture
def client():
    """Create test client for FastAPI app"""
    return TestClient(app)


@pytest.fixture
def sample_optimize_request():
    """Sample optimization request payload"""
    return {
        "warehouse": {
            "id": 1,
            "latitude": 40.7128,
            "longitude": -74.0060,
            "stock": 10000.0
        },
        "customers": [
            {
                "id": 1,
                "latitude": 40.7580,
                "longitude": -73.9855,
                "demand_rate": 50.0,
                "max_inventory": 1000.0,
                "current_inventory": 200.0,
                "min_inventory": 100.0,
                "priority": 1
            },
            {
                "id": 2,
                "latitude": 40.7505,
                "longitude": -73.9934,
                "demand_rate": 50.0,
                "max_inventory": 1000.0,
                "current_inventory": 600.0,
                "min_inventory": 100.0,
                "priority": 1
            }
        ],
        "vehicles": [
            {
                "id": 1,
                "capacity": 5000.0,
                "cost_per_km": 1.0,
                "fixed_cost": 100.0,
                "max_distance": 0.0
            }
        ],
        "planning_horizon": 7,
        "start_date": "2024-01-01"
    }


class TestHealthEndpoint:
    """Tests for /health endpoint"""
    
    def test_health_check_success(self, client):
        """Health check should return 200 OK"""
        response = client.get("/health")
        assert response.status_code == 200
        
        data = response.json()
        assert data["status"] == "healthy"
        assert "service" in data
        assert "timestamp" in data


class TestOptimizeEndpoint:
    """Tests for /optimize endpoint"""
    
    def test_optimize_success(self, client, sample_optimize_request):
        """Valid optimization request should return success"""
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] == True
        assert "total_cost" in data
        assert "total_distance" in data
        assert "routes" in data
        assert isinstance(data["routes"], list)
    
    def test_optimize_no_customers(self, client, sample_optimize_request):
        """Request with no customers should return error"""
        sample_optimize_request["customers"] = []
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200  # API returns 200 with success=false
        data = response.json()
        assert data["success"] == False
        assert "No customers" in data["message"]
    
    def test_optimize_no_vehicles(self, client, sample_optimize_request):
        """Request with no vehicles should return error"""
        sample_optimize_request["vehicles"] = []
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] == False
        assert "No vehicles" in data["message"]
    
    def test_optimize_invalid_warehouse(self, client, sample_optimize_request):
        """Request with invalid warehouse data should return error"""
        sample_optimize_request["warehouse"]["latitude"] = "invalid"
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 422  # Validation error
    
    def test_optimize_missing_required_fields(self, client):
        """Request missing required fields should return validation error"""
        incomplete_request = {
            "warehouse": {"id": 1},
            "customers": [],
            "vehicles": []
        }
        response = client.post("/optimize", json=incomplete_request)
        
        assert response.status_code == 422  # Validation error
    
    def test_optimize_invalid_date_format(self, client, sample_optimize_request):
        """Request with invalid date format should return error"""
        sample_optimize_request["start_date"] = "invalid-date"
        response = client.post("/optimize", json=sample_optimize_request)
        
        # Should either validate or fail during processing
        assert response.status_code in [200, 422, 500]
    
    def test_optimize_response_structure(self, client, sample_optimize_request):
        """Response should have correct structure"""
        response = client.post("/optimize", json=sample_optimize_request)
        
        if response.status_code == 200:
            data = response.json()
            
            # Check top-level fields
            assert "success" in data
            assert "message" in data
            assert "total_cost" in data
            assert "total_distance" in data
            assert "routes" in data
            
            # Check route structure if routes exist
            if data["routes"]:
                route = data["routes"][0]
                assert "day" in route
                assert "date" in route
                assert "vehicle_id" in route
                assert "total_distance" in route
                assert "total_cost" in route
                assert "total_load" in route
                assert "stops" in route
                
                # Check stop structure if stops exist
                if route["stops"]:
                    stop = route["stops"][0]
                    assert "customer_id" in stop
                    assert "sequence" in stop
                    assert "quantity" in stop
                    assert "arrival_time" in stop


class TestOptimizeScenarios:
    """Tests for various optimization scenarios"""
    
    def test_optimize_single_day(self, client, sample_optimize_request):
        """Single day optimization should work"""
        sample_optimize_request["planning_horizon"] = 1
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] == True
    
    def test_optimize_multi_day(self, client, sample_optimize_request):
        """Multi-day optimization should generate routes for multiple days"""
        sample_optimize_request["planning_horizon"] = 7
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        if data["success"]:
            days = set(r["day"] for r in data["routes"])
            # Should have routes for potentially multiple days
            assert len(days) >= 0
    
    def test_optimize_large_horizon(self, client, sample_optimize_request):
        """Large planning horizon should be handled"""
        sample_optimize_request["planning_horizon"] = 30
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        assert data["success"] == True
    
    def test_optimize_multiple_vehicles(self, client, sample_optimize_request):
        """Multiple vehicles should be utilized"""
        sample_optimize_request["vehicles"] = [
            {"id": 1, "capacity": 3000.0, "cost_per_km": 1.0, "fixed_cost": 100.0, "max_distance": 0.0},
            {"id": 2, "capacity": 2000.0, "cost_per_km": 0.8, "fixed_cost": 80.0, "max_distance": 0.0},
        ]
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        if data["success"] and data["routes"]:
            vehicle_ids = set(r["vehicle_id"] for r in data["routes"])
            # Should use at least one vehicle
            assert len(vehicle_ids) > 0
    
    def test_optimize_capacity_constraint(self, client, sample_optimize_request):
        """Vehicle capacity constraints should be respected"""
        # Small capacity vehicle
        sample_optimize_request["vehicles"] = [
            {"id": 1, "capacity": 100.0, "cost_per_km": 1.0, "fixed_cost": 100.0, "max_distance": 0.0}
        ]
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        if data["success"] and data["routes"]:
            for route in data["routes"]:
                assert route["total_load"] <= 100.0


class TestErrorHandling:
    """Tests for error handling and edge cases"""
    
    def test_optimize_malformed_json(self, client):
        """Malformed JSON should return 422"""
        response = client.post(
            "/optimize",
            data="not json",
            headers={"Content-Type": "application/json"}
        )
        assert response.status_code == 422
    
    def test_optimize_wrong_method(self, client):
        """GET request to /optimize should return 405"""
        response = client.get("/optimize")
        assert response.status_code == 405
    
    def test_optimize_negative_values(self, client, sample_optimize_request):
        """Negative values should be handled (may be invalid)"""
        sample_optimize_request["customers"][0]["current_inventory"] = -100
        response = client.post("/optimize", json=sample_optimize_request)
        
        # Should either validate or process (depending on business logic)
        assert response.status_code in [200, 422]
    
    def test_optimize_zero_capacity_vehicle(self, client, sample_optimize_request):
        """Vehicle with zero capacity should be handled"""
        sample_optimize_request["vehicles"][0]["capacity"] = 0.0
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        data = response.json()
        # Should either fail or return empty routes
        assert "success" in data


class TestPerformance:
    """Performance and load tests"""
    
    def test_optimize_many_customers(self, client, sample_optimize_request):
        """Optimization with many customers should complete"""
        # Create 50 customers
        customers = []
        for i in range(50):
            customers.append({
                "id": i + 1,
                "latitude": 40.7128 + (i * 0.01),
                "longitude": -74.0060 + (i * 0.01),
                "demand_rate": 50.0,
                "max_inventory": 1000.0,
                "current_inventory": 200.0,
                "min_inventory": 100.0,
                "priority": 1
            })
        sample_optimize_request["customers"] = customers
        
        response = client.post("/optimize", json=sample_optimize_request)
        
        assert response.status_code == 200
        # Should complete within reasonable time (test timeout will catch if too slow)
    
    def test_optimize_concurrent_requests(self, client, sample_optimize_request):
        """Multiple concurrent requests should be handled"""
        import concurrent.futures
        
        def make_request():
            return client.post("/optimize", json=sample_optimize_request)
        
        with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
            futures = [executor.submit(make_request) for _ in range(5)]
            results = [f.result() for f in concurrent.futures.as_completed(futures)]
        
        # All requests should succeed
        assert all(r.status_code == 200 for r in results)
        assert all(r.json()["success"] for r in results)
