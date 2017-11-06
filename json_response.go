package csnotes

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	// An array of models that are returned from the database. For GET
	// endpoints, this will be the requested resources. For POST/PUT/DELETE
	// endpoints, this will be the affected resources.
	Models []Model `json:"models"`

	// An object mapping field names to their errors. This is used for
	// validating input.
	Fields map[string]string `json:"fields"`

	// An array of errors. A successful operation will return an empty array,
	// but any errors will be appended to this array.
	Errors []string `json:"errors"`
}

// NewJSONResponse creates a new JSON response struct with initialized slices
// and maps.
func NewJSONResponse() JSONResponse {
	return JSONResponse {
		Models: []Model{},
		Fields: map[string]string{},
		Errors: []string{},
	}
}

// Respond is a helper function for sending a properly formatted JSON response
// for any handler function. If the response could not be serialized for any
// reason, a 500 error is written.
func (jr *JSONResponse) Respond(w http.ResponseWriter) {
	// Marshal the response into a JSON string.
	res, err := json.Marshal(jr)
	if err != nil {
		http.Error(w, "Could not create JSON response.", http.StatusInternalServerError)
		return
	}

	// Write the JSON string to the response writer.
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
