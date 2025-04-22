package deserializer

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

func TestLoginDecode_Decode(t *testing.T) {
	mux := http.ServeMux{}

	mux.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		deserialize := NewLoginDecode()
		if err := deserialize.Decode(r); err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{"login": "approve"})
	})

	dataForRequest := []struct {
		body      string
		resCode   int
		resRegexp string
		msg       string
	}{
		{
			body:      `{"password":"qwert12345"}`,
			resCode:   http.StatusOK,
			resRegexp: `{"login":"approve"}`,
			msg:       `valid decode`,
		},
		{
			body:      `{"password":"qwer"}`,
			resCode:   http.StatusBadRequest,
			resRegexp: `{"error":"login decode error - {password:short lenght}"}`,
			msg:       `invalid decode`,
		},
	}

	for _, test := range dataForRequest {
		req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(test.body)))
		require.NoError(t, err, fmt.Sprintf("request create error - %v", err))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.resCode, w.Code, "status code not equal "+test.msg)
		assert.Regexp(t, test.resRegexp, w.Body.String(), "other body from response "+test.msg)
	}
}

func TestTaskDecode_Decode(t *testing.T) {
	mux := http.ServeMux{}

	mux.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		deserialize := NewTaskDecode()
		if err := deserialize.Decode(r); err != nil {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{"task": "approve"})
	})

	dataForRequest := []struct {
		body      string
		resCode   int
		resRegexp string
		msg       string
	}{
		{
			body:      `{"date":"20240201","title":"Summarize","comment":"my comment","repeat":"d 5"}`,
			resCode:   http.StatusOK,
			resRegexp: `{"task":"approve"}`,
			msg:       `valid decode`,
		},
		{
			body:      `{"id":"no numeric","date":"ffff","title":"","comment":"my comment","repeat":"sgfgasoigadgiohadioghdwioghwdijoghwiojghwijghiadjghjdklsngdjsfghueiroqwghjwdighveioqghdjiwbviqejghvkjwdbvieqwghijenvbijqweghiqejnbviodjfvhbqeioghvdklajfbqeioghqdiojvbqdiojghiqeghqeojvhiqepufvhipjvhqeighiqeprhvwdiofjvhiqepughvqjipevhqeiopghqiepghqiejpghiqefjghvqeigv"}`,
			resCode:   http.StatusUnprocessableEntity,
			resRegexp: `{"error":"taskdecode: error - {date:invalid date format},{id:not numeric},{repeat:length exceeded},{title:empty}"}`,
			msg:       `invalid decode`,
		},
	}

	for _, test := range dataForRequest {
		req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(test.body)))
		require.NoError(t, err, fmt.Sprintf("request create error - %v", err))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		assert.Equal(t, test.resCode, w.Code, "status code not equal "+test.msg)
		assert.Regexp(t, test.resRegexp, w.Body.String(), "other body from response "+test.msg)
	}
}
