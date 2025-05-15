package repositories

import (
	"database/sql"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
