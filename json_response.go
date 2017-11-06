package csnotes

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	// The HTTP status code that will be set upon response. By default, this
	// is 200.
	StatusCode int `json:"-"`

	// This will be the main error message if there was a non-200 response.
	ErrorMessage string `json:"-"`

	// An array of models that are returned from the database. For GET
	// endpoints, this will be the requested resources. For POST/PUT/DELETE
	// endpoints, this will be the affected resources.
	Models []interface{} `json:"models"`

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
		StatusCode: 200,
		Models: []interface{}{},
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

	// If the response code is not 200, write an error instead with a plain
	// body containing the first error message in the response struct.
	if jr.StatusCode != 200 {
		http.Error(w, jr.ErrorMessage, jr.StatusCode)
		return
	}

	// If there were no serious errors, write the JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
