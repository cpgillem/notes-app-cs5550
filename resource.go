package csnotes

import (
	"database/sql"
	"fmt"
	"strings"
	
	_ "github.com/go-sql-driver/mysql"
)

// Resource defines the data common to all resources, including the database
// connection, ID, and table name.
type Resource struct {
	DB *sql.DB `json:"-"`
	ID int64 `json:"id"`
	Table string `json:"-"`
}

// Select runs a SELECT statement on the database.
// cols is a slice of strings representing what columns you want to pull.
// ptrs is a slice of pointers to variables in which to store the results.
func (r *Resource) Select(cols []string, ptrs ...interface{}) error {
	query := fmt.Sprintf("SELECT %v FROM %v WHERE id=?", strings.Join(cols, ", "), r.Table)

	// Query the database and store all data into the pointers corresponding
	// to the database columns desired.
	err := r.DB.QueryRow(query, r.ID).Scan(ptrs...)

	return err
}

// Sync either inserts a new record into the database or updates an existing one.
// cols is a slice of strings representing which columns you would like to save to.
// vals is a slice of variables that contain data to save.
func (r *Resource) Sync(cols []string, vals ...interface{}) error {
	var err error
	var query string
	
	// If the resource does not exist, the query will be an INSERT statement.
	// If it already exists, it will be an UPDATE statement.
	if r.ID == 0 {
		query = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", r.Table, strings.Join(cols, ", "), "?" + strings.Repeat(", ?", len(vals) - 1))

		res, err := r.DB.Exec(query, vals...)
		if err != nil {
			return err
		}

		// Give the resource the proper ID.
		r.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		var updateCols []string
		for _, c := range cols {
			updateCols = append(updateCols, c + "=?")
		}

		query = fmt.Sprintf("UPDATE %v SET %v WHERE id=%v", r.Table, strings.Join(updateCols, ", "), r.ID)
		_, err = r.DB.Exec(query, vals...)
		if err != nil {
			return err
		}
	}

	return err
}

// Delete removes the resource from the database, if it exists. The ID field
// is set to zero. If Save is run after this command, it will create a new
// record for the resource again.
func (r *Resource) Delete() error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id=?", r.Table)
	_, err := r.DB.Exec(query, r.ID)

	if err == nil {
		r.ID = 0
	}

	return err
}
