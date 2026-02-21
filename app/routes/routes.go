package routes

import (
	"log"
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

	// Cloudinary Service (optional: graceful fallback if not configured)
	cloudinarySvc, err := services.NewCloudinaryService(cfg)
	if err != nil {
		log.Printf("⚠️  Cloudinary not configured: %v — image upload endpoints will be unavailable", err)
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
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
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

		// ─── Protected routes (JWT required) ───────────────────────────────────────
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
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

			// ─── Image Upload Routes (Cloudinary) ─────────────────────────────────
			// Only registered if Cloudinary credentials are provided in .env
			if cloudinarySvc != nil {
				uploadHandler := handlers.NewUploadHandler(cloudinarySvc)

				upload := protected.Group("/upload")
				{
					// POST /api/upload/news-image → Author + Admin
					// Upload a featured image for a news article.
					// After upload, take the returned "url" and set it as news.image_url.
					upload.POST("/news-image",
						middleware.RequireAnyRole([]string{"admin", "author"}),
						uploadHandler.UploadNewsImage,
					)

					// POST /api/upload/match-preview → CX + Admin
					// Upload a pre-match promo/preview photo before the game kicks off.
					// After upload, take the returned "url" and set it as fixture.preview_image.
					upload.POST("/match-preview",
						middleware.RequirePermission("publish_news"),
						uploadHandler.UploadMatchPreview,
					)

					// POST /api/upload/match-photo → CX + Admin
					// Upload action/gallery photos after the match.
					// After upload, take the returned "url" and add it to fixture.match_photos JSON array.
					upload.POST("/match-photo",
						middleware.RequirePermission("publish_news"),
						uploadHandler.UploadMatchPhoto,
					)

					// DELETE /api/upload/image → Admin only
					// Remove an image from Cloudinary by public_id.
					upload.DELETE("/image",
						middleware.RequireRole("admin"),
						uploadHandler.DeleteImage,
					)
				}
			}

			// ─── Admin Only Routes ─────────────────────────────────────────────────
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
