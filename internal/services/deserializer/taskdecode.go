// taskdecode - rules for decode Task object from http.request
package deserializer

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrServicesFiledEmpty - if the field can't be empty, use this error
var ErrServicesFiledEmpty = errors.New("empty")

// ErrServicesFiledLengthExceeded - field exceeded the maximum range
var ErrServicesFiledLengthExceeded = errors.New("length exceeded")

// ErrServicesWrongID - ID consists of more than just numbers
var ErrServicesWrongID = errors.New("not numeric")

// TaskValidtor - rules for deserialize object 'TaskModel'
// task - set fields
type TaskDecode struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`

	task model.TaskModel `json:"-"`
}

// NewTaskDecode - if need add default property
func NewTaskDecode() *TaskDecode {
	return &TaskDecode{}
}

// Model - return task
// called without a pointer for save 'task'
func (td TaskDecode) Model() model.TaskModel {
	return td.task
}

// Decode - deserialize object TaskModel from Request
//
// for check all property 'TaskModel' and create full error list if exist bad data
// use map - common.Message look (../../pkg/common/common.go)
//
// for set td.task.Date use func - algoNextDate
// important: set td.task.Repeat before call member of TaskDecode 'setDate' because 'repeat' use in it
func (td *TaskDecode) Decode(r *http.Request, algoNextDate func(time.Time, string, string) (string, error)) error {
	if err := common.DecodeJSON(r, td); err != nil {
		return err
	}
	msgErr := make(common.Message)
	taskID := uint(0)
	if idSTR := td.ID; idSTR != "" {
		if id, err := strconv.Atoi(idSTR); err != nil {
			msgErr["id"] = ErrServicesWrongID.Error()
		} else {
			taskID = uint(id)
		}
	}
	titleLen := len(td.Title)
	if titleLen < 1 {
		msgErr["title"] = ErrServicesFiledEmpty.Error()
	}
	if titleLen > model.TaskTitleLen {
		msgErr["title"] = ErrServicesFiledLengthExceeded.Error()
	}
	if len(td.Comment) > model.TaskCommentLen {
		msgErr["comment"] = ErrServicesFiledLengthExceeded.Error()
	}
	if len(td.Repeat) > model.TaskRepeatLen {
		msgErr["repeat"] = ErrServicesFiledLengthExceeded.Error()
	}
	taskDate, err := td.executeDate(algoNextDate)
	if err != nil {
		msgErr["date"] = err.Error()
	}
	if len(msgErr) != 0 {
		return fmt.Errorf("task decode error - %v", msgErr)
	}
	td.task.ID = taskID
	td.task.Date = taskDate
	td.task.Title = td.Title
	td.task.Comment = td.Comment
	td.task.Repeat = td.Repeat
	return nil
}

// executeDate - metod of TaskDecode find executeble date
//
// 1 or 2 or 2.1 or 2.2 or 3
//
// 1. date is not specified, today's date is taken
// 2. date is less than today's date
// 2.1. repeat == "" date = now
// 2.2. find date with use 'nextDate'
// 3. return td.Date
//
// 'nextDate' - selected algorithm - execute if 'date' less 'now' and 't.Repeat' not empty
func (td *TaskDecode) executeDate(nextDate func(time.Time, string, string) (string, error)) (string, error) {
	if nextDate == nil {
		return "", model.ErrModelAlgorithmNextDateIsNULL
	}
	now := common.ReduceTimeToDay(time.Now().UTC())
	date := td.Date
	if date == "" {
		date = now.Format(model.DateFormat)
	}
	dateToTime, err := time.Parse(model.DateFormat, date)
	if err != nil {
		return "", services.ErrServicesInvalidDate
	}
	if dateToTime.UTC().Before(now.UTC()) {
		repeat := td.Repeat
		if repeat == "" {
			return now.Format(model.DateFormat), nil
		}
		if date, err = nextDate(now, date, repeat); err != nil {
			return "", err
		}
		return date, nil
	}
	return date, nil
}
