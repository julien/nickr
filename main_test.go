package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// var t *testing.T
	// if *port != "" {
	// 	t.Errorf("got %v want 8080", *port)
	// }

	os.Exit(m.Run())
}

func TestLoadCollection(t *testing.T) {
	err := loadCollection()
	if err != nil {
		t.Errorf("got %v", err)
	}

	if collection == nil {
		t.Errorf("got %v expected collection not to be nil\n", err)
	}
}

func TestIndexHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := collectionHandler()
	h.ServeHTTP(w, req)

	if w.Body == nil {
		t.Errorf("Body is nil")
	}

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}
}

func TestNonExistingHandler(t *testing.T) {

	req, _ := http.NewRequest("GET", "/nonExisting", nil)
	w := httptest.NewRecorder()
	h := collectionHandler()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %v want 400", w.Code)
	}

}
