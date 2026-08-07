package main

import (
	database "konsera-backend/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("[ENV] Error loading .env file")
	}

	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("[Database] Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("[Database] Failed to ping database:", err)
	}

	log.Println("[Database] Database Connect")

	defer db.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	log.Fatal("[Server] Failed to start server:", router.Run(":8080"))
}
