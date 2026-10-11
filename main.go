package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	route "github.com/napster-mayu/service-inspector/internal/http"
	"github.com/napster-mayu/service-inspector/internal/service"
)

func main() {
	serverPort := ":8080"

	svc := service.NewService()

	r := route.GetRoutes(svc)

	srv := &http.Server{
		Addr:              serverPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Cancelled on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Println("listening on ", serverPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		stop() // a second Ctrl+C now force-kills
		log.Println("shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}

	// Close DB, flush logs, stop workers, etc.
	log.Println("bye")
}
