// logincase - logic of servicing 'LoginModel'
package usecase

import (
	"context"
	"errors"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/serializer"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

var ErrCaseLoginNotFound = errors.New("login not found")

type LoginService interface {
	services.LoginValidPasswordCase
}

// multiLogin - contain all LoginModel interfaces
type MultiLogin interface {
	model.LoginRead
}

type loginService struct {
	loginRepo MultiLogin
}

func NewLoginService(store MultiLogin) LoginService {
	return loginService{loginRepo: store}
}

// UserExist - member of loginService
// checking on exists login in application
func (l loginService) UserExist(ctx context.Context) (bool, error) {
	return l.loginRepo.PasswordExist(ctx), nil
}

// CreateToken - member of loginService
//
// create hash password from 'login' -> compare passwords with login inside of application
// create jwt.Token -> create object with Token for Response
func (l loginService) CreateToken(
	ctx context.Context,
	login model.LoginModel) (*serializer.TokenResponse, error) {
	login.Password = common.HashData(login.Password)
	if !l.loginRepo.ValidLogin(ctx, login) {
		return nil, ErrCaseLoginNotFound
	}
	serialize := serializer.TokenEncode{Content: "Task Access"}
	token, err := serialize.Response()
	if err != nil {
		return nil, services.ErrServicesInternalError
	}
	return token, nil
}
