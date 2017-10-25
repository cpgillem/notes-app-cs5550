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
}

func TestResourceDelete(t *testing.T) {
}
