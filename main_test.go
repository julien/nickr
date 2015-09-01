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

func TestIndexHandler(t *testing.T) {

	http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}

}
