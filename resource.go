package main

import (
	"database/sql"
	"fmt"
	"strings"
	
	_ "github.com/go-sql-driver/mysql"
)

type Resource struct {
	DB *sql.DB
	ID int64
	Table string
	Columns map[string]interface{}
}

func (r *Resource) Load() error {
	var cols []string
	var ptrs []interface{}

	for k, v := range r.Columns {
		cols = append(cols, k)
		ptrs = append(ptrs, v)
	}

	query := fmt.Sprintf("SELECT %v FROM %v WHERE id=?", strings.Join(cols, ", "), r.Table)

	err := r.DB.QueryRow(query, r.ID).Scan(ptrs...)

	return err
}

func (r *Resource) Save() error {
	return nil
}

func (r *Resource) Delete() error {
	return nil
}
