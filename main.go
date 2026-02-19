package main

import (
	"log"

	"webuye-sportif/app/config"
	"webuye-sportif/app/database"
	"webuye-sportif/app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	database.Connect(cfg)

	// Initialize Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, database.DB, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
