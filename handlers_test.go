package csnotes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

/* Test the handler for the index route. Should return a success status code. */
func TestGetIndex(t *testing.T) {
	// Send a request to the index route.
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetIndex(nil))

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Wrong status code. Got %v want %v", status, http.StatusOK)
	}
}

func TestGetUser(t *testing.T) {
	// Mock some data.
	db, ids, _ := SeededTestDB()
	defer TearDownDbTest(db)

	context := Context {
		DB: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{id}", http.HandlerFunc(GetUser(&context)))
	server := httptest.NewServer(router)
	defer server.Close()

	// Send the request.
	url := fmt.Sprintf("%s/%d", server.URL, ids["user.nonadmin"])
	res, err := http.Get(url)

	if err != nil {
		t.Fatal(err)
	}

	// Read the response and ensure that the data matches.
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	str := string(body)

	AssertContains(str, "nonadmin", t)
	AssertContains(str, "false", t)
}
