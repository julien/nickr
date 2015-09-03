package data

import "testing"

func TestUserHasNickname(t *testing.T) {

	u := &User{"tester", []string{"bob"}, "profile.png"}
	if u.HasNickname("foo") {
		t.Errorf("expected false got true")
	}
}
