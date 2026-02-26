package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared/utils"
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	httptransport "github.com/RubenPari/clear-songs/internal/infrastructure/transport/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize environment and DI
	utils.LoadEnvVariables()

	// Switch to release mode in production based on env
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}

	// Initialize Database with Pooling
	if errConnectDb := postgres.Init(); errConnectDb != nil {
		log.Printf("WARNING: Database initialization failed: %v", errConnectDb)
	}

	// Setup Gin Router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "http://127.0.0.1:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Routes
	httptransport.SetUpRoutesRefactored(router, container)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
		// Good practice: enforce timeouts for server
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run server in a goroutine so that it doesn't block
	go func() {
		log.Println("Starting server on :3000")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Properly close database connections
	if postgres.Db != nil {
		sqlDB, err := postgres.Db.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}

	log.Println("Server exiting")
}
