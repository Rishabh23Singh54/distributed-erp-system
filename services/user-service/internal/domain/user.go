// Here, we will define the data structures
package domain // Domain Data Driven Design

import "time"

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name`
	Email     string    `db:"email"`
	Password  string    `db:password` // hashed
	RoleID    int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
