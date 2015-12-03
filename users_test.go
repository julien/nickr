package main

import (
	"fmt"
	"testing"
	"time"
)

const fbTestsURL = "https://nickr.firebaseio.com/tests/"

func TestUserHasNickname(t *testing.T) {

	nicknames := []*Nickname{&Nickname{"nickname1", "nickname1.png"}}

	u := &User{"tester", nicknames}
	b := &Nickname{"foo", "foo.png"}

	if u.HasNickname(b) {
		t.Errorf("expected false got true")
	}
}

func TestNewUsers(t *testing.T) {
	u := NewUsers(fbTestsURL + "test-new/users/")
	_, err := u.All()
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
}

func TestAdd(t *testing.T) {
	u := NewUsers(fbTestsURL + "test-add/users/")

	// unique name for tester each time this test is run
	name := fmt.Sprintf("tester-for-add-%d", time.Now().UnixNano())

	nicknames := []*Nickname{&Nickname{"tester", "tester.png"}}
	v := &User{name, nicknames}

	if err := u.Add(v); err != nil {
		t.Errorf("error adding user: %v\n", err)
	}

}

func TestGet(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-get/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"tester", "tester.png"})

	v := &User{"tester", nicknames}
	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	id := u.GetUserID("tester")
	if id == "" {
		t.Errorf("error getting user id\n")
	}

	b, err := u.Get(id)
	if err != nil {
		t.Errorf("error getting user\n")
	}

	if b.Name != v.Name {
		t.Errorf("expected both users to be equal\n")
	}
}

func TestGetByName(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-get-by-name/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"tester", "tester.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"tester", nicknames}

	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	b, err := u.GetByName("tester")
	if err != nil {
		t.Errorf("error getting user by name: %v\n", err)
	}
	if b.Name != v.Name {
		t.Errorf("expected users to be the same\n")
	}
}

func TestGetByID(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-get-by-id/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"tester", "tester.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"tester", nicknames}

	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	// get the ID
	if id := u.GetUserID("tester"); id == "" {
		t.Errorf("couldn't get user id by name\n")
	} else {
		b := u.GetByID(id)
		if b == nil {
			t.Errorf("error getting user by id\n")
		}
		if b.Name != "tester" {
			t.Errorf("expected both users to be equal\n")
		}
	}

}

// func TestUpdate(t *testing.T) {
// 	u := NewUsers(fbTestsURL + "tests-update/users/")
//
// 	v := &User{"tester", []string{"tester", "Ned Flanders"}, ""}
// 	// add the user don't care if it errors, we need to Get it later
// 	u.Add(v)
//
// 	// get the ID
// 	if id := u.GetUserID("tester"); id == "" {
// 		t.Errorf("couldn't get user id by name\n")
// 	} else {
// 		b := &User{"bob", []string{"snoop dog"}, ""}
//
// 		c, err := u.Update(id, b)
// 		if err != nil {
// 			t.Errorf("error updating user: %v\n", err)
// 		}
// 		if c.Name != "bob" {
// 			t.Errorf("expected user name to be bob\n")
// 		}
// 	}
// }

func TestAddExisting(t *testing.T) {
	u := NewUsers(fbTestsURL + "tests-add/existing/users/")

	nicknames := []*Nickname{&Nickname{"tester", "tester.png"}}
	v := &User{"tester", nicknames}
	// add a user
	u.Add(v)

	id := u.GetUserID("tester")
	if id == "" {
		t.Errorf("expected a user id")
	}

	q := &User{"tester", nicknames}
	if err := u.Add(q); err == nil {
		t.Errorf("expected an error (existing user): %v\n", err)
	}
}

func TestDelete1(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-delete-1/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"tester", "tester.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"tester", nicknames}
	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	if err := u.Delete("tester"); err != nil {
		t.Errorf("error deleting user: %v\n", err)
	}

	if err := u.Delete("non_existing"); err == nil {
		t.Errorf("expected an error while deleting a non existing user\n")
	}

}

func TestDelete2(t *testing.T) {

	u := NewUsers("")

	if err := u.Delete(""); err == nil {
		t.Errorf("expected an error%v\n")
	}

}

func TestDelete3(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-delete-3/users/")

	if err := u.Delete("tester"); err == nil {
		t.Errorf("expected an error%v\n")
	}

}
