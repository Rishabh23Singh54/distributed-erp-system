package main

import (
	"log"
	"os"

	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/db"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/handlers"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/repository"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/services"
	"github.com/gofiber/fiber/v2" // Handles HTTP routing and serving
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dbUrl == "" || jwtSecret == "" {
		log.Fatal("Missing environment variables")
	}

	pool := db.ConnectDB(dbUrl)
	defer pool.Close() // Schedules the pool.Close() method to be called when the surrounding function (main) exits, ensuring the database connections are properly released

	// Dependency injection
	authRepo := repository.NewAuthRepository(pool)
	authService := services.NewAuthService(authRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New() // Creates a new Fiber application instance, which acts as a core web server

	api := app.Group("/api/v1/auth") // Creates a route group. All subsequent routes defined using 'api' will be prefixed with "/api/v1/auth"

	// POST Requests
	api.Post("/signup", authHandler.Signup)
	api.Post("/login", authHandler.Login)

	log.Println("Auth Service running on port 8080") // Logs a message
	log.Fatal(app.Listen(":8080"))                   // Starts the Fiber application to listen for incoming HTTP connections on port 8080
	// app.Listen returns an error if the server fails (e.g. port is already in use) and log.Fatal ensures the error is printed and the program terminates if listening fails
}
