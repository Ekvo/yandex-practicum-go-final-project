// loginencode - rules for encode jwt.Token after valid Login
package serializer

import "github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenEncode struct {
	Content string
}

func (t TokenEncode) Response() (TokenResponse, error) {
	token, err := common.TokenGenerator(t.Content)
	return TokenResponse{Token: token}, err
}
