package csnotes

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Resource 
	Name sql.NullString `json:"name"`
	Username string `json:"username"`
	Admin bool `json:"admin"`
}

func NewUser(db *sql.DB) (u User) {
	return User {
		Resource: Resource {
			DB: db,
			Table: "users",
		},
	}
}

func LoadUser(id int64, db *sql.DB) (u User, err error) {
	u = NewUser(db)
	u.ID = id
	err = u.Load()

	return
}

// ValidateUser takes a username and password and attempts to load a 
// user model from this information. If the user could not be found, or if
// the password is incorrect, an error is returned.
func ValidateUser(username, password string, db *sql.DB) (u User, err error) {
	// Create a new user model with an ID. If the user was not found,
	// return an empty user and and error.
	u = NewUser(db)
	row := db.QueryRow("SELECT id FROM users WHERE username=?", username)
	err = row.Scan(&u.ID)
	if err != nil {
		return u, err
	}

	// Validate the user's password. If the password is not valid, do not load 
	// the model but return an error.
	valid, err := CheckPassword(u.ID, password, db)
	if !valid || err != nil {
		return u, err
	}

	// Load the user model.
	err = u.Load()
	if err != nil {
		return u, err
	}

	return
}

func (u *User) Load() error {
	return u.Select([]string{"username", "name", "admin"}, &u.Username, &u.Name, &u.Admin)
}

func (u *User) Save() error {
	return u.Sync([]string{"username", "name", "admin"}, u.Username, u.Name, u.Admin)
}

func (u *User) Notes() (ns []Note, err error) {
	rows, err := u.DB.Query("SELECT id FROM notes WHERE user_id = ?", u.ID)
	ns = []Note{}

	defer rows.Close()
	for rows.Next() {
		var nID int64

		err := rows.Scan(&nID)
		if err != nil {
			continue
		}

		n := Note {
			Resource: Resource {
				ID: nID,
				DB: u.DB,
				Table: "notes",
			},
		}

		err = n.Load()
		if err != nil {
			continue
		}

		ns = append(ns, n)
	}

	return 
}

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

type Tag struct {
	Resource
	Title string
}

func (t *Tag) Load() error {
	return t.Select([]string{"title"}, &t.Title)
}

func (t *Tag) Save() error {
	return t.Sync([]string{"title"}, t.Title)
}

func (t *Tag) Notes() (ns []Note, err error) {
	// Create an empty slice of notes.
	ns = []Note{}

	// Query for notes with this tag's ID.
	rows, err := t.DB.Query("SELECT note_id FROM note_tag WHERE tag_id=?", t.ID)
	defer rows.Close()

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
