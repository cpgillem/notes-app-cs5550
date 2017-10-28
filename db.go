package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// SetUpDB creates all the tables necessary for the app.
func SetUpDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE users (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				name	VARCHAR(191) NOT NULL UNIQUE,
				admin	BOOLEAN DEFAULT FALSE NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`CREATE TABLE notes (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				title	VARCHAR(191) NOT NULL,
				content TEXT,
				time	DATETIME,
				user_id	INT(10) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE tags (
				id		INT(10) NOT NULL UNIQUE AUTO_INCREMENT,
				title	VARCHAR(191) NOT NULL,
				PRIMARY KEY (id)
			)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE note_tag (
				note_id	INT(10) NOT NULL,
				tag_id	INT(10) NOT NULL,
				PRIMARY KEY (note_id, tag_id)
			)`)
	if err != nil {
		return err
	}

	return nil
}

// TearDownDB clears the database of all tables that the app uses.
func TearDownDB(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS notes")
	db.Exec("DROP TABLE IF EXISTS tags")
	db.Exec("DROP TABLE IF EXISTS note_tag")
}
