package repository

import (
	"context"

	"github.com/Rishabh23Singh54/distributed-erp-system/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password, role_id, is_active, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users(id, name, email, password, role_id, is_active, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, %6, NOW(), NOW(), NOW())`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.RoleID, user.IsActive)
	return err
}
