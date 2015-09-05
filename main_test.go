package nickr

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// func TestIndexHandler(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	w := httptest.NewRecorder()
//
// 	h := collectionHandler()
// 	h.ServeHTTP(w, req)
//
// 	if w.Body == nil {
// 		t.Errorf("Body is nil")
// 	}
//
// 	if w.Code != http.StatusOK {
// 		t.Errorf("got %v want 200", w.Code)
// 	}
// }
//
// func TestNonExistingHandler(t *testing.T) {
//
// 	req, _ := http.NewRequest("GET", "/nonExisting", nil)
// 	w := httptest.NewRecorder()
// 	h := collectionHandler()
// 	h.ServeHTTP(w, req)
//
// 	if w.Code != http.StatusNotFound {
// 		t.Errorf("got %v want 404", w.Code)
// 	}
//
// }
