package csnotes

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetNotes(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Get the logged in user's ID.
		currentUserID, _, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not get logged in user."
			return
		}

		// Create and load a user model.
		u, err := LoadUser(currentUserID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load user data."
			return
		}

		// Load the notes from the user model.
		ns, err := u.Notes()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load user's notes."
			return
		}

		// Add the notes to the response.
		for _, n := range ns {
			resp.Models = append(resp.Models, n)
		}
	}
}

// GetNote retrieves a note from a user. If the logged in user is not admin,
// they will only be able to retrieve a note that's theirs.
func GetNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Get the note ID.
		nID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the note's existence.
		if e, err := CheckExistence(nID, "notes", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "Note not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify note's existence."
			}
			return
		}
		
		// Load the note model.
		n, err := LoadNote(nID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load note."
			return
		}

		// Get the logged in user's data.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// If the currently logged in user does not own the note, or is not
		// admin, access will be denied to the note.
		if !currentUserAdmin && currentUserID != n.UserID {
			resp.StatusCode = 403
			resp.ErrorMessage ="Access denied."
			return
		}

		// Add the note model to the response.
		resp.Models = append(resp.Models, n)
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
			resp.StatusCode = 500
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

// PutNote updates a note's data. If the user doesn't own the note, and is not
// an admin, they will be denied access.
func PutNote(context *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve the form values.
		title := r.FormValue("title")
		content := r.FormValue("content")
		time := r.FormValue("time")

		// Make sure the title is not empty.
		if len(title) == 0 {
			resp.Fields["title"] = "Title must not be empty."
			return
		}

		// Retrieve the note ID.
		nID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the note's existence.
		if e, err := CheckExistence(nID, "notes", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "Note not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify note's existence."
			}
			return
		}

		// Load the note model.
		n, err := LoadNote(nID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Note not found."
			return
		}

		// Retrieve the logged in user's data.
		userID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 403
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// If the user does not own this note, and isn't an admin, deny access.
		if !currentUserAdmin && userID != n.UserID {
			resp.StatusCode = 403
			resp.ErrorMessage = "Access denied."
			return
		}

		// Update the note's values.
		n.Title = title

		n.Content.String = content
		n.Content.Valid = len(content) > 0

		n.Time.String = time
		n.Time.Valid = len(time) > 0

		// Save the note.
		err = n.Save()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not save note."
			return
		}

		// Add the newly updated note to the response.
		resp.Models = append(resp.Models, n)
	}
}

// DeleteNote removes a note from the database as long as the user is either
// admin or the owner of the note.
func DeleteNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Retrieve the ID from the URL.
		nID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the note's existence.
		if e, err := CheckExistence(nID, "notes", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "Note not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify note's existence."
			}
			return
		}

		// Load the data of the user that's logged in.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// Create a model for the note from the ID.
		n, err := LoadNote(nID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve note."
			return
		}

		// Make sure the user is either admin or the owner of the note.
		if n.UserID != currentUserID && !currentUserAdmin {
			resp.StatusCode = 403
			resp.ErrorMessage = "Must be admin or owner of note."
			return
		}

		// Delete the note.
		err = n.Delete()
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not delete note."
			return
		}

		// Add the old note's data to the response.
		resp.Models = append(resp.Models, n)
	}
}

func GetNoteTags(context *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Get the note ID from the URL.
		nID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Check for the existence of the note.
		if e, err := CheckExistence(nID, "notes", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "Note not found."
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify note's existence."
			}
			return
		}

		// Attempt to load the note.
		n, err := LoadNote(nID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load note."
			return
		}

		// Get the logged in user's data.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not get logged in user."
			return
		}

		// Verify that the user owns the note, or is admin.
		if !currentUserAdmin && currentUserID != n.UserID {
			resp.StatusCode = 403
			resp.ErrorMessage = "Access denied."
		}

		// Retrieve the note's tags.
		ts, err := n.Tags()
		if err != nil {
			fmt.Println(err)
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load tags."
			return
		}

		// Add the tags to the response.
		for _, t := range ts {
			resp.Models = append(resp.Models, t)
		}
	}
}

// PostNoteTags attaches a tag to a note through an intermediate table. The tag
// ID is sent through a form parameter.
func PostNoteTag(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create a response.
		resp := NewJSONResponse()
		defer resp.Respond(w)
		
		// Retrieve the note ID.
		nID, ok := GetURLID(r, &resp)
		if !ok {
			return
		}

		// Retrieve the tag ID.
		tIDform := r.FormValue("tag_id")
		
		// Ensure a tag ID was given.
		if len(tIDform) == 0 {
			resp.Fields["tag_id"] = "No tag ID given."
			return
		}

		// Ensure that the tag ID is a valid int.
		tIDint, err := strconv.Atoi(tIDform)
		if err != nil {
			resp.Fields["tag_id"] = "Invalid tag ID."
			return
		}
		tID := int64(tIDint)

		// Make sure the note exists.
		if e, err := CheckExistence(nID, "notes", context.DB); !e {
			if err == nil {
				resp.StatusCode = 404
				resp.ErrorMessage = "Note not found."
				return
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify existence of note."
				return
			}
		}

		// Load the note model.
		n, err := LoadNote(nID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load note."
			return
		}

		// Make sure the tag exists.
		if e, err := CheckExistence(tID, "tags", context.DB); !e {
			if err == nil {
				resp.Fields["tag_id"] = "Tag does not exist."
				return
			} else {
				resp.StatusCode = 500
				resp.ErrorMessage = "Could not verify existence of tag."
				return
			}
		}

		// Load the tag as a model.
		t, err := LoadTag(tID, context.DB)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not load tag."
			return
		}

		// Retrieve the logged in user.
		currentUserID, currentUserAdmin, err := context.LoggedInUser(r)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not retrieve logged in user."
			return
		}

		// Make sure the user is either admin or owns the note and the tag.
		if !currentUserAdmin && (currentUserID != t.UserID || currentUserID != n.UserID) {
			resp.StatusCode = 403
			resp.ErrorMessage = "Access Denied."
			return
		}

		// Add the tag to the note.
		err = n.AddTag(t.ID)
		if err != nil {
			resp.StatusCode = 500
			resp.ErrorMessage = "Could not add tag to note."
			return
		}

		// Add the tag model to the response.
		resp.Models = append(resp.Models, t)
	}
}
