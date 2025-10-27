package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5" // A popular external library for creating and parsing JWTs in Go
)

func GenerateJWT(userID, secret string) (string, error) { // Function to create a new JWT which is used for authentication and authorization
	claims := jwt.MapClaims{ // Creates a map to hold the custom and standard claims (data payload) of the JWT
		"user_id": userID,                                // Custom claim: Identifer to embed in token's claims (Subject of the token)
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Standard claim: Sets the token's expiration time to 24 hours from now
		// Unix() converts the time.Time value into a Unix timestamp (seconds since epoch), which is the standard format for JWT claims
		"iat": time.Now().Unix(), // Standard claims: Sets the time the token was issued
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Creates a new JWT structure
	// jwt.SigningMethodHS256 specifies the algorithm used to sign the token (HMAC with SHA-256), a common symmetric signing method
	// claims: The payload defined above
	return token.SignedString([]byte(secret)) // Signs the token using the provided secret key. The secret must be converted to a byte slice.
	// This step generates the final, encoded JWT string (header, payload, signature)
}

func ValidateJWT(tokenString, secret string) (*jwt.Token, error) { // Defines a function to parse and validate an incoming JWT string
	// tokenString: The JWT string received (e.g. from an HTTP Authorization header)
	// secret: The cryptographic key used to verify the token's signature
	// Returns the validated token structure (containing the claims) and an error status
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // jwt.Parse attempts to decode the token string and verify its signature
		// The second argument is a KeyFunc, a function required to provide the secret key needed for verification
		// Validation check (Implicitly): It's highly recommended practice to explicitly check the signing method inside the KeyFunc for added security
		return []byte(secret), nil // The KeyFunc simply returns the secret key (converted to a byte slice) that the parser should use to verify the signature of the token
	})
}
