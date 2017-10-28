package main

import (
	"database/sql"
	"fmt"
	"strings"
	
	_ "github.com/go-sql-driver/mysql"
)

// Model implements a function for saving to the database, and a function for
// loading from it. They may use the helper methods Select, Sync, or Delete on
// a resource, as needed.
type Model interface {
	Save() error
	Load() error
}

// Resource defines the data common to all resources, including the database
// connection, ID, and table name.
type Resource struct {
	DB *sql.DB
	ID int64
	Table string
	// Columns is a map of strings to pointers, which are the final locations of 
	// each value returned from the database. Usually these pointers will reference
	// fields of a struct.
	Columns map[string]interface{}
}

func (r *Resource) Select(cols []string, ptrs ...interface{}) error {
	// Build a query that selects the columns with the names stored in Columns.
	query := fmt.Sprintf("SELECT %v FROM %v WHERE id=?", strings.Join(cols, ", "), r.Table)

	// Query the database and store all data into the pointers corresponding
	// to the database columns desired.
	err := r.DB.QueryRow(query, r.ID).Scan(ptrs...)

	return err
}

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
	return nil
}
