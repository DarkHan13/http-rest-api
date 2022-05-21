package models

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:             "user@example.org",
		DecryptedPassword: "password",
	}
}
