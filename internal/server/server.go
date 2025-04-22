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
	"syscall"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
)

// ServerTimeoutShut - use in 'ListenAndServeAndShut' look down
const serverTimeoutShut = 10 * time.Second

type Srv struct {
	*http.Server
}

func InitSRV(cfg *config.Config, router http.Handler) Srv {
	return Srv{
		Server: &http.Server{
			Addr:    net.JoinHostPort("", cfg.ServerPort),
			Handler: router,
		},
	}
}

// ListenAndServeAndShut - application wait SIGTERM or SYGINT - signal to inform that it is time to shutdovn
func (s Srv) ListenAndServeAndShut() error {
	go func() {
		log.Print("server: listen and serve - start\n")
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
		return fmt.Errorf("server: HTTP shutdown error - %w", err)
	}
	log.Print("server: shutdown complete\n")
	return nil
}
