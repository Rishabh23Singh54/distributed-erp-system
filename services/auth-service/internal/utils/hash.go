package utils // Contains utility functions that can be reused throughout the application

import "golang.org/x/crypto/bcrypt" // bcrypt is a robust, widely-used and slow (computationally expensive) hashing algorithm designed specifically for securing passwords, making it resistant to brute-force attacks

func HashPassword(password string) (string, error) { // Function to convert a plaintext password into a secure hash
	// It returns the hashed password as a string and an error if the hashing fails
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// bcrypt.GenerateFromPassword performs the actual hashing process
	// []byte(password) converts the plaintext to byte slice
	// bcrypt.DefaultCost specifies the computational complexity (the "cost" or number of rounds) used for hashing.
	// A higher cost makes the function slower, which increases security against offline brute force attacks
	return string(bytes), err // Converts the resulting byte slice into a string for storage (in db) and returns the error status
}

func CheckPassword(hashedPassword, password string) bool { // Function to compare a stored hash against a user provided plaintext password during login
	// Returns boolean true if the password matches the hash, false otherwise
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) // Safely compares a bcrypt hash with a plaintext password
	// Handles the complex task of extracting the salt and cost from the hash, rehashing the plaintext, and then comparing the results, all in cryptographically secure, time constant manner
	return err == nil // Retuns true if error is nil (password matched) and false if it returns an error (mismatch or invalid hash)
}
