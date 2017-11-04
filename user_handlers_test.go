package csnotes

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestPostUser(t *testing.T) {
	// Mock data.
	db := SetUpDbTest()
	defer TearDownDbTest(db)

	context := Context {
		DB: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc(PostUser(&context)))
	server := httptest.NewServer(router)
	defer server.Close()

	// Create and send a request.
	_, err := http.PostForm(server.URL, url.Values{"username": {"newtestuser"}, "password": {"testpassword"}})
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the user was inserted into the database.
	row := db.QueryRow("SELECT id FROM users WHERE username='newtestuser'")
	var id int64
	err = row.Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUser(t *testing.T) {
	// Mock data.
	db, ids, _ := SeededTestDB()
	defer TearDownDbTest(db)

	// Create surrounding context.
	context := Context {
		DB: db,
	}
	router := mux.NewRouter()
	router.HandleFunc("/{id}", http.HandlerFunc(UpdateUser(&context))).Methods("PUT")
	server := httptest.NewServer(router)
	defer server.Close()

	// Create and send an update request.
	body := []byte(`{
		"name": "Test User"
	}`)
	url := fmt.Sprintf("%s/%d", server.URL, ids["user.nonadmin"])
	req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))
	client := &http.Client {
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the response was valid and the user is updated.
	AssertEqual(200, res.StatusCode, t)

	row := db.QueryRow("SELECT name FROM users WHERE id=?", ids["user.nonadmin"])
	var name sql.NullString
	err = row.Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual("Test User", name, t)
}
