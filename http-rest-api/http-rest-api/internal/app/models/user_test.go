package models_test

import (
	"testing"

	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := models.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.Password)
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *models.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *models.User {
				return models.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty Email",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid Email",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty Password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.DecryptedPassword = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short Password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.DecryptedPassword = "123"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
