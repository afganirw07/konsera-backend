package main

import (
	"log"

	"github.com/joho/godotenv"

	"konsera-backend/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env file not found")
	}

	app, err := server.New()
	if err != nil {
		log.Fatal("[Server] Failed to initialize:", err)
	}

	defer app.DB.Close()

	log.Println("[Database] Database connected")
	log.Println("[Server] Starting server on :8080")

	if err := app.Run(); err != nil {
		log.Fatal("[Server] Failed to start:", err)
	}
}