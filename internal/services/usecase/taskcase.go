// taskcase - biz logic of create, read, update, delete 'model.TaskModel'
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/nextdate"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/serializer"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

var (
	// ErrCaseTaskAlreadyExist - if add Task with ID
	ErrCaseTaskAlreadyExist = errors.New("task already exist")

	ErrCaseTaskNotFound = errors.New("task not found")

	// ErrCaseTaskZeroID - for read, update, delete and 'done' task
	ErrCaseTaskZeroID = errors.New("task ID is zero")

	// ErrCaselAlgorithmNextDateIsNULL - nextDate is nil in 'TaskService'
	ErrCaselAlgorithmNextDateIsNULL = errors.New("algorithm not selected")
)

// contain all business logic of task
type TaskService interface {
	services.TaskCreateCase
	services.TaskReadCase
	services.TaskUpdateCase
	services.TaskDeleteCase
	services.TaskDoneCase
}

// multiTask - contain all TaskModel interfaces
type MultiTask interface {
	model.TaskCreate
	model.TaskRead
	model.TaskUpdate
	model.TaskDelete
}

type taskService struct {
	// taskRepo - work with store
	taskRepo MultiTask

	// nextDate - algorithm for create next date to task
	nextDate nextdate.NextDateFunc
}

func NewTaskService(cfg *config.Config, store MultiTask) (TaskService, error) {
	nextDate, err := setNextDate(cfg.TaskNextDate)
	if err != nil {
		return nil, err
	}
	return taskService{
		taskRepo: store,
		nextDate: nextDate,
	}, nil
}

// return - algotithm by name
func setNextDate(name string) (nextdate.NextDateFunc, error) {
	if name == "nextdate" {
		return nextdate.NextDate, nil
	}
	return nil, ErrCaselAlgorithmNextDateIsNULL
}

// CreateTask - member of taskService
//
// 1. if create with ID -> check in database 'FindOneTask' -> ID exist -> error
// 2. find execute date see below 'executeDate(date, repeat string) (string, error)'
// 3. add in database task and get ID
// 4. return TaskIDResponse
func (ts taskService) CreateTask(
	ctx context.Context,
	task model.TaskModel) (*serializer.TaskIDResponse, error) {
	if id := task.ID; id != 0 {
		_, err := ts.taskRepo.FindOneTask(ctx, id)
		if err == nil {
			return nil, ErrCaseTaskAlreadyExist
		}
		if !errors.Is(err, database.ErrDataBaseNotFound) {
			return nil, services.ErrServicesInternalError
		}
	}
	date, err := ts.executeDate(task.Date, task.Repeat)
	if err != nil {
		if errors.Is(err, nextdate.ErrNextDateInvalidDate) ||
			errors.Is(err, nextdate.ErrNextDateWrongRepeat) {
			return nil, err
		}
		return nil, services.ErrServicesInternalError
	}
	task.Date = date
	id, err := ts.taskRepo.SaveOneTask(ctx, task)
	if err != nil {
		return nil, services.ErrServicesInternalError
	}
	serizlize := serializer.TaskIDEncode{ID: id}
	return serizlize.Response(), nil
}

// executeDate - metod of taskService find executeble date where task create or update
//
// 1 or 2.1 or 2.2 or 3
//
// 1. 'date' is not specified, 'now' is taken
// 2. 'date' is less than 'now'
// 2.1. repeat == "" date = now
// 2.2. find date with use 'nextDate'
// 3. return td.Date
//
// 'nextDate' - selected algorithm - execute if 'date' less 'now' and 't.Repeat' not empty
func (ts taskService) executeDate(date, repeat string) (string, error) {
	now := common.ReduceTimeToDay(time.Now())
	if date == "" {
		date = now.Format(model.DateFormat)
	}
	dateToTime, err := time.Parse(model.DateFormat, date)
	if err != nil {
		return "", nextdate.ErrNextDateInvalidDate
	}
	dateAfterNow := false
	if dateToTime.UTC().Before(now.UTC()) {
		dateAfterNow = true
	}
	if repeat == "" {
		if dateAfterNow {
			return now.Format(model.DateFormat), nil
		}
		return date, nil
	}
	nextDate, err := ts.nextDate(now, date, repeat)
	if err != nil {
		return "", err
	}
	if dateAfterNow {
		date = nextDate
	}
	return date, nil
}

