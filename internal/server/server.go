// server - wrapper to &gttp.Server
package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ServerTimeoutShut = 10 * time.Second

type Srv struct {
	*http.Server
}

func NewSrvWihtHTTPServer(server *http.Server) Srv {
	return Srv{server}
}

func InitSRV(r http.Handler) Srv {
	return NewSrvWihtHTTPServer(&http.Server{
		Addr:    ":" + os.Getenv("TODO_PORT"),
		Handler: r,
	})
}

func (s Srv) ListenAndServeAndShut(timeShut time.Duration) {
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: HTTP server error: %v", err)
		}
		log.Print("server: stopped serving\n")
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeShut)
	defer shutdownRelease()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server: HTTP shutdown error: %v", err)
	}
	log.Print("server: shutdown complete\n")
}
