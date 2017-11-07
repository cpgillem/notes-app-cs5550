package csnotes

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func CreateRouter(context *Context) *mux.Router {
	router := mux.NewRouter()

	// Create middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return context.VerifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Public Routes
	router.HandleFunc("/", GetIndex(context))
	router.HandleFunc("/login", PostLogin(context)).Methods("POST")
	//router.HandleFunc("/logout", GetLogout(context)).Methods("GET")

	router.HandleFunc("/user", GetUsers(context)).Methods("GET")
	router.HandleFunc("/user", PostUser(context)).Methods("POST")
	router.HandleFunc("/user/{id}", GetUser(context)).Methods("GET")
	router.HandleFunc("/user/{id}", PutUser(context)).Methods("PUT")
	router.HandleFunc("/user/{id}/notes", GetUserNotes(context)).Methods("GET")

	router.HandleFunc("/note", PostNote(context)).Methods("POST")
	router.HandleFunc("/note/{id}", GetNote(context)).Methods("GET")

	// Authenticated Routes
	router.Handle("/note/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(GetNote(context)),
	))

	return router
}
