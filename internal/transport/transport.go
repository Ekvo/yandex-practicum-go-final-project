// transport - wrapper to '*http.ServeMux' and 'server.Srv'
package transport

import (
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/server"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/usecase"
)

type Transport struct {
	*http.ServeMux
	server.Srv
}

func NewTransport(cfg *config.Config) Transport {
	pathWeb = cfg.PathFilesWeb
	mux := http.NewServeMux()
	return Transport{ServeMux: mux, Srv: server.InitSRV(cfg, mux)}
}

// path to dir: ./web (look from main folder)
// set with help config.Config
var pathWeb = ""

type shedulerCase interface {
	usecase.LoginService
	usecase.AuthService
	usecase.TaskService
}

// Routes - logic of application routes
func (r Transport) Routes(sheduler shedulerCase) {
	muxTask := NewHandlerModel().apiRoutes(sheduler)

	r.Handle("/", http.FileServer(http.Dir(pathWeb)))
	r.Handle("/api/", http.StripPrefix("/api", muxTask))
}

// Start - set all routes and 'ListenAndServe' see (/internal/server/server.go)
func (r Transport) Start(sheduler shedulerCase) error {
	r.Routes(sheduler)
	return r.ListenAndServeAndShut()
}
