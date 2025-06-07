package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect establishes a connection to the specified SQLite database using GORM.
// It takes a string parameter 'database' which is the path to the SQLite database file.
// It returns a pointer to a gorm.DB instance representing the database connection.
// If the connection fails, the function will panic with the message "failed to connect database".
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./db/database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
