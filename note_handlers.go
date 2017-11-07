package csnotes

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "logged in")
	}
}

// PostNotes is a handler for creating a new note. It will add the note to
// the user that is logged in, unless they are an admin and the user_id field
// is filled.
func PostNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve form values.
		title := r.FormValue("title")
		content := r.FormValue("content")
		time := r.FormValue("time")
		readUserID := r.FormValue("user_id")

		// Perform validation on the form values.
		if len(title) == 0 {
			resp.Fields["title"] = "Title must be specified."
		}

		if len(resp.Fields) > 0 {
			return
		}

		// Retrieve the logged in user's data.
		userID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 403
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// If the user ID was specified, and the logged in user is admin, the
		// user ID will be set. If not, it will default to the current user's.
		if len(readUserID) > 0 {
			// Make sure the logged in user is admin.
			if !currentUserAdmin {
				resp.StatusCode = 403
				resp.ErrorMessage = "Could not add a note to this user."
				return
			}

			// Attempt to read the ID as an int64.
			convUserID, err := strconv.Atoi(readUserID)
			if err != nil {
				resp.Fields["user_id"] = "Improper user ID."
				return
			}
			userID = int64(convUserID)
		}

		// Create a new note model.
		n := NewNote(context.DB)

		// Set the values.
		n.Title = title

		n.Content.String = content
		n.Content.Valid = len(content) > 0

		n.Time.String = time
		n.Time.Valid = len(time) > 0

		n.UserID = userID

		// Save the new note.
		err = n.Save()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not save note."
			return
		}

		// Add the note model to the response.
		resp.Models = append(resp.Models, n)
	}
}
