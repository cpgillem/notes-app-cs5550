package csnotes

import (
	"database/sql"
	
)

// Context is a struct that contains the webapp's global variables.
type Context struct {
	DB *sql.DB
}
