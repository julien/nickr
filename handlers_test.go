package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRoot(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if len(w.Body.String()) == 0 {
		t.Errorf("expected a body\n")
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("got %v want 404", w.Code)
	}

}

func TestGetUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if len(w.Body.String()) == 0 {
		t.Errorf("expected a body\n")
	}

	if w.Header().Get("Content-Type") == "" {
		t.Errorf("expected a content type header\n")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected an application/json content type header\n")
	}

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}

}

func TestHandleGetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/frank", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if len(w.Body.String()) == 0 {
		t.Errorf("expected a body\n")
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("got %v want 404", w.Code)
	}

}

func TestOptions(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/users/frank", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}

}

func TestPostOK(t *testing.T) {
	postData := []byte("{\"name\": \"tester\", \"nicknames\": [\"nickname1\"]}")

	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(postData))

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	ok := w.Code == http.StatusCreated || w.Code == http.StatusBadRequest

	if !ok {
		t.Errorf("got %v want 201 or 400", w.Code)
	}

	if w.Body.String() == "" {
		t.Errorf("expected a response, got %s\n", w.Body.String())
	}
}

func TestPostBadJSON(t *testing.T) {

	postData := []byte("{\"name\": \"\"}")

	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(postData))

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %v want 400", w.Code)
	}

	if w.Body.String() == "" {
		t.Errorf("expected a response, got %s\n", w.Body.String())
	}
}

func TestHandlePutRequest(t *testing.T) {
	postData := []byte("{\"name\": \"tester\", \"nicknames\": [\"nickname1\"]}")

	req, _ := http.NewRequest("PUT", "/users", bytes.NewReader(postData))

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	ok := w.Code == http.StatusOK || w.Code == http.StatusBadRequest || w.Code == http.StatusNotFound

	if !ok {
		t.Errorf("got %v want 201 or 400", w.Code)
	}

	t.Logf("%s\n", w.Body.String())

	if w.Body.String() == "" {
		t.Errorf("expected a response, got %s\n", w.Body.String())
	}
}

func TestHandlePutNotFoundRequest(t *testing.T) {
	postData := []byte("{\"name\": \"T\"}")

	req, _ := http.NewRequest("PUT", "/users", bytes.NewReader(postData))

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %v want 404", w.Code)
	}

	t.Logf("%s\n", w.Body.String())

	if w.Body.String() == "" {
		t.Errorf("expected a response, got %s\n", w.Body.String())
	}
}

func TestHandleDeleteRequest(t *testing.T) {
	postData := []byte("{\"name\": \"tester\"}")

	req, _ := http.NewRequest("DELETE", "/users", bytes.NewReader(postData))

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	ok := w.Code == http.StatusOK || w.Code == http.StatusBadRequest || w.Code == http.StatusNotFound

	if !ok {
		t.Errorf("got %v want 200 or 400", w.Code)
	}
}
