package database

import (
	"log"
	"webuye-sportif/app/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	// 1. Permissions
	permissions := []models.Permission{
		{Name: "manage_users", Description: "Full user management"},
		{Name: "assign_roles", Description: "Assign roles to users"},
		{Name: "crud_players", Description: "Manage players"},
		{Name: "crud_fixtures", Description: "Manage fixtures"},
		{Name: "update_scores", Description: "Update match scores"},
		{Name: "manage_league", Description: "Manage league table"},
		{Name: "crud_jerseys", Description: "Manage jerseys in store"},
		{Name: "manage_orders", Description: "Manage store orders"},
		{Name: "manage_tickets", Description: "Manage ticket booking"},
		{Name: "manage_membership", Description: "Manage membership plans"},
		{Name: "view_donations", Description: "View all donations"},
		{Name: "publish_news", Description: "Publish or unpublish news"},
		{Name: "delete_news", Description: "Delete news articles"},
		{Name: "system_settings", Description: "Manage system configuration"},
		{Name: "create_news", Description: "Create news articles"},
		{Name: "edit_own_news", Description: "Edit own news drafts"},
		{Name: "upload_images", Description: "Upload images"},
		{Name: "view_news", Description: "View news articles"},
		{Name: "buy_merch", Description: "Purchase from store"},
		{Name: "buy_tickets", Description: "Purchase tickets"},
		{Name: "subscribe_membership", Description: "Subscribe to plans"},
		{Name: "make_donations", Description: "Donate to club"},
	}

	for i := range permissions {
		DB.FirstOrCreate(&permissions[i], models.Permission{Name: permissions[i].Name})
	}

	// 2. Roles
	// helper to get permissions by names
	getPerms := func(names []string) []models.Permission {
		var ps []models.Permission
		DB.Where("name IN ?", names).Find(&ps)
		return ps
	}

	roles := []models.Role{
		{
			Name:        "admin",
			Description: "Full system access",
			Permissions: permissions, // Admin gets everything
		},
		{
			Name:        "author",
			Description: "News content creator",
			Permissions: getPerms([]string{"create_news", "edit_own_news", "upload_images", "view_news"}),
		},
		{
			Name:        "cx",
			Description: "Content Executive",
			Permissions: getPerms([]string{
				"publish_news", "edit_own_news", "delete_news",
				"manage_league", "update_scores", "crud_jerseys",
				"view_news", "upload_images",
			}),
		},
		{
			Name:        "user",
			Description: "Fan account",
			Permissions: getPerms([]string{"view_news", "buy_merch", "buy_tickets", "subscribe_membership", "make_donations"}),
		},
	}

	for i := range roles {
		// Use a map to avoid updating permissions if role exists, or handle it properly
		var existing models.Role
		if err := DB.Preload("Permissions").Where("name = ?", roles[i].Name).First(&existing).Error; err != nil {
			DB.Create(&roles[i])
		} else {
			// Update permissions of existing role
			DB.Model(&existing).Association("Permissions").Replace(roles[i].Permissions)
			roles[i] = existing
		}
	}

	// 3. Admin User
	adminRole := roles[0] // Assuming admin is first in slice

	// Re-fetch to be sure
	DB.Where("name = ?", "admin").First(&adminRole)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		FullName: "System Admin",
		Username: "admin",
		Email:    "admin@sportif.com",
		Phone:    "0700000000",
		Password: string(hashedPassword),
		RoleID:   adminRole.ID,
	}
	DB.FirstOrCreate(&admin, models.User{Username: "admin"})

	// 4. Sample Players
	players := []models.Player{
		{Name: "John Doe", Position: "Striker", JerseyNumber: 9, Nationality: "Kenyan", Age: 25},
		{Name: "Jane Smith", Position: "Midfielder", JerseyNumber: 10, Nationality: "Kenyan", Age: 23},
	}
	for _, p := range players {
		DB.FirstOrCreate(&p, models.Player{Name: p.Name})
	}

	log.Println("Seeding completed")
}