// ReadTask - member of taskService
//
// 1. check ID by zero
// 2. find task by ID
// 3. create TaskResponse
func (ts taskService) ReadTask(
	ctx context.Context,
	id uint) (*serializer.TaskResponse, error) {
	if id == 0 {
		return nil, ErrCaseTaskZeroID
	}
	task, err := ts.taskRepo.FindOneTask(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDataBaseNotFound) {
			return nil, ErrCaseTaskNotFound
		}
		return nil, services.ErrServicesInternalError
	}
	serialize := serializer.TaskEncode{TaskModel: task}
	return serialize.Response(), nil
}

// UpdateTask - member of taskService
//
// 1. check ID by zero
// 2. find execute date use - 'executeDate'
// 3. update task by ID in database
func (ts taskService) UpdateTask(ctx context.Context, task model.TaskModel) error {
	id := task.ID
	if id == 0 {
		return ErrCaseTaskZeroID
	}
	date, err := ts.executeDate(task.Date, task.Repeat)
	if err != nil {
		if errors.Is(err, nextdate.ErrNextDateInvalidDate) ||
			errors.Is(err, nextdate.ErrNextDateWrongRepeat) {
			return err
		}
		return services.ErrServicesInternalError
	}
	task.Date = date
	if err := ts.taskRepo.NewDataTask(ctx, task); err != nil {
		if errors.Is(err, database.ErrDataBaseNotFound) {
			return ErrCaseTaskNotFound
		}
		return services.ErrServicesInternalError
	}
	return nil
}

// DeleteTask - member of taskService
//
// 1. check ID by zero
// 2. delete task by ID
func (ts taskService) DeleteTask(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrCaseTaskZeroID
	}
	if err := ts.taskRepo.ExpirationTask(ctx, id); err != nil {
		if errors.Is(err, database.ErrDataBaseNotFound) {
			return ErrCaseTaskNotFound
		}
		return services.ErrServicesInternalError
	}
	return nil
}

// DoneTask - member of taskService
//
// 1. check ID by zero
// 2. find task by ID
// 3. processing the task
//
//	3.1 find execute date see bellow 'updateDateAfterDone(date, repeat string) (string, error)'
//
// 3.2.1 task done -> delete task from database by ID
// 3.2.2 update task by ID in database
func (ts taskService) DoneTask(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrCaseTaskZeroID
	}
	task, err := ts.taskRepo.FindOneTask(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDataBaseNotFound) {
			return ErrCaseTaskNotFound
		}
		return services.ErrServicesInternalError
	}
	date, err := ts.updateDateAfterDone(task.Date, task.Repeat)
	if err != nil {
		if errors.Is(err, model.ErrModelTaskDone) {
			if err := ts.taskRepo.ExpirationTask(ctx, id); err != nil {
				return services.ErrServicesInternalError
			}
			return nil
		}
		return services.ErrServicesInternalError
	}
	task.Date = date
	if err := ts.taskRepo.NewDataTask(ctx, task); err != nil {
		return services.ErrServicesInternalError
	}
	return nil
}

//	updateDateAfterDone - metod of taskService used only in 'DoneTask'
//
// rules:
// 1. repeat - empty -> task done -> delete
// 2. update the date using 'nextDate' algorithm
// 2.1 if now Before t.Date -> now = t.Date
func (ts taskService) updateDateAfterDone(date, repeat string) (string, error) {
	if repeat == "" {
		return "", model.ErrModelTaskDone
	}
	oldDate, err := time.Parse(model.DateFormat, date)
	if err != nil {
		return "", err
	}
	now := common.ReduceTimeToDay(time.Now())
	if now.UTC().Before(oldDate.UTC()) {
		now = oldDate
	}
	return ts.nextDate(now, date, repeat)
}

// ReadTaskList - member of taskService
//
// 1. find task list by 'entity.TaskProperty' look (/internal/services/entity/taskproperty.go)
// 2. create TaslListResponse
func (ts taskService) ReadTaskList(
	ctx context.Context,
	property *entity.TaskProperty) (*serializer.TaslListResponse, error) {
	tasks, err := ts.taskRepo.FindTaskList(ctx, property)
	if err != nil {
		return nil, services.ErrServicesInternalError
	}
	serialize := serializer.TaskListEncode{Tasks: tasks}
	return serialize.Response(), nil
}
