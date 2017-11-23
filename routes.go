package csnotes

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateRouter(context *Context) *mux.Router {
	router := mux.NewRouter()
	api := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

	// Create middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return context.VerifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Public assets
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	// Public Routes (non-GET)
	router.HandleFunc("/login", PostLogin(context)).Methods("POST")
	//router.HandleFunc("/logout", PostLogout(context)).Methods("POST")

	// User Routes
	api.HandleFunc("/user", GetUsers(context)).Methods("GET")
	api.HandleFunc("/user", PostUser(context)).Methods("POST")
	api.HandleFunc("/user/{id}", GetUser(context)).Methods("GET")
	api.HandleFunc("/user/{id}", PutUser(context)).Methods("PUT")
	api.HandleFunc("/user/{id}/note", GetUserNotes(context)).Methods("GET")

	// Note Routes
	api.HandleFunc("/note", GetNotes(context)).Methods("GET")
	api.HandleFunc("/note", PostNote(context)).Methods("POST")
	api.HandleFunc("/note/{id}", GetNote(context)).Methods("GET")
	api.HandleFunc("/note/{id}", PutNote(context)).Methods("PUT")
	api.HandleFunc("/note/{id}", DeleteNote(context)).Methods("DELETE")
	//router.HandleFunc("/note/{id}/tag", GetNoteTags(context)).Methods("GET")
	//router.HandleFunc("/note/{id}/tag", PostNoteTag(context)).Methods("POST")

	// Authenticated API Routes
	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(api),
	))

	return router
}
