package csnotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
			// TODO: Possibly use a generated struct that contains all necessary data
			"CustomUserInfo": struct {
				ID int64
			} {user.ID},
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

// GetIndex handles requests for the main page of the site.
func GetIndex(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	}
}

// GetUser retrieves a user and their data by id.
func GetUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if idVar, ok := vars["id"]; ok {
			id, err := strconv.Atoi(idVar)
			if err != nil {
				http.Error(w, "URL does not contain a user ID.", http.StatusNotFound)
				return
			}

			// Send the user struct with the password redacted.
			// The id should be included along with the rest of the variables.
			u, err := LoadUser(int64(id), context.DB)
			if err != nil {
				// Return a 404.
				http.Error(w, "Could not find user.", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(u)
			w.Write(j)
		} else {
			// Return a 404.
			http.Error(w, "Could not find user ID.", http.StatusNotFound)
		}
	}
}

func PostUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		password := r.FormValue("password")

		// Validate the input data
		if len(name) < 8 {
			// TODO: Error message in response
			return
		}

		if len(password) < 8 {
			return
		}

		// Store the new user.

	}
}

func GetNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "logged in")
	}
}
