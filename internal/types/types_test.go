package types

import (
	"testing"
	"time"
)

func TestNewMockData(t *testing.T) {
	mockData := NewMockData()

	if mockData == nil {
		t.Fatal("NewMockData() returned nil")
	}

	if len(mockData.StatusResponses) == 0 {
		t.Error("StatusResponses should not be empty")
	}

	if len(mockData.TrafficResponses) == 0 {
		t.Error("TrafficResponses should not be empty")
	}
}

func TestGenerateStatusResponse(t *testing.T) {
	mockData := NewMockData()
	response := mockData.GenerateStatusResponse()

	if response.Status == "" {
		t.Error("Status should not be empty")
	}

	if response.Version == "" {
		t.Error("Version should not be empty")
	}

	if response.Uptime <= 0 {
		t.Error("Uptime should be positive")
	}

	// Check that timestamp is recent
	now := time.Now()
	if response.Timestamp.After(now) {
		t.Error("Timestamp should not be in the future")
	}

	if response.Timestamp.Before(now.Add(-time.Second)) {
		t.Error("Timestamp should be recent")
	}
}

func TestGenerateTrafficResponse(t *testing.T) {
	mockData := NewMockData()
	request := TrafficRequest{
		TrafficType: "incoming",
		Volume:      100.5,
		Priority:    "medium",
		Metadata: map[string]interface{}{
			"source":      "service-a",
			"destination": "service-b",
		},
	}

	response := mockData.GenerateTrafficResponse("service123", request)

	if response.ID == "" {
		t.Error("ID should not be empty")
	}

	if response.ServiceID != "service123" {
		t.Errorf("ServiceID should be 'service123', got '%s'", response.ServiceID)
	}

	if response.Status == "" {
		t.Error("Status should not be empty")
	}

	if response.Message == "" {
		t.Error("Message should not be empty")
	}

	// Check that timestamp is recent
	now := time.Now()
	if response.Timestamp.After(now) {
		t.Error("Timestamp should not be in the future")
	}

	if response.Timestamp.Before(now.Add(-time.Second)) {
		t.Error("Timestamp should be recent")
	}
}

func TestTrafficRequestValidation(t *testing.T) {
	validRequest := TrafficRequest{
		TrafficType: "incoming",
		Volume:      100.5,
		Priority:    "medium",
	}

	if validRequest.TrafficType == "" {
		t.Error("TrafficType should not be empty")
	}

	if validRequest.Volume < 0 {
		t.Error("Volume should be non-negative")
	}

	if validRequest.Priority == "" {
		t.Error("Priority should not be empty")
	}
}

func TestStatusResponseValidation(t *testing.T) {
	validResponse := StatusResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    3600,
	}

	if validResponse.Status == "" {
		t.Error("Status should not be empty")
	}

	if validResponse.Version == "" {
		t.Error("Version should not be empty")
	}

	if validResponse.Uptime < 0 {
		t.Error("Uptime should be non-negative")
	}
}

func TestTrafficResponseValidation(t *testing.T) {
	validResponse := TrafficResponse{
		ID:        "traffic-123",
		ServiceID: "service123",
		Status:    "accepted",
		Message:   "Request processed",
		Timestamp: time.Now(),
	}

	if validResponse.ID == "" {
		t.Error("ID should not be empty")
	}

	if validResponse.ServiceID == "" {
		t.Error("ServiceID should not be empty")
	}

	if validResponse.Status == "" {
		t.Error("Status should not be empty")
	}

	if validResponse.Message == "" {
		t.Error("Message should not be empty")
	}
}
