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
	model.TaskCreate // C
	model.TaskRead   // R
	model.TaskUpdate // U
	model.TaskDelete // D :)
}

// TaskModelRoutes - create a group for the object "model.TaskModel"
func (h HandlerModel) taskRoutes(db multiTask) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /task", TaskRetrive(db))
	mux.HandleFunc("POST /task", TaskNew(db))
	mux.HandleFunc("PUT /task", TaskChange(db))
	mux.HandleFunc("DELETE /task", TaskRemove(db))
	mux.HandleFunc("POST /task/done", TaskDone(db))

	mux.HandleFunc("GET /tasks", TaskRetriveList(db))

	mux.HandleFunc("GET /nextdate", TestNextDate)
	return mux
}
