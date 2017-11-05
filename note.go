package csnotes

import (
	"database/sql"
)

type Note struct {
	Resource
	Title string
	Content sql.NullString
	Time sql.NullString
	UserID int64
}

func (n *Note) Load() error {
	return n.Select([]string{"title", "content", "time", "user_id"}, &n.Title, &n.Content, &n.Time, &n.UserID)
}

func (n *Note) Save() error {
	return n.Sync([]string{"title", "content", "time", "user_id"}, n.Title, n.Content.String, n.Time.String, n.UserID)
}

func (n *Note) User() (u User, err error) {
	// Create an unloaded model for the user.
	u = User {
		Resource: Resource {
			ID: n.UserID,
			DB: n.DB,
			Table: "users",
		},
	}

	// Define err as the result of loading the user from their ID.
	err = u.Load()

	return
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
