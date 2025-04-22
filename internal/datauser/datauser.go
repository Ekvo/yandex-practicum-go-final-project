// datauser - store for user password
//
// password set from the .env file during application startup, see 'NewUserData(cfg *config.Config) UserData'
package datauser

import (
	"context"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
)

type UserData struct {
	Password string
}

// cfg.User.Password - may be empty -> application does not require authorization
func NewUserData(cfg *config.Config) UserData {
	return UserData{Password: cfg.UserPassword}
}

func (ul UserData) ValidLogin(_ context.Context, login model.LoginModel) bool {
	return ul.Password == login.Password
}

func (ul UserData) PasswordExist(_ context.Context) bool {
	return len(ul.Password) > 0
}
