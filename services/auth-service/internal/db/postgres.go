package db // This package name suggests the package's primary role is managing database connectivity

import (
	"context" // The 'context' package is used to handle timeouts and cancellation signals during database operations
	"log"     // The 'log' package is used for logging status messages and critical errors (like connection failures)

	"github.com/jackc/pgx/v5/pgxpool" // A highly regarded PostgreSQL specific connection pool in Go.
	// This package is often preferred over the standard 'database/sql' for pure PostgreSQL applications due to better performance and features
)

func ConnectDB(connString string) *pgxpool.Pool { // Establishes and tests database connections
	// connString string: Takes the database connection string (DSN) as input
	// pgxpool.Pool: Returns a pointer to a connection pool object, which manages multiple database connections
	pool, err := pgxpool.New(context.Background(), connString) // Attempts to create a new connection pool using the provided connection string
	// context.Background() is used before establishing the pool connection is typically a background setup task
	if err != nil { // If pool creation fails (e.g. malformed connection string, driver issue), logs the error and terminates the application
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = pool.Ping(context.Background()) // Attempts to send a quick verification message (a 'ping') to the database server using one of the pool's connections to ensure the server is reachable and responsive
	if err != nil {                       // If the ping fails (e.g. database is offline, firewall issue), logs the error and terminates the application
		log.Fatalf("DB ping failed: %v", err) // This check ensures a live connection before proceeding with the rest of the application setup
	}

	log.Println("Connected to PostgreSQL") // Logs a successful connection message to the console
	return pool                            // Returns the successfully configured and tested connection pool object
}
