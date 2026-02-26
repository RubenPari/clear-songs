/**
 * Clear Songs Backend - Main Entry Point (Refactored)
 * 
 * This is the refactored main entry point that uses Dependency Injection
 * instead of global variables, following Clean Architecture principles.
 * 
 * @package main
 * @author Clear Songs Development Team
 */
package main

import (
	"log"
	"time"

	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/transport/http"
	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/**
 * Main function - Application entry point (Refactored)
 * 
 * Initializes the application using Dependency Injection:
 * 1. Loads environment variables
 * 2. Creates DI Container with all dependencies
 * 3. Sets up routes with injected dependencies
 * 4. Initializes database (optional)
 * 5. Starts HTTP server
 */
func main() {
	// Set Gin to debug mode for development
	gin.SetMode(gin.DebugMode)
	
	// Create a Gin router with default middleware
	server := gin.Default()

	/**
	 * CORS Configuration
	 */
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "http://127.0.0.1:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	/**
	 * Load Environment Variables
	 * 
	 * Must be done before creating DI container as it needs env vars
	 */
	utils.LoadEnvVariables()

	/**
	 * Initialize Dependency Injection Container
	 * 
	 * This creates all dependencies (repositories, use cases, etc.)
	 * and makes them available through the container.
	 */
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	log.Println("DI Container initialized successfully")

	/**
	 * Set Up Routes with DI
	 * 
	 * Routes are configured with dependencies injected from the container.
	 * This eliminates the need for global variables.
	 */
	http.SetUpRoutesRefactored(server, container)

	/**
	 * Database Connection (Optional)
	 * 
	 * Database is optional - application can function without it.
	 * Only backup functionality will be disabled.
	 */
	if errConnectDb := postgres.Init(); errConnectDb != nil {
		log.Printf("WARNING: Database initialization failed: %v", errConnectDb)
		log.Println("WARNING: Application will continue without database. Backup functionality disabled.")
	}

	/**
	 * Start HTTP Server
	 */
	log.Println("Starting server on :3000")
	if err := server.Run(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
