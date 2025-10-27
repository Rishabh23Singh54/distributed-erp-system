package handlers // Contains the HTTP request heandlers (controllers) for the user-service

import (
	"context"
	"net/http"

	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(us *services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: us}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	uidLoc := c.Locals("user_id")
	if uidLoc == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthenticated"})
	}
	userID := uidLoc.(string)
	user, err := h.userService.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unauthenticated"})
	}
	if user == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	user.password = ""
	return c.Status(http.StatusOK).JSON(user)
}
