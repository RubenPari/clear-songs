/**
 * Clear Songs Backend - Main Entry Point
 * 
 * This is the main entry point for the Clear Songs backend server built with Go and Gin.
 * The application provides a RESTful API for managing Spotify playlists and tracks,
 * including features for bulk deletion, playlist management, and track organization.
 * 
 * Architecture:
 * - Gin web framework for HTTP routing and middleware
 * - PostgreSQL database for data persistence
 * - Redis cache for performance optimization
 * - Spotify OAuth 2.0 for authentication
 * 
 * Server Configuration:
 * - Runs on port 3000
 * - CORS enabled for frontend communication (localhost:4200)
 * - Debug mode enabled for development
 * 
 * Initialization Order:
 * 1. Configure CORS middleware
 * 2. Load environment variables
 * 3. Initialize Spotify OAuth
 * 4. Set up API routes
 * 5. Initialize cache manager
 * 6. Connect to database
 * 7. Start HTTP server
 * 
 * @package main
 * @author Clear Songs Development Team
 */
package main

import (
	"log"
	"time"

	cacheManager "github.com/RubenPari/clear-songs/src/cache"
	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/routes"
	"github.com/RubenPari/clear-songs/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/**
 * Main function - Application entry point
 * 
 * Initializes and starts the HTTP server with all necessary middleware,
 * routes, and services. The function performs the following operations:
 * 
 * 1. Sets Gin to debug mode for development (detailed error messages)
 * 2. Creates a default Gin router with logging and recovery middleware
 * 3. Configures CORS to allow frontend requests
 * 4. Loads environment variables from .env file
 * 5. Initializes Spotify OAuth configuration
 * 6. Sets up all API routes
 * 7. Initializes Redis cache manager
 * 8. Connects to PostgreSQL database
 * 9. Starts the HTTP server on port 3000
 * 
 * Error Handling:
 * - Database connection failures cause application panic
 * - Server startup failures cause application panic
 * 
 * @func main
 */
func main() {
	// Set Gin to debug mode for development
	// In production, this should be set to gin.ReleaseMode
	gin.SetMode(gin.DebugMode)
	
	// Create a Gin router with default middleware:
	// - Logger: Logs HTTP requests
	// - Recovery: Recovers from panics and returns 500 error
	server := gin.Default()

	/**
	 * CORS (Cross-Origin Resource Sharing) Configuration
	 * 
	 * Allows the Angular frontend (running on localhost:4200) to make
	 * requests to this backend API. CORS is essential for web applications
	 * where frontend and backend run on different origins.
	 * 
	 * Configuration:
	 * - AllowOrigins: Frontend URLs that can access the API
	 * - AllowMethods: HTTP methods allowed (GET, POST, PUT, DELETE, OPTIONS)
	 * - AllowHeaders: Request headers that can be sent
	 * - ExposeHeaders: Response headers that can be read by frontend
	 * - AllowCredentials: Allows cookies and authentication headers
	 * - MaxAge: How long preflight requests can be cached (12 hours)
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
	 * Environment and OAuth Initialization
	 * 
	 * These must be initialized before setting up routes because:
	 * - Routes depend on environment variables (API keys, database URLs)
	 * - OAuth endpoints need Spotify client credentials
	 */
	utils.LoadEnvVariables()
	utils.InitOAuth()

	/**
	 * Route Setup
	 * 
	 * Registers all API endpoints including:
	 * - Authentication routes (/auth/login, /auth/callback, /auth/logout)
	 * - Track management routes (/track/summary, /track/by-artist, /track/by-range)
	 * - Playlist management routes (/playlist/delete-tracks, etc.)
	 */
	routes.SetUpRoutes(server)

	/**
	 * Cache Manager Initialization
	 * 
	 * Initializes Redis cache for:
	 * - Caching API responses to reduce database queries
	 * - Storing session data
	 * - Improving application performance
	 */
	cacheManager.Init()

	/**
	 * Database Connection
	 * 
	 * Connects to PostgreSQL database. If connection fails, the application
	 * continues without database (backup functionality will be disabled).
	 * 
	 * The database stores:
	 * - Track backup data (optional)
	 * - User information (optional)
	 * 
	 * Note: The application can function without a database, but backup
	 * features will not be available.
	 */
	if errConnectDb := database.Init(); errConnectDb != nil {
		log.Printf("WARNING: Database initialization failed: %v", errConnectDb)
		log.Println("WARNING: Application will continue without database. Backup functionality disabled.")
	}

	/**
	 * Start HTTP Server
	 * 
	 * Starts the HTTP server on port 3000. If the server fails to start
	 * (e.g., port already in use), the application panics.
	 * 
	 * The server will listen for incoming HTTP requests and route them
	 * to the appropriate handlers based on the route configuration.
	 */
	if server.Run(":3000") != nil {
		panic("Error starting server")
	}
}
