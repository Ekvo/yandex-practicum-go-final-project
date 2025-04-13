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

// TaskModelRoutes - create a group for the object(s) 'model.TaskModel','model.LoginModel'
func (h HandlerModel) apiRoutes(db multiTask) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signin", Login)

	mux.HandleFunc("GET /task", autorization.AuthZ(TaskRetrive(db)))
	mux.HandleFunc("POST /task", autorization.AuthZ(TaskNew(db)))
	mux.HandleFunc("PUT /task", autorization.AuthZ(TaskChange(db)))
	mux.HandleFunc("DELETE /task", autorization.AuthZ(TaskRemove(db)))
	mux.HandleFunc("POST /task/done", autorization.AuthZ(TaskDone(db)))

	mux.HandleFunc("GET /tasks", autorization.AuthZ(TaskRetriveList(db)))

	mux.HandleFunc("GET /nextdate", TestNextDate)
	return mux
}
