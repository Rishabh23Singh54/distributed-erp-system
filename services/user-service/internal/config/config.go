package config // The name 'config' suggests this package is responsible for loading and managing application configuration settings

import (
	"fmt" // The 'fmt' package is used for formatted I/O, specifically here for constructing the database connection string (DSN)
	"os"  // The 'os' package provides an interface to operating system functionality, primarily used here to access environment variables
)

type Config struct { // The struct holds all the configuration settings needed for database connection
	DBHost     string // Hostname or IP address of the database server
	DBPort     string // Port number where the database server is listening
	DBUser     string // Username for the database connection
	DBPassword string // Password for the database connection
	DBName     string // Name of the specific database instance to connect to
}

func LoadConfig() *Config { // Function to load the config values
	return &Config{
		DBHost:     getEnv("DB_HOST", "janis-software-postgres-database-do-user-17670141-0.k.db.ondigitalocean.com"), // If DB_Host found in getEnv then it's value, else localhost
		DBPort:     getEnv("DB_PORT", "25060"),                                                                       // If DB_Port found in getEnv then it's value, else 5432 (standard postgres PORT)
		DBUser:     getEnv("DB_User", "doadmin"),                                                                     // If DB_User found in getEnv then it's value, else postgres
		DBPassword: getEnv("DB_Password", "AVNS_cSBdfrfFUA3x699eaiP"),                                                // If DB_Password found in getEnv then it's value, else postgres
		DBName:     getEnv("DB_Name", "erp_user_service"),                                                            // If DB_Name found in getEnv then it's value, else erp_user_service
	}
}

func getEnv(key, fallback string) string { // A utility function to safely retrieve environment variables with a fallback (default) value
	if value, exists := os.LookupEnv(key); exists { // os.LookupEnv(key) attempts to find the environment variable named 'key'. It returns a value and a boolean 'exists' indicating whether the variable was set
		return value // If exists, return value
	}
	return fallback // If environment variable is not set, the provided 'fallback' value is returned
}

func (c *Config) BuildDSN() string { // Method attached to config struct (receiver 'c'). It constructs the Database Source Name (DSN)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
