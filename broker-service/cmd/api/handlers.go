package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string        `json:"action"`
	Auth   AuthPayload   `json:"auth,omitempty"`
	Log    LoggerPayload `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggerPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) BrokerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Hit the broker")
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(responseWriter, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(responseWriter http.ResponseWriter, request *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(responseWriter, request, &requestPayload)
	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(responseWriter, requestPayload.Auth)
	case "log":
		app.logItem(responseWriter, requestPayload.Log)
	default:
		app.errorJson(responseWriter, errors.New("unknown action"))
	}
}

func (app *Config) logItem(responseWriter http.ResponseWriter, logPayload LoggerPayload) {
	// insert data
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(responseWriter, errors.New("error calling log service -- status no StatusAccepted"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged!"

	app.writeJson(responseWriter, http.StatusAccepted, payload)
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(responseWriter http.ResponseWriter, authPayload AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	// call the service
	request, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(responseWriter, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(responseWriter, errors.New("error calling auth service"))
		return
	}

	// create authPayload variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(responseWriter, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(responseWriter, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(responseWriter, http.StatusAccepted, payload)
}
