// server - wrapper to &http.Server
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// ServerTimeoutShut - use in 'ListenAndServeAndShut' look down
const serverTimeoutShut = 10 * time.Second

type Srv struct {
	*http.Server
}

func NewSrvWihtHTTPServer(server *http.Server) Srv {
	return Srv{Server: server}
}

func InitSRV(r http.Handler) Srv {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(8000)
	}
	return NewSrvWihtHTTPServer(&http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: r,
	})
}

// ListenAndServeAndShut - application wait SIGTERM or SYGINT - signal to inform that it is time to shutdovn
func (s Srv) ListenAndServeAndShut() error {
	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: HTTP server error - %v", err)
		}
		log.Print("server: stopped serving\n")
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), serverTimeoutShut)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("server: HTTP shutdown error - %v", err)
	}
	log.Print("server: shutdown complete\n")
	return nil
}
