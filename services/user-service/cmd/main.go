package main // Entry point, required

import (
	"context"      // The 'context' package is used to carry cancellation, timeouts, and deadlines, ensuring database operations don't hang indefinitely
	"database/sql" // The standard Go interface for interacting with SQL databases
	"fmt"          // Used for formatted printing (e.g. to the console)
	"log"          // Used for logging errors and messages, typically to standard error, often used for critical failures like log.Fatalf
	"time"         // Used for time related operations, specifically for setting timeouts and connection lifetime

	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/config"     // Imports config package
	"github.com/Rishabh23Singh54/distributed-erp-system/services/user-service/internal/repository" // Imports repository package
	_ "github.com/jackc/pgx/v5/stdlib"                                                             // This is a blank import (using the underscore `_`). It imports the PostgreSQL driver (pgx) and executes its package level initialization code
	// Without making any of its functions or types directly accessible. This is necessary so the driver registers itself with the standard 'database/sql' package
)

func main() {
	dsn := config.LoadConfig().BuildDSN() // Loading from config

	db, err := sql.Open("pgx", dsn) // Attempts to open a database connection pool
	// pgx is the name the imported pgx driver registers itself under
	// dsn is the connection string built earlier
	// db is the connection pool object; err captures any error during initialization
	if err != nil {
		log.Fatalf("failed to open DB:%v", err)
	}

	db.SetMaxOpenConns(25)                 // Configures the connection pool: limits the total number of open database connections to 25
	db.SetMaxIdleConns(25)                 // Configures the connection pool: limits the number of idle (unused but kept open) connections to 25
	db.SetConnMaxLifetime(5 * time.Minute) // Configures the connection pool: sets a maximum time a connection can be reused (5 minutes). After this time, it's closed and new connections are opened as needed

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Creates a new context with 5 second timeout, ensuring the subsequent Ping operation doesn't wait forever
	// context.Background() is the root context
	// cancel is the function used to manually stop the operation early
	defer cancel() // Schedules the 'cancel' function to be executed just before the surrounding function (main) returns. This releases the context's resources

	if err := db.PingContext(ctx); err != nil { // Attempts to verify the connection to the database is alive and reachable using the created context
		log.Fatalf("Failed to ping DB: %v", err) // If the ping fails, log the error and terminate
	}

	fmt.Println("Connected to PostgreSQL successfully") // Prints a confirmation message to the console after a successful connection and ping

	userRepo := repository.NewUserRepository(db) // Initializes the UserRepository struct, passing the open database connection pool 'db' to it
	roleRepo := repository.NewRoleRepository(db) // Initializes the RoleRepository struct, passing the open database connection pool 'db' to it

	fmt.Printf("UserRepo and RoleRepo initialized: %+v,%+v\n", userRepo, roleRepo) // The '%+v' format specifier prints the struct with its field names, useful for debugging
}
