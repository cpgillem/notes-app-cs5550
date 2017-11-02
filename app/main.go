package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cpgillem/csnotes"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
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
	}

	// Define the routes.
	router := csnotes.CreateRouter(&context)

	// Define a server object.
	server := &http.Server {
		Handler:		router,
		Addr:			"127.0.0.1:" + port,
		WriteTimeout:	15 * time.Second,
		ReadTimeout:	15 * time.Second,
	}

	// Start the server.
	log.Fatal(server.ListenAndServe())
}
