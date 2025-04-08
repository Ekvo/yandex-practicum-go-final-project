// handler - routes group to http.ServeMux
package transport

import (
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
)

type HandlerModel struct{}

func NewHandlerModel() HandlerModel {
	return HandlerModel{}
}

// multiTask - contains all interface for "model.TaskModel"
type multiTask interface {
	model.TaskCreate
	model.TaskFind
}

// TaskModelRoutes - create a group for the object "model.TaskModel"
func (h HandlerModel) TaskModelRoutes(db multiTask) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /nextdate", TestNextDate)
	mux.HandleFunc("POST /task", TaskNew(db))
	mux.HandleFunc("GET /tasks", TaskList(db))
	return mux
}
