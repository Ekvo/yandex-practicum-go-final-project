// handler - routes group to http.ServeMux
package transport

import (
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/autorization"
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

	mux.HandleFunc("POST /signin", Login)

	mux.HandleFunc("GET /task", autorization.Autorization(TaskRetrive(db)))
	mux.HandleFunc("POST /task", autorization.Autorization(TaskNew(db)))
	mux.HandleFunc("PUT /task", autorization.Autorization(TaskChange(db)))
	mux.HandleFunc("DELETE /task", autorization.Autorization(TaskRemove(db)))
	mux.HandleFunc("POST /task/done", autorization.Autorization(TaskDone(db)))

	mux.HandleFunc("GET /tasks", autorization.Autorization(TaskRetriveList(db)))

	mux.HandleFunc("GET /nextdate", TestNextDate)
	return mux
}
