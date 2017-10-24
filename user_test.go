package main

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	testDB := SetUpDbTest()
	defer TearDownDbTest(testDB)

	user := NewUser(testDB)
	if user.Db != testDB {
		t.Fatal("The database connections did not match.")
	}
}

func TestLoadUser(t *testing.T) {
	testDB := SetUpDbTest()
	defer TearDownDbTest(testDB)

	// Manually create a user to test with.
	_, err := testDB.Exec("INSERT INTO users (id, name, admin) VALUES (1, 'test', false)")
	if err != nil {
		t.Fatal(err)
	}

	// Load the user.
	user, err := LoadUser(testDB, 1)

	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Errorf("Expected id 1, got %v", user.ID)
	}

	if user.Name != "test" {
		t.Errorf("Expected name 'test', got %v", user.Name)
	}

	if user.Admin != false {
		t.Errorf("Expected admin false, got %v", user.Admin)
	}
}

func TestSaveNew(t *testing.T) {
	testDB := SetUpDbTest()
	defer TearDownDbTest(testDB)

	// Test one with all values defined.
	user := User {
		Db: testDB,
		Name: "test",
		Admin: true,
	}

	// This should store successfully.
	err := user.Save()
	
	if err != nil {
		t.Fatal(err)
	}

	// Query the database. The Save function should have stored the id.
	var name string
	var admin bool

	err = testDB.QueryRow("SELECT name, admin FROM users WHERE id=?", user.ID).Scan(&name, &admin)

	if err != nil {
		t.Fatal(err)
	}

	if name != user.Name {
		t.Errorf("Expected %v, got %v", user.Name, name)
	}

	if admin != user.Admin {
		t.Errorf("Expected %v, got %v", user.Admin, admin)
	}
}

func TestSave(t *testing.T) {
	testDB := SetUpDbTest()
	defer TearDownDbTest(testDB)

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
		Db: testDB,
		ID: id,
		Name: "name-updated",
		Admin: false,
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

	if name != user.Name {
		t.Errorf("Expected %v, got %v", user.Name, name)
	}

	if admin != user.Admin {
		t.Errorf("Expected %v, got %v", user.Admin, admin)
	}
}
