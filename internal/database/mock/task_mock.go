package mock

import (
	"context"
	"sort"
	"strings"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
)

type MockTaskStore struct {
	id    *uint
	tasks map[uint]model.TaskModel
}

func NewMockTaskStore() MockTaskStore {
	return MockTaskStore{
		id:    new(uint),
		tasks: make(map[uint]model.TaskModel),
	}
}

func (s MockTaskStore) incrementID() {
	*s.id++
}

func (s MockTaskStore) SaveOneTask(_ context.Context, data any) (uint, error) {
	newTask := data.(model.TaskModel)
	s.incrementID()
	id := *s.id
	newTask.ID = id
	s.tasks[id] = newTask
	return id, nil
}

func (s MockTaskStore) FindOneTask(_ context.Context, data any) (model.TaskModel, error) {
	id := data.(uint)
	task, ex := s.tasks[id]
	if !ex {
		return model.TaskModel{}, database.ErrDataBaseNotFound
	}
	return task, nil
}

func (s MockTaskStore) NewDataTask(_ context.Context, data any) error {
	updateTask := data.(model.TaskModel)
	id := updateTask.ID
	if _, ex := s.tasks[id]; !ex {
		return database.ErrDataBaseNotFound
	}
	s.tasks[id] = updateTask
	return nil
}

func (s MockTaskStore) ExpirationTask(ctx context.Context, data any) error {
	taskID := data.(uint)
	if _, ex := s.tasks[taskID]; !ex {
		return database.ErrDataBaseNotFound
	}
	delete(s.tasks, taskID)
	return nil
}

func (s MockTaskStore) FindTaskList(_ context.Context, data any) ([]model.TaskModel, error) {
	property := data.(*entity.TaskProperty)
	var arrOfTask []model.TaskModel

	word := property.PassWord()
	date := property.PassDate().UTC().Format(model.DateFormat)
	for _, task := range s.tasks {
		if property.IsWord() {
			if strings.Contains(task.Title, word) || strings.Contains(task.Comment, word) {
				arrOfTask = append(arrOfTask, task)
			}
		} else if property.IsDate() {
			if task.Date == date {
				arrOfTask = append(arrOfTask, task)
			}
		} else {
			arrOfTask = append(arrOfTask, task)
		}
	}
	limit := property.PassLimit()
	if len(arrOfTask) > int(limit) {
		return arrOfTask[:limit], nil
	}
	sort.Slice(arrOfTask, func(i, j int) bool {
		return arrOfTask[i].Date < arrOfTask[j].Date
	})
	return arrOfTask, nil
}
