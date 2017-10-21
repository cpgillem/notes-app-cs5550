package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
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

	// Define the routes.
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)

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
