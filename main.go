package main

import (
	"context"
	"log"
	"strconv"

	"webuye-sportif/app/config"
	"webuye-sportif/app/database"
	"webuye-sportif/app/routes"
	worker "webuye-sportif/app/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to MySQL
	database.Connect(cfg)

	// Connect to Redis (optional — app works without it, sessions just won't be whitelisted)
	redisDB, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		redisDB = 0
	}
	ctx := context.Background()
	if err = database.ConnectRedis(ctx, cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, redisDB); err != nil {
		log.Printf("Redis unavailable: %v — session whitelisting disabled", err)
	} else {
		defer database.CloseRedisConn()
		// Start background worker that clears Redis every 30 min (debug/cleanup)
		redisWorker := worker.NewRedisWorker(database.Redis)
		redisWorker.Start(ctx)
		defer redisWorker.Stop()
	}

	// Initialize Gin
	r := gin.Default()

	// Setup routes — pass database.Redis (will be nil if Redis failed to connect)
	routes.SetupRoutes(r, database.DB, cfg, database.Redis)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
