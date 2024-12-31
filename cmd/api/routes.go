package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/tasks", app.listTasksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.showTaskHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.deleteTaskHandler)

	return router
}
