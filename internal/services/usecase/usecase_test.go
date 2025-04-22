package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database/mock"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/datauser"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/jwtsign"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/nextdate"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/serializer"
	"github.com/Ekvo/yandex-practicum-go-final-project/tests"
)

func Test_All_Usecase(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	cfg, err := config.NewConfig(filepath.Join("..", "..", "..", "init", ".env"))
	requires.NoError(err, fmt.Sprintf("usecase_test: config error - %v - should be no error", err))
	requires.NoError(jwtsign.NewSecretKey(cfg), "usecase_test: secret key error")

	//--------------------------------------------------------------------------------------
	log.Print("test - Task Service\n")
	//--------------------------------------------------------------------------------------

	var ( // for compare type
		nilPtrTaskIDResponse *serializer.TaskIDResponse = nil
		nilPtrTaskResponse   *serializer.TaskResponse   = nil
	)

	var dataForTaskService = []struct {
		description string
		init        func(ctx context.Context, ts TaskService, data any) (any, error)
		ctxTimeOut  time.Duration
		data        any
		expectedRes any
		err         error
		msg         string
	}{
		{ // 1
			description: `task create no error`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.CreateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				Date:    "20251003",
				Title:   "first",
				Comment: "ololo",
				Repeat:  "d 1",
			},
			expectedRes: &serializer.TaskIDResponse{ID: `^[1-9][0-9]*$`},
			err:         nil,
			msg:         `should return *TaskResponse and error is nil`,
		},
		{ // 2
			description: `task create wrong already exist`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.CreateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				ID:      1,
				Date:    "20251003",
				Title:   "second",
				Comment: "ololo",
				Repeat:  "d 1",
			},
			expectedRes: nilPtrTaskIDResponse,
			err:         ErrCaseTaskAlreadyExist,
			msg:         `should return nil and error`,
		},
		{ // 3
			description: `task create wrong data format`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.CreateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				Date:  "2025.10.03",
				Title: "third",
			},
			expectedRes: nilPtrTaskIDResponse,
			err:         nextdate.ErrNextDateInvalidDate,
			msg:         `should return nil and error`,
		},
		{ // 4
			description: `task create no error, only title in TaskModel`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.CreateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				Title: "fourh",
			},
			expectedRes: &serializer.TaskIDResponse{ID: `^[1-9][0-9]*$`},
			err:         nil,
			msg:         `should return *TaskResponse and error is nil`,
		},
		{ // 5
			description: `task Read no error`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.ReadTask(ctx, data.(uint))
			},
			ctxTimeOut: 100 * time.Second,
			data:       uint(1),
			expectedRes: &serializer.TaskResponse{
				Title:   "first",
				Comment: "ololo",
				Repeat:  "d 1",
			},
			err: nil,
			msg: `should return *TaskResponse and error is nil`,
		},
		{ // 6
			description: `task Read not found`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.ReadTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(1_000_000),
			expectedRes: nilPtrTaskResponse,
			err:         ErrCaseTaskNotFound,
			msg:         `should return nil and error`,
		},
		{ // 7
			description: `task Update valid`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.UpdateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				ID:      1,
				Date:    "20261003",
				Title:   "fifth",
				Comment: "abcd",
				Repeat:  "w 3,4,5",
			},
			expectedRes: nil,
			err:         nil,
			msg:         `should return nil and error is nil`,
		},
		{ // 8
			description: `task Update not found`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.UpdateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				ID:      1_000_000,
				Date:    "20261003",
				Title:   "fifth",
				Comment: "abcd",
				Repeat:  "w 3,4,5",
			},
			expectedRes: nil,
			err:         ErrCaseTaskNotFound,
			msg:         `should return nil and error is Not found`,
		},
		{ // 9
			description: `task Update bad data`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.UpdateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				ID:    1,
				Date:  "sagfssdf",
				Title: "sixth",
			},
			expectedRes: nil,
			err:         nextdate.ErrNextDateInvalidDate,
			msg:         `should return nil and error is nil`,
		},
		{ // 10
			description: `task Update bad repeat`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.UpdateTask(ctx, data.(model.TaskModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.TaskModel{
				ID:     1,
				Title:  "sixth",
				Repeat: "d rtv",
			},
			expectedRes: nil,
			err:         nextdate.ErrNextDateWrongRepeat,
			msg:         `should return nil and error not found`,
		},
		{ // 11
			description: `task done valid`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DoneTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(1),
			expectedRes: nil,
			err:         nil,
			msg:         `should update date res is nil and error nil`,
		},
		{ // 12
			description: `task done not found`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DoneTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(1_000_000),
			expectedRes: nil,
			err:         ErrCaseTaskNotFound,
			msg:         `wrong done task res is nil and error - not found`,
		},
		{ // 13
			description: `task list`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return ts.ReadTaskList(ctx, data.(*entity.TaskProperty))
			},
			ctxTimeOut: 100 * time.Second,
			data:       entity.NewTaskProperty("", 123),
			expectedRes: &serializer.TaslListResponse{
				// ID,Date set from Response
				TasksResp: []serializer.TaskResponse{
					{
						ID:      "1",
						Title:   "fifth",
						Comment: "abcd",
						Repeat:  "w 3,4,5",
					},
					{
						ID:    "2",
						Title: "fourh",
					},
				},
			},
			err: nil,
			msg: `task list res is list with lenght - 2 and error - nil`,
		},
		{ // 14
			description: `task done valid`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DoneTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(2),
			expectedRes: nil,
			err:         nil,
			msg:         `should update date res is nil and error nil`,
		},
		{ // 15
			description: `task delete not found`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DeleteTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(2),
			expectedRes: nil,
			err:         ErrCaseTaskNotFound,
			msg:         `wrong delete task, res is nil and error - not found`,
		},
		{ // 16
			description: `task delete`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DeleteTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(1),
			expectedRes: nil,
			err:         nil,
			msg:         `valid delete task, res and err is nil`,
		},
		{ // 17
			description: `task delete not found`,
			init: func(ctx context.Context, ts TaskService, data any) (any, error) {
				return nil, ts.DeleteTask(ctx, data.(uint))
			},
			ctxTimeOut:  100 * time.Second,
			data:        uint(1),
			expectedRes: nil,
			err:         ErrCaseTaskNotFound,
			msg:         `wrong delete task, res is nil and error - not found`,
		},
	}

	ctx := context.Background()

	taskService, err := NewTaskService(cfg, mock.NewMockTaskStore())
	requires.NoError(err, fmt.Sprintf("usecase_test: task service error - %v - should be no error", err))

	for i, test := range dataForTaskService {
		log.Printf("\t%d %s", i+1, test.description)

		ctx, cancel := context.WithTimeout(ctx, test.ctxTimeOut)
		defer cancel()

		res, err := test.init(ctx, taskService, test.data)

		asserts.ErrorIs(test.err, err, "erros no equal "+test.msg)

		switch v := res.(type) {
		case *serializer.TaskIDResponse:
			numeric, ok := test.expectedRes.(*serializer.TaskIDResponse)
			requires.True(ok, "false!!! "+test.msg)
			if numeric == nil {
				asserts.Nil(v, "should be nil "+test.msg)
			} else {
				requires.NotNil(v, "should be not nil "+test.msg)
				asserts.Regexp(numeric.ID, v.ID, "should be numeric "+test.msg)
			}
		case *serializer.TaskResponse:
			expectedTask, ok := test.expectedRes.(*serializer.TaskResponse)
			requires.True(ok, "false!!! "+test.msg)
			if expectedTask == nil {
				asserts.Nil(v, "should be nil "+test.msg)
			} else {
				requires.NotNil(v, "should be not nil "+test.msg)
				expectedTask.ID = v.ID
				expectedTask.Date = v.Date // see 'executeDate' ./taskcase.go
				asserts.Equal(*expectedTask, *v, "TaskResponse not equal "+test.msg)
			}
		case *serializer.TaslListResponse:
			expectedTasks, ok := test.expectedRes.(*serializer.TaslListResponse)
			requires.True(ok, "false!!! "+test.msg)
			if expectedTasks == nil {
				asserts.Nil(v, "should be nil "+test.msg)
			} else {
				requires.NotNil(v, "should be not nil "+test.msg)
				requires.Len(v.TasksResp, len(expectedTasks.TasksResp), "arrays of TaslListResponse not equal "+test.msg)
				sort.Slice(expectedTasks.TasksResp, func(i, j int) bool {
					return expectedTasks.TasksResp[i].ID < expectedTasks.TasksResp[j].ID
				})
				sort.Slice(expectedTasks.TasksResp, func(i, j int) bool {
					return v.TasksResp[i].ID < v.TasksResp[j].ID
				})
				for i, _ := range expectedTasks.TasksResp {
					expectedTasks.TasksResp[i].Date = v.TasksResp[i].Date
				}
				asserts.Equal(*expectedTasks, *v, "compare TaslListResponse - faild "+test.msg)
			}
		case nil:
			asserts.Nil(test.expectedRes, "should be nil "+test.msg)
		default:
			requires.FailNow("incorrect type in switch RES " + test.msg)
		}
	}

	//--------------------------------------------------------------------------------------
	log.Print("test - Login Service\n")
	//--------------------------------------------------------------------------------------

	var nulPtrTokenResponse *serializer.TokenResponse = nil

	var dataForLoginService = []struct {
		description string
		init        func(ctx context.Context, l LoginService, data any) (any, error)
		ctxTimeOut  time.Duration
		data        any
		expectedRes any
		err         error
		msg         string
	}{
		{ // 1
			description: `valid login`,
			init: func(ctx context.Context, l LoginService, data any) (any, error) {
				return l.CreateToken(ctx, data.(model.LoginModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.LoginModel{
				Password: `qwert12345`,
			},
			expectedRes: &serializer.TokenResponse{
				Token: `[a-zA-Z0-9-_.]{148}`,
			},
			err: nil,
			msg: `correct login return token, error is nil`,
		},
		{ // 2
			description: `invalid login alien password`,
			init: func(ctx context.Context, l LoginService, data any) (any, error) {
				return l.CreateToken(ctx, data.(model.LoginModel))
			},
			ctxTimeOut: 100 * time.Second,
			data: model.LoginModel{
				Password: `qwerty123456`,
			},
			expectedRes: nulPtrTokenResponse,
			err:         ErrCaseLoginNotFound,
			msg:         `wrong login return nil, error - not found`,
		},
		{ // 3
			description: `login exist password`,
			init: func(ctx context.Context, l LoginService, _ any) (any, error) {
				return l.UserExist(ctx)
			},
			ctxTimeOut:  100 * time.Second,
			expectedRes: true,
			err:         nil,
			msg:         `login exist return true, error is nil`,
		},
	}

	loginService := NewLoginService(datauser.NewUserData(cfg))

	for i, test := range dataForLoginService {
		log.Printf("\t%d %s", i+1, test.description)

		ctx, cancel := context.WithTimeout(ctx, test.ctxTimeOut)
		defer cancel()

		res, err := test.init(ctx, loginService, test.data)

		asserts.ErrorIs(test.err, err, "erros no equal "+test.msg)

		switch v := res.(type) {
		case *serializer.TokenResponse:
			expectedToken, ok := test.expectedRes.(*serializer.TokenResponse)
			requires.True(ok, "false!!! "+test.msg)
			if expectedToken == nil {
				asserts.Nil(v, "should be nil "+test.msg)
			} else {
				requires.NotNil(v, "should be not nil "+test.msg)
				asserts.Regexp(expectedToken.Token, v.Token, "regexp compare should be valid "+test.msg)
			}
		case bool:
			expected, ok := test.expectedRes.(bool)
			requires.True(ok, "false!!! "+test.msg)
			asserts.Equal(expected, v, "should be equal "+test.msg)
		default:
			requires.FailNow("incorrect type in switch RES " + test.msg)
		}
	}

	var dataForAuthService = []struct {
		description string
		cookie      *http.Cookie
		err         error
		msg         string
	}{
		{ // 1
			description: `valid auth`,
			cookie: &http.Cookie{
				Name:    `token`,
				Value:   tests.Token,
				Expires: time.Now().Add(7 * 24 * time.Hour),
			},
			err: nil,
			msg: `approve auth, error - is nil`,
		},
		{ // 2
			description: `wrong auth coolkie not contain "token"`,
			cookie: &http.Cookie{
				Name:    `alien`,
				Value:   tests.Token,
				Expires: time.Now().Add(7 * 24 * time.Hour),
			},
			err: services.ErrServicesInternalError,
			msg: `invalid auth, error - internal `,
		},
		{ // 3
			description: `wrong auth coolkie with broken token`,
			cookie: &http.Cookie{
				Name:    `token`,
				Value:   `some line`,
				Expires: time.Now().Add(7 * 24 * time.Hour),
			},
			err: services.ErrServicesInternalError,
			msg: `invalid auth, error - internal `,
		},
		{ //4
			description: `wrong auth coolkie with old token`,
			cookie: &http.Cookie{
				Name: `token`,
				// old token
				Value:   `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250ZW50IjoiVGFzayBBY2Nlc3MiLCJleHBsb3JhdGlvbiI6MTc0NDk2NDI2NH0.5KCkfZJwbFoB--pLNGe-Qat8_hVKGtUvMv5TC3zHkn8`,
				Expires: time.Now().Add(7 * 24 * time.Hour),
			},
			err: jwt.ErrTokenExpired,
			msg: `invalid auth, error - internal `,
		},
	}

	authService := NewAuthService()

	for i, test := range dataForAuthService {
		log.Printf("\t%d %s", i+1, test.description)

		req, err := http.NewRequest(http.MethodGet, "/empty", nil)
		requires.NoError(err, fmt.Sprintf("request create error - %v", err))
		req.AddCookie(test.cookie)

		err = authService.AuthZ(req)

		asserts.ErrorIs(test.err, err, "unexpected error "+test.msg)
	}
}
