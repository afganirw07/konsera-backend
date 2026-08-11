package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	config "konsera-backend/internal/config"
	database "konsera-backend/internal/database"
	userHandler "konsera-backend/internal/handler/user"
	userRepository "konsera-backend/internal/repository/user"
	email "konsera-backend/internal/services/email"
	userService "konsera-backend/internal/services/user"
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

	// SERVICE
	userService := userService.NewUserService(userRepo, emailService, db.DB)

	// HANDLER
	userHandler := userHandler.NewUserHandler(userService)

	// ROUTER

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authGroup := router.Group("/auth")
	{
		authGroup.POST(
			"/register",
			userHandler.CreateUser,
		)
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
