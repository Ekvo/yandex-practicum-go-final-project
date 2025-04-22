// logindecode - rules for decode Login object from http.Request
package deserializer

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrServicesFiledLengthShort - field is shorter than minimum range
var ErrServicesFiledLengthShort = errors.New("short lenght")

type LoginDecode struct {
	Password string `json:"password"`

	login model.LoginModel `json:"-"`
}

func NewLoginDecode() *LoginDecode {
	return &LoginDecode{}
}

func (ld *LoginDecode) Model() model.LoginModel {
	return ld.login
}

// Decode - deserialize object LoginModel from Request
func (ld *LoginDecode) Decode(r *http.Request) error {
	if err := common.DecodeJSON(r, ld); err != nil {
		return err
	}
	msgErr := common.Message{}
	pass := ld.Password
	if len(pass) < model.MinPassword {
		msgErr["password"] = ErrServicesFiledLengthShort.Error()
	}
	if len(pass) > model.MaxPassword {
		msgErr["password"] = ErrServicesFiledLengthExceeded.Error()
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("login decode error - %s", msgErr.String())
	}
	ld.login.Password = common.HashData(ld.Password)
	return nil
}
