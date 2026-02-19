package database

import (
	"log"
	"webuye-sportif/app/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	// 1. Admin User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		FullName: "System Admin",
		Username: "admin",
		Email:    "admin@sportif.com",
		Phone:    "0700000000",
		Password: string(hashedPassword),
		Role:     "admin",
	}
	// Using GORM for seeding is fine as it's a utility, but I'll use Raw for consistency if desired.
	// Actually, FirstOrCreate is very useful here. I'll stick to it unless raw is strictly required for seeding too.
	DB.FirstOrCreate(&admin, models.User{Username: "admin"})

	// 2. Sample Players
	players := []models.Player{
		{Name: "John Doe", Position: "Striker", JerseyNumber: 9, Nationality: "Kenyan", Age: 25},
		{Name: "Jane Smith", Position: "Midfielder", JerseyNumber: 10, Nationality: "Kenyan", Age: 23},
	}
	for _, p := range players {
		DB.FirstOrCreate(&p, models.Player{Name: p.Name})
	}

	// 3. Sample Membership Plans
	plans := []models.MembershipPlan{
		{Name: "Gold", Price: 1000, DurationMonths: 12, Benefits: "Free tickets, Jersey discount"},
		{Name: "Silver", Price: 500, DurationMonths: 6, Benefits: "Free tickets"},
	}
	for _, p := range plans {
		DB.FirstOrCreate(&p, models.MembershipPlan{Name: p.Name})
	}

	log.Println("Seeding completed")
}
