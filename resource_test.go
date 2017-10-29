package main

import (
	"testing"
)

// A generic model that makes use of a Resource as an anonymous field.
type testModel struct {
	Resource
	Name string
	Admin bool
}

func (m *testModel) Load() error {
	err := m.Select([]string{"name", "admin"}, &m.Name, &m.Admin)

	return err
}

func (m *testModel) Save() error {
	err := m.Sync([]string{"name", "admin"}, m.Name, m.Admin)

	return err
}

func TestResourceSelect(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Insert a resource manually.
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES ('admin', true)")
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the ID.
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Define a resource manually.
	user := testModel {
		Resource: Resource {
			DB: db,
			ID: id,
			Table: "users",
		},
	}

	user.Columns = map[string]interface{} {
		"name": &user.Name,
		"admin": &user.Admin,
	}

	// Load the rest of the resource's data.
	err = user.Load()
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the resource data was loaded.
	AssertEqual("admin", user.Name, t)
	AssertEqual(true, user.Admin, t)
}

func TestResourceSyncNew(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	user := testModel {
		Resource: Resource {
			DB: db,
			Table: "users",
		},
	}

	// Set the values for the user.
	user.Name = "admin"
	user.Admin = true

	// Save the resource.
	err := user.Save()
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the id was set after storing the model.
	AssertUnequal(0, user.ID, t)

	// Query the database for the model.
	var retrievedName string
	var retrievedAdmin bool
	err = db.QueryRow("SELECT name, admin FROM users WHERE id=?", user.ID).Scan(&retrievedName, &retrievedAdmin)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the values are correct.
	AssertEqual("admin", retrievedName, t)
	AssertEqual(true, retrievedAdmin, t)
}

func TestResourceSyncExisting(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Insert the data manually.
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES ('nonadmin', false)")
	if err != nil {
		t.Fatal(err)
	}

	// Get the ID.
	newID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Make a model.
	user := testModel {
		Resource: Resource {
			DB: db,
			Table: "users",
			ID: newID,
		},
		Name: "admin",
		Admin: true,
	}

	// Sync the model.
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the database was updated accordingly.
	var retrievedName string
	var retrievedAdmin bool
	err = db.QueryRow("SELECT name, admin FROM users WHERE id=?", user.ID).Scan(&retrievedName, &retrievedAdmin)

	AssertEqual("admin", retrievedName, t)
	AssertEqual(true, retrievedAdmin, t)
}

func TestResourceDelete(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Insert a model manually.
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES ('test', false)")
	if err != nil {
		t.Fatal(err)
	}

	// Get the ID and create a model.
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	user := testModel {
		Resource: Resource {
			DB: db,
			ID: id,
			Table: "users",
		},
		Name: "test",
		Admin: false,
	}

	// Delete the user
	err = user.Delete()
	if err != nil {
		t.Fatal(err)
	}

	// The model's ID should have been reset to 0.
	if user.ID != 0 {
		t.Errorf("Expected ID to be 0, received %v.", user.ID)
	}

	// Make sure the user does not exist in the database.
	rows, err := db.Query("SELECT * FROM users WHERE id=?", id)
	defer rows.Close()

	for rows.Next() {
		t.Fatalf("A row was found for ID %v.", id)
	}
}
