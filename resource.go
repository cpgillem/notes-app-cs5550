package main

import (
	"database/sql"
	"fmt"
	"strings"
	
	_ "github.com/go-sql-driver/mysql"
)

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

// Load runs a select statement on the database and scans the row data into
// the variables pointed to from Columns.
func (r *Resource) Load() error {
	var cols []string
	var ptrs []interface{}

	// Columns is a map of string keys (column names) to pointers (to fields).
	for k, v := range r.Columns {
		cols = append(cols, k)
		ptrs = append(ptrs, v)
	}

	// Build a query that selects the columns with the names stored in Columns.
	query := fmt.Sprintf("SELECT %v FROM %v WHERE id=?", strings.Join(cols, ", "), r.Table)

	// Query the database and store all data into the pointers corresponding
	// to the database columns desired.
	err := r.DB.QueryRow(query, r.ID).Scan(ptrs...)

	return err
}

// Save runs an insert or update statement on the database, depending on
// whether the resource already exists. If it didn't previously, the ID is
// stored after the insert statement is run.
func (r *Resource) Save() error {
	return nil
}

// Delete removes the resource from the database, if it exists. The ID field
// is set to zero. If Save is run after this command, it will create a new
// record for the resource again.
func (r *Resource) Delete() error {
	return nil
}
