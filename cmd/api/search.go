package main

import (
	"net/http"

	"github.com/thomascastle/tarsk/internal/data"
)

func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Description *string        `json:"description"`
		Done        *bool          `json:"done"`
		From        int64          `json:"from"`
		Priority    *data.Priority `json:"priority"`
		Size        int64          `json:"size"`
	}

	e := app.readJSON(w, r, &input)
	if e != nil {
		app.badRequestResponse(w, r, e)
		return
	}

	results, e := app.search.Query(
		r.Context(),
		data.SearchParams{
			Description: input.Description,
			Done:        input.Done,
			From:        input.From,
			Priority:    input.Priority,
			Size:        input.Size,
		},
	)
	if e != nil {
		app.serverErrorResponse(w, r, e)
		return
	}

	e = app.writeJSON(
		w,
		http.StatusOK,
		envelope{"tasks": results.Tasks, "total": results.Total},
		nil,
	)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}
