package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/domain"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	query := `INSERT INTO roles (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	now := time.Now()

	err := r.db.QueryRowContext(ctx, query, role.Name, now, now).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE id=$1`
	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT id, name, created_at, updated_at FROM roles ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*domain.Role{}
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name=$1, updated_at=$2 WHERE id=$3`
	now := time.Now()

	res, err := r.db.ExecContext(ctx, query, role.Name, now, role.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM roles WHERE id=$1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}
