package csnotes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PostUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve the form values.
		username := r.FormValue("username")
		name := r.FormValue("name")
		password := r.FormValue("password")

		// Validate the input data.
		if len(username) < 8 {
			resp.Fields["username"] = "Username must be longer than 8 characters."
		}

		if exists, err := CheckUsernameExists(username, context.DB); exists {
			if err != nil {
				http.Error(w, "Could not check for username existence.", http.StatusInternalServerError)
			}
			resp.Fields["username"] = "Username already exists."
		}

		if len(password) < 8 {
			resp.Fields["password"] = "Password must be longer than 8 characters."
		}

		// If one or more fields were invalid, respond early.
		if len(resp.Fields) > 0 {
			return
		}

		// Create a user model.
		u := NewUser(context.DB)

		// Set the new model's data.
		if len(name) > 0 {
			u.Name.String = name
			u.Name.Valid = true
		}
		u.Username = username

		// Save the model.
		err := u.Save()
		if err != nil {
			resp.Errors = append(resp.Errors, "Could not save new user.")
			return
		}

		// Hash and store the user's password.
		err = StorePassword(u.ID, password, context.DB)
		if err != nil {
			resp.Errors = append(resp.Errors, "Could not store password.")

			// Try to delete the user, since the password is required.
			u.Delete()

			return
		}

		// If all was successful, the response will include the user model.
		resp.Models = append(resp.Models, u)
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
