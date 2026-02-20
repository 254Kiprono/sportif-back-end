package routes

import (
	"webuye-sportif/app/config"
	"webuye-sportif/app/handlers"
	"webuye-sportif/app/middleware"
	"webuye-sportif/app/repository"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
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
	authService := services.NewAuthService(userRepo, roleRepo, cfg)

	playerService := services.NewPlayerService(playerRepo)
	storeService := services.NewStoreService(storeRepo)
	fixtureService := services.NewFixtureService(fixtureRepo)
	leagueService := services.NewLeagueService(leagueRepo)
	newsService := services.NewNewsService(newsRepo)
	ticketService := services.NewTicketService(ticketRepo)
	membershipService := services.NewMembershipService(membershipRepo)
	donationService := services.NewDonationService(donationRepo)

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
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Players
		api.GET("/players", playerHandler.GetAll)

		// Fixtures
		api.GET("/fixtures", fixtureHandler.GetAll)

		// League
		api.GET("/league", leagueHandler.GetTable)

		// News
		api.GET("/news", newsHandler.GetAll)
		api.GET("/news/:id", newsHandler.GetByID)

		// Store (Public)
		// api.GET("/store/jerseys", storeHandler.GetJerseys) // To be implemented if needed

		// Donations
		api.POST("/donations", donationHandler.Donate)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// User specific
			protected.POST("/store/order", storeHandler.PlaceOrder)
			protected.POST("/tickets/purchase", ticketHandler.Purchase)
			protected.POST("/memberships/subscribe", membershipHandler.Subscribe)
			protected.POST("/donations/member", donationHandler.Donate)

			// Admin/Content Routes
			// /api/news/create → Author + Admin
			protected.POST("/news", middleware.RequireAnyRole([]string{"admin", "author"}), newsHandler.Create)

			// /api/news/publish → CX + Admin
			protected.PUT("/news/:id/publish", middleware.RequirePermission("publish_news"), newsHandler.Update) // Assuming CX has this permission

			// /api/league/update → CX + Admin
			protected.PUT("/league/:id", middleware.RequirePermission("edit_league"), leagueHandler.UpdateEntry)

			// /api/store/jerseys/update → CX + Admin
			protected.PUT("/store/jerseys/:id", middleware.RequirePermission("manage_store"), storeHandler.Update)

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
				admin.GET("/admin/users", authHandler.GetAllUsers) // Assuming this exists or will be added
			}
		}
	}
}
