package main

import (
	"log"
	"os"

	"dynamiccontrol/internal/opa"
	"dynamiccontrol/internal/router"
	"dynamiccontrol/internal/validator"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Dynamic Control Plane Server...")

	// Initialize components
	policyManager := opa.NewPolicyManager()
	schemaValidator := validator.NewSchemaValidator()
	routeManager := router.NewRouteManager(policyManager, schemaValidator)

	// Load policies
	policiesDir := "policies"
	if err := policyManager.LoadPolicies(policiesDir); err != nil {
		log.Printf("Warning: Failed to load policies: %v", err)
	} else {
		log.Printf("Loaded policies: %v", policyManager.ListLoadedPolicies())
	}

	// Load route configuration
	configPath := "config/routes.json"
	if err := routeManager.LoadConfig(configPath); err != nil {
		log.Fatalf("Failed to load route configuration: %v", err)
	}

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "dynamic-control-plane",
		})
	})

	// Register dynamic routes
	if err := routeManager.RegisterRoutes(router); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Add info endpoint
	router.GET("/info", func(c *gin.Context) {
		config := routeManager.GetConfig()
		policies := policyManager.ListLoadedPolicies()

		c.JSON(200, gin.H{
			"service":  "Dynamic Control Plane",
			"version":  "1.0.0",
			"routes":   len(config.Routes),
			"policies": policies,
			"endpoints": []string{
				"GET /health - Health check",
				"GET /info - Service information",
				"GET /v1/status - Service status",
				"POST /v1/services/:serviceId/traffic - Traffic management",
			},
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
