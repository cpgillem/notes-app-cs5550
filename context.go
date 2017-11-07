package csnotes

import (
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Context is a struct that contains the webapp's global variables.
type Context struct {
	DB *sql.DB
	VerifyKey *rsa.PublicKey
	SignKey *rsa.PrivateKey
}

// LoggedInUser parses a JWT and returns a user ID and admin status, if 
// possible.
func (c *Context) LoggedInUser(r *http.Request) (uID int64, admin bool, err error) {
	authHeader := r.Header.Get("Authorization")

	// Check for an authorization header.
	if len(authHeader) < 8 {
		err = errors.New("Authorization header not found.")
		return 
	}
	
	// Check for a token.
	if strings.Index(authHeader, "Bearer ") != 0 {
		err = errors.New("No token found in authorization header.")
		return
	}

	// Read and parse the token.
	tokenString := authHeader[7:]
	token, err := jwt.Parse(tokenString, func (*jwt.Token) (interface{}, error) {
		return c.VerifyKey, nil
	})
	if err != nil {
		return
	}

	// Make sure the token is valid.
	if !token.Valid {
		err = errors.New("Token invalid.")
		return
	}

	// Make sure the claims are a valid map.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("Could not get token claims.")
		return
	}

	// Extract the user ID from the claims.
	id, ok := claims["user_id"].(float64)
	if !ok {
		err = errors.New("Could not load user data from token.")
		return
	}
	fmt.Println(claims)
	uID = int64(id)

	// Extract the admin status.
	admin, ok = claims["user_admin"].(bool)
	if !ok {
		err = errors.New("Could not load user data from token.")
		return
	}

	return
}
