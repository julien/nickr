package main

import (
	"errors"
	"fmt"

	"github.com/julien/nickr/Godeps/_workspace/src/github.com/melvinmt/firebase"
)

// Users is the main data model.
type User struct {
	Name      string   `json:"name"`
	Nicknames []string `json:"nicknames"`
	Picture   string   `json:"picture"`
}

// HasNickname returns a boolean indicating if a User has a nickname or not.
func (u *User) HasNickname(nick string) bool {
	for _, v := range u.Nicknames {
		if v == nick {
			return true
		}
	}
	return false
}

// Users is a collection of User objects, that uses Firebase as it's persistance layer.
type Users struct {
	data  map[string]User
	fbURL string // Firebase URL
}

// NewUsers returns a new collection of "User" objects.
func NewUsers(fbURL string) *Users {
	return &Users{fbURL: fbURL}
}

// All returns all the users in this collection.
func (u *Users) All() (map[string]User, error) {
	ref := firebase.NewReference(u.fbURL)
	if err := ref.Value(&u.data); err != nil {
		return nil, err
	}
	return u.data, nil
}

// Add adds a new user to this collection.
func (u *Users) Add(usr *User) error {
	ref := firebase.NewReference(u.fbURL)

	if usr.Name == "" || len(usr.Nicknames) < 1 {
		return errors.New("Invalid user")
	}

	existing, err := u.GetByName(usr.Name)
	if err != nil {
		return err
	}

	if existing != nil {
		return errors.New("existing user")
	} else {
		if err := ref.Push(usr); err != nil {
			return err
		}
	}
	return nil
}

// Get returns a user given an id.
func (u *Users) Get(id string) (*User, error) {
	ref := firebase.NewReference(u.fbURL + id).Export(false)
	usr := &User{}
	if err := ref.Value(usr); err != nil {
		return nil, err
	}
	return usr, nil
}

// GetByName returns a user given a name.
func (u *Users) GetByName(name string) (*User, error) {
	if u.data == nil {
		if _, err := u.All(); err != nil {
			return nil, err
		}
	}

	for _, v := range u.data {

		fmt.Printf("looking up user: %v\n", v.Name == name)

		if v.Name == name {
			fmt.Printf("found user: %v\n", v)
			return &v, nil
		}
	}
	return nil, nil
}

// GetByID returns a user given an id.
func (u *Users) GetByID(id string) *User {
	if u.data == nil {
		u.All()
	}
	for k, v := range u.data {
		if k == id {
			return &v
		}
	}
	return nil
}

// GetUserID returns a user's id given a name.
func (u *Users) GetUserID(name string) string {
	if u.data == nil {
		if _, err := u.All(); err != nil {
			return ""
		}
	}
	for k, v := range u.data {
		if v.Name == name {
			return k
		}
	}
	return ""
}

// Update a user given an id and a "User" object.
func (u *Users) Update(id string, v *User) (*User, error) {
	if usr := u.GetByID(id); usr != nil {

		if v.Name != "" && usr.Name != v.Name {
			usr.Name = v.Name
		}

		for i := 0; i < len(v.Nicknames); i++ {
			if !usr.HasNickname(v.Nicknames[i]) {
				usr.Nicknames = append(usr.Nicknames, v.Nicknames[i])
			}
		}

		if v.Picture != "" && usr.Picture != v.Picture {
			usr.Picture = v.Picture
		}

		ref := firebase.NewReference(u.fbURL + id)

		if err := ref.Write(usr); err != nil {
			return nil, err
		}

		return usr, nil
	}

	return nil, nil
}

// Delete a user with the specified name.
func (u *Users) Delete(name string) error {
	if id := u.GetUserID(name); id != "" {
		ref := firebase.NewReference(u.fbURL + id)
		if err := ref.Delete(); err != nil {
			return err
		}
		// also delete from the map
		delete(u.data, id)
	} else {
		return errors.New("user not found")
	}
	return nil
}
