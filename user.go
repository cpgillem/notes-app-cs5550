package csnotes

import (
	"database/sql"
)

type User struct {
	Resource 
	Name sql.NullString `json:"name"`
	Username string `json:"username"`
	Admin bool `json:"admin"`
}

// CheckUsernameExists checks for a username in the database. If it exists,
// the function returns true.
func CheckUsernameExists(username string, db *sql.DB) (bool, error) {
	// Query the database, and if there is any row, return true.
	rows, err := db.Query("SELECT username FROM users WHERE username=?", username)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	return rows.Next(), nil
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

func LoadAllUsers(db *sql.DB) (us []User, err error) {
	// Initalize the slice.
	us = []User{}

	// Query the database for users.
	rows, err := db.Query("SELECT id, name, username, admin FROM users")
	if err != nil {
		return
	}

	// Load the users into the slice.
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name sql.NullString
		var username string
		var admin bool

		// Scan the data. If the user data could not be scanned, do not add
		// a new model.
		err = rows.Scan(&id, &name, &username, &admin)
		if err != nil {
			continue
		}

		// Create a new model.
		u := NewUser(db)
		u.ID = id
		u.Name = name
		u.Username = username
		u.Admin = admin

		// Add the user model.
		us = append(us, u)
	}

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

		err = rows.Scan(&nID)
		if err != nil {
			continue
		}

		n, err := LoadNote(nID, u.DB)
		if err != nil {
			continue
		}

		ns = append(ns, n)
	}

	return 
}

// Tags retrieves all the tags belonging to this user.
func (u *User) Tags() (ts []Tag, err error) {
	// Query the database for tags.
	rows, err := u.DB.Query("SELECT id FROM tags WHERE user_id = ?", u.ID)
	ts = []Tag{}

	defer rows.Close()
	for rows.Next() {
		var tID int64

		// Scan the ID from the tag row. If not possible, do not add the tag.
		err = rows.Scan(&tID)
		if err != nil {
			continue
		}

		// Create a new tag model. If not possible, do not add the tag.
		t, err := LoadTag(tID, u.DB)
		if err != nil {
			continue
		}

		// Add the tag to the slice.
		ts = append(ts, t)
	}

	return
}
