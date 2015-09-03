package utils

import "testing"

func TestContainsString(t *testing.T) {

	src := make([]string, 0)
	src = append(src, "one")
	src = append(src, "two")
	src = append(src, "three")

	if ContainsString(src, "four") {
		t.Errorf("expected false")
	}
	if !ContainsString(src, "three") {
		t.Errorf("expected true")
	}
}

func TestFromJSON(t *testing.T) {
	c := Collection{}

	_, err := c.FromJSON("non_existing.json")
	if err == nil {
		t.Errorf("got nil expected an error\n")
	}

	_, err = c.FromJSON("../data.json")
	if err != nil {
		t.Errorf("got %v expected no error\n", err)
	}

	if ok := c.Get("nedFlanders"); ok == nil {
		t.Errorf("got %v expected non nil\n", ok)
	}
}

func TestAdd(t *testing.T) {
	c := Collection{}
	c.Add("one", []string{"one"})
	c.Add("two", []string{"two"})

	oldlen := len(c.Get("one"))

	c.Add("one", []string{"one"})

	newlen := len(c.Get("one"))

	if oldlen != newlen {
		t.Error("expected item length not to change")
	}

	c.Add("three", []string{"three"})

	if c.Get("three") == nil {
		t.Errorf("expected a value got nil")
	}
}

func TestGet(t *testing.T) {
	c := Collection{}
	if c.Get("non_existing") != nil {
		t.Errorf("expected nil got %v\n", c.Get("non_existing"))
	}

	s := []string{"foo", "bar"}
	c.Add("non_existing", s)

	if c.Get("non_existing") == nil {
		t.Errorf("got nil expected something\n")
	}
}

func TestSet(t *testing.T) {
	c := Collection{}

	c.Add("non_existing", []string{"foo", "bar"})

	if c.Get("non_existing") == nil {
		t.Errorf("got nil expected something\n")
	}

	c.Set("non_existing", []string{"baz", "qux"})
	k := c.Get("non_existing")

	if k == nil {
		t.Errorf("got nil expected something\n")
	}

	if !ContainsString(k, "baz") {
		t.Errorf("expected baz to be present")
	}
}

func TestDelete(t *testing.T) {
	c := Collection{}

	c.Add("non_existing", []string{"foo", "bar"})

	if c.Get("non_existing") == nil {
		t.Errorf("got nil expected something\n")
	}

	c.Delete("non_existing")
	if c.Get("non_existing") != nil {
		t.Errorf("expected nil\n")
	}
}

func TestFlush(t *testing.T) {
	c := Collection{}
	c.Add("first", []string{"one", "two"})
	c.Add("second", []string{"three", "four"})
	c.Add("third", []string{"five", "six"})

	err := c.Flush("../jizzle\\drizzle/foojson")
	if err == nil {
		t.Errorf("error: %v\n", err)
	}

	err = c.Flush("../tests.json")
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	_, err = c.FromJSON("../tests.json")
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	if c.Get("first") == nil {
		t.Errorf("got nil expected something")
	}
}
