package csnotes

import (
	"github.com/gorilla/mux"
)

func CreateRouter(context *Context) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", GetIndex(context))
	router.HandleFunc("/user/{id}", GetUser(context))

	return router
}
