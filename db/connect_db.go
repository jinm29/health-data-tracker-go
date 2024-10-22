package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/Fidel-wole/wearable-integration/db/sqlc" // Import the generated sqlc package
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var Queries *db.Queries // This will hold the instance of Queries
var DB *sql.DB            // Database connection

// InitDB initializes the database connection and the SQLC Queries instance
func InitDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("DB_URL")
	
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	// Initialize the Queries instance with the connected DB
	Queries = db.New(DB)

	log.Println("Successfully connected to the database")
}

// GetQueries returns the initialized Queries instance
func GetQueries() *db.Queries {
	return Queries
}

// GetDB returns the initialized database connection
func GetDB() *sql.DB {
	return DB
}
