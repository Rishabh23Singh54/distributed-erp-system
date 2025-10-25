package services // This package typically contains the application's core business logic, coordinating data access (repository) and utility functions

import (
	"context"
	"errors"

	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/models"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/repository"
	"github.com/Rishabh23Singh54/distributed-erp-system/services/auth-service/internal/utils"
	"github.com/google/uuid"
)

type AuthService struct { // Defines the primary service structure for handling all authentication related business logic
	repo      *repository.AuthRepository // A pointer to the database repository, providing methods for data persistence
	jwtSecret string                     // The secret key required to sign and verify JWTs
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string) *AuthService { // Constructor function
	return &AuthService{ // Returns a pointer to the AuthService, injecting the required dependencies (repository and secret)
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Signup registers a new user
func (s *AuthService) Signup(ctx context.Context, name, email, password, roleID string) (string, error) { // Method attached to AuthService ('s')
	// Check if user exists
	existingUser, _ := s.repo.FindUserByEmail(ctx, email)
	if existingUser != nil {
		return "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		RoleID:   roleID,
		IsActive: true,
	}

	// Save user
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
