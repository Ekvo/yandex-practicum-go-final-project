// query - describes requests to database
package database

import (
	"context"
	"errors"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
)

// ErrDatabseAlreadyExist - if unique columns already exist in the storage
var ErrDatabseAlreadyExist = errors.New("resource already exists")

func (s Source) NewTables(ctx context.Context, tables ...string) error {
	newTables := func(ctx context.Context) error {
		for _, table := range tables {
			_, err := s.store.Tx.ExecContext(ctx, table)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return s.store.Transaction(ctx, newTables)
}

// SaveOneTask - Implements the 'model.taskModel' interface - 'TaskCreate
// use -> Transaction(ctx fucn(ctx)error)error look (./transaction.go)
//
// write task to database:
// check on null 'comment' if comment="" write null
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
			newTask.Date,                             // 1
			newTask.Title,                            // 2
			WhenEmptyStringThenNULL(newTask.Comment), // 3
			newTask.Repeat,                           // 4
		).Scan(&newTask.ID)
		if err != nil {
			return ErrDatabseAlreadyExist
		}
		return nil
	}
	return newTask.ID, s.store.Transaction(ctx, createTask)
}

// WhenEmptyStringThenNULL - if field VARCHAR or TEXT in database can be NULL
func WhenEmptyStringThenNULL(s string) *string {
	if len(s) != 0 {
		return &s
	}
	return nil
}
