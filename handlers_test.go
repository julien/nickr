package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if len(w.Body.String()) == 0 {
		t.Errorf("expected a body\n")
	}

	if w.Body.String() != "nickr" {
		t.Errorf("expected the body to be nickr\n")
	}

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}

}

func TestHandleRequest2(t *testing.T) {
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

func TestHandleRequest3(t *testing.T) {
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

func TestHandleRequest4(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/users/frank", nil)

	w := httptest.NewRecorder()
	h := handleRequest()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}

}
