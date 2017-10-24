package main

import (
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

// User represents a user of the site.
type User struct {
	Db *sql.DB
	ID int64
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
	err = db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Admin)

	if err != nil {
		user = nil
	}

	return
}

// LoadUsers loads a slice of Users from the database, given a query fragment.
// Query fragments are the part which follows WHERE, unless you want all models.
// Examples of query fragments:
// - id < 10 AND id > 2
// - name = "?" (succeeded by a string)
// - (empty string)
func LoadUsers(db *sql.DB, query string, vars ...interface{}) (users []User, err error) {
	// If the query is not empty, add WHERE to the beginning.
	queryString := query
	if len(query) > 0 {
		queryString = "WHERE " + queryString
	}

	// Set users to empty slice.
	users = []User{}

	rows, err := db.Query("SELECT id, name, admin FROM users " + queryString, vars...)

	if rows != nil {
		defer rows.Close() 
	}

	if err == nil {
		// If the query worked, load the results into models and return them.
		for rows.Next() {
			user := NewUser(db)
			if err := rows.Scan(&user.ID, &user.Name, &user.Admin); err == nil {
				users = append(users, *user)
			}
		}
	}

	return
}

// Save persists any recently changed data to the database.
func (u *User) Save() error {
	err := error(nil)
	
	if u.ID != 0 {
		// Update a user if they already exist.
		_, err = u.Db.Exec(`UPDATE users SET name=?, admin=? WHERE id=?`, u.Name, u.Admin, u.ID)
	} else {
		// If a user doesn't exist yet, insert a new record and save the
		// new ID.
		res, err := u.Db.Exec(`INSERT INTO users (name, admin) VALUES (?, ?)`, u.Name, u.Admin)
		if err == nil { 
			u.ID, err = res.LastInsertId()
		}
	}

	return err
}
