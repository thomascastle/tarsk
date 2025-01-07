package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/thomascastle/tarsk/internal/data"
	"github.com/thomascastle/tarsk/internal/messaging"
	"github.com/thomascastle/tarsk/internal/search"
	"github.com/thomascastle/tarsk/internal/structuredlog"
)

type application struct {
	logger  *structuredlog.Logger
	indexer *search.TaskIndexer
}

func main() {
	logger := structuredlog.New(os.Stdout, structuredlog.LevelInfo)

	s_client, e := search.NewClient()
	if e != nil {
		logger.Fatal(e, nil)
	}

	logger.Info("search client created", nil)

	app := &application{
		logger:  logger,
		indexer: search.NewTaskIndexer(s_client),
	}

	e = app.serve()
	if e != nil {
		logger.Fatal(e, nil)
	}
}

func (app *application) serve() error {
	r_client, e := messaging.NewClient()
	if e != nil {
		return e
	}

	pubsub := r_client.PSubscribe(context.Background(), "tasks.*")

	errorShuttingDown := make(chan error)

	go func() {
		signalQuitting := make(chan os.Signal, 1)
		signal.Notify(signalQuitting, syscall.SIGINT, syscall.SIGTERM)
		s := <-signalQuitting // blocks until a signal is received

		app.logger.Info("server shutting down...", map[string]string{"signal": s.String()})

		e := pubsub.Close()
		if e != nil {
			errorShuttingDown <- e
		}

		errorShuttingDown <- nil
	}()

	app.logger.Info("server started", nil)

	_, e = pubsub.Receive(context.Background())
	if e != nil {
		return e
	}

	messagePublished := pubsub.Channel()

	go func() {
		for message := range messagePublished {
			app.logger.Info("message received on: "+message.Channel, nil)

			switch message.Channel {
			case "tasks.event.created", "tasks.event.updated":
				var task data.Task
				if e := json.NewDecoder(strings.NewReader(message.Payload)).Decode(&task); e != nil {
					app.logger.Info("invalid message: "+e.Error(), nil)
					continue
				}
				if e := app.indexer.Index(context.Background(), task); e != nil {
					app.logger.Info("failed to index the task: "+e.Error(), nil)
				}
			case "tasks.event.deleted":
				var id string
				if e := json.NewDecoder(strings.NewReader(message.Payload)).Decode(&id); e != nil {
					app.logger.Info("invalid message: "+e.Error(), nil)
					continue
				}
				if e := app.indexer.Delete(context.Background(), id); e != nil {
					app.logger.Info("failed to delete the task: "+e.Error(), nil)
				}
			}
		}

		app.logger.Info("No more message to consume!", nil)
	}()

	e = <-errorShuttingDown
	if e != nil {
		return e
	}

	app.logger.Info("server stopped", nil)

	return nil
}
