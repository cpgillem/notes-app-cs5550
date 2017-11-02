package csnotes

import (
	"database/sql"
	"strings"
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

// SeededTestDB Creates a database object referring to a database seeded with
// random test data. This is a way to approximate isolation in unit tests for
// non-database functions, since this function is the only point of possible 
// failure.
// This function returns a map of "useful names" to ID's, in order to expedite
// retrieval of relevant data. The syntax is ["(table name).(name/title/etc)."]
// e.g. ids["user.nonadmin"] => 1
//   or ids["note.note1"] => 2
//   or ids["tag.tag1"] => 3
func SeededTestDB() (db *sql.DB, ids map[string]int64) {
	db = SetUpDbTest()
	ids = map[string]int64{}

	// Users
	res, _ := db.Exec("INSERT INTO users (name, admin) VALUES (?, ?)", "nonadmin", false)
	ids["user.nonadmin"], _ = res.LastInsertId()

	res, _ = db.Exec("INSERT INTO users (name, admin) VALUES (?, ?)", "admin", true)
	ids["user.admin"], _ = res.LastInsertId()

	// Notes
	res, _ = db.Exec("INSERT INTO notes (title, content, time, user_id) VALUES (?, ?, ?, ?)",
		"note1", "content", "2017-01-01 12:00", ids["user.nonadmin"])
	ids["note.note1"], _ = res.LastInsertId()

	res, _ = db.Exec("INSERT INTO notes (title, content, time, user_id) VALUES (?, ?, ?, ?)",
		"note2", "content", "2017-02-01 12:00", ids["user.nonadmin"])
	ids["note.note2"], _ = res.LastInsertId()

	// Tags
	res, _ = db.Exec("INSERT INTO tags (title) VALUES (?)", "tag1")
	ids["tag.tag1"], _ = res.LastInsertId()

	res, _ = db.Exec("INSERT INTO tags (title) VALUES (?)", "tag2")
	ids["tag.tag2"], _ = res.LastInsertId()

	// Attach tags to notes. Note 1 will be tagged with tag 1, etc.
	db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)",
		ids["note.note1"], ids["tag.tag1"])

	db.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)",
		ids["note.note2"], ids["tag.tag2"])

	return 
}

func AssertEqual(expected interface{}, received interface{}, t *testing.T) {
	if expected != received {
		t.Errorf("Expected %v, received %v.", expected, received)
	}
}

func AssertUnequal(unexpected interface{}, received interface{}, t *testing.T) {
	if unexpected == received {
		t.Errorf("Did not expect %v to equal %v", received, unexpected)
	}
}

func AssertContains(haystack string, needle string, t *testing.T) {
	if !strings.Contains(haystack, needle) {
		t.Errorf("Could not find %s within %s.", needle, haystack)
	}
}
