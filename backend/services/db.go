package services

import (
	"canny-clone/repositories"
	"canny-clone/utils"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() *sql.DB {
	config := utils.GetConfig()
	
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	log.Println("Successfully connected to database")
	
	// Set the database connection in the repository
	repositories.SetDB(db)
	
	return db
}
