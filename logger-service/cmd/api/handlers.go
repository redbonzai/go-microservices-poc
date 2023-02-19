package main

import (
	"fmt"
	"log-Service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLogHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var requestPayload JSONPayload
	_ = app.readJson(responseWriter, request, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		fmt.Printf("Error inserting log entry: %s", err)
		app.errorJson(responseWriter, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Log entry created",
	}

	fmt.Printf("Log entry created: %v", response)

	app.writeJson(responseWriter, http.StatusAccepted, response)
}
