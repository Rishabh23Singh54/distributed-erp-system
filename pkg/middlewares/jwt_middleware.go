package middlewares

import (
	"strings" // Standard 'strings' package for manipulating strings, used here to parse Authorization header

	"github.com/gofiber/fiber/v2"  // Fiber web framework, providing the Handler type and Context object for middleware functionality
	"github.com/golang-jwt/jwt/v5" // JWT library for parsing and validating the token
)

func JWTMiddleware(secret string) fiber.Handler { // Factory function that takes the JWT signing secret and returns a Fiber middleware function
	return func(c *fiber.Ctx) error { // Returns the actual middleware function which executes before the intended route handler
		authHeader := c.Get("Authorization") // Rrtrieves the value of 'Authorization' HTTP header from the request
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"}) // Returns 401 if header is missing terminating the request chain
		}

		parts := strings.Split(authHeader, " ")                       // Splits the header value (expected format: "Bearer <token>") by the space character
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" { // Must be exactly two parts.. The first part should be "Bearer"
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header"})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Attempts to parse the JWT string and verify its signature
			return []byte(secret), nil // Returns the secret key used for signature verification
		})
		if err != nil || !token.Valid {
			// err != nil covers errors like signature mismatch, malformed token or parsing issues
			// !token.Valid covers built in claim checks, like expired ('exp') or not before ('nbf') times
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok { // If the token is valid, attempt to extract the claims (the payload data)
			if userID, exists := claims["user_id"].(string); exists { // Checks if the custom claim "user_id" exists and is a string
				c.Locals("user_id", userID) // Stores the extracted userID in Fiber's Locals context. This makes the userID easily accessible to subsequent handlers in the request chain without having to reparse the token
			}
		}
		return c.Next() // If the token is successfully validated and user information is stored in Locals, c.Next() calls the next handler in the route chain (the actual business logic)
	}
}
