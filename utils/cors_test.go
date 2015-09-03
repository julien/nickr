package utils

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

func TestAddCORSHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h := AddCORSHeaders(dummyHandler())
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "X-Requested-With" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET,DELETE,PATCH,POST,PUT" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("got %v want 204", w.Code)
	}
}
