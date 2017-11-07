package csnotes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetURLResource retrieves the ID of a resource requested via URL.
func GetURLID(req *http.Request, resp *JSONResponse) int64 {
	// Retrieve the ID.
	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		resp.StatusCode = 404
		resp.ErrorMessage = "No ID specified."
		return 0
	}

	// Convert the ID to an int.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.StatusCode = 404
		resp.ErrorMessage = "Improper ID."
		return 0
	}

	return int64(id)
}
