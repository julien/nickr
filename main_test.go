package main

import (
	"io"
	"os"
	"testing"
)

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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
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
