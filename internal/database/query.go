// query - describes requests to database
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrDatabaseNotFound = mark error - if object in base not exist
var ErrDataBaseNotFound = errors.New("resource not found")

// SaveOneTask - Implements the 'model.taskModel' interface - 'TaskCreate
// use -> Transaction(ctx fucn(ctx)error)error look (./transaction.go)
//
// write task to database:
// return unique ID of new Task if no error
func (s Source) SaveOneTask(ctx context.Context, data any) (uint, error) {
	newTask := data.(model.TaskModel)
	createTask := func(ctx context.Context) error {
		err := s.store.Tx.QueryRowContext(ctx, `
INSERT INTO scheduler (date,
                       title,
                       comment,
                       repeat)
VALUES ($1, $2, $3, $4)
RETURNING id;`,
			newTask.Date,    // 1
			newTask.Title,   // 2
			newTask.Comment, // 3 // if empty need write null, but _test_ need ""
			newTask.Repeat,  // 4
		).Scan(&newTask.ID)
		return err
	}
	return newTask.ID, s.store.Transaction(ctx, createTask)
}

// FindOneTask - get ID from 'date' and return solo task if exist
func (s Source) FindOneTask(ctx context.Context, data any) (model.TaskModel, error) {
	taskID := data.(uint)
	row := s.store.DB.QueryRowContext(ctx, `
SELECT *
FROM scheduler
WHERE id = $1
LIMIT 1;`, taskID)
	task, err := scanTask[*sql.Row](row)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.TaskModel{}, ErrDataBaseNotFound
	}
	return task, err
}

func scanTask[T common.ScanSQL](r T) (model.TaskModel, error) {
	var task model.TaskModel
	err := r.Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	return task, err
}

// NewDataTask - update task in database except id and scan id to check existence
// use -> Transaction(ctx fucn(ctx)error)error
func (s Source) NewDataTask(ctx context.Context, data any) error {
	newTask := data.(model.TaskModel)
	updateTask := func(ctx context.Context) error {
		id := uint(0)
		err := s.store.Tx.QueryRowContext(ctx, `
UPDATE scheduler
SET date    = $2,
    title   = $3,
    comment = $4,
    repeat  = $5
WHERE id = $1
RETURNING id;`,
			newTask.ID,      //1
			newTask.Date,    //2
			newTask.Title,   //3
			newTask.Comment, //4
			newTask.Repeat,  //5
		).Scan(&id)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return ErrDataBaseNotFound
		}
		return err
	}
	return s.store.Transaction(ctx, updateTask)
}

// ExpirationTask - remove task by ID with scan id
// use -> Transaction(ctx fucn(ctx)error)error
func (s Source) ExpirationTask(ctx context.Context, data any) error {
	taskID := data.(uint)
	deleteTask := func(ctx context.Context) error {
		id := uint(0)
		err := s.store.Tx.QueryRowContext(ctx, `
DELETE
from scheduler
WHERE id = $1
RETURNING id;`, taskID).Scan(&id)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return ErrDataBaseNotFound
		}
		return err
	}
	return s.store.Transaction(ctx, deleteTask)
}

// FindTaskList - get 'ptr' of type 'services.TaskProperty' from 'data' look (internal/services/taskproperty.go)
//
// add a command line to 'strings.Builder' if we find any characteristic from "TaskProperty" then append 'args'
// args        - pass to sql.QueryContext
// numberOfArg - marks the argument number in the query string
func (s Source) FindTaskList(ctx context.Context, data any) ([]model.TaskModel, error) {
	property := data.(*services.TaskProperty)
	query := strings.Builder{}
	args := make([]any, 0, 2)
	numberOfArg := 1

	query.WriteString("SELECT * FROM scheduler")
	if property.IsWord() {
		query.WriteString("\nWHERE title LIKE $1 OR comment LIKE $1\nORDER BY date ASC")
		args = append(args, fmt.Sprintf(`%%%s%%`, property.PassWord()))
		numberOfArg++
	} else if property.IsDate() {
		query.WriteString("\nWHERE date = $1")
		args = append(args, property.PassDate().UTC().Format(model.DateFormat))
		numberOfArg++
	}
	query.WriteString(fmt.Sprintf("\nLIMIT $%d;", numberOfArg))
	args = append(args, property.PassLimit())

	rows, err := s.store.DB.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	return scanTaskList(rows)
}

func scanTaskList(rows *sql.Rows) ([]model.TaskModel, error) {
	var tasks []model.TaskModel
	for rows.Next() {
		task, err := scanTask[*sql.Rows](rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}
