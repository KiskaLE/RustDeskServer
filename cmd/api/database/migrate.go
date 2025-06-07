package database

import (
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	// create database if not exist
	// createDatabase()
	// Migrate the schema
	println("Startin db migration")
	db.AutoMigrate(Computers{})
	println("Migration done")

}
