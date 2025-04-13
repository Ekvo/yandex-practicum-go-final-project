package transport

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
	"github.com/Ekvo/yandex-practicum-go-final-project/tests"
)

type mockStore struct {
	id    *uint
	tasks map[uint]model.TaskModel
}

func NewmockStore() mockStore {
	return mockStore{
		id:    new(uint),
		tasks: make(map[uint]model.TaskModel),
	}
}

func (s mockStore) incrementID() {
	*s.id++
}

func (s mockStore) SaveOneTask(_ context.Context, data any) (uint, error) {
	newTask := data.(model.TaskModel)
	s.incrementID()
	id := *s.id
	newTask.ID = id
	s.tasks[id] = newTask
	return id, nil
}

func (s mockStore) FindOneTask(_ context.Context, data any) (model.TaskModel, error) {
	id := data.(uint)
	task, ex := s.tasks[id]
	if !ex {
		return model.TaskModel{}, database.ErrDataBaseNotFound
	}
	return task, nil
}

func (s mockStore) NewDataTask(_ context.Context, data any) error {
	updateTask := data.(model.TaskModel)
	id := updateTask.ID
	if _, ex := s.tasks[id]; !ex {
		return database.ErrDataBaseNotFound
	}
	s.tasks[id] = updateTask
	return nil
}

func (s mockStore) ExpirationTask(ctx context.Context, data any) error {
	taskID := data.(uint)
	if _, ex := s.tasks[taskID]; !ex {
		return database.ErrDataBaseNotFound
	}
	delete(s.tasks, taskID)
	return nil
}

func (s mockStore) FindTaskList(_ context.Context, data any) ([]model.TaskModel, error) {
	property := data.(*services.TaskProperty)
	var arrOfTask []model.TaskModel

	if property.IsWord() {
		word := property.PassWord()
		for _, task := range s.tasks {
			if strings.Contains(task.Title, word) ||
				strings.Contains(task.Comment, word) {
				arrOfTask = append(arrOfTask, task)
			}
		}
	} else if property.IsDate() {
		date := property.PassDate().UTC().Format(model.DateFormat)
		for _, task := range s.tasks {
			if task.Date == date {
				arrOfTask = append(arrOfTask, task)
			}
		}
	}
	limit := property.PassLimit()
	if len(arrOfTask) > int(limit) {
		return arrOfTask[:limit], nil
	}
	return arrOfTask, nil
}

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
		resRegexp:   `{"error":"invalid password"}`,
		msg:         `login with invalid password, status 403, return JSON error`,
	},
	{ //3
		description: `wrong login (short password)`,
		method:      http.MethodPost,
		url:         `/api/signin`,
		body:        `{"password":"qwert"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"login decode error - {password : short lenght}"`,
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
		description: `new task invalid`,
		method:      http.MethodPost,
		url:         `/api/task`,
		body:        `{"date":"fff","title":"", "comment":"what it is","repeat":"p 765"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   ``,
		msg:         `bad task, status 422, return JSON error`,
	},
	{ //7
		description: `get task valid`,
		method:      http.MethodGet,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{"id":"1","date":"[0-9]{8}","title":"Summarize","comment":"my comment","repeat":"d 5"}`,
		msg:         `find task, status 200, return JSON TaskResponse`,
	},
	{ //8
		description: `task not exist`,
		method:      http.MethodGet,
		url:         `/api/task?id=3`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"resource not found"}`,
		msg:         `task not exist, status 404, return JSON error`,
	},
	{ //9
		description: `task find with incorrect param`,
		method:      http.MethodGet,
		url:         `/api/task?id=alien`,
		body:        ``,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"invalid param"}`,
		msg:         `bad id in param, status 400, return JSON error`,
	},
	{ //10
		description: `task change`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"id":"1","date": "20240201","title": "new title","comment": "change comment","repeat": "d 5"}`,
		resCode:     http.StatusOK,
		resRegexp:   ``,
		msg:         `update task, status 200, return empty JSON`,
	},
	{ //11
		description: `wrong update task`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"date":"fff","title":"T","comment":"what it is","repeat":"p 765"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"task decode error - {date : invalid date}"}`,
		msg:         `bad task, status 422, return JSON error`,
	},
	{ //12
		description: `wrong update task`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"date":"20240201","title":"new title","comment":"change comment","repeat":"d 5"}`,
		resCode:     http.StatusUnprocessableEntity,
		resRegexp:   `{"error":"not numeric"}`,
		msg:         `wrong update task, status 422, return JSON error`,
	},
	{ //13
		description: `wrong update task not found ID`,
		method:      http.MethodPut,
		url:         `/api/task`,
		body:        `{"id":"3","date":"20240201","title":"new title 2","comment":"","repeat":""}`,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"resource not found"}`,
		msg:         `task ID not found, status 422, return JSON error`,
	},
	{ //14
		description: `task list valid`,
		method:      http.MethodGet,
		url:         `/api/tasks?search=arize`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{"tasks":\[{"id":"2","date":"\d{8}","title":"Summarize","comment":"","repeat":""}]}`,
		msg:         `tasks list find, status 200, return JSON with 2 task`,
	},
	{ //15
		description: `valid task done`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is done, status 200, return empty JSON`,
	},
	{ //16
		description: `valid task done with delete`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=2`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is done delete from store, status 200, return empty JSON`,
	},
	{ //17
		description: `task done not found`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=2`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"resource not found"}`,
		msg:         `task not found, status 404, return JSON error`,
	},
	{ //18
		description: `wrong task done (invalid ID)`,
		method:      http.MethodPost,
		url:         `/api/task/done?id=ght`,
		body:        ``,
		resCode:     http.StatusBadRequest,
		resRegexp:   `{"error":"invalid param"}`,
		msg:         `wrong task done, status 400, return JSON error`,
	},
	{ //19
		description: `valid delete task`,
		method:      http.MethodDelete,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusOK,
		resRegexp:   `{}`,
		msg:         `task is deleted, status 200, return empty JSON`,
	},
	{ //20
		description: `valid delete task (not found)`,
		method:      http.MethodDelete,
		url:         `/api/task?id=1`,
		body:        ``,
		resCode:     http.StatusNotFound,
		resRegexp:   `{"error":"resource not found"}`,
		msg:         `task not exist, status 404, return JSON error`,
	},
	{ //21
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

	err := godotenv.Load("../../init/.env")
	requires.NoError(err, fmt.Sprintf("transport: no .env file error - %v", err))
	common.SecretKey = os.Getenv("TODO_SECRET_KEY")
	requires.NotEmpty(common.SecretKey, "transport: SecretKey is empty")

	store := NewmockStore()
	r := NewTransport(http.NewServeMux())
	r.Routes(store)

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
}
