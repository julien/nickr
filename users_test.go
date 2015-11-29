package main

import (
	"fmt"
	"testing"
	"time"
)

const fbTestsURL = "https://nickr.firebaseio.com/tests/"

func TestUserHasNickname(t *testing.T) {
	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"Gabbo", "gabbo.png"})

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

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"2", "2.png"})

	v := &User{name, nicknames}

	if err := u.Add(v); err != nil {
		t.Errorf("error adding user: %v\n", err)
	}

}

func TestGet(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-get/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"Gabbo", "gabbo.png"})

	v := &User{"joe", nicknames}
	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	id := u.GetUserID("joe")
	if id == "" {
		t.Errorf("error getting user id\n")
	}

	b, err := u.Get(id)
	if err != nil {
		t.Errorf("error getting user\n")
	}

	fmt.Printf("Got user: %v\n", b)
	if b.Name != v.Name {
		t.Errorf("expected both users to be equal\n")
	}
}

func TestGetByName(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-get-by-name/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"Gabbo", "gabbo.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"joe", nicknames}

	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	b, err := u.GetByName("joe")
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
	nicknames = append(nicknames, &Nickname{"Gabbo", "gabbo.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"joe", nicknames}

	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	// get the ID
	if id := u.GetUserID("joe"); id == "" {
		t.Errorf("couldn't get user id by name\n")
	} else {
		b := u.GetByID(id)
		if b == nil {
			t.Errorf("error getting user by id\n")
		}
		if b.Name != "joe" {
			t.Errorf("expected both users to be equal\n")
		}
	}

}

// func TestUpdate(t *testing.T) {
// 	u := NewUsers(fbTestsURL + "tests-update/users/")
//
// 	v := &User{"joe", []string{"Gabbo", "Ned Flanders"}, ""}
// 	// add the user don't care if it errors, we need to Get it later
// 	u.Add(v)
//
// 	// get the ID
// 	if id := u.GetUserID("joe"); id == "" {
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

func TestDelete(t *testing.T) {

	u := NewUsers(fbTestsURL + "tests-delete/users/")

	var nicknames []*Nickname
	nicknames = append(nicknames, &Nickname{"Gabbo", "gabbo.png"})
	nicknames = append(nicknames, &Nickname{"Ned Flanders", "ned_flanders.png"})

	v := &User{"joe", nicknames}
	// add the user don't care if it errors, we need to Get it later
	u.Add(v)

	if err := u.Delete("joe"); err != nil {
		t.Errorf("error deleting user: %v\n", err)
	}

	if err := u.Delete("non_existing"); err == nil {
		t.Errorf("expected an error while deleting a non existing user\n")
	}

}
