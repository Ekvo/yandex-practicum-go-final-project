// task - describes the Task object and its implementing interfaces
package model

import (
	"context"
	"errors"
)

var ErrModelAlgorithmNextDateIsNULL = errors.New("algorithm not selected")

// ErrModelTaskDone - mark task - completed
var ErrModelTaskDone = errors.New("task done")

// max lenght to 'TaskModel' Fields
const (
	TaskTitleLen   = 255
	TaskCommentLen = 2048
	TaskRepeatLen  = 128
)

// in this format database store 'date' in type VARCHAR(8)
const DateFormat = "20060102"

type TaskModel struct {
	ID uint

	// task date in format '20060102'
	// max 8 characters
	Date string

	// not empty
	// max 255 characters
	Title string

	// max 2048 characters
	Comment string

	// containing the repetition rules for the task.
	// max 128 characters,
	Repeat string
}

// TaskCreate - save a task to storage, and return a unique ID for the new task
type TaskCreate interface {
	SaveOneTask(ctx context.Context, data any) (uint, error)
}

// TaskRead - read task from store
type TaskRead interface {
	FindOneTask(ctx context.Context, data any) (TaskModel, error)
	FindTaskList(ctx context.Context, data any) ([]TaskModel, error)
}

// TaskUpdate - write new data for a specific task
type TaskUpdate interface {
	NewDataTask(ctx context.Context, data any) error
}

// TaskDelete - delete task from store, if not exist -> error
type TaskDelete interface {
	ExpirationTask(ctx context.Context, data any) error
}
