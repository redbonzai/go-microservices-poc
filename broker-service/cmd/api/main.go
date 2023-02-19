package main

import (
	"fmt"
	"log"
	"net/http"
)

const brokerPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", brokerPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", brokerPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
