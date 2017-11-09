package csnotes

import (
	"testing"
	"database/sql"
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

type testNoteModel struct {
	Resource
	Title string
	Content sql.NullString
}

func (m *testNoteModel) Load() error {
	return m.Select([]string{"title", "content"}, &m.Title, &m.Content)
}

func (m *testNoteModel) Save() error {
	return m.Sync([]string{"title", "content"}, m.Title, m.Content)
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

	// Make sure the user does not exist in the database.
	rows, err := db.Query("SELECT * FROM users WHERE id=?", id)
	defer rows.Close()

	for rows.Next() {
		t.Fatalf("A row was found for ID %v.", id)
	}
}

func TestResourceWithNulls(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Create a new resource and give it null content.
	res, err := db.Exec("INSERT INTO notes (title, user_id) VALUES ('test', 1)")
	if err != nil {
		t.Fatal(err)
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create a model for this resource.
	note := testNoteModel {
		Resource: Resource {
			ID: id,
			DB: db,
			Table: "notes",
		},
	}

	// Try to load its data.
	err = note.Load()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual("test", note.Title, t)
	AssertEqual(false, note.Content.Valid, t)
}
