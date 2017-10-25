package main

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// SetUpDbTest sets up the database tables.
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

// TearDownDbTest tears down the database tables, removing all data.
func TearDownDbTest(testDB *sql.DB) {
	defer testDB.Close()
	TearDownDB(testDB)
}

func AssertEqual(expected interface{}, received interface{}, t *testing.T) {
	if expected != received {
		t.Errorf("Expected %v, received %v.", expected, received)
	}
}
