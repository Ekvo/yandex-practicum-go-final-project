// autorization - describes the midlweare function responsible for authorization validation
package autorization

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// AuthZ (Autorization)
//  1. check password from .env file
//
// password exist  -> 2 -> 3
// password !exist -> 3
//
//  2. SecretKey                           - empty -> error          - end
//     get token.(string) from ReadCookie  - if error - clean cookie - end
//     parse jwt.Token from token.(string) - if error - clean cookie - end
//     get content with check exploration in ReceiveValueFromToken - if error - clean cookie - end
//     write to log line content
//
//  3. call next(w,r)
func AuthZ(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password != "" {
			value, err := common.ReadCookie(r, "token")
			if err != nil {
				common.CleanCookie(w, r)
				common.EncodeJSON(w, http.StatusUnauthorized, common.Message{"error": err.Error()})
				return
			}
			token, err := jwt.Parse(value, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrHashUnavailable
				}
				if common.SecretKey == "" {
					return nil, common.ErrCommonEmptySecretKey
				}
				return []byte(common.SecretKey), nil
			})
			if err != nil {
				common.CleanCookie(w, r)
				common.EncodeJSON(w, http.StatusUnauthorized, common.Message{"error": err.Error()})
				return
			}
			content, err := common.ReceiveValueFromToken[string](token, "content")
			if err != nil {
				common.CleanCookie(w, r)
				common.EncodeJSON(w, http.StatusUnauthorized, common.Message{"error": err.Error()})
				return
			}
			log.Print(content)
		}
		next(w, r)
	})
}
