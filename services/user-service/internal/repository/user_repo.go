package repository // The name 'repository' suggests this package is responsible for handling data access and persistence, following the Repository pattern

import (
	"context"      // The 'context' package provides a way to carry deadlines, cancellation signals, and other request-scoped values across API boundaries. It's crucial for well-behaved database and networked applications in Go.
	"database/sql" // The 'database/sql' package is the standard Go interface for interacting with SQL databases.
	"errors"       // The standard 'errors' package is used for creating and handling errors.
	"time"         // The 'time' package provides functionality for measuring and displaying time. It is used here to get the current timestamp for 'created_at' and 'updated_at'.

	// The following import is assumed from the previous context, even if slightly simplified here:
	// "github.com/Rishabh23Singh54/distributed-erp-system/services/user-service" // Higher-level package import (often for interfaces)
	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/domain" // This imports the 'domain' package, which contains the core data structures like the User struct.
)

type UserRepository struct { // Defines a struct named 'UserRepository', the concrete implementation of the Repository pattern for the User entity.
	db *sql.DB // This field holds a reference to the standard library's database connection pool. It is the primary tool used by the repository methods to interact with the database.
}

func NewUserRepository(db *sql.DB) *UserRepository { // This is a constructor function (a convention in Go, not mandatory). It takes a pointer to a database connection pool (*sql.DB) as an argument.
	return &UserRepository{db: db} // It initializes and returns a pointer to a new UserRepository struct, embedding the provided database connection.
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) { // This is a method attached to the UserRepository struct (indicated by receiver 'r'). It implements the logic to persist a new user record in the database.
	// ctx context.Context: Allows the operation to be canceled or timed out externally.
	// user *domain.User: The user data to be created.
	// returns (*domain.User, error): The created user (updated with ID and timestamps) or an error.
	query := `INSERT INTO users (name, email, password, role_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at` // This variable holds the SQL query string for inserting a new row into the 'users' table. $1, $2, etc., are placeholders for values that will be provided safely.
	// RETURNING id, created_at, updated_at: This part is specific to PostgreSQL/certain SQL dialects and tells the database to return the generated ID and the timestamps right after insertion.
	now := time.Now() // Capture the current timestamp for 'created_at' and 'updated_at'.

	err := r.db.QueryRowContext(ctx, query, // This executes the SQL query. QueryRowContext is used because we expect exactly one row back (the RETURNING clause).
		// ctx: Passes the context to the database driver.
		// query: The SQL string.
		user.Name, // These arguments are the values to replace $1, $2, etc., respectively.
		user.Email,
		user.Password,
		user.RoleID,
		now, // $5: created_at
		now, // $6: updated_at
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt) // .Scan(...) is chained to QueryRowContext. It reads the returned values from the database (id, created_at, updated_at) and assigns them to the fields of the 'user' struct.
	// The & symbol is crucial; it passes pointers to the struct fields so their values can be updated.

	if err != nil { // Standard Go error handling. If the database operation failed (e.g., duplicate email constraint), the error is returned immediately.
		return nil, err
	}
	return user, nil // If successful, the 'user' struct now contains the generated ID and timestamps. It is returned along with a nil error, indicating success.
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) { // Retrieves a single user record based on their email address.
	query := `SELECT id, name, email, password, role_id, created_at, updated_at FROM users WHERE email=$1 AND deleted_at IS NULL` // SQL to select all non-deleted user fields where the email matches the placeholder $1.
	var user domain.User                                                                                                          // Declare a variable 'user' of type domain.User to hold the result.

	err := r.db.QueryRowContext(ctx, query, email).Scan( // Executes the query and attempts to scan the result row into the 'user' struct fields.
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // Checks if the error specifically indicates no rows were found.
			return nil, nil // If no user is found, returns (nil user, nil error), which is common for "not found" scenarios in repositories.
		}
		return nil, err // Returns any other database error.
	}
	return &user, nil // Returns a pointer to the found user.
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) { // Retrieves a single user record based on their ID.
	query := `SELECT id, name, email, password, role_id, created_at, updated_at FROM users WHERE id=$1 AND deleted_at IS NULL` // SQL to select all non-deleted user fields where the ID matches the placeholder $1.
	var user domain.User                                                                                                       // Variable to hold the retrieved user data.

	err := r.db.QueryRowContext(ctx, query, id).Scan( // Executes the query and attempts to scan the result row into the 'user' struct fields.
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // Checks for the "no rows found" error.
			return nil, nil // Returns (nil user, nil error) if the user is not found.
		}
		return nil, err // Returns any other database error.
	}
	return &user, nil // Returns a pointer to the found user.
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) { // Retrieves a list (page) of user records.
	query := `SELECT id, name, email, password, role_id, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY id ASC LIMIT $1 OFFSET $2` // SQL for paginated listing of non-deleted users, ordered by ID. $1 is LIMIT, $2 is OFFSET.

	rows, err := r.db.QueryContext(ctx, query, limit, offset) // Executes the query expected to return multiple rows.
	if err != nil {
		return nil, err // Returns an error if the query fails to execute.
	}
	defer rows.Close() // Ensures that the result set is closed when the function exits, preventing resource leaks.

	users := []*domain.User{} // Initializes an empty slice to store pointers to User structs.

	for rows.Next() { // Iterates through the rows returned by the query.
		var user domain.User // Declares a new user struct for each row.
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err // If scanning a row fails, return the error.
		}
		users = append(users, &user) // Appends a pointer to the newly scanned user to the slice.
	}

	// Note: It's good practice to check rows.Err() here for any error that occurred during iteration.
	// if err := rows.Err(); err != nil { return nil, err }

	return users, nil // Returns the slice of users and a nil error.
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error { // Updates an existing user record in the database.
	query := `UPDATE users SET name=$1, email=$2, password=$3, role_id=$4, updated_at=$5 WHERE id=$6 AND deleted_at IS NULL` // SQL to update specific fields for a non-deleted user based on their ID ($6).
	now := time.Now()                                                                                                        // Capture the current time for the 'updated_at' field.

	res, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.RoleID, now, user.ID) // Executes the update query. ExecContext is used for operations that don't return rows (INSERT, UPDATE, DELETE).
	if err != nil {
		return err // Returns an error if the query execution fails.
	}

	rowsAffected, err := res.RowsAffected() // Gets the number of rows modified by the UPDATE statement.
	if err != nil {
		return err // Returns an error if retrieving the rows affected count fails.
	}

	if rowsAffected == 0 {
		return errors.New("User not found") // If zero rows were affected, the user ID was likely not found or was already soft-deleted. Returns a custom error.
	}
	return nil // Returns nil on successful update.
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error { // Performs a **soft delete** on a user record by setting the 'deleted_at' timestamp.
	query := `UPDATE users SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL` // SQL for soft-deleting: sets 'deleted_at' to current time for the user ID ($2), but only if it hasn't been deleted already.
	now := time.Now()                                                            // Capture the current time to mark the deletion.

	res, err := r.db.ExecContext(ctx, query, now, id) // Executes the update query (soft delete).
	if err != nil {
		return err // Returns an error if the query execution fails.
	}

	rowsAffected, err := res.RowsAffected() // Gets the number of rows modified.
	if err != nil {
		return err // Returns an error if retrieving the rows affected count fails.
	}

	if rowsAffected == 0 {
		return errors.New("User not found") // If zero rows were affected, the user was not found or was already soft-deleted. Returns an error.
	}
	return nil // Returns nil on successful soft delete.
}
