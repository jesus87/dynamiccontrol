# Dynamic Control Plane in Go

A lightweight prototype of a dynamic, policy-enforced control plane using Go and OPA (Rego). This project demonstrates how to build a flexible API gateway that loads routes from configuration and validates requests using OPA policies.

## Features

- **Dynamic Route Loading**: Routes are defined in JSON configuration and loaded at startup
- **OPA/Rego Integration**: Request validation using Open Policy Agent policies
- **JSON Schema Validation**: Request and response validation against JSON schemas
- **Mock Responses**: Simulated backend responses for demonstration
- **Policy Testing**: Unit tests for all Rego policies
- **RESTful API**: Clean HTTP endpoints with proper status codes

## Project Structure

```
dynamiccontrol/
├── cmd/server/
│   └── main.go                 # Main application entry point
├── config/
│   └── routes.json             # Route configuration
├── internal/
│   ├── opa/
│   │   └── policy_manager.go   # OPA policy management
│   ├── router/
│   │   └── route_manager.go    # Dynamic route management
│   ├── types/
│   │   └── types.go           # Data structures and types
│   └── validator/
│       └── schema_validator.go # JSON schema validation
├── policies/
│   ├── status_policy.rego      # Status endpoint policy
│   ├── status_policy.rego.test # Status policy tests
│   ├── traffic_policy.rego     # Traffic endpoint policy
│   ├── traffic_policy.rego.test # Traffic policy tests
│   ├── service_policy.rego     # Service validation policy
│   └── service_policy.rego.test # Service policy tests
├── go.mod                      # Go module dependencies
└── README.md                   # This file
```

## Prerequisites

- Go 1.21 or later
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/jesus87/dynamiccontrol.git
cd dynamiccontrol
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

4. Execute endpoint tests
```bash
cd examples
test_endpoints.sh
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

## Configuration

### Route Configuration (`config/routes.json`)

Routes are defined in JSON format with the following structure:

```json
{
  "routes": [
    {
      "routeName": "/v1/status",
      "method": "GET",
      "requestSchema": {},
      "responseSchema": {
        "type": "object",
        "properties": {
          "status": {"type": "string"},
          "timestamp": {"type": "string"},
          "version": {"type": "string"},
          "uptime": {"type": "number"}
        }
      },
      "policies": ["status_policy"]
    }
  ]
}
```

### OPA Policies

Policies are written in Rego and stored in the `policies/` directory. Each policy file should:

1. Define an `allow` rule that returns a boolean
2. Include proper input validation
3. Have corresponding test files (`.rego.test`)

## API Endpoints

### Health Check
```bash
GET /health
```
Returns service health status.

### Service Information
```bash
GET /info
```
Returns information about the service, loaded routes, and policies.

### Status Endpoint
```bash
GET /v1/status
```
Returns service status information. Validated by `status_policy`.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": 3600
}
```

### Traffic Management
```bash
POST /v1/services/:serviceId/traffic
```

**Request Body:**
```json
{
  "trafficType": "incoming",
  "volume": 100.5,
  "priority": "medium",
  "metadata": {
    "source": "service-a",
    "destination": "service-b",
    "protocol": "http"
  }
}
```

**Response:**
```json
{
  "id": "traffic-20240101120000",
  "serviceId": "service123",
  "status": "accepted",
  "message": "Traffic request processed successfully",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Testing

### Running Go Tests
```bash
go test ./...
```

### Running OPA Policy Tests
```bash
# Install OPA if not already installed
curl -L -o opa https://openpolicyagent.org/downloads/latest/opa_linux_amd64
chmod +x opa

# Test all policies
opa test policies/ --verbose
```

### Testing with curl

1. **Health Check:**
```bash
curl http://localhost:8080/health
```

2. **Service Info:**
```bash
curl http://localhost:8080/info
```

3. **Status Endpoint:**
```bash
curl http://localhost:8080/v1/status
```

4. **Traffic Management:**
```bash
curl -X POST http://localhost:8080/v1/services/service123/traffic \
  -H "Content-Type: application/json" \
  -d '{
    "trafficType": "incoming",
    "volume": 100.5,
    "priority": "medium",
    "metadata": {
      "source": "service-a",
      "destination": "service-b",
      "protocol": "http"
    }
  }'
```

## Policy Examples

### Status Policy (`policies/status_policy.rego`)
```rego
package status_policy

default allow = false

allow if {
    input.method == "GET"
    input.path == "/v1/status"
}
```

### Traffic Policy (`policies/traffic_policy.rego`)
```rego
package traffic_policy

default allow = false

allow if {
    input.method == "POST"
    input.path = startswith("/v1/services/")
    input.path = endswith("/traffic")
    
    input.body.trafficType in ["incoming", "outgoing", "internal"]
    input.body.volume >= 0
    input.body.priority in ["low", "medium", "high", "critical"]
}
```

## Architecture

### Components

1. **Policy Manager**: Loads and evaluates OPA/Rego policies
2. **Route Manager**: Handles dynamic route registration and request processing
3. **Schema Validator**: Validates requests and responses against JSON schemas
4. **Mock Data**: Provides simulated responses for endpoints

### Request Flow

1. **Route Matching**: Request is matched to configured route
2. **Schema Validation**: Request body is validated against JSON schema
3. **Policy Evaluation**: OPA policies are evaluated with request data
4. **Response Generation**: Mock response is generated and validated
5. **Response Return**: Validated response is returned to client

## Extending the Project

### Adding New Routes

1. Add route definition to `config/routes.json`
2. Create corresponding Rego policy in `policies/`
3. Add policy tests in `policies/*.rego.test`
4. Restart the server

### Adding New Policies

1. Create `.rego` file in `policies/` directory
2. Define `allow` rule with appropriate logic
3. Create corresponding `.rego.test` file
4. Reference policy name in route configuration

### Custom Response Generation

Modify the `MockData` struct in `internal/types/types.go` to add custom response generation logic.

## Error Handling

The application provides comprehensive error handling:

- **400 Bad Request**: Invalid JSON or schema validation failures
- **403 Forbidden**: Policy evaluation denies the request
- **404 Not Found**: Route not found
- **500 Internal Server Error**: Server-side errors

## Security Considerations

- All requests are validated against JSON schemas
- OPA policies provide fine-grained access control
- Input sanitization and validation
- Proper HTTP status codes for different error conditions

## Performance

- OPA policies are pre-compiled for efficient evaluation
- JSON schema validation is optimized
- Minimal memory footprint
- Fast request processing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under jesus87 permission