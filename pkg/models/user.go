package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	RoleID    string    `json:"role_id" db:"role_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}

//  `json` --> Data Transfer/Serialization --> Standard library's `encoding/json` package --> Tells the API how to expose or consume this data over the network (HTTP). It controls the field names and serialization rules in JSON format.
//  `db` --> Data Persistence/Mapping --> Database libraries (e.g., `sqlx`, `pgx`) --> Tells the repository layer how to map the struct fields to the exact column names in the SQL database table.

// Even when the names are identical (e.g., `json:"id"` and `db:"id"`), both tags must be present because the JSON library does not read the `db` tag, and the database library does not read the `json` tag. They are independent instructions for independent components of your application.
