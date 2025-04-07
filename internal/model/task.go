// task - describe the Task object and its implementing interfaces
package modelmodel

import (
	"context"
	"time"
)

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
	// no more than 8 characters
	Date string

	// no more than 255 characters
	Title string

	// no more than 2048 characters
	Comment string

	// containing the repetition rules for the task.
	// no more than 128 characters,
	Repeat string
}

// TaskCreate - save a task to storage, and return a unique ID for the new task
type TaskCreate interface {
	SaveOneTask(ctx context.Context, data any) (uint, error)
}

// SetDate - metod of TaskModel find new executeble date to Task
//
// algoNewDate - selected algorithm - execute if 'date' less 'now' and 't.Repeat' not empty
func (t *TaskModel) SetDate(date string, algoNextDate func(time.Time, string, string) (string, error)) error {
	now := time.Now().UTC()
	if date == "" {
		t.Date = now.Format(DateFormat)
		return nil
	}
	dateToTime, err := time.Parse(DateFormat, date)
	if err != nil {
		return err
	}
	if dateToTime.UTC().Before(now.UTC()) {
		if t.Repeat == "" || algoNextDate == nil {
			t.Date = now.Format(DateFormat)
			return nil
		}
		if t.Date, err = algoNextDate(now, date, t.Repeat); err != nil {
			return err
		}
		return nil
	}
	t.Date = date
	return nil
}

// TaskComplete - check before remove task from the base
//
// no 'repeat' and 'date' before or equal time.Now() -> task must be deleted from base
func (t *TaskModel) TaskComplete() (bool, error) {
	timeDate, err := time.Parse(DateFormat, t.Date)
	if err != nil {
		return false, err
	}
	return t.Repeat == "" && !timeDate.UTC().After(time.Now().UTC()), nil
}
