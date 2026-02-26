package database

import (
	"fmt"
	"log"
	"time"

	"webuye-sportif/app/config"
	"webuye-sportif/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	DB = db
	fmt.Println("Database connected successfully")

	AutoMigrate()
}

func AutoMigrate() {
	// Temporarily disable foreign key checks to allow structural changes
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	// Drop tables with incompatible foreign keys for Player ID (UUID vs UINT)
	// This is necessary because changing the type of a referenced column is not allowed while constraints exist.
	// Since data needs to be re-entered anyway with the new ID format, dropping is safest.
	DB.Migrator().DropTable("match_events", "lineup_players", "lineups")

	err := DB.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.Player{},

		&models.Fixture{},
		&models.LeagueTable{},
		&models.News{},
		&models.Jersey{},
		&models.Order{},
		&models.OrderItem{},
		&models.Ticket{},
		&models.TicketOrder{},
		&models.MembershipPlan{},
		&models.MembershipOrder{},
		&models.Donation{},
		&models.Sponsor{},
		&models.Fan{},
		&models.Payment{},
		&models.MatchEvent{},
		&models.Lineup{},
		&models.LineupPlayer{},
	)

	// Re-enable foreign key checks
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed")
	Seed()
}
