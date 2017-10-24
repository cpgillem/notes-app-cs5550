package main

import (
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

// User represents a user of the site.
type User struct {
	Db *sql.DB
	Id int64
	Name string
	Admin bool
}

// NewUser creates a new user from scratch.
func NewUser(db *sql.DB) *User {
	// Create a new user with defaults. This will not be saved if any required
	// fields are still nil.
	return &User {
		Db: db,
	}
}

// LoadUser creates a User instance from data in the database, given 
// a primary key.
func LoadUser(db *sql.DB, id int64) (user *User, err error) {
	// Create a new empty user model.
	user = NewUser(db)

	// Use the database connection to retrieve the user. If the user does not
	// exist, an error will be returned instead.
	err = db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.Id, &user.Name, &user.Admin)

	if err != nil {
		user = nil
	}

	return
}

// Save persists any recently changed data to the database.
func (u *User) Save() error {
	err := error(nil)
	
	if u.Id != 0 {
		// Update a user if they already exist.
		_, err = u.Db.Exec(`UPDATE users SET name=?, admin=? WHERE id=?`, u.Name, u.Admin, u.Id)
	} else {
		// If a user doesn't exist yet, insert a new record and save the
		// new ID.
		res, err := u.Db.Exec(`INSERT INTO users (name, admin) VALUES (?, ?)`, u.Name, u.Admin)
		if err == nil { 
			u.Id, err = res.LastInsertId()
		}
	}

	return err
}
