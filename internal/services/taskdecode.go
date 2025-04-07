// taskdecode - rules for decode Task object from http.request
package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrServicesFiledEmpty - if the field can't be empty, use this error
var ErrServicesFiledEmpty = errors.New("empty")

// ErrServicesFiledLengthExceeded - field exceeded the maximum range
var ErrServicesFiledLengthExceeded = errors.New("length exceeded")

// TaskValidtor - rules for deserialize object 'TaskModel'
// task - set fields
type TaskDecode struct {
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
// important: set td.task.Repeat before call member of TaskModel 'SetDate' because 'repeat' use in member
func (td *TaskDecode) Decode(r *http.Request, algoNextDate func(time.Time, string, string) (string, error)) error {
	if err := common.DecodeJSON(r, td); err != nil {
		return err
	}
	msgErr := make(common.Message)
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
	// td.Date == "" - use in member of TaskModel 'SetDate'
	if _, err := time.Parse(model.DateFormat, td.Date); err != nil && td.Date != "" {
		msgErr["date"] = ErrServicesInvalidDate.Error()
	}
	if len(msgErr) != 0 {
		return fmt.Errorf("task decode error - %v", msgErr)
	}
	td.task.Title = td.Title
	td.task.Comment = td.Comment
	td.task.Repeat = td.Repeat
	if err := td.task.SetDate(td.Date, algoNextDate); err != nil {
		return err
	}
	return nil
}
