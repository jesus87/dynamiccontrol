package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"dynamiccontrol/internal/opa"
	"dynamiccontrol/internal/types"
	"dynamiccontrol/internal/validator"

	"github.com/gin-gonic/gin"
)

// RouteManager handles dynamic route registration and management
type RouteManager struct {
	config          *types.RoutesConfig
	policyManager   *opa.PolicyManager
	schemaValidator *validator.SchemaValidator
	mockData        *types.MockData
}

// NewRouteManager creates a new route manager
func NewRouteManager(policyManager *opa.PolicyManager, schemaValidator *validator.SchemaValidator) *RouteManager {
	return &RouteManager{
		policyManager:   policyManager,
		schemaValidator: schemaValidator,
		mockData:        types.NewMockData(),
	}
}

// LoadConfig loads the route configuration from JSON file
func (rm *RouteManager) LoadConfig(configPath string) error {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.RoutesConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	rm.config = &config
	log.Printf("Loaded %d routes from configuration", len(config.Routes))
	return nil
}

// RegisterRoutes registers all routes from the configuration
func (rm *RouteManager) RegisterRoutes(router *gin.Engine) error {
	if rm.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	for _, route := range rm.config.Routes {
		if err := rm.registerRoute(router, route); err != nil {
			log.Printf("Failed to register route %s: %v", route.RouteName, err)
			continue
		}
		log.Printf("Registered route: %s %s", route.Method, route.RouteName)
	}

	return nil
}

// registerRoute registers a single route
func (rm *RouteManager) registerRoute(router *gin.Engine, route types.RouteConfig) error {
	switch route.Method {
	case "GET":
		router.GET(route.RouteName, rm.createHandler(route))
	case "POST":
		router.POST(route.RouteName, rm.createHandler(route))
	case "PUT":
		router.PUT(route.RouteName, rm.createHandler(route))
	case "DELETE":
		router.DELETE(route.RouteName, rm.createHandler(route))
	default:
		return fmt.Errorf("unsupported HTTP method: %s", route.Method)
	}

	return nil
}

// createHandler creates a Gin handler for a route
func (rm *RouteManager) createHandler(route types.RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract headers
		headers := make(map[string]string)
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}

		// Handle different HTTP methods
		switch route.Method {
		case "GET":
			rm.handleGET(c, route, headers)
		case "POST":
			rm.handlePOST(c, route, headers)
		default:
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method not allowed",
			})
		}
	}
}

// handleGET handles GET requests
func (rm *RouteManager) handleGET(c *gin.Context, route types.RouteConfig, headers map[string]string) {
	// Create policy input
	input := opa.CreatePolicyInput("GET", route.RouteName, headers, nil)

	// Evaluate policies
	policyResult, err := rm.policyManager.EvaluatePolicies(route.Policies, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Policy evaluation error: %v", err),
		})
		return
	}

	if !policyResult.Allowed {
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Request denied by policy: %s", policyResult.Error),
		})
		return
	}

	// Generate mock response based on route
	var response interface{}
	switch route.RouteName {
	case "/v1/status":
		statusResponse := rm.mockData.GenerateStatusResponse()
		response = statusResponse

		// Validate response against schema
		validationResult := rm.schemaValidator.ValidateStatusResponse(statusResponse)
		if !validationResult.Valid {
			log.Printf("Response validation failed: %v", validationResult.Errors)
		}
	default:
		// Generic response for other GET routes
		response = gin.H{
			"message": "GET request processed successfully",
			"route":   route.RouteName,
			"method":  route.Method,
		}
	}

	c.JSON(http.StatusOK, response)
}

// handlePOST handles POST requests
func (rm *RouteManager) handlePOST(c *gin.Context, route types.RouteConfig, headers map[string]string) {
	// Parse request body
	var requestBody interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid JSON: %v", err),
		})
		return
	}

	// Validate request against schema
	validationResult := rm.schemaValidator.ValidateRequest(route.RequestSchema, requestBody)
	if !validationResult.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Request validation failed",
			"details": validator.FormatValidationErrors(validationResult.Errors),
		})
		return
	}

	// Create policy input
	input := opa.CreatePolicyInput("POST", route.RouteName, headers, requestBody)

	// Evaluate policies
	policyResult, err := rm.policyManager.EvaluatePolicies(route.Policies, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Policy evaluation error: %v", err),
		})
		return
	}

	if !policyResult.Allowed {
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Request denied by policy: %s", policyResult.Error),
		})
		return
	}

	// Generate mock response based on route
	var response interface{}
	switch {
	case strings.Contains(route.RouteName, "/services/") && strings.HasSuffix(route.RouteName, "/traffic"):
		// Extract service ID from actual path
		pathParts := strings.Split(c.Request.URL.Path, "/")
		if len(pathParts) >= 4 {
			serviceID := pathParts[3] // /v1/services/service123/traffic

			// Parse traffic request
			var trafficRequest types.TrafficRequest
			if requestBytes, err := json.Marshal(requestBody); err == nil {
				json.Unmarshal(requestBytes, &trafficRequest)
			}

			trafficResponse := rm.mockData.GenerateTrafficResponse(serviceID, trafficRequest)
			response = trafficResponse

			// Validate response against schema
			responseValidation := rm.schemaValidator.ValidateTrafficResponse(trafficResponse)
			if !responseValidation.Valid {
				log.Printf("Response validation failed: %v", responseValidation.Errors)
			}
		}
	default:
		// Generic response for other POST routes
		response = gin.H{
			"message": "POST request processed successfully",
			"route":   route.RouteName,
			"method":  route.Method,
			"data":    requestBody,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetConfig returns the current route configuration
func (rm *RouteManager) GetConfig() *types.RoutesConfig {
	return rm.config
}

// GetMockData returns the mock data instance
func (rm *RouteManager) GetMockData() *types.MockData {
	return rm.mockData
}
