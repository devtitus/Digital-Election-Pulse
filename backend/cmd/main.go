package main

import (
	"log"

	"election-pulse-backend/db"
	"election-pulse-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to Database
	db.Connect()

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())

	// Routes
	handlers.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
