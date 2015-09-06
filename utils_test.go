package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func TestAddCORS(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/", nil)

	w := httptest.NewRecorder()
	h := AddCORS(dummyHandler(), "", "X-Requested-With", "GET,POST,PUT,DELETE")
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "X-Requested-With" {
		t.Errorf("got %v want X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET,POST,PUT,DELETE" {
		t.Errorf("got %v want GET,POST,PUT,DELETE", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("got %v want 204", w.Code)
	}
}
