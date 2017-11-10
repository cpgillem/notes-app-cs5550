package csnotes

import (
	"database/sql"
)

type Tag struct {
	Resource
	Title string `json:"title"`
	UserID int64 `json:"-"`
}

func (t *Tag) Load() error {
	return t.Select([]string{"title", "user_id"}, &t.Title, &t.UserID)
}

func (t *Tag) Save() error {
	return t.Sync([]string{"title", "user_id"}, t.Title, t.UserID)
}

func (t *Tag) Notes() (ns []Note, err error) {
	// Create an empty slice of notes.
	ns = []Note{}

	// Query for notes with this tag's ID.
	rows, err := t.DB.Query("SELECT note_id FROM note_tag WHERE tag_id=?", t.ID)
	defer rows.Close()

	// If there was an error in the query, return nothing.
	if err != nil {
		return ns, err
	}

	for rows.Next() {
		// Get the ID.
		var nID int64
		rows.Scan(&nID)
		if err != nil {
			continue
		}

		// Create a note model.
		n := Note {
			Resource: Resource {
				ID: nID,
				DB: t.DB,
				Table: "notes",
			},
		}
		err = n.Load()
		if err != nil {
			continue
		}

		// Append the model.
		ns = append(ns, n)
	}

	return
}

func (t *Tag) User() (u User, err error) {
	// Load the user from the ID.
	return LoadUser(t.UserID, t.DB)
}

func NewTag(db *sql.DB) (t Tag) {
	return Tag {
		Resource: Resource {
			DB: db,
			Table: "tags",
		},
	}
}

func LoadTag(id int64, db *sql.DB) (t Tag, err error) {
	t = NewTag(db)
	t.ID = id
	err = t.Load()

	return
}
