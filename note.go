package csnotes

import (
	"database/sql"
)

type Note struct {
	Resource
	Title string `json:"title"`
	Content sql.NullString `json:"content"`
	Time sql.NullString	`json:"time"`
	UserID int64 `json:"user_id"`
}

// NewNote creates a new note model with no ID or any fields set.
func NewNote(db *sql.DB) (n Note) {
	return Note {
		Resource: Resource {
			DB: db,
			Table: "notes",
		},
	}
}

// LoadNote attempts to load a note's fields from the database, given its ID.
func LoadNote(id int64, db *sql.DB) (n Note, err error) {
	n = NewNote(db)
	n.ID = id
	err = n.Load()

	return
}

func (n *Note) Load() error {
	return n.Select([]string{"title", "content", "time", "user_id"}, &n.Title, &n.Content, &n.Time, &n.UserID)
}

func (n *Note) Save() error {
	return n.Sync([]string{"title", "content", "time", "user_id"}, n.Title, n.Content.String, n.Time.String, n.UserID)
}

func (n *Note) User() (u User, err error) {
	// Load the tag's user from their ID.
	return LoadUser(n.UserID, n.DB)
}

func (n *Note) Tags() (ts []Tag, err error) {
	// Create empty slice of tags.
	ts = []Tag{}

	// Query for tags.
	rows, err := n.DB.Query("SELECT tag_id FROM note_tag WHERE note_id=?", n.ID)
	defer rows.Close()

	// If there was an error in the query, return nothing.
	if err != nil {
		return ts, err
	}

	for rows.Next() {
		var tID int64
		rows.Scan(&tID)
		if err != nil {
			// Skip this row if there was a problem scanning it.
			continue
		}

		t := Tag {
			Resource: Resource {
				ID: tID,
				DB: n.DB,
				Table: "tags",
			},
		}
		// TODO: Implement lazy-loading so this isn't done every time
		err = t.Load()
		if err != nil {
			continue
		}

		// Append the model.
		ts = append(ts, t)
	}

	return
}

// AddTag attaches a tag to this note.
func (n *Note) AddTag(id int64) error {
	// Insert a new row in the note_tag table.
	_, err := n.DB.Exec("INSERT INTO note_tag (note_id, tag_id) VALUES (?, ?)", n.ID, id)
	return err
}

// RemoveTag detaches a tag from this note.
func (n *Note) RemoveTag(id int64) error {
	// Remove a row from the note_tag table.
	_, err := n.DB.Exec("DELETE FROM note_tag WHERE note_id=? AND tag_id=?", n.ID, id)
	return err
}
