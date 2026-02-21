package routes

import (
	"log"
	"webuye-sportif/app/config"
	"webuye-sportif/app/handlers"
	"webuye-sportif/app/middleware"
	"webuye-sportif/app/repository"
	"webuye-sportif/app/services"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, rdb *redis.Client) {
	// Root health check (outside /api)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "message": "Sportif Backend is running"})
	})

	// Repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	fixtureRepo := repository.NewFixtureRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	newsRepo := repository.NewNewsRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	donationRepo := repository.NewDonationRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, roleRepo, cfg, rdb)
	playerService := services.NewPlayerService(playerRepo)
	storeService := services.NewStoreService(storeRepo)
	fixtureService := services.NewFixtureService(fixtureRepo)
	leagueService := services.NewLeagueService(leagueRepo)
	newsService := services.NewNewsService(newsRepo)
	ticketService := services.NewTicketService(ticketRepo)
	membershipService := services.NewMembershipService(membershipRepo)
	donationService := services.NewDonationService(donationRepo)

	// Storage Service (optional: graceful fallback if not configured)
	storageSvc, err := services.NewStorageService(cfg)
	if err != nil {
		log.Printf("Backblaze B2 not configured: %v — image upload endpoints will be unavailable", err)
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	playerHandler := handlers.NewPlayerHandler(playerService)
	storeHandler := handlers.NewStoreHandler(storeService)
	fixtureHandler := handlers.NewFixtureHandler(fixtureService)
	leagueHandler := handlers.NewLeagueHandler(leagueService)
	newsHandler := handlers.NewNewsHandler(newsService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	membershipHandler := handlers.NewMembershipHandler(membershipService)
	donationHandler := handlers.NewDonationHandler(donationService)

	// API Groups
	api := r.Group("/api")
	{
		// API health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "up", "api": "v1"})
		})

		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(cfg, rdb), authHandler.Logout)
		}

		// Players (public)
		api.GET("/players", playerHandler.GetAll)

		// Fixtures (public)
		api.GET("/fixtures", fixtureHandler.GetAll)

		// League (public)
		api.GET("/league", leagueHandler.GetTable)

		// News (public)
		api.GET("/news", newsHandler.GetAll)
		api.GET("/news/:id", newsHandler.GetByID)

		// Donations (public)
		api.POST("/donations", donationHandler.Donate)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg, rdb))
		{
			// Fan / User actions
			protected.POST("/store/order", storeHandler.PlaceOrder)
			protected.POST("/tickets/purchase", ticketHandler.Purchase)
			protected.POST("/memberships/subscribe", membershipHandler.Subscribe)
			protected.POST("/donations/member", donationHandler.Donate)

			// Author + Admin: Create news articles
			protected.POST("/news", middleware.RequireAnyRole([]string{"admin", "author"}), newsHandler.Create)

			// CX + Admin (via permission): content management actions
			protected.PUT("/news/:id/publish", middleware.RequirePermission("publish_news"), newsHandler.Update)
			protected.PUT("/league/:id", middleware.RequirePermission("edit_league"), leagueHandler.UpdateEntry)
			protected.PUT("/store/jerseys/:id", middleware.RequirePermission("manage_store"), storeHandler.Update)

			// Image Upload Routes (B2 S3)
			if storageSvc != nil {
				uploadHandler := handlers.NewUploadHandler(storageSvc)

				upload := protected.Group("/upload")
				{
					upload.POST("/news-image",
						middleware.RequireAnyRole([]string{"admin", "author"}),
						uploadHandler.UploadNewsImage,
					)

					upload.POST("/match-preview",
						middleware.RequirePermission("publish_news"),
						uploadHandler.UploadMatchPreview,
					)

					upload.POST("/match-photo",
						middleware.RequirePermission("publish_news"),
						uploadHandler.UploadMatchPhoto,
					)

					upload.DELETE("/image",
						middleware.RequireRole("admin"),
						uploadHandler.DeleteImage,
					)
				}
			}

			// Admin Only Routes
			admin := protected.Group("/")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.POST("/players", playerHandler.Create)
				admin.PUT("/players/:id", playerHandler.Update)
				admin.DELETE("/players/:id", playerHandler.Delete)

				admin.POST("/fixtures", fixtureHandler.Create)
				admin.PUT("/fixtures/:id/score", fixtureHandler.UpdateScore)

				admin.POST("/league", leagueHandler.CreateEntry)

				admin.GET("/admin/news", newsHandler.GetAllAdmin)
				admin.PUT("/news/:id", newsHandler.Update)
				admin.DELETE("/news/:id", newsHandler.Delete)

				admin.POST("/tickets", ticketHandler.Create)
				admin.POST("/memberships/plans", membershipHandler.CreatePlan)
				admin.GET("/donations", donationHandler.GetAll)

				// /api/admin/users → Admin only
				admin.GET("/admin/users", authHandler.GetAllUsers)
			}
		}
	}
}
