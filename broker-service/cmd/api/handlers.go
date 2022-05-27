package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse {
		Error: false,
		Message: "Broker service is running",
	}

	_ = app.writeJson(w, http.StatusOK, payload)

	//out, _ := json.MarshalIndent(payload, "", "\t")
	//w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.WriteHeader(http.StatusAccepted)
	//w.Write(out)
}
