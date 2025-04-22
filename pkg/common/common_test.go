package common

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func Test_Message_String(t *testing.T) {
	msg := Message{
		"id":      "1111",
		"title":   "some woed",
		"error":   "nil",
		"created": "empty",
	}
	msgLine := msg.String()
	assert.Equal(t, `{created:empty},{error:nil},{id:1111},{title:some woed}`, msgLine)
	assert.NotEqual(t, `{error:nil},{created:empty},{id:1111},{title:some woed}`, msgLine)
}

func Test_Abs(t *testing.T) {
	asserts := assert.New(t)
	var nums = []struct {
		num int
		res int
	}{{1, 1}, {-1, 1}, {-567, 567}, {123, 123}}
	for _, test := range nums {
		ans := Abs(test.num)
		asserts.Equal(test.res, ans)
	}
}

func Test_HashData(t *testing.T) {
	letters := []string{"1", "2", "3", "4", "a", "b", "s", "y", "i"}
	n := len(letters)
	for i := 0; i < 100; i++ {
		rand.Shuffle(n, func(i, j int) {
			letters[i], letters[j] = letters[j], letters[i]
		})
		password := strings.Join(letters, "")
		hashedPassword := HashData(password)

		assert.NotEqual(t, password, hashedPassword)
		assert.Regexp(t, `[0-9a-f]{64}`, hashedPassword)
	}
}

func Test_DecodeJSON_EncodeJSON(t *testing.T) {
	type user struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Age     uint   `json:"age"`
	}

	dataForRequest := []struct {
		description string
		body        string
		resCode     int
		respRegexp  string
		msg         string
	}{
		{
			description: `valid user`,
			body:        `{"name":"Alex","surname":"","age":26}`,
			resCode:     http.StatusOK,
			respRegexp:  `{"user":"approve"}`,
			msg:         `valid Decode and Encode`,
		},
		{
			description: `wrong user`,
			body:        `{"name":"Alex","surname":"","age":26,"avp":"alien"}`,
			resCode:     http.StatusBadRequest,
			respRegexp:  `{"error":"json: unknown field \\"avp\\""}`,
			msg:         `invalid Decode and valid Encode`,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		var u user
		if err := DecodeJSON(r, &u); err != nil {
			EncodeJSON(w, http.StatusBadRequest, Message{"error": err.Error()})
			return
		}
		EncodeJSON(w, http.StatusOK, Message{"user": "approve"})
	})

	for i, test := range dataForRequest {
		log.Printf("\t%d %s", i+1, test.description)

		req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(test.body)))
		require.NoError(t, err, "NewRequest error")
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.resCode, w.Code, "http status code should be equal "+test.msg)
		assert.Regexp(t, test.respRegexp, w.Body.String(), "body from Response not equal "+test.msg)
	}
}
