package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	e := decoder.Decode(dst)
	if e != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(e, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(e, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(e, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(e, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(e.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(e.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case e.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(e, &invalidUnmarshalError):
			panic(e)
		default:
			return e
		}
	}

	e = decoder.Decode(&struct{}{})
	if e != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	data_JSON, e := json.MarshalIndent(data, "", "\t")
	if e != nil {
		return e
	}

	data_JSON = append(data_JSON, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data_JSON)

	return nil
}

func prioritize(priority string) string {
	if priority == "" {
		return "none"
	}

	return priority
}

func routeParam(r *http.Request, name string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(name)
}
