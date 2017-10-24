package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/* Test the handler for the index route. Should return a success status code. */
func TestGetIndex(t *testing.T) {
	// Send a request to the index route.
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetIndex)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Wrong status code. Got %v want %v", status, http.StatusOK)
	}
}

func TestGetUser(t *testing.T) {

}
