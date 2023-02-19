package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/**
 * @api {get} /broker Get broker status
 */
func (app *Config) readJson(responseWriter http.ResponseWriter, request *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte
	request.Body = http.MaxBytesReader(responseWriter, request.Body, int64(maxBytes))
	dec := json.NewDecoder(request.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body myst have only a single JSON value")
	}

	return nil
}

/**
 * @api {post} /broker/:id/:action/:value Set broker status
 */
func (app *Config) writeJson(responseWriter http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			responseWriter.Header()[key] = value
		}
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.WriteHeader(status)
	_, err = responseWriter.Write(out)
	if err != nil {
		return err
	}

	return nil
}

/**
 * @api {get} /broker Get broker error status
 */
func (app *Config) errorJson(responseWriter http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(responseWriter, statusCode, payload)
}
