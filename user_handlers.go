package csnotes

import (
	"fmt"
	"net/http"
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

// GetUsers retrieves all users and returns them as an array of user models.
func GetUsers(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Make sure the logged in user is available.
		_, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// Make sure the logged in user is an admin.
		if !currentUserAdmin {
			resp.StatusCode = 403
			resp.ErrorMessage = "Only admins can retrieve all users."
			return
		}

		// Retrieve all the user models.
		users, err := LoadAllUsers(context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve users."
			return
		}

		// Add the user models to the response.
		for _, m := range users {
			resp.Models = append(resp.Models, m)
		}
	}
}

func GetUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve the ID.
		uID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the user's existence.
		if e, err := CheckExistence(uID, "users", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "User not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify user's existence."
			}
			return
		}

		// Retrieve a user model.  
		u, err := LoadUser(uID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load user."
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
		uID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the user's existence.
		if e, err := CheckExistence(uID, "users", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "User not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify user's existence."
			}
			return
		}

		// Retrieve the user model.
		u, err := LoadUser(int64(uID), context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load user."
			return
		}

		// Retrieve the logged in user ID.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 403
			resp.ErrorMessage = "Could not retrieve logged in user."
			fmt.Println(err)
			return
		}

		// Ensure that the logged in user is allowed to modify the specified 
		// user.
		if currentUserID != u.ID && !currentUserAdmin {
			resp.StatusCode = 403
			resp.ErrorMessage = "Access denied. Must be admin or logged in as this user."
			return
		}

		// Set the new values.
		// NOTE: These values will not take effect until the user logs out and back
		// in, since the JWT isn't updated. This could be remedied in a future
		// iteration if necessary.
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

// GetUserNotes retrieves all the notes belonging to a user.
func GetUserNotes(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create a response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Get the user ID from the URL.
		uID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the user's existence.
		if e, err := CheckExistence(uID, "users", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "User not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify user's existence."
			}
			return
		}

		// Retrieve the logged in user.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 403
			resp.ErrorMessage = "Could not get logged in user."
			return
		}

		// Make sure the logged in user is allowed to see the notes.
		if currentUserID != uID && !currentUserAdmin {
			resp.StatusCode = 403
			resp.ErrorMessage = "Access denied."
			return
		}

		// Load the user model.
		u, err := LoadUser(uID, context.DB)
		if err != nil {
			resp.StatusCode = 404
			resp.ErrorMessage = "User not found."
			return
		}

		// Get the user's notes.
		ns, err := u.Notes()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load notes."
			return
		}

		// Add the notes to the response.
		for _, n := range ns {
			resp.Models = append(resp.Models, n)
		}
	}
}
