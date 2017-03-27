package users

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_GetUser_ShouldReturnUser_WhenLookingForExistingUser(t *testing.T) {
	username := "foo"
	password := "1234"

	user, err := GetUser(username, password)
	if err != nil {
		t.Errorf(`GetUser("%s", "%s") failed: %v`, username, password, err)
		return
	}
	if user == nil {
		t.Errorf(`GetUser("%s", "%s") could not find a user for the given username and password: %v`, username, password, err)
		return
	}

	if user.Username != username || user.Password != password {
		t.Errorf(`GetUser("%s", "%s") = %s, returned a User struct not properly filled: %v`, username, password, spew.Sprint(user), err)
		return
	}
}
