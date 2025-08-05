package types

import (
	"time"
)

// RouteConfig represents the configuration for a single route
type RouteConfig struct {
	RouteName      string                 `json:"routeName"`
	Method         string                 `json:"method"`
	RequestSchema  map[string]interface{} `json:"requestSchema"`
	ResponseSchema map[string]interface{} `json:"responseSchema"`
	Policies       []string               `json:"policies"`
}

// RoutesConfig represents the complete routes configuration
type RoutesConfig struct {
	Routes []RouteConfig `json:"routes"`
}

// StatusResponse represents the response for the status endpoint
type StatusResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    int64     `json:"uptime"`
}

// TrafficRequest represents the request payload for traffic endpoint
type TrafficRequest struct {
	TrafficType string                 `json:"trafficType"`
	Volume      float64                `json:"volume"`
	Priority    string                 `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TrafficResponse represents the response for the traffic endpoint
type TrafficResponse struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"serviceId"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// PolicyResult represents the result of a policy evaluation
type PolicyResult struct {
	Allowed bool   `json:"allowed"`
	Error   string `json:"error,omitempty"`
}

// ValidationResult represents the result of request validation
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	Details string   `json:"details,omitempty"`
}

// MockData provides mock responses for endpoints
type MockData struct {
	StatusResponses  map[string]StatusResponse
	TrafficResponses map[string]TrafficResponse
}

// NewMockData creates a new instance of MockData with default values
func NewMockData() *MockData {
	return &MockData{
		StatusResponses: map[string]StatusResponse{
			"default": {
				Status:    "healthy",
				Timestamp: time.Now(),
				Version:   "1.0.0",
				Uptime:    3600,
			},
		},
		TrafficResponses: map[string]TrafficResponse{
			"default": {
				ID:        "traffic-123",
				ServiceID: "service-123",
				Status:    "accepted",
				Message:   "Traffic request processed successfully",
				Timestamp: time.Now(),
			},
		},
	}
}

// GenerateTrafficResponse creates a mock traffic response
func (md *MockData) GenerateTrafficResponse(serviceID string, request TrafficRequest) TrafficResponse {
	return TrafficResponse{
		ID:        generateID(),
		ServiceID: serviceID,
		Status:    "accepted",
		Message:   "Traffic request processed successfully",
		Timestamp: time.Now(),
	}
}

// GenerateStatusResponse creates a mock status response
func (md *MockData) GenerateStatusResponse() StatusResponse {
	return StatusResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    3600,
	}
}

// generateID generates a simple ID for mock responses
func generateID() string {
	return "traffic-" + time.Now().Format("20060102150405")
}
