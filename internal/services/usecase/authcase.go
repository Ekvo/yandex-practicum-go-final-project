// authcase - biz logic of autorization
package usecase

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/jwtsign"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

type AuthService interface {
	services.AutorizationCase
}

type authService struct{}

func NewAuthService() AuthService {
	return authService{}
}

// AuthZ - implemet of 'AutorizationCase interface' look (/internal/services/services.go)
//
// take cookie by key -> get jwt.token -> receive data from 'jwt.MapClaims' -> mark in 'log' string line"
func (a authService) AuthZ(r *http.Request) error {
	value, err := common.ReadCookie(r, "token")
	if err != nil {
		return common.ErrCookieEmptyKey
	}
	token, err := jwtsign.TokenRetrieve(value)
	if err != nil {
		return services.ErrServicesInternalError
	}
	content, err := jwtsign.ReceiveValueFromToken[string](token, "content")
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return err
		}
		return services.ErrServicesInternalError
	}
	log.Printf("authcase: %s", content)
	return nil
}
