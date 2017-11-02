package csnotes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetIndex handles requests for the main page of the site.
func GetIndex(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	}
}

// GetUser retrieves a user and their data by id.
func GetUser(db *sql.DB) http.HandlerFunc {
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
			// TODO: DB can't be a global variable, this makes testing impossible.
			u, err := LoadUser(int64(id), db)
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

