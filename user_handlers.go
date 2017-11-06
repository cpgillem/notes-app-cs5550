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

func PostUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate the input data
		if len(username) < 8 {
			// TODO: Validation should be in 200 response.
			http.Error(w, "Username must be longer than 8 characters.", http.StatusInternalServerError)
			return
		}

		if len(password) < 8 {
			http.Error(w, "Password must be longer than 8 characters.", http.StatusInternalServerError)
			return
		}

		// Create a user model.
		u := NewUser(context.DB)
		u.Username = username
		err := u.Save()
		if err != nil {
			http.Error(w, "Could not save new user.", http.StatusInternalServerError)
			return
		}

		// Hash and store the user's password.
		err = StorePassword(u.ID, password, context.DB)
		if err != nil {
			http.Error(w, "Could not store password.", http.StatusInternalServerError)
			return
		}
	}
}

func GetUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Retrieve the ID.
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "No ID specified.", http.StatusNotFound)
			return
		}

		// Convert the ID to an int.
		uID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Improper ID.", http.StatusNotFound)
			return
		}

		// Retrieve a user model.  
		u, err := LoadUser(int64(uID), context.DB)
		if err != nil {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		// Return this user's data as a JSON response.
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(u)
		w.Write(j)
	}
}

func PutUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Retrieve and validate the given data.
		name := r.FormValue("name")

		if len(name) == 0 {
			// TODO: invalid data
			return
		}

		// Retrieve the user ID.
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "No ID specified.", http.StatusNotFound)
			return
		}

		// Convert the ID to an int.
		uID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Improper ID.", http.StatusNotFound)
			return
		}

		// Retrieve the user model.
		_, err = LoadUser(int64(uID), context.DB)
		if err != nil {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		// Store the new values.
		// TODO: make sure user is allowed to change the data.

	}
}
