package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// Define a global database connection.
var db *sql.DB

// GetIndex handles requests for the main page of the site.
func GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}

// GetUser retrieves a user and their data by id.
func GetUser(w http.ResponseWriter, r *http.Request) {
	
}

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
	_, err := sql.Open("mysql", "notes_app:notes_app@/notes_app")
	if err != nil {
		panic(err)
	}

	// Define the routes.
	router := mux.NewRouter()
	router.HandleFunc("/", GetIndex)

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
