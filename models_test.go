package main

import (
	"testing"
	"database/sql"
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
		Content: sql.NullString{String: "content", Valid: true},
		Time: sql.NullString{String: "2017-10-01 12:00", Valid: true},
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

// TestUserNotes ensures that you can retrieve all notes owned by a user.
func TestUserNotes(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Insert a user.
	res, err := db.Exec("INSERT INTO users (name, admin) VALUES ('test', false)")
	if err != nil {
		t.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create a user model.
	user := User {
		Resource: Resource {
			ID: id, 
			DB: db,
			Table: "users",
		},
		Name: "name",
		Admin: false,
	}

	// Insert some notes.
	_, err = db.Exec("INSERT INTO notes (title, user_id) VALUES ('title1', ?)", id)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO notes (title, user_id) VALUES ('title2', ?)", id)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the note models.
	notes, err := user.Notes()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual("title1", notes[0].Title, t)
	AssertEqual("title2", notes[1].Title, t)
}
