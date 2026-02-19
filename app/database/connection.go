package database

import (
	"fmt"
	"log"

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

	DB = db
	fmt.Println("Database connected successfully")

	AutoMigrate()
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Player{},
		&models.Fixture{},
		&models.LeagueTable{},
		&models.News{},
		&models.Jersey{},
		&models.Order{},
		&models.OrderItem{},
		&models.Ticket{},
		&models.MembershipPlan{},
		&models.MembershipOrder{},
		&models.Donation{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed")
	Seed()
}
