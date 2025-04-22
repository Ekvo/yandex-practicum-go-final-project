package transport

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database/mock"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/datauser"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/jwtsign"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/usecase"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
	"github.com/Ekvo/yandex-practicum-go-final-project/tests"
)

var dataForRequest = []struct {
	description string
	method      string
	url         string
	body        string
	resCode     int
	resRegexp   string
	msg         string
}{
	{ //1
		description: `login valid`,
		method:      http.MethodPost,
		url:         `/api/signin`,
		body:        `{"password":"qwert12345"}`,
		resCode:     http.StatusOK,
		resRegexp:   `{"token":"[a-zA-Z0-9-_.]{148}"}`,
		msg:         `valid login, status 200, return token`,
	},
	{ //2
		description: `wrong login (bad password)`,
		method:      http.MethodPost,
		url:         `/api/signin`,
		body:        `{"password":"123456789"}`,
		resCode:     http.StatusForbidden,
		resRegexp:   `{"error":"login not found"}`,
		msg:         `login with invalid password, status 403, return JSON error`,
	},
	{ //3
		description: `wrong login (short password)`,
		method:      http.MethodPost,
		url:         `/api/signin`,
		body:        `{"password":"qwert"}`,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"login decode error - {password:short lenght}"`,
		msg:         `login with short password, status 422, return JSON error`,
	},
	{ //4
		description: `new task valid`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"date":"20240201","title": "Summarize","comment": "my comment","repeat": "d 5"}`,
		resCode:     http.StatusCreated,
		resRegexp:   `{"id":"1"}`,
		msg:         `save new task, status 201, return ID`,
	},
	{ //5
		description: `new task valid`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"date":"20240203","title":"Summarize","comment":"","repeat": ""}`,
		resCode:     http.StatusCreated,
		resRegexp:   `{"id":"2"}`,
		msg:         `save new task, status 201, return ID`,
	},
	{ //6
		description: `new task with ID invalid`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"id":"2","date":"20240203","title":"Summarize","comment":"","repeat": ""}`,
		resCode:     http.StatusConflict,
		resRegexp:   `{"error":"task already exist"}`,
		msg:         `task with ID = 2 is exist, status 409, return error`,
	},
	{ //7
		description: `new task invalid bad repeat`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"date":"20240203","title":"Summarize","comment":"","repeat": "k 1"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"invalid repeat data"}`,
		msg:         `wrong task bad repeat, status 422, return error`,
	},
	{ //8
		description: `new task invalid`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"date":"fff","title":"", "comment":"what it is","repeat":"p 765"}`,
		resCode:     http.StatusBadRequest,
		resRegexp:   ``,
		msg:         `bad task, status 422, return JSON error`,
	},
	{ //9
		description: `get task valid`,
		method:      http.MethodGet,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{"id":"1","date":"[0-9]{8}","title":"Summarize","comment":"my comment","repeat":"d 5"}`,
		msg:         `find task, status 200, return JSON TaskResponse`,
	},
	{ //10
		description: `task not exist`,
		method:      http.MethodGet,
		url:         `/api/task?id=3`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"task not found"}`,
		msg:         `task not exist, status 404, return JSON error`,
	},
	{ //11
		description: `task find with incorrect param`,
		method:      http.MethodGet,
		url:         `/api/task?id=alien`,
		body:        ``,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"invalid param"}`,
		msg:         `bad id in param, status 400, return JSON error`,
	},
	{ //12
		description: `task find with incorrect param`,
		method:      http.MethodGet,
		url:         `/api/task?id=0`,
		body:        ``,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"task ID is zero"}`,
		msg:         `bad id in param, status 400, return JSON error`,
	},
	{ //13
		description: `task change`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"id":"1","date": "20240201","title": "new title","comment": "change comment","repeat": "d 5"}`,
		resCode:     http.StatusOK,
		resRegexp:   ``,
		msg:         `update task, status 200, return empty JSON`,
	},
	{ //14
		description: `wrong update task`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"date":"fff","title":"T","comment":"what it is","repeat":"p 765"}`,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"taskdecode: error - {date:invalid date format}"}`,
		msg:         `bad task, status 422, return JSON error`,
	},
	{ //15
		description: `wrong update task`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"date":"20240201","title":"new title","comment":"change comment","repeat":"d 5"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"task ID is zero"}`,
		msg:         `wrong update task, status 422, return JSON error`,
	},
	{ //16
		description: `wrong update task not found ID`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"id":"3","date":"20240201","title":"new title 2","comment":"","repeat":""}`,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"task not found"}`,
		msg:         `task ID not found, status 422, return JSON error`,
	},
	{ //17
		description: `task change with bad repeat`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"id":"1","date": "20240201","title": "new title","comment": "change comment","repeat": "d 401"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"invalid repeat data"}`,
		msg:         `wrong update task, status 422, return error`,
	},
	{ //18
		description: `task list valid`,
		method:      http.MethodGet,
		url:         `/api/tasks?search=arize`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{"tasks":\[{"id":"2","date":"\d{8}","title":"Summarize","comment":"","repeat":""}]}`,
		msg:         `tasks list find, status 200, return JSON with 2 task`,
	},
	{ //19
		description: `valid task done`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is done, status 200, return empty JSON`,
	},
	{ //20
		description: `valid task done with delete`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=2`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is done delete from store, status 200, return empty JSON`,
	},
	{ //21
		description: `task done not found`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=2`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"task not found"}`,
		msg:         `task not found, status 404, return JSON error`,
	},
	{ //22
		description: `wrong task done (invalid ID)`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=ght`,
		body:        ``,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"invalid param"}`,
		msg:         `wrong task done, status 400, return JSON error`,
	},
	{ //23
		description: `valid delete task`,
		method:      http.MethodDelete,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is deleted, status 200, return empty JSON`,
	},
	{ //24
		description: `valid delete task (not found)`,
		method:      http.MethodDelete,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"task not found"}`,
		msg:         `task not exist, status 404, return JSON error`,
	},
	{ //25
		description: `wrong task delete (invalid ID)`,
		method:      http.MethodDelete,
		url:         `/api/task?id=ght`,
		body:        ``,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"invalid param"}`,
		msg:         `wrong task delete, status 400, return JSON error`,
	},
}

