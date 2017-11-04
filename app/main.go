package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cpgillem/csnotes"
	"github.com/urfave/negroni"
	"github.com/dgrijalva/jwt-go"

	_ "github.com/go-sql-driver/mysql"
)

const (
	PRIVATE_KEY_PATH = "./keys/app.rsa"
	PUBLIC_KEY_PATH = "./keys/app.rsa.pub"
)

func main() {
	// Load the RSA keys.
	signRaw, err := ioutil.ReadFile(PRIVATE_KEY_PATH)
	if err != nil {
		panic(err)
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signRaw)
	if err != nil {
		panic(err)
	}

	verifyRaw, err := ioutil.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		panic(err)
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyRaw)
	if err != nil {
		panic(err)
	}

	// Validate command-line arguments.
	if len(os.Args) < 2 {
		log.Fatal("First argument must be a port number.")
	}

	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Port number must be numerical. [%v]", err)
	}

	// Read command line arguments or any config data.
	port := os.Args[1]

	// Setup the database connection.
	db, err := sql.Open("mysql", "notes_app:notes_app@/notes_app")
	if err != nil {
		panic(err)
	}

	// Create a context variable to pass around.
	context := csnotes.Context {
		DB: db,
		SignKey: signKey,
		VerifyKey: verifyKey,
	}

	// Define the routes.
	router := csnotes.CreateRouter(&context)
	n := negroni.Classic()
	n.UseHandler(router)

	// Define a server object.
	server := &http.Server {
		Handler:		n,
		Addr:			"127.0.0.1:" + port,
		WriteTimeout:	15 * time.Second,
		ReadTimeout:	15 * time.Second,
	}

	// Start the server.
	log.Fatal(server.ListenAndServe())
}
