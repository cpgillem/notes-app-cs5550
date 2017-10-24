package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// SetUpDB creates all the tables necessary for the app.
func SetUpDB(db *sql.DB) {
	db.Exec(`CREATE TABLE users (
				id	INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				name	VARCHAR(191) NOT NULL UNIQUE,
				admin	BOOLEAN DEFAULT FALSE NOT NULL,
				PRIMARY KEY (id)
			)`)
}

// TearDownDB clears the database of all tables that the app uses.
func TearDownDB(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS users")
}
