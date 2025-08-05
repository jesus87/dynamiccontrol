package status_policy

import future.keywords.if
import future.keywords.in

# Default to deny
default allow = false

# Allow GET requests to status endpoint
allow if {
    input.method == "GET"
    input.path == "/v1/status"
}

# Validate that the request has proper headers
allow if {
    input.method == "GET"
    input.path == "/v1/status"
    input.headers["Content-Type"] == "application/json"
}

# Additional validation for request parameters
validate_request if {
    input.method == "GET"
    input.path == "/v1/status"
    
    # Check if any query parameters are provided (optional)
    not input.query
}

# Helper function to check if request is from authorized source
is_authorized if {
    input.headers["Authorization"] != null
}

# Enhanced validation for authorized requests
allow if {
    input.method == "GET"
    input.path == "/v1/status"
    is_authorized
} 