package optimizer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHealthCheck tests optimizer health check
func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		serverStatus   int
		serverResponse string
		wantErr        bool
	}{
		{
			name:           "healthy service",
			serverStatus:   http.StatusOK,
			serverResponse: `{"status":"healthy"}`,
			wantErr:        false,
		},
		{
			name:           "unhealthy service",
			serverStatus:   http.StatusInternalServerError,
			serverResponse: `{"error":"internal error"}`,
			wantErr:        true,
		},
		{
			name:           "service unavailable",
			serverStatus:   0, // Server not started
			serverResponse: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.serverStatus > 0 {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.serverStatus)
					w.Write([]byte(tt.serverResponse))
				}))
				defer server.Close()
			}

			url := "http://localhost:9999"
			if server != nil {
				url = server.URL
			}

			client := NewClient(url)
			err := client.HealthCheck()

			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestOptimize tests optimization request
func TestOptimize(t *testing.T) {
	tests := []struct {
		name           string
		request        *OptimizeRequest
		serverResponse OptimizeResponse
		serverStatus   int
		wantErr        bool
		validateResult func(*OptimizeResponse) bool
	}{
		{
			name: "successful optimization",
			request: &OptimizeRequest{
				Warehouse: WarehouseData{
					ID:        1,
					Latitude:  40.7128,
					Longitude: -74.0060,
					Stock:     10000,
				},
				Customers: []CustomerData{
					{
						ID:               1,
						Latitude:         40.7580,
						Longitude:        -73.9855,
						DemandRate:       50,
						MaxInventory:     1000,
						CurrentInventory: 200,
						MinInventory:     100,
						Priority:         1,
					},
				},
				Vehicles: []VehicleData{
					{
						ID:          1,
						Capacity:    5000,
						CostPerKm:   1.0,
						FixedCost:   100,
						MaxDistance: 0,
					},
				},
				PlanningHorizon: 7,
				StartDate:       "2024-01-01",
			},
			serverResponse: OptimizeResponse{
				Success:       true,
				Message:       "Optimization complete",
				TotalCost:     500.0,
				TotalDistance: 100.0,
				Routes:        []RouteResult{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			validateResult: func(resp *OptimizeResponse) bool {
				return resp.Success && resp.TotalCost == 500.0
			},
		},
		{
			name: "optimizer returns error",
			request: &OptimizeRequest{
				Warehouse:       WarehouseData{ID: 1, Latitude: 40.7128, Longitude: -74.0060},
				Customers:       []CustomerData{},
				Vehicles:        []VehicleData{{ID: 1, Capacity: 5000}},
				PlanningHorizon: 7,
				StartDate:       "2024-01-01",
			},
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
		},
		{
			name: "invalid JSON response",
			request: &OptimizeRequest{
				Warehouse:       WarehouseData{ID: 1, Latitude: 40.7128, Longitude: -74.0060},
				Customers:       []CustomerData{{ID: 1, Latitude: 40.0, Longitude: -74.0}},
				Vehicles:        []VehicleData{{ID: 1, Capacity: 5000}},
				PlanningHorizon: 1,
				StartDate:       "2024-01-01",
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK && !tt.wantErr {
					json.NewEncoder(w).Encode(tt.serverResponse)
				} else {
					w.Write([]byte("invalid json"))
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			result, err := client.Optimize(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Optimize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				if !tt.validateResult(result) {
					t.Error("Optimize() result validation failed")
				}
			}
		})
	}
}

// TestOptimizeTimeout tests timeout handling
func TestOptimizeTimeout(t *testing.T) {
	// Create server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Minute) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req := &OptimizeRequest{
		Warehouse:       WarehouseData{ID: 1, Latitude: 40.7128, Longitude: -74.0060},
		Customers:       []CustomerData{{ID: 1, Latitude: 40.0, Longitude: -74.0}},
		Vehicles:        []VehicleData{{ID: 1, Capacity: 5000}},
		PlanningHorizon: 1,
		StartDate:       "2024-01-01",
	}

	_, err := client.Optimize(req)
	if err == nil {
		t.Error("Optimize() should timeout but didn't")
	}
}

// TestOptimizeRequestMarshaling tests request JSON marshaling
func TestOptimizeRequestMarshaling(t *testing.T) {
	req := &OptimizeRequest{
		Warehouse: WarehouseData{
			ID:        1,
			Latitude:  40.7128,
			Longitude: -74.0060,
			Stock:     10000,
		},
		Customers: []CustomerData{
			{
				ID:               1,
				Latitude:         40.0,
				Longitude:        -74.0,
				DemandRate:       50,
				MaxInventory:     1000,
				CurrentInventory: 500,
				MinInventory:     100,
				Priority:         1,
			},
		},
		Vehicles: []VehicleData{
			{
				ID:          1,
				Capacity:    5000,
				CostPerKm:   1.0,
				FixedCost:   100,
				MaxDistance: 0,
			},
		},
		PlanningHorizon: 7,
		StartDate:       "2024-01-01",
	}

	// Test that request can be marshaled to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Test that JSON can be unmarshaled back
	var unmarshaled OptimizeRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.PlanningHorizon != req.PlanningHorizon {
		t.Errorf("PlanningHorizon = %v, want %v", unmarshaled.PlanningHorizon, req.PlanningHorizon)
	}
}
