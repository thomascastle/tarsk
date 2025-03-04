package main

import (
	"net/http"
)

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, e error) {
	app.errorResponse(w, r, http.StatusBadRequest, e.Error())
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request, e error) {
	app.errorResponse(w, r, http.StatusBadRequest, e.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	e := app.writeJSON(
		w,
		status,
		envelope{"error": message},
		nil,
	)
	if e != nil {
		app.logError(r, e)
		w.WriteHeader(500)
	}
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) resourceNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) logError(r *http.Request, e error) {
	app.logger.Error(e, map[string]string{"request_method": r.Method, "request_url": r.URL.String()})
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, e error) {
	app.logError(r, e)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
