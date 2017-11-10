package csnotes

import (
	"net/http"
)

// GetTags retrieves all tags owned by the logged in user.
func GetTags(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Create the response.
		resp := NewJSONResponse()
		defer resp.Respond(w)

		// Get the logged in user's data.

		// Load the user's data into a model.

		// Retrieve the user's tags.

		// Add the tags to the response.
	}
}
