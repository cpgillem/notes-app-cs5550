package main

import (
	"testing"
)

// TestNoteUser ensures that the correct user can be retrieved from a note 
// model.
func TestNoteUser(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Insert a user.
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES ('test', false)")
	if err != nil {
		t.Fatal(err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create a note model.
	note := Note {
		Resource: Resource {
			DB: db,
			Table: "notes",
		},
		Title: "title",
		Content: "content",
		Time: "2017-10-01 12:00",
		UserID: userID,
	}

	// Get a model for the user that owns the note.
	user, err := note.User()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual("test", user.Name, t)
	AssertEqual(false, user.Admin, t)
}
