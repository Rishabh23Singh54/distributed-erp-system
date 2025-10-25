package handlers // HTTP request handlers (controllers) that receive and respond to client requests. Interface between web framework and business logic (services)

import (
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/services"
	"github.com/gofiber/fiber/v2" // Imports the Fiber web framework, used for handling HTTP requests, routing and sending responses
)

type AuthHandler struct {
	AuthService *services.AuthService // A pointer to the AuthService which injects the business logic dependency. The handler uses the service to perform operations
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService} // Initializes and returns a pointer to the AuthHandler, injecting the AuthService dependency
}

// These structs define the expected JSON payload format for incoming HTTP requests
type SignupRequest struct {
	Name     string `json:"name"` // `json:"name"`: The struct tag ensures the Fiber BodyParser correctly maps the JSON field "name" from the request body to the Go struct field name
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req SignupRequest                      // Variable to hold the parsed request data
	if err := c.BodyParser(&req); err != nil { // Uses Fiber's BodyParser to automatically unmarshal the incoming JSON request body into the 'req' struct
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"}) // If parsing fiels (e.g. malformed JSON), returns a 400 Bad Request status with an error message
	}

	token, err := h.AuthService.Signup(c.Context(), req.Name, req.Email, req.Password, req.RoleID) // Calling Signup located in AuthService (core business logic)
	// c.Context() provides a context (with potential request deadlines/cancellation) to the service layer
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) // If service returns an error ("User already exists"), returns a 400 request with specific error message
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token}) // If successful, returns a 201 Created status and sends the newly generated JWT token back to the client
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := h.AuthService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
