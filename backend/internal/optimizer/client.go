package optimizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Optimization can take time
		},
	}
}

// OptimizeRequest represents the request to the optimizer service
type OptimizeRequest struct {
	Warehouse  WarehouseData   `json:"warehouse"`
	Customers  []CustomerData  `json:"customers"`
	Vehicles   []VehicleData   `json:"vehicles"`
	PlanningHorizon int        `json:"planning_horizon"`
	StartDate  string          `json:"start_date"`
}

type WarehouseData struct {
	ID        int64   `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Stock     float64 `json:"stock"`
}

type CustomerData struct {
	ID               int64   `json:"id"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	DemandRate       float64 `json:"demand_rate"`
	MaxInventory     float64 `json:"max_inventory"`
	CurrentInventory float64 `json:"current_inventory"`
	MinInventory     float64 `json:"min_inventory"`
	Priority         int     `json:"priority"`
}

type VehicleData struct {
	ID          int64   `json:"id"`
	Capacity    float64 `json:"capacity"`
	CostPerKm   float64 `json:"cost_per_km"`
	FixedCost   float64 `json:"fixed_cost"`
	MaxDistance float64 `json:"max_distance"`
}

// OptimizeResponse represents the response from the optimizer service
type OptimizeResponse struct {
	Success       bool          `json:"success"`
	Message       string        `json:"message"`
	TotalCost     float64       `json:"total_cost"`
	TotalDistance float64       `json:"total_distance"`
	Routes        []RouteResult `json:"routes"`
}

type RouteResult struct {
	Day           int          `json:"day"`
	Date          string       `json:"date"`
	VehicleID     int64        `json:"vehicle_id"`
	TotalDistance float64      `json:"total_distance"`
	TotalCost     float64      `json:"total_cost"`
	TotalLoad     float64      `json:"total_load"`
	Stops         []StopResult `json:"stops"`
}

type StopResult struct {
	CustomerID  int64   `json:"customer_id"`
	Sequence    int     `json:"sequence"`
	Quantity    float64 `json:"quantity"`
	ArrivalTime string  `json:"arrival_time"`
}

// HealthCheck checks if the optimizer service is available
func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("optimizer service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("optimizer service returned status %d", resp.StatusCode)
	}
	return nil
}

// Optimize sends the optimization request to the Python service
func (c *Client) Optimize(req *OptimizeRequest) (*OptimizeResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/optimize",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call optimizer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("optimizer returned status %d", resp.StatusCode)
	}

	var result OptimizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

