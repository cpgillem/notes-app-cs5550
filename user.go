package main

import (
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

// User represents a user of the site.
type User struct {
	db *sql.DB
	id int64
	name string
	admin bool
}

// NewUser creates a new user from scratch.
func NewUser(db *sql.DB) *User {
	// Create a new user with defaults. This will not be saved if any required
	// fields are still nil.
	return &User {
		db: db,
	}
}

// LoadUser creates a User instance from data in the database, given 
// a primary key.
func LoadUser(db *sql.DB, id int64) (user *User, err error) {
	// Create a new empty user model.
	user = NewUser(db)

	// Use the database connection to retrieve the user. If the user does not
	// exist, an error will be returned instead.
	err = db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.id, &user.name, &user.admin)

	if err != nil {
		user = nil
	}

	return
}

// Save persists any recently changed data to the database.
func (u *User) Save() error {
	err := error(nil)
	
	if u.id != 0 {
		// Update a user if they already exist.
		_, err = u.db.Exec(`UPDATE users SET name=?, admin=? WHERE id=?`, u.name, u.admin, u.id)
	} else {
		// If a user doesn't exist yet, insert a new record and save the
		// new ID.
		res, err := u.db.Exec(`INSERT INTO users (name, admin) VALUES (?, ?)`, u.name, u.admin)
		if err == nil { 
			u.id, err = res.LastInsertId()
		}
	}

	return err
}
