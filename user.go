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
