package http

import (
	"github.com/Rishabh23Singh54/distributed-erp-system/pkg/middlewares"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/api/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, jwtSecret string, ph *handlers.ProfileHandler) {
	api := app.Group("/api/v1")
	userGroup := api.Group("/users")
	userGroup.Use(middlewares.JWTMiddleware(jwtSecret))
	userGroup.Get("/profile", ph.GetProfile)
}
