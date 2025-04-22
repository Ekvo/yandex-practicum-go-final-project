package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
)

func newTask() model.TaskModel {
	return model.TaskModel{
		Date:    "20251003",
		Title:   "first",
		Comment: "ololo",
		Repeat:  "d 1",
	}
}

func updateTask() model.TaskModel {
	return model.TaskModel{
		Date:    "30000102",
		Title:   "first, not it's second",
		Comment: "ololo pololo",
		Repeat:  "m 3,2,1 1",
	}
}

var LastID = uint(0)

var dataForQuery = []struct {
	descriptiom string
	init        func(ctx context.Context, s Source, data any) (any, error)
	ctxTimeOut  time.Duration
	data        any
	expectedRes any
	err         error
	msg         string
}{
	{
		descriptiom: `update task - wrong ID`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return nil, s.NewDataTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        updateTask(),
		expectedRes: nil,
		err:         ErrDataBaseNotFound,
		msg:         `invalid update task, error - not dound`,
	},
	{
		descriptiom: `new task - valid`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return s.SaveOneTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        newTask(),
		expectedRes: uint(0),
		err:         nil,
		msg:         `create new task, no error`,
	},
	{
		descriptiom: `find task - valid`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return s.FindOneTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        uint(0),
		expectedRes: newTask(),
		err:         nil,
		msg:         `find a task, no error`,
	},
	{
		descriptiom: `wrong find task`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return s.FindOneTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        uint(1000),
		expectedRes: model.TaskModel{},
		err:         ErrDataBaseNotFound,
		msg:         `invalid find task,error - not found`,
	},
	{
		descriptiom: `update task - valid`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return nil, s.NewDataTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        updateTask(),
		expectedRes: nil,
		err:         nil,
		msg:         `update task, no error`,
	},
	{
		descriptiom: `task list - valid by date`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return s.FindTaskList(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        entity.NewTaskProperty("02.01.3000", 123),
		expectedRes: []model.TaskModel{updateTask()},
		err:         nil,
		msg:         `find array task, no error`,
	},
	{
		descriptiom: `task list - valid by word`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return s.FindTaskList(ctx, data)
		},
		ctxTimeOut: 100 * time.Second,
		// first, not it's second -> first, no(t i)t's second
		data:        entity.NewTaskProperty("t i", 123),
		expectedRes: []model.TaskModel{updateTask()},
		err:         nil,
		msg:         `find array task, no error`,
	},
	{
		descriptiom: `task delete - valid`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return nil, s.ExpirationTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        uint(0),
		expectedRes: nil,
		err:         nil,
		msg:         `delete task, no error`,
	},
	{
		descriptiom: `task delete - invalid ID not exist`,
		init: func(ctx context.Context, s Source, data any) (any, error) {
			return nil, s.ExpirationTask(ctx, data)
		},
		ctxTimeOut:  100 * time.Second,
		data:        uint(0),
		expectedRes: nil,
		err:         ErrDataBaseNotFound,
		msg:         `wrong delete task, error - task not exist`,
	},
}

func TestDataBase(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	cfg, err := config.NewConfig(filepath.Join("..", "..", "init", ".env"))
	requires.NoError(err, fmt.Sprintf("database_test: config error - %v", err))
	cfg.DataBaseDataSourceName = filepath.Join("..", "..", cfg.DataBaseDataSourceName)

	db, err := InitDB(cfg)
	requires.NoError(err, "database_test: DB Open error")
	defer func() {
		err := db.Close()
		asserts.NoError(err, "database_test: DB Close error")
	}()
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `DELETE FROM scheduler;`)
	requires.NoError(err, "delete data from table error")

	source := NewSource(db)

	for i, test := range dataForQuery {
		log.Printf("\t%d %s", i+1, test.descriptiom)

		// create data
		var data any
		switch v := test.data.(type) {
		case uint:
			data = v + LastID
		case model.TaskModel:
			v.ID = LastID
			data = v
		case *entity.TaskProperty:
			data = v
		default:
			requires.FailNow("unexpected data for query " + test.msg)
		}

		ctx, cancel := context.WithTimeout(ctx, test.ctxTimeOut)
		defer cancel()

		res, err := test.init(ctx, source, data)
		asserts.ErrorIs(test.err, err, "erros no equal "+test.msg)

		// check result
		switch v := res.(type) {
		case uint:
			LastID = v
			asserts.Equal(reflect.TypeOf(test.expectedRes), reflect.TypeOf(LastID), "ist not a uint type ", test.msg)
		case model.TaskModel:
			expectedTask, ok := test.expectedRes.(model.TaskModel)
			requires.True(ok, "false!!! "+test.msg)
			if !errors.Is(err, ErrDataBaseNotFound) {
				expectedTask.ID = LastID
			}
			asserts.Equal(expectedTask, v, "compare task is failed "+test.msg)
		case []model.TaskModel:
			expectedArrTask, ok := test.expectedRes.([]model.TaskModel)
			requires.True(ok, "false!!! "+test.msg)
			requires.Len(v, len(expectedArrTask), "arrays of TaskModel not equal "+test.msg)
			if len(expectedArrTask) > 0 {
				//("ORDER BY date ASC")
				sort.Slice(expectedArrTask, func(i, j int) bool {
					return expectedArrTask[i].Date < expectedArrTask[j].Date
				})
				for i := range expectedArrTask {
					expectedArrTask[i].ID = v[i].ID
				}
			}
			asserts.Equal(expectedArrTask, v, "compare task is failed "+test.msg)
		case nil:
			asserts.Nil(test.expectedRes, "res from query is nil "+test.msg)
		default:
			requires.FailNow("incorrect type in switch RES " + test.msg)
		}
	}
}
