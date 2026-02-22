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
	// Root health check (handles both GET and HEAD for docker/wget compatibility)
	health := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "message": "Sportif Backend is running"})
	}
	r.GET("/health", health)
	r.HEAD("/health", health)

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
	sponsorRepo := repository.NewSponsorRepository(db)
	fanRepo := repository.NewFanRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

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
	sponsorService := services.NewSponsorService(sponsorRepo)
	fanService := services.NewFanService(fanRepo)
	paymentService := services.NewPaymentService(paymentRepo)

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
	sponsorHandler := handlers.NewSponsorHandler(sponsorService)
	fanHandler := handlers.NewFanHandler(fanService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// API Groups — Using root "/" because Nginx typically strips the "/api" prefix
	// before passing to the backend container.
	api := r.Group("/")
	{
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

		// Store (public)
		api.GET("/store/jerseys", storeHandler.GetJerseys)

		// News (public)
		api.GET("/news", newsHandler.GetAll)
		api.GET("/news/:id", newsHandler.GetByID)

		// Tickets (public)
		api.GET("/tickets", ticketHandler.GetAll)

		// Memberships (public)
		api.GET("/memberships/plans", membershipHandler.GetPlans)

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
			protected.PUT("/league/:id", middleware.RequirePermission("manage_league"), leagueHandler.UpdateEntry)
			protected.PUT("/store/jerseys/:id", middleware.RequirePermission("crud_jerseys"), storeHandler.Update)

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

			// Staff Routes (permissions)
			{
				protected.POST("/players", middleware.RequirePermission("crud_players"), playerHandler.Create)
				protected.PUT("/players/:id", middleware.RequirePermission("crud_players"), playerHandler.Update)
				protected.DELETE("/players/:id", middleware.RequirePermission("crud_players"), playerHandler.Delete)

				protected.POST("/fixtures", middleware.RequirePermission("crud_fixtures"), fixtureHandler.Create)
				protected.PUT("/fixtures/:id/score", middleware.RequirePermission("update_scores"), fixtureHandler.UpdateScore)

				protected.POST("/league", middleware.RequirePermission("manage_league"), leagueHandler.CreateEntry)

				protected.GET("/admin/news", middleware.RequirePermission("publish_news"), newsHandler.GetAllAdmin)
				protected.PUT("/news/:id", middleware.RequirePermission("publish_news"), newsHandler.Update)
				protected.DELETE("/news/:id", middleware.RequirePermission("delete_news"), newsHandler.Delete)

				protected.POST("/tickets", middleware.RequirePermission("manage_tickets"), ticketHandler.Create)
				protected.GET("/store/orders", middleware.RequirePermission("manage_orders"), storeHandler.GetOrders)
				protected.POST("/memberships/plans", middleware.RequirePermission("manage_membership"), membershipHandler.CreatePlan)
				protected.GET("/donations", middleware.RequirePermission("view_donations"), donationHandler.GetAll)
			}

			// Admin Only Routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", authHandler.GetAllUsers)
				admin.POST("/users", authHandler.Create)
				admin.GET("/sponsors", sponsorHandler.GetAll)
				admin.POST("/sponsors", sponsorHandler.Create)
				admin.PUT("/sponsors/:id", sponsorHandler.Update)
				admin.DELETE("/sponsors/:id", sponsorHandler.Delete)

				admin.GET("/fans", fanHandler.GetAll)
				admin.POST("/fans", fanHandler.Create)
				admin.PUT("/fans/:id", fanHandler.Update)
				admin.DELETE("/fans/:id", fanHandler.Delete)

				admin.GET("/payments", paymentHandler.GetAll)
				admin.POST("/payments", paymentHandler.Create)
				admin.PUT("/payments/:id", paymentHandler.Update)
				admin.DELETE("/payments/:id", paymentHandler.Delete)
			}
		}
	}
}

//Tetsing
