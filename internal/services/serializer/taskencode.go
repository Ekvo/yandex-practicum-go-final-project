// taskencode - rules for encode Task object
package serializer

import (
	"strconv"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
)

// TaskResponse - task properties for writing to http.ResponseWriter
type TaskResponse struct {
	ID      string `json:"id"` // need "-" in my opinion
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"` // need omitempty
	Repeat  string `json:"repeat"`  // need omitempty
}

type TaskEncode struct {
	model.TaskModel
}

// create a TaskResponse according to given rules
func (te TaskEncode) Response() TaskResponse {
	taskresponse := TaskResponse{
		ID:      strconv.Itoa(int(te.ID)),
		Date:    te.Date,
		Title:   te.Title,
		Comment: te.Comment,
		Repeat:  te.Repeat,
	}
	return taskresponse
}

type TaskListEncode struct {
	Tasks []model.TaskModel
}

// create a TaskResponse list
func (tle TaskListEncode) Response() []TaskResponse {
	arrTaskresponse := make([]TaskResponse, len(tle.Tasks))
	for i, task := range tle.Tasks {
		arrTaskresponse[i] = TaskEncode{task}.Response()
	}
	return arrTaskresponse
}
