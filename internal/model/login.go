// login - describes the Login object
package model

import "context"

// password lenght [8,64]
const (
	MinPassword = 8
	MaxPassword = 64
)

type LoginModel struct {
	Password string
}

// LoginRead - check for and compare user password
type LoginRead interface {
	ValidLogin(ctx context.Context, login LoginModel) bool
	PasswordExist(ctx context.Context) bool
}
