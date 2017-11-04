package csnotes

import (
	"crypto/rsa"
	"database/sql"
)

// Context is a struct that contains the webapp's global variables.
type Context struct {
	DB *sql.DB
	VerifyKey *rsa.PublicKey
	SignKey *rsa.PrivateKey
}
