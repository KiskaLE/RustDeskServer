package db

import (
	"database/sql"
	"log"
)

func MigrateDatabase() {
	// create database if not exist
	createDatabase()
	db := Connect()
	// Migrate the schema
	println("Startin db migration")
	db.AutoMigrate(Computers{})
	println("Migration done")

}

func createDatabase() {
	println("Check if database exist")
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create empty database
	sqlStmt := `
    create table aTable(field1 int); drop table aTable;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
