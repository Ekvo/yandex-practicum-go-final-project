package autorization

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
	"github.com/Ekvo/yandex-practicum-go-final-project/tests"
)

func init() {
	if err := godotenv.Load("../../../init/.env"); err != nil {
		log.Printf("autorization_test: no .env file - %v", err)
	}
	common.SecretKey = os.Getenv("TODO_SECRET_KEY")
	if common.SecretKey == "" {
		log.Printf("autorization_test: SecretKey is empty")
	}
}

func TestAuthZ(t *testing.T) {
	mux := http.ServeMux{}

	mux.HandleFunc("GET /test", AuthZ(func(w http.ResponseWriter, r *http.Request) {
		common.EncodeJSON(w, http.StatusOK, common.Message{"auth": "approve"})
	}))
	cookie := http.Cookie{
		Name:  "token",
		Value: tests.Token,
	}

	dataForRequest := []struct {
		cookie    bool
		resCode   int
		resRegexp string
		msg       string
	}{
		{
			cookie:    true,
			resCode:   http.StatusOK,
			resRegexp: `{"auth":"approve"}`,
			msg:       `valid auth`,
		},
		{
			cookie:    false,
			resCode:   http.StatusForbidden,
			resRegexp: `{"error":"http: named cookie not present"}`,
			msg:       `wrong auth with JSON error`,
		},
	}

	for _, test := range dataForRequest {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err, fmt.Sprintf("request create error - %v", err))
		if test.cookie {
			req.AddCookie(&cookie)
		}

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.resCode, w.Code, "status code not equal "+test.msg)
		assert.Regexp(t, test.resRegexp, w.Body.String(), "other body from response "+test.msg)

	}

}
