package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"dynamiccontrol/internal/types"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidator handles JSON schema validation
type SchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemas: make(map[string]*gojsonschema.Schema),
	}
}

// ValidateRequest validates a request against its schema
func (sv *SchemaValidator) ValidateRequest(schema map[string]interface{}, data interface{}) *types.ValidationResult {
	if schema == nil || len(schema) == 0 {
		return &types.ValidationResult{
			Valid: true,
		}
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return &types.ValidationResult{
			Valid:  false,
			Errors: []string{fmt.Sprintf("Invalid schema: %v", err)},
		}
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return &types.ValidationResult{
			Valid:  false,
			Errors: []string{fmt.Sprintf("Invalid data: %v", err)},
		}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	documentLoader := gojsonschema.NewBytesLoader(dataBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return &types.ValidationResult{
			Valid:  false,
			Errors: []string{fmt.Sprintf("Validation error: %v", err)},
		}
	}

	if result.Valid() {
		return &types.ValidationResult{
			Valid: true,
		}
	}

	errors := make([]string, 0, len(result.Errors()))
	for _, err := range result.Errors() {
		errors = append(errors, err.String())
	}

	return &types.ValidationResult{
		Valid:   false,
		Errors:  errors,
		Details: fmt.Sprintf("Validation failed with %d errors", len(errors)),
	}
}

// ValidateResponse validates a response against its schema
func (sv *SchemaValidator) ValidateResponse(schema map[string]interface{}, data interface{}) *types.ValidationResult {
	return sv.ValidateRequest(schema, data)
}

// ValidateTrafficRequest validates a traffic request
func (sv *SchemaValidator) ValidateTrafficRequest(request types.TrafficRequest) *types.ValidationResult {
	// Define the schema for traffic request
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"trafficType": map[string]interface{}{
				"type": "string",
				"enum": []string{"incoming", "outgoing", "internal"},
			},
			"volume": map[string]interface{}{
				"type":    "number",
				"minimum": 0,
			},
			"priority": map[string]interface{}{
				"type": "string",
				"enum": []string{"low", "medium", "high", "critical"},
			},
			"metadata": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"source": map[string]interface{}{
						"type": "string",
					},
					"destination": map[string]interface{}{
						"type": "string",
					},
					"protocol": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
		"required": []string{"trafficType", "volume", "priority"},
	}

	return sv.ValidateRequest(schema, request)
}

// ValidateStatusRequest validates a status request (empty for GET)
func (sv *SchemaValidator) ValidateStatusRequest() *types.ValidationResult {
	// Status endpoint is GET, so no request body validation needed
	return &types.ValidationResult{
		Valid: true,
	}
}

// ValidateTrafficResponse validates a traffic response
func (sv *SchemaValidator) ValidateTrafficResponse(response types.TrafficResponse) *types.ValidationResult {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type": "string",
			},
			"serviceId": map[string]interface{}{
				"type": "string",
			},
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"accepted", "rejected", "pending"},
			},
			"message": map[string]interface{}{
				"type": "string",
			},
			"timestamp": map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
		},
		"required": []string{"id", "serviceId", "status", "message", "timestamp"},
	}

	return sv.ValidateResponse(schema, response)
}

// ValidateStatusResponse validates a status response
func (sv *SchemaValidator) ValidateStatusResponse(response types.StatusResponse) *types.ValidationResult {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"healthy", "degraded", "unhealthy"},
			},
			"timestamp": map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
			"version": map[string]interface{}{
				"type": "string",
			},
			"uptime": map[string]interface{}{
				"type": "number",
			},
		},
		"required": []string{"status", "timestamp", "version", "uptime"},
	}

	return sv.ValidateResponse(schema, response)
}

// FormatValidationErrors formats validation errors for better readability
func FormatValidationErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	var formattedErrors []string
	for _, err := range errors {
		// Remove common prefixes for cleaner error messages
		cleanErr := strings.TrimPrefix(err, "(root): ")
		cleanErr = strings.TrimPrefix(cleanErr, "I[0,")
		formattedErrors = append(formattedErrors, cleanErr)
	}

	return strings.Join(formattedErrors, "; ")
}
