// middlweare - describes the midlweare function
package transport

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// rulesForAuthZ - set of rules for 'AuthZ(authCase rulesForAuthZ, next http.HandlerFunc) http.HandlerFunc'
type rulesForAuthZ interface {
	services.AutorizationCase

	services.LoginValidPasswordCase
}

// AuthZ (Autorization)
//  1. check password (services.LoginValidPassword)
//
// password exist  -> (services.Autorization)
// password !exist -> call  next(w,r)
//
//  2. password exist:
//     use bizness logic of Autorization look (internal/services/usecase/authcase.go)
//     if error clean all cookie (MaxAge = -1) and see status code
//     no error -> call next(w,r)
func AuthZ(authCase rulesForAuthZ, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := authCase.UserExist(context.Background())
		if err != nil {
			err = services.ErrServicesInternalError
			common.EncodeJSON(w, http.StatusInternalServerError, common.NewError(err))
			return
		}
		if !ok {
			next(w, r)
			return
		}

		if err := authCase.AuthZ(r); err != nil {
			code := 0
			if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				//jwt.ErrTokenExpired or common.ErrCookieEmptyKey
				code = http.StatusUnauthorized
			}
			common.CleanCookie(w, r)
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		next(w, r)
	}
}
