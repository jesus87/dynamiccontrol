package traffic_policy

import future.keywords.if
import future.keywords.in

# Default to deny
default allow = false

# Allow POST requests to traffic endpoint with basic validation
allow if {
    input.method == "POST"
    startswith(input.path, "/v1/services/")
    endswith(input.path, "/traffic")
    
    # Basic validation
    input.body.trafficType in ["incoming", "outgoing", "internal"]
    input.body.volume >= 0
    input.body.priority in ["low", "medium", "high", "critical"]
} 