package jwtsign

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
)

func Test_TokenGenerator_ReceiveValueFromToken(t *testing.T) {
	cfg, err := config.NewConfig("../../../init/.env")
	require.NoError(t, err, fmt.Sprintf("jwt_test: config error - %v", err))
	err = NewSecretKey(cfg)
	require.NoError(t, err, fmt.Sprintf("jwt_test: NewSecretKey error - %v", err))
	content := "Task Access"

	tokenLine, err := TokenGenerator(content)
	require.NoError(t, err, "err - should be nil")
	require.Regexp(t, `[a-zA-Z0-9-_.]{148}`, tokenLine, "bad token")

	token, err := TokenRetrieve(tokenLine)
	require.NoError(t, err, fmt.Sprintf("jwt.Parse error - %v - should be nil", err))

	contentFromToken, err := ReceiveValueFromToken[string](token, "content")
	require.NoError(t, err, fmt.Sprintf("receive content error - %v - should be nil", err))
	require.Equal(t, content, contentFromToken, "content from token no euql start content")

	contentFromToken, err = ReceiveValueFromToken[string](token, "alien")
	assert.ErrorIs(t, err, jwt.ErrInvalidKey, "should be equal errors")
}
