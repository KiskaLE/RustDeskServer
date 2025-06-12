package database

import (
	"log"

	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	// create database if not exist
	// createDatabase()
	// Migrate the schema
	log.Println("Startin db migration")
	db.AutoMigrate(Computers{}, Users{})
	log.Println("Migration done")

}
