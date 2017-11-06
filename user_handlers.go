package csnotes

import (
	"fmt"
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
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not check for username existence."
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
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not save new user."
			return
		}

		// Hash and store the user's password.
		err = StorePassword(u.ID, password, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not store password."

			// Try to delete the user, since the password is required.
			u.Delete()

			return
		}

		// If all was successful, the response will include the user model.
		fmt.Println(u)
		resp.Models = append(resp.Models, u)
	}
}

func GetUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve the ID.
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			resp.StatusCode = 404
			resp.ErrorMessage = "No ID specified."
			return
		}

		// Convert the ID to an int.
		uID, err := strconv.Atoi(idStr)
		if err != nil {
			resp.StatusCode = 404
			resp.ErrorMessage = "Improper ID."
			return
		}

		// Retrieve a user model.  
		u, err := LoadUser(int64(uID), context.DB)
		if err != nil {
			resp.StatusCode = 404
			resp.ErrorMessage = "User not found."
			return
		}

		// Return this user's data.
		resp.Models = append(resp.Models, u)
	}
}

func PutUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create a new response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve any new data.
		name := r.FormValue("name")

		// Retrieve the user ID.
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			resp.StatusCode = 404
			resp.ErrorMessage = "No ID specified."
			return
		}

		// Convert the ID to an int.
		uID, err := strconv.Atoi(idStr)
		if err != nil {
			resp.StatusCode = 404
			resp.ErrorMessage = "Improper ID."
			return
		}
		
		// Make sure the user actually exists.
		if exists, _ := CheckExistence(int64(uID), "users", context.DB); !exists {
			resp.StatusCode = 404
			resp.ErrorMessage = "User does not exist."
			return
		}

		// Retrieve the user model.
		u, err := LoadUser(int64(uID), context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load user."
			return
		}

		// Set the new values.
		// TODO: make sure user is allowed to change the data.
		u.Name.String = name
		u.Name.Valid = len(name) > 0

		// Store the new values in the database.
		err = u.Save()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not store new data."
			return
		}

		// Add the user model to the response.
		resp.Models = append(resp.Models, u)
	}
}
