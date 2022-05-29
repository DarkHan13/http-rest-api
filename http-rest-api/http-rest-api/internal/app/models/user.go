package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

//User Model
type User struct {
	Id                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	DecryptedPassword string `json:"-"`
	Password          string `json:"-"`
}

//Before Create User
func (u *User) BeforeCreate() error {
	if len(u.DecryptedPassword) > 0 {
		enc, err := EncryptString(u.DecryptedPassword)
		if err != nil {
			return err
		}
		u.Password = enc
	}
	return nil
}

func (u *User) Sanitize() {
	u.DecryptedPassword = ""
}

//Compare Passwords
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

//Validate User
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.DecryptedPassword, validation.By(RequiredIf(u.Password == "")), validation.Length(6, 30)),
	)
}

func EncryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
