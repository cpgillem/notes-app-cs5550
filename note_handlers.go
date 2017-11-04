package csnotes

import (
	"fmt"
	"net/http"
)

func GetNote(context *Context) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "logged in")
	}
}
