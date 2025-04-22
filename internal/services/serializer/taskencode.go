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
func (te TaskEncode) Response() *TaskResponse {
	taskResponse := TaskResponse{
		ID:      strconv.FormatUint(uint64(te.ID), 10),
		Date:    te.Date,
		Title:   te.Title,
		Comment: te.Comment,
		Repeat:  te.Repeat,
	}
	return &taskResponse
}

// TaskIDResponse - properties for ID of task to write to http.ResponseWriter
type TaskIDResponse struct {
	ID string `json:"id"`
}

type TaskIDEncode struct {
	ID uint
}

// create a task ID response
func (tIDe TaskIDEncode) Response() *TaskIDResponse {
	taskIDResponse := TaskIDResponse{
		ID: strconv.FormatUint(uint64(tIDe.ID), 10),
	}
	return &taskIDResponse
}

type TaslListResponse struct {
	TasksResp []TaskResponse `json:"tasks"`
}

type TaskListEncode struct {
	Tasks []model.TaskModel
}

// create a 'taskResponse' list
func (tle TaskListEncode) Response() *TaslListResponse {
	arrTaskResponse := make([]TaskResponse, 0, len(tle.Tasks))
	for _, task := range tle.Tasks {
		arrTaskResponse = append(arrTaskResponse, *TaskEncode{task}.Response())
	}
	return &TaslListResponse{TasksResp: arrTaskResponse}
}
