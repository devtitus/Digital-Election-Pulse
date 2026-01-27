package main

import (
	"log"
	"os"

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

	// --- CRITICAL FIX FOR RENDER START ---
	// Get the PORT from Render's environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default fallback for local development
	}

	// Start server on the dynamic port
	log.Println("Server is starting on port " + port)
	log.Fatal(app.Listen(":" + port))
	// --- CRITICAL FIX END ---
}
