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
	dbUrl := getEnv("DATABASE_URL", "postgres://doadmin:AVNS_cSBdfrfFUA3x699eaiP@janis-software-postgres-database-do-user-17670141-0.k.db.ondigitalocean.com:25060/erp_user_service?sslmode=require")
	jwtSecret := getEnv("JWT_SECRET", "jwt_supersecret")
	if dbUrl == "" {
		log.Fatal("Missing DB URL")
	}
	if jwtSecret == "" {
		log.Fatal("Missing JWT Secret")
	}
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
	api.Post("/logout", authHandler.Logout)
	log.Println("Auth Service running on port 8080") // Logs a message
	log.Fatal(app.Listen(":8080"))                   // Starts the Fiber application to listen for incoming HTTP connections on port 8080
	// app.Listen returns an error if the server fails (e.g. port is already in use) and log.Fatal ensures the error is printed and the program terminates if listening fails
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
