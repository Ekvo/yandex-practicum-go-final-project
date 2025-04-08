// taskencode - rules for encode Task object
package services

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

type TaskListEncode struct {
	Tasks []model.TaskModel
}

// Response - create a TaskResponse list according to given rules
func (te TaskListEncode) Response() []TaskResponse {
	arrTaskresponse := make([]TaskResponse, len(te.Tasks))
	for i, task := range te.Tasks {
		arrTaskresponse[i] = TaskResponse{
			ID:      strconv.Itoa(int(task.ID)),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}
	}
	return arrTaskresponse
}
