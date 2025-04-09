// transport - wrapper to '*http.ServeMux' and 'server.Srv'
package transport

import (
	"net/http"
	"os"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/server"
)

type Transport struct {
	*http.ServeMux
	server.Srv
}

func NewTransport(mux *http.ServeMux) Transport {
	return Transport{ServeMux: mux, Srv: server.InitSRV(mux)}
}

// Routes - logic of application routes
func (t Transport) Routes(db multiTask) {
	mux := t.ServeMux
	muxTask := NewHandlerModel().taskRoutes(db)

	mux.Handle("/", http.FileServer(http.Dir(os.Getenv("TODO_WEB"))))
	mux.Handle("/api/", http.StripPrefix("/api", muxTask))
}

func (t Transport) Run(db multiTask) error {
	t.Routes(db)
	return t.ListenAndServeAndShut()
}
