package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

type FakeReader struct {
	content string
	done    bool
}

func (r *FakeReader) Read(p []byte) (int, error) {
	if r.done {
		return 0, io.EOF
	}
	for i, b := range []byte(r.content) {
		p[i] = b
	}
	r.done = true
	return len(r.content), nil
}

func TestAddCORS1(t *testing.T) {
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

func TestAddCORS2(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/", nil)

	w := httptest.NewRecorder()
	h := AddCORS(dummyHandler(), "", "", "")
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "" {
		t.Errorf("got %v want X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("got %v want 204", w.Code)
	}
}

func TestAddCORS3(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/", nil)

	w := httptest.NewRecorder()
	h := AddCORS(dummyHandler(), "", "FOO", "BAAAAAR")
	h.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("got %v want *", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "FOO" {
		t.Errorf("got %v want X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Code != http.StatusNoContent {
		t.Errorf("got %v want 204", w.Code)
	}
}

func TestBodyToByte(t *testing.T) {

	r := &FakeReader{"tester", false}

	b, err := bodyToByte(r)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	if len(b) == 0 {
		t.Errorf("expected reader to have data\n")
	}
}

func TestBodyToUser1(t *testing.T) {
	r := &FakeReader{"{\"name\": \"tester\"}", false}

	b, err := bodyToUser(r)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	if len(b.Name) == 0 {
		t.Errorf("expected user to have a name\n")
	}

	if b.Name != "tester" {
		t.Errorf("expected user name to be tester, got: %v\n", b.Name)
	}
}

func TestBodyToUser2(t *testing.T) {
	r := &FakeReader{"error", false}

	b, err := bodyToUser(r)
	if err == nil {
		t.Errorf("expected an error: %v\n", err)
	}

	if b != nil {
		t.Errorf("expected a nil user\n")
	}

}

func TestEncodeJSON(t *testing.T) {
	v := struct {
		Name string `json:"name"`
	}{
		"tester",
	}

	if _, err := encodeJSON(v); err != nil {
		t.Errorf("got %v\n", err)
	}
}

func TestDecodeJSON(t *testing.T) {
	type tester struct {
		Name string `json:"name"`
	}

	v := &tester{}
	b := []byte("{\"name\": \"tester\"}")

	if err := decodeJSON(b, &v); err != nil {
		t.Errorf("got %v\n", err)
	}

	if v.Name != "tester" {
		t.Errorf("expected v.Name to be tester, got %v\n", v.Name)
	}
}

func TestDecodeJSONError(t *testing.T) {
	type tester struct {
		Name string `json:"name"`
	}

	v := &tester{}
	b := []byte("{\"name\"}")

	if err := decodeJSON(b, &v); err == nil {
		t.Errorf("expected error, got %v\n", err)
	}
}
