package main

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUser(t *testing.T) {
	// Open a database connection. This presumes that the testing database has
	// been created and that the user has access.
	testDB, err := sql.Open("mysql", "notes_app:notes_app@/notes_app_testing")
	if err != nil {
		t.Fatal(err)
	}
	
	if testDB == nil {
		t.Fatal("Could not create the database connection.")
	}

	// All this test needs to do is make sure that the original database
	// connection is passed to the new user model.
	t.Run("NewUser", func(t *testing.T) {
		user := NewUser(testDB)
		if user.db != testDB {
			t.Fatal("The database connections did not match.")
		}
	})

	// Tests creating a user and loading the data into a model struct.
	t.Run("LoadUser", func(t *testing.T) {
		SetUpDB(testDB)
		defer TearDownDB(testDB)

		// Manually create a user to test with.
		_, err = testDB.Exec("INSERT INTO users (id, name, admin) VALUES (1, 'test', false)")
		if err != nil {
			t.Fatal(err)
		}

		// Load the user.
		user, err := LoadUser(testDB, 1)

		if err != nil {
			t.Fatal(err)
		}

		if user.id != 1 {
			t.Errorf("Expected id 1, got %v", user.id)
		}

		if user.name != "test" {
			t.Errorf("Expected name 'test', got %v", user.name)
		}

		if user.admin != false {
			t.Errorf("Expected admin false, got %v", user.admin)
		}
	})

	// Ensures that all the data from a new user model will save.
	t.Run("SaveNew", func(t *testing.T) {
		SetUpDB(testDB)
		defer TearDownDB(testDB)

		// Test one with all values defined.
		user := User {
			db: testDB,
			name: "test",
			admin: true,
		}

		// This should store successfully.
		err := user.Save()
		
		if err != nil {
			t.Fatal(err)
		}

		// Query the database. The Save function should have stored the id.
		var name string
		var admin bool

		err = testDB.QueryRow("SELECT name, admin FROM users WHERE id=?", user.id).Scan(&name, &admin)

		if err != nil {
			t.Fatal(err)
		}

		if name != user.name {
			t.Errorf("Expected %v, got %v", user.name, name)
		}

		if admin != user.admin {
			t.Errorf("Expected %v, got %v", user.admin, admin)
		}
	})

	// This will test the modification of a user that already exists in the 
	// database.
	t.Run("Save", func(t *testing.T) {
		SetUpDB(testDB)
		defer TearDownDB(testDB)

		// Create a sample user and retrieve its id.
		res, err := testDB.Exec("INSERT INTO users (name, admin) VALUES ('test', true)")
		if err != nil {
			t.Fatal(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			t.Fatal(err)
		}

		// Create a corresponding model manually.
		user := User {
			db: testDB,
			id: id,
			name: "name-updated",
			admin: false,
		}

		// Save the model.
		err = user.Save()

		if err != nil {
			t.Fatal(err)
		}

		// Ensure that the new information is saved.
		var name string
		var admin bool
		err = testDB.QueryRow("SELECT name, admin FROM users WHERE id = ?", id).Scan(&name, &admin)

		if err != nil {
			t.Fatal(err)
		}

		if name != user.name {
			t.Errorf("Expected %v, got %v", user.name, name)
		}

		if admin != user.admin {
			t.Errorf("Expected %v, got %v", user.admin, admin)
		}
	})
}
