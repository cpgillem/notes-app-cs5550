package csnotes

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
	res, err := db.Exec("INSERT INTO users (username, admin) VALUES ('test', false)")
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

	AssertEqual("test", user.Username, t)
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
		Username: "name",
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

func TestTagNotes(t *testing.T) {
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	// Prepare the database.
	res, err := db.Exec("INSERT INTO tags (title) VALUES ('tag')")
	tID, err := res.LastInsertId()
	res, err = db.Exec("INSERT INTO users (name) VALUES ('test')")
	uID, err := res.LastInsertId()
	res, err = db.Exec("INSERT INTO notes (title, user_id) VALUES ('title1', ?)", uID)
	n1ID, err := res.LastInsertId()
	res, err = db.Exec("INSERT INTO notes (title, user_id) VALUES ('title2', ?)", uID)
	n2ID, err := res.LastInsertId()
	_, err = db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)", n1ID, tID)
	_, err = db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)", n2ID, tID)

	if err != nil {
		t.Fatal(err)
	}

	// Create mock model.
	tag := Tag {
		Resource: Resource {
			ID: tID,
			DB: db,
			Table: "tags",
		},
		Title: "tag",
	}

	// Retrieve the tag's notes.
	ns, err := tag.Notes()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual("title1", ns[0].Title, t)
	AssertEqual("title2", ns[1].Title, t)
}

func TestNoteTags(t *testing.T) {
	db, ids, err := SeededTestDB()
	defer TearDownDbTest(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock model.
	note := Note {
		Resource: Resource {
			ID: ids["note.note1"],
			DB: db,
			Table: "notes",
		},
	}

	// Retrieve the note model's tags.
	ts, err := note.Tags()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual(2, len(ts), t)
	AssertEqual("tag1", ts[0].Title, t)
	AssertEqual("tag2", ts[1].Title, t)
}

// TestUserTags retrieves tags that belong to a user
func TestUserTags(t *testing.T) {
	db, ids, err := SeededTestDB()
	defer TearDownDbTest(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock user model.
	user := User {
		Resource: Resource {
			ID: ids["user.nonadmin"],
			DB: db,
			Table: "users",
		},
	}

	// Retrieve the user's tags.
	ts, err := user.Tags()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual(2, len(ts), t)
	AssertEqual("tag1", ts[0].Title, t)
	AssertEqual("tag2", ts[1].Title, t)
}

func TestTagUser(t *testing.T) {
	// Create seeded database.
	db, ids, err := SeededTestDB()
	defer TearDownDbTest(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock tag model.
	tag := Tag {
		Resource: Resource {
			ID: ids["tags.tag1"],
			DB: db,
			Table: "tags",
		},
		Title: "tag1",
		UserID: ids["user.nonadmin"],
	}

	// Load the user from the tag.
	user, err := tag.User()
	if err != nil {
		t.Fatal(err)
	}

	AssertEqual(ids["user.nonadmin"], user.ID, t)
	AssertEqual("nonadmin", user.Username, t)
}

func TestNoteAddTag(t *testing.T) {
	// Create seeded database.
	db, ids, err := SeededTestDB()
	defer TearDownDbTest(db)
	if err != nil {
		t.Fatal(err)
	}

	// Add a new tag to the table for the nonadmin user.
	res, err := db.Exec("INSERT INTO tags (title, user_id) VALUES ('note3', ?)", ids["user.nonadmin"])
	if err != nil {
		t.Fatal(err)
	}
	
	// Get the ID of the tag.
	tID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock note model.
	note := Note {
		Resource: Resource {
			ID: ids["note.note1"],
			DB: db,
			Table: "notes",
		},
	}

	// Add the tag to one of the nonadmin user's notes.
	err = note.AddTag(tID)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure an entry now exists in the note_tag table.
	rows, err := db.Query("SELECT * FROM note_tag WHERE note_id=? AND tag_id=?", note.ID, tID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Errorf("No entry for note %d and tag %d.", note.ID, tID)
	}
}

func TestNoteRemoveTag(t *testing.T) {
	// Create seeded database.
	db, ids, err := SeededTestDB()
	defer TearDownDbTest(db)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock note model.
	note := Note {
		Resource: Resource {
			ID: ids["note.note1"],
			DB: db,
			Table: "notes",
		},
	}

	// Remove a tag from one of the nonadmin user's notes.
	err = note.RemoveTag(ids["tag.tag1"])
	if err != nil {
		t.Fatal(err)
	}

	// Make sure there is no entry in the note_tag table.
	rows, err := db.Query("SELECT * FROM note_tag WHERE note_id=? AND tag_id=?", note.ID, ids["tag.tag1"])
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Errorf("Entry for note %d and tag %d not deleted.", note.ID, ids["tag.tag1"])
	}
}
