package csnotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetIndex handles requests for the main page of the site.
func GetIndex(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	}
}

// GetUser retrieves a user and their data by id.
func GetUser(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if idVar, ok := vars["id"]; ok {
			id, err := strconv.Atoi(idVar)
			if err != nil {
				http.Error(w, "URL does not contain a user ID.", http.StatusNotFound)
				return
			}

			// Send the user struct with the password redacted.
			// The id should be included along with the rest of the variables.
			u, err := LoadUser(int64(id), context.DB)
			if err != nil {
				// Return a 404.
				http.Error(w, "Could not find user.", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(u)
			w.Write(j)
		} else {
			// Return a 404.
			http.Error(w, "Could not find user ID.", http.StatusNotFound)
		}
	}
}