func TestRoutes(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	cfg, err := config.NewConfig(filepath.Join("..", "..", "init", ".env"))
	requires.NoError(err, fmt.Sprintf("transport_test: config error - %v - should be no error", err))
	requires.NoError(jwtsign.NewSecretKey(cfg), "transport_test: secret key error")

	//---------------------------------------------------------------------------------------
	log.Print("test of routes\n") // test routes
	//---------------------------------------------------------------------------------------

	type (
		mockTaskCase interface {
			services.TaskCreateCase
			services.TaskReadCase
			services.TaskUpdateCase
			services.TaskDeleteCase
			services.TaskDoneCase
		}

		mockSheduler struct {
			mockTaskCase

			services.AutorizationCase

			services.LoginValidPasswordCase
		}
	)

	taskCase, err := usecase.NewTaskService(cfg, mock.NewMockTaskStore())
	requires.NoError(err, fmt.Sprintf("transport_test: task service error - %v - should be no error", err))

	sheduler := mockSheduler{
		mockTaskCase:           taskCase,
		AutorizationCase:       usecase.NewAuthService(),
		LoginValidPasswordCase: usecase.NewLoginService(datauser.NewUserData(cfg)),
	}

	r := NewTransport(cfg)
	r.Routes(sheduler)

	cookie := http.Cookie{
		Name:  "token",
		Value: tests.Token,
	}

	for i, test := range dataForRequest {
		log.Printf("\t%d %s", i+1, test.description)

		req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer([]byte(test.body)))
		requires.NoError(err, fmt.Sprintf("request create error - %v", err))
		if test.method == http.MethodPost || test.method == http.MethodPut {
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		}
		req.AddCookie(&cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(test.resCode, w.Code, "status code not equal "+test.msg)
		asserts.Regexp(test.resRegexp, w.Body.String(), "other body from response "+test.msg)
	}

	//---------------------------------------------------------------------------------------
	log.Print("test of midlweare\n") // test midlweare
	//---------------------------------------------------------------------------------------

	type mockAuthLogin struct {
		services.AutorizationCase
		services.LoginValidPasswordCase
	}

	dataForAuth := []struct {
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
			resCode:   http.StatusInternalServerError,
			resRegexp: `{"error":"internal error"}`,
			msg:       `wrong auth with JSON error`,
		},
	}

	authLogin := mockAuthLogin{
		AutorizationCase:       usecase.NewAuthService(),
		LoginValidPasswordCase: usecase.NewLoginService(datauser.NewUserData(cfg)),
	}

	mux := http.ServeMux{}

	mux.HandleFunc("GET /test", AuthZ(
		authLogin,
		func(w http.ResponseWriter, r *http.Request) {
			common.EncodeJSON(w, http.StatusOK, common.Message{"auth": "approve"})
		}))

	for _, test := range dataForAuth {
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
