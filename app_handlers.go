package csnotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GetLogin should take the user's credentials and create a
// JSON web token if the authentication was successful.
func PostLogin(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate the data.
		if len(username) == 0 {
			http.Error(w, "No username.", http.StatusInternalServerError)
			return
		}

		if len(password) == 0 {
			http.Error(w, "No password.", http.StatusInternalServerError)
			return
		}

		// Validate the credentials against the database.
		user, err := ValidateUser(username, password, context.DB)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not authenticate user.", http.StatusForbidden)
			return
		}

		// Create a token.
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims {
			"iss": "admin",
			"exp": time.Now().Add(time.Minute * 20).Unix(),
			"user_id": user.ID,
			"user_admin": user.Admin,
		})

		tokenString, err := token.SignedString(context.SignKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		response := struct {
			Token string `json:"token"`
		} {tokenString}

		// Turn the token into a json string.
		json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
