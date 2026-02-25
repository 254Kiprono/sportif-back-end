package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"webuye-sportif/app/config"
	"webuye-sportif/app/database"
	"webuye-sportif/app/loggers"
	"webuye-sportif/app/middleware"
	"webuye-sportif/app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Set Timezone
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Africa/Nairobi"
	}

	var loc *time.Location
	if tz == "Africa/Nairobi" {
		loc = time.FixedZone("EAT", 3*60*60)
	} else {
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("Warning: Failed to load timezone %s: %v", tz, err)
			loc = time.UTC
		}
	}

	time.Local = loc
	log.Printf("Timezone set to %s", tz)

	// Initialize logger before any background workers start
	env := os.Getenv("APP_ENV")
	if env == "" {
		if gin.Mode() == gin.ReleaseMode {
			env = "production"
		} else {
			env = "development"
		}
	}
	if err := loggers.InitLogger(env); err != nil {
		log.Printf("Logger init failed: %v", err)
	}

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
	}

	// Initialize Gin (using New to configure SkipPaths for Health Check logs)
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	r.Use(gin.Recovery())

	// Use CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes — pass database.Redis (will be nil if Redis failed to connect)
	routes.SetupRoutes(r, database.DB, cfg, database.Redis)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
