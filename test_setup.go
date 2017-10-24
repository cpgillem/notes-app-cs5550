package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func SetUpDbTest() *sql.DB {
	// Open a database connection. This presumes that the testing database has
	// been created and that the user has access.
	newDB, err := sql.Open("mysql", "notes_app:notes_app@/notes_app_testing")
	if err != nil {
		panic(err)
	}
	SetUpDB(newDB)

	return newDB
}

func TearDownDbTest(testDB *sql.DB) {
	defer testDB.Close()
	TearDownDB(testDB)
}
