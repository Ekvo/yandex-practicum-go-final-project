// taskdecode - rules for decode Task object from http.request
package deserializer

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

var (
	// ErrServicesFiledEmpty - if the field can't be empty, use this error
	ErrServicesFiledEmpty = errors.New("empty")

	// ErrServicesFiledLengthExceeded - field exceeded the maximum range
	ErrServicesFiledLengthExceeded = errors.New("length exceeded")

	// ErrServicesWrongID - ID consists of more than just numbers
	ErrServicesWrongID = errors.New("not numeric")

	// ErrBizInvalidDate - wrong format of date
	ErrServicesInvalidDate = errors.New("invalid date format")
)

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
func (td *TaskDecode) Decode(r *http.Request) error {
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
	date := td.Date
	if date != "" {
		if _, err := time.Parse(model.DateFormat, date); err != nil {
			msgErr["date"] = ErrServicesInvalidDate.Error()
		}
	}
	if len(msgErr) != 0 {
		return fmt.Errorf("task decode error - %s", msgErr.String())
	}
	td.task.ID = taskID
	td.task.Date = date
	td.task.Title = td.Title
	td.task.Comment = td.Comment
	td.task.Repeat = td.Repeat
	return nil
}
