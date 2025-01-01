package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	errorShuttingDown := make(chan error)

	go func() {
		signalQuitting := make(chan os.Signal, 1)
		signal.Notify(signalQuitting, syscall.SIGINT, syscall.SIGTERM)
		s := <-signalQuitting // blocks until a signal is received

		log.Println("The server is gracefully shutting down...", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		e := server.Shutdown(ctx)
		if e != nil {
			errorShuttingDown <- e
		}

		errorShuttingDown <- nil
	}()

	log.Printf("The server started at :%d", app.config.port)

	e := server.ListenAndServe()
	if !errors.Is(e, http.ErrServerClosed) {
		return e
	}

	e = <-errorShuttingDown
	if e != nil {
		return e
	}

	log.Println("The server stopped")

	return nil
}
