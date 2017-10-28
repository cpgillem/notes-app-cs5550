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

func TestResourceLoad(t *testing.T) {
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

func TestResourceSave(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Test saving a new resource that doesn't yet exist in the database.
	t.Run("New", func(t *testing.T) {
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
	})
}

func TestResourceDelete(t *testing.T) {
}
