package service_policy

import future.keywords.if
import future.keywords.in

# Default to deny
default allow = false

# Allow all POST requests to traffic endpoint
allow if {
    input.method == "POST"
    startswith(input.path, "/v1/services/")
    endswith(input.path, "/traffic")
} 