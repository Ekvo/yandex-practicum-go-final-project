// handler - routes group to http.ServeMux
package transport

import (
	"net/http"
)

type HandlerModel struct{}

func NewHandlerModel() HandlerModel {
	return HandlerModel{}
}

// apiRoutes - create a group for the object(s) 'model.TaskModel','model.LoginModel'
func (h HandlerModel) apiRoutes(sheduler shedulerCase) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signin", Login(sheduler))

	mux.HandleFunc("GET /task", AuthZ(sheduler, TaskRetrieve(sheduler)))
	mux.HandleFunc("POST /task", AuthZ(sheduler, TaskNew(sheduler)))
	mux.HandleFunc("PUT /task", AuthZ(sheduler, TaskChange(sheduler)))
	mux.HandleFunc("DELETE /task", AuthZ(sheduler, TaskRemove(sheduler)))
	mux.HandleFunc("POST /task/done", AuthZ(sheduler, TaskDone(sheduler)))

	mux.HandleFunc("GET /tasks", AuthZ(sheduler, TaskRetriveList(sheduler)))

	mux.HandleFunc("GET /nextdate", TestNextDate)
	return mux
}
