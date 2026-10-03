package server

import (
	"database/sql"

	config "konsera-backend/internal/config"
	database "konsera-backend/internal/database"
	checkoutHandler "konsera-backend/internal/handler/checkout"
	communityHandler "konsera-backend/internal/handler/community"
	eventHandler "konsera-backend/internal/handler/event"
	notificationHandler "konsera-backend/internal/handler/notification"
	organizerHandler "konsera-backend/internal/handler/organizer"
	promoHandler "konsera-backend/internal/handler/promo"
	ticketHandler "konsera-backend/internal/handler/ticket"
	userHandler "konsera-backend/internal/handler/user"
	venueHandler "konsera-backend/internal/handler/venue"
	appMiddleware "konsera-backend/internal/middleware"
	checkoutRepository "konsera-backend/internal/repository/checkout"
	communityRepository "konsera-backend/internal/repository/community"
	eventRepository "konsera-backend/internal/repository/event"
	notificationRepository "konsera-backend/internal/repository/notification"
	organizerRepository "konsera-backend/internal/repository/organizer"
	promoRepository "konsera-backend/internal/repository/promo"
	ticketRepository "konsera-backend/internal/repository/ticket"
	userRepository "konsera-backend/internal/repository/user"
	venueRepository "konsera-backend/internal/repository/venue"
	checkoutService "konsera-backend/internal/services/checkout"
	communityService "konsera-backend/internal/services/community"
	email "konsera-backend/internal/services/email"
	eventService "konsera-backend/internal/services/event"
	notificationService "konsera-backend/internal/services/notification"
	organizerService "konsera-backend/internal/services/organizer"
	promoService "konsera-backend/internal/services/promo"
	ticketService "konsera-backend/internal/services/ticket"
	userService "konsera-backend/internal/services/user"
	venueService "konsera-backend/internal/services/venue"

	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router *gin.Engine
	DB     *sql.DB

	EmailService *email.Service
}

