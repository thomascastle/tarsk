package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/thomascastle/tarsk/internal/data"
	"github.com/thomascastle/tarsk/internal/validator"
)

func (app *application) listTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, e := app.repositories.Tasks.Select()
	if e != nil {
		app.serverErrorResponse(w, r, e)
		return
	}

	e = app.writeJSON(w, http.StatusOK, envelope{"tasks": tasks}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Description string        `json:"description"`
		DueAt       time.Time     `json:"due_at"`
		Priority    data.Priority `json:"priority"`
		StartedAt   time.Time     `json:"started_at"`
	}

	e := app.readJSON(w, r, &input)
	if e != nil {
		app.badRequestResponse(w, r, e)
		return
	}

	task := &data.Task{
		Description: input.Description,
		DueAt:       input.DueAt,
		Priority:    prioritize(input.Priority),
		StartedAt:   input.StartedAt,
	}

	v := validator.New()
	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	e = app.repositories.Tasks.Insert(task)
	if e != nil {
		app.serverErrorResponse(w, r, e)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%s", task.ID))

	e = app.writeJSON(w, http.StatusCreated, envelope{"task": task}, headers)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := routeParam(r, "id")

	task, e := app.repositories.Tasks.SelectOne(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrorRecordNotFound):
			app.resourceNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
		return
	}

	e = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := routeParam(r, "id")

	task, e := app.repositories.Tasks.SelectOne(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrorRecordNotFound):
			app.resourceNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
		return
	}

	var input struct {
		Description *string        `json:"description"`
		Done        *bool          `json:"done"`
		DueAt       time.Time      `json:"due_at"`
		Priority    *data.Priority `json:"priority"`
		StartedAt   time.Time      `json:"started_at"`
	}

	e = app.readJSON(w, r, &input)
	if e != nil {
		app.badRequestResponse(w, r, e)
		return
	}

	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Done != nil {
		task.Done = *input.Done
	}
	if !input.DueAt.IsZero() {
		task.DueAt = input.DueAt
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if !input.StartedAt.IsZero() {
		task.StartedAt = input.StartedAt
	}

	v := validator.New()
	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	e = app.repositories.Tasks.Update(task)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrorEditConflict):
			app.editConflictResponse(w, r, e)
		default:
			app.serverErrorResponse(w, r, e)
		}
		return
	}

	e = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := routeParam(r, "id")

	e := app.repositories.Tasks.Delete(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrorRecordNotFound):
			app.resourceNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
		return
	}

	e = app.writeJSON(w, http.StatusOK, envelope{"message": "The task has been deleted successfully."}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}
