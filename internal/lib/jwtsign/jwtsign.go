//	jwtsign - container for 'secretKey'
//
// 'secretKey' - use for create, parse - jwt.Token
package jwtsign

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
)

var (
	// ErrJwtSetSecretKeyAfterRun - secret k
	ErrJwtSetSecretKeyAfterRun = errors.New("secret key already exists")

	// ErrJwtEmptySecretKey - use only in
	ErrJwtEmptySecretKey = errors.New("secret key not found")
)

// SecretKey -  key for jwt.Token
var secretKey = ""

// NewSecretKey - call in Run -> during application startup
func NewSecretKey(cfg *config.Config) error {
	if secretKey == "" {
		secretKey = cfg.JWTSecretKey
		return nil
	}
	return ErrJwtSetSecretKeyAfterRun
}

// TokenRetrieve - wrapper for jwt.Parse
func TokenRetrieve(value string) (*jwt.Token, error) {
	return jwt.Parse(value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		if secretKey == "" {
			return nil, ErrJwtEmptySecretKey
		}
		return []byte(secretKey), nil
	})
}

// time exploration for jwt.Token see 'TokenGenerator'
const tokenLife = 7 * 24 * time.Hour

// TokenGenerator - create jwt token using specific key
// set time of exploration in claims
func TokenGenerator(content string) (string, error) {
	if secretKey == "" {
		return "", ErrJwtEmptySecretKey
	}
	claims := jwt.MapClaims{
		"content":     content,
		"exploration": time.Now().UTC().Add(tokenLife).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secretKey))
}

// ReceiveValueFromToken - get value by key from jwt.Token
//
// 1. check time exploration -> if Expired - error
// 2. get value by key
func ReceiveValueFromToken[T any](token *jwt.Token, key string) (T, error) {
	var obj T
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return obj, jwt.ErrTokenInvalidClaims
	}
	exploration, ok := claims["exploration"].(float64)
	if !ok {
		return obj, jwt.ErrInvalidKey
	}
	if int64(exploration) < time.Now().UTC().Unix() {
		return obj, jwt.ErrTokenExpired
	}
	value, ok := claims[key].(T)
	if !ok {
		return obj, jwt.ErrInvalidKey
	}
	return value, nil
}