func New() (*Server, error) {
	// CONFIG
	emailConfig := config.LoadEmailConfig()

	// DATABASE
	db, err := database.ConnectDatabase()
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// EMAIL
	emailService := email.NewService(
		email.Config{
			SMTPHost:     emailConfig.SMTPHost,
			SMTPPort:     emailConfig.SMTPPort,
			SMTPUsername: emailConfig.SMTPUsername,
			SMTPPassword: emailConfig.SMTPPassword,
			FromName:     emailConfig.FromName,
			FromEmail:    emailConfig.FromEmail,
		},
	)

	// REPOSITORY
	userRepo := userRepository.NewUserRepository(db.DB)
	organizerRepo := organizerRepository.NewOrganizerRepository(db.DB)
	venueRepo := venueRepository.NewVenueRepository(db.DB)

	// SERVICE
	userService := userService.NewUserService(userRepo, emailService, db.DB)
	organizerService := organizerService.NewOrganizerService(organizerRepo)
	venueService := venueService.NewVenueService(venueRepo)

	// HANDLER
	userHandler := userHandler.NewUserHandler(userService)
	organizerHandler := organizerHandler.NewOrganizerHandler(organizerService)
	venueHandler := venueHandler.NewVenueHandler(venueService)
	eventRepo := eventRepository.NewEventRepository(db.DB)
	eventSvc := eventService.NewEventService(eventRepo)
	eventH := eventHandler.NewEventHandler(eventSvc)
	ticketRepo := ticketRepository.NewTicketRepository(db.DB)
	ticketSvc := ticketService.NewTicketService(ticketRepo)
	ticketH := ticketHandler.NewTicketHandler(ticketSvc)
	checkoutRepo := checkoutRepository.NewRepository(db.DB)
	checkoutSvc := checkoutService.NewService(checkoutRepo)
	checkoutH := checkoutHandler.NewHandler(checkoutSvc)
	promoRepo := promoRepository.NewRepository(db.DB)
	promoSvc := promoService.NewService(promoRepo)
	promoH := promoHandler.NewHandler(promoSvc)
	notificationRepo := notificationRepository.NewRepository(db.DB)
	notificationSvc := notificationService.NewService(notificationRepo)
	notificationH := notificationHandler.NewHandler(notificationSvc)
	communityRepo := communityRepository.NewRepository(db.DB)
	communitySvc := communityService.NewService(communityRepo)
	communityH := communityHandler.NewHandler(communitySvc)

	// ROUTER

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	router.Use(appMiddleware.SecurityHeaders())
	loginRateLimiter := appMiddleware.NewRateLimiter(5, time.Minute)
	authRateLimiter := appMiddleware.NewRateLimiter(20, time.Minute)
	organizerRateLimiter := appMiddleware.NewRateLimiter(60, time.Minute)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", loginRateLimiter.Middleware(), userHandler.Login)
		authGroup.POST("/register", authRateLimiter.Middleware(), userHandler.CreateUser)
		authGroup.POST("/verify-otp", authRateLimiter.Middleware(), userHandler.VerifyOTP)
		authGroup.POST("/verify-otp/:profile_id/:code", authRateLimiter.Middleware(), userHandler.VerifyOTPParams)
		authGroup.POST("/resend-otp", authRateLimiter.Middleware(), userHandler.ResendOTP)
		authGroup.POST("/users/preferences", appMiddleware.Auth(), userHandler.CreateUserPreference)
	}

	organizerGroup := router.Group("/organizers")
	organizerGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		organizerGroup.POST("", appMiddleware.RequireRole("customer", "organizer", "admin_event", "super_admin"), organizerHandler.Create)
		organizerGroup.GET("", organizerHandler.List)
		organizerGroup.GET("/:organizer_id", organizerHandler.Get)
		organizerGroup.PUT("/:organizer_id", appMiddleware.RequireRole("organizer", "admin_event", "super_admin"), organizerHandler.Update)
		organizerGroup.DELETE("/:organizer_id", appMiddleware.RequireRole("organizer", "admin_event", "super_admin"), organizerHandler.Delete)
		organizerGroup.POST("/:organizer_id/verifications", appMiddleware.RequireRole("organizer", "admin_event", "super_admin"), organizerHandler.CreateVerification)
		organizerGroup.GET("/:organizer_id/verifications", organizerHandler.ListVerifications)
		organizerGroup.GET("/:organizer_id/verifications/:verification_id", organizerHandler.GetVerification)
		organizerGroup.PUT("/:organizer_id/verifications/:verification_id", appMiddleware.RequireRole("admin_event", "super_admin"), organizerHandler.UpdateVerification)
		organizerGroup.DELETE("/:organizer_id/verifications/:verification_id", appMiddleware.RequireRole("admin_event", "super_admin"), organizerHandler.DeleteVerification)
	}

	venueGroup := router.Group("/venues")
	venueGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	venueMutationRoles := appMiddleware.RequireRole("organizer", "admin_event", "super_admin")
	{
		venueGroup.POST("", venueMutationRoles, venueHandler.Create)
		venueGroup.GET("", venueHandler.List)
		venueGroup.GET("/:venue_id", venueHandler.Get)
		venueGroup.PUT("/:venue_id", venueMutationRoles, venueHandler.Update)
		venueGroup.DELETE("/:venue_id", venueMutationRoles, venueHandler.Delete)
		venueGroup.POST("/:venue_id/sections", venueMutationRoles, venueHandler.CreateSection)
		venueGroup.GET("/:venue_id/sections", venueHandler.ListSections)
		venueGroup.GET("/:venue_id/sections/:section_id", venueHandler.GetSection)
		venueGroup.PUT("/:venue_id/sections/:section_id", venueMutationRoles, venueHandler.UpdateSection)
		venueGroup.DELETE("/:venue_id/sections/:section_id", venueMutationRoles, venueHandler.DeleteSection)
		venueGroup.POST("/:venue_id/sections/:section_id/seats", venueMutationRoles, venueHandler.CreateSeat)
		venueGroup.GET("/:venue_id/sections/:section_id/seats", venueHandler.ListSeats)
		venueGroup.GET("/:venue_id/sections/:section_id/seats/:seat_id", venueHandler.GetSeat)
		venueGroup.PUT("/:venue_id/sections/:section_id/seats/:seat_id", venueMutationRoles, venueHandler.UpdateSeat)
		venueGroup.DELETE("/:venue_id/sections/:section_id/seats/:seat_id", venueMutationRoles, venueHandler.DeleteSeat)
	}

	eventGroup := router.Group("/events")
	eventGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	eventMutationRoles := appMiddleware.RequireRole("organizer", "admin_event", "super_admin")
	{
		eventGroup.POST("", eventMutationRoles, eventH.Create)
		eventGroup.GET("", eventH.Events)
		eventGroup.GET("/:event_id", eventH.Get)
		eventGroup.PUT("/:event_id", eventMutationRoles, eventH.Update)
		eventGroup.DELETE("/:event_id", eventMutationRoles, eventH.Delete)
		eventGroup.POST("/:event_id/sessions", eventMutationRoles, eventH.CreateSession)
		eventGroup.GET("/:event_id/sessions", eventH.Sessions)
		eventGroup.GET("/:event_id/sessions/:session_id", eventH.GetSession)
		eventGroup.PUT("/:event_id/sessions/:session_id", eventMutationRoles, eventH.UpdateSession)
		eventGroup.DELETE("/:event_id/sessions/:session_id", eventMutationRoles, eventH.DeleteSession)
	}

	categoryGroup := router.Group("/event-categories")
	categoryGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		categoryGroup.POST("", eventMutationRoles, eventH.CreateCategory)
		categoryGroup.GET("", eventH.Categories)
		categoryGroup.GET("/:category_id", eventH.GetCategory)
		categoryGroup.PUT("/:category_id", eventMutationRoles, eventH.UpdateCategory)
		categoryGroup.DELETE("/:category_id", eventMutationRoles, eventH.DeleteCategory)
	}

	artistGroup := router.Group("/artists")
	artistGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		artistGroup.POST("", eventMutationRoles, eventH.CreateArtist)
		artistGroup.GET("", eventH.Artists)
		artistGroup.GET("/:artist_id", eventH.GetArtist)
		artistGroup.PUT("/:artist_id", eventMutationRoles, eventH.UpdateArtist)
		artistGroup.DELETE("/:artist_id", eventMutationRoles, eventH.DeleteArtist)
	}

	lineupGroup := router.Group("/sessions/:session_id/artists")
	lineupGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		lineupGroup.POST("", eventMutationRoles, eventH.AddArtist)
		lineupGroup.GET("", eventH.Lineup)
		lineupGroup.PUT("/:event_artist_id", eventMutationRoles, eventH.UpdateLineup)
		lineupGroup.DELETE("/:event_artist_id", eventMutationRoles, eventH.DeleteLineup)
	}

	ticketTierGroup := router.Group("/events/:event_id/ticket-tiers")
	ticketTierGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		ticketTierGroup.POST("", eventMutationRoles, ticketH.CreateTier)
		ticketTierGroup.GET("", ticketH.ListTiers)
		ticketTierGroup.GET("/:tier_id", ticketH.GetTier)
		ticketTierGroup.PUT("/:tier_id", eventMutationRoles, ticketH.UpdateTier)
		ticketTierGroup.DELETE("/:tier_id", eventMutationRoles, ticketH.DeleteTier)
	}

	inventoryGroup := router.Group("/ticket-tiers/:tier_id/inventory")
	inventoryGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		inventoryGroup.POST("", eventMutationRoles, ticketH.CreateInventory)
		inventoryGroup.GET("", ticketH.ListInventory)
		inventoryGroup.GET("/:event_session_id", ticketH.GetInventory)
		inventoryGroup.PUT("/:event_session_id", eventMutationRoles, ticketH.UpdateInventory)
		inventoryGroup.DELETE("/:event_session_id", eventMutationRoles, ticketH.DeleteInventory)
	}

	ticketGroup := router.Group("/tickets")
	ticketGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		ticketGroup.GET("", ticketH.ListOwned)
		ticketGroup.GET("/:ticket_id", ticketH.GetOwned)
		ticketGroup.PUT("/:ticket_id", ticketH.UpdateOwned)
		ticketGroup.POST("/:ticket_id/transfers", ticketH.CreateTransfers)
	}

	transferGroup := router.Group("/ticket-transfers")
	transferGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		transferGroup.GET("", ticketH.Transfers)
		transferGroup.PUT("/:transfer_id", ticketH.ResolveTransfer)
	}

	checkInGroup := router.Group("/check-ins")
	checkInGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware(), appMiddleware.RequireRole("gatekeeper", "admin_event", "super_admin"))
	{
		checkInGroup.POST("", ticketH.CreateCheckIn)
		checkInGroup.GET("", ticketH.ListCheckIns)
	}

	checkoutGroup := router.Group("/checkout")
	checkoutGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		checkoutGroup.GET("/carts", checkoutH.Carts)
		checkoutGroup.GET("/carts/:cart_id", checkoutH.Cart)
		checkoutGroup.POST("/carts", checkoutH.CreateCart)
		checkoutGroup.PUT("/carts/:cart_id", checkoutH.UpdateCart)
		checkoutGroup.DELETE("/carts/:cart_id", checkoutH.DeleteCart)
		checkoutGroup.POST("/bookings", checkoutH.CreateBooking)
		checkoutGroup.GET("/bookings", checkoutH.Bookings)
		checkoutGroup.GET("/bookings/:booking_id", checkoutH.Booking)
		checkoutGroup.PUT("/bookings/:booking_id", checkoutH.UpdateBooking)
		checkoutGroup.DELETE("/bookings/:booking_id", checkoutH.DeleteBooking)
		checkoutGroup.GET("/payment-methods", checkoutH.Methods)
		checkoutGroup.POST("/payments", checkoutH.CreatePayment)
		checkoutGroup.GET("/payments", checkoutH.Payments)
		checkoutGroup.PUT("/payments/:payment_id", appMiddleware.RequireRole("admin_event", "super_admin"), checkoutH.UpdatePayment)
		checkoutGroup.POST("/refunds", checkoutH.CreateRefund)
		checkoutGroup.GET("/refunds", checkoutH.Refunds)
		checkoutGroup.PUT("/refunds/:refund_id", appMiddleware.RequireRole("admin_event", "super_admin"), checkoutH.UpdateRefund)
		checkoutGroup.DELETE("/refunds/:refund_id", appMiddleware.RequireRole("admin_event", "super_admin"), checkoutH.DeleteRefund)
	}

	paymentMethodGroup := router.Group("/payment-methods")
	paymentMethodGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware(), appMiddleware.RequireRole("admin_event", "super_admin"))
	{
		paymentMethodGroup.POST("", checkoutH.CreateMethod)
		paymentMethodGroup.PUT("/:payment_method_id", checkoutH.UpdateMethod)
		paymentMethodGroup.DELETE("/:payment_method_id", checkoutH.DeleteMethod)
	}

	promoGroup := router.Group("/promo-codes")
	promoGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		promoGroup.GET("", promoH.List)
		promoGroup.GET("/:promo_id", promoH.Get)
		promoGroup.POST("", eventMutationRoles, promoH.Create)
		promoGroup.PUT("/:promo_id", eventMutationRoles, promoH.Update)
		promoGroup.DELETE("/:promo_id", eventMutationRoles, promoH.Delete)
	}

	checkoutPromoGroup := router.Group("/checkout/bookings/:booking_id/promo")
	checkoutPromoGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	checkoutPromoGroup.POST("", promoH.Apply)

	communityGroup := router.Group("/community")
	communityGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		communityGroup.GET("/reviews", communityH.Reviews)
		communityGroup.GET("/reviews/:review_id", communityH.GetReview)
		communityGroup.POST("/reviews", communityH.CreateReview)
		communityGroup.PUT("/reviews/:review_id", communityH.UpdateReview)
		communityGroup.DELETE("/reviews/:review_id", communityH.DeleteReview)
		communityGroup.GET("/favorites", communityH.Favorites)
		communityGroup.POST("/favorites", communityH.CreateFavorite)
		communityGroup.DELETE("/favorites/:event_id", communityH.DeleteFavorite)
	}

	financeAdminGroup := router.Group("/admin")
	financeAdminGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware(), appMiddleware.RequireRole("admin_event", "super_admin"))
	{
		financeAdminGroup.GET("/commissions", communityH.Commissions)
		financeAdminGroup.POST("/commissions", communityH.CreateCommission)
		financeAdminGroup.PUT("/commissions/:commission_id", communityH.UpdateCommission)
		financeAdminGroup.DELETE("/commissions/:commission_id", communityH.DeleteCommission)
		financeAdminGroup.GET("/payouts", communityH.Payouts)
		financeAdminGroup.POST("/payouts", communityH.CreatePayout)
		financeAdminGroup.PUT("/payouts/:payout_id", communityH.UpdatePayout)
		financeAdminGroup.DELETE("/payouts/:payout_id", communityH.DeletePayout)
		financeAdminGroup.GET("/audit-logs", communityH.Audits)
		financeAdminGroup.POST("/audit-logs", communityH.CreateAudit)
		financeAdminGroup.DELETE("/audit-logs/:audit_id", communityH.DeleteAudit)
	}

	notificationTemplateGroup := router.Group("/notification-templates")
	notificationTemplateGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	notificationAdminRoles := appMiddleware.RequireRole("admin_event", "super_admin")
	{
		notificationTemplateGroup.GET("", notificationH.ListTemplates)
		notificationTemplateGroup.GET("/:template_id", notificationH.GetTemplate)
		notificationTemplateGroup.POST("", notificationAdminRoles, notificationH.CreateTemplate)
		notificationTemplateGroup.PUT("/:template_id", notificationAdminRoles, notificationH.UpdateTemplate)
		notificationTemplateGroup.DELETE("/:template_id", notificationAdminRoles, notificationH.DeleteTemplate)
	}

	notificationGroup := router.Group("/notifications")
	notificationGroup.Use(appMiddleware.Auth(), organizerRateLimiter.Middleware())
	{
		notificationGroup.GET("", notificationH.List)
		notificationGroup.GET("/:notification_id", notificationH.Get)
		notificationGroup.POST("", notificationAdminRoles, notificationH.Create)
		notificationGroup.PUT("/:notification_id", notificationH.Update)
		notificationGroup.DELETE("/:notification_id", notificationH.Delete)
	}

	return &Server{
		Router:       router,
		DB:           db.DB,
		EmailService: emailService,
	}, nil
}

func (s *Server) Run() error {
	return s.Router.Run(":8080")
}
