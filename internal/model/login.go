// login - describes the Login object
package model

import "errors"

var ErrModelsLoginInvalidPassword = errors.New("invalid password")

// password lenght [8,64]
const (
	MinPassword = 8
	MaxPassword = 64
)

type LoginModel struct {
	Password string
}

// ValidPassword - compare password from .env file "TODO_PASSWORD"
func (lm LoginModel) ValidPassword(password string) bool {
	return lm.Password == password
}
