package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	webPort = 80
)

type Config struct{}

func main() {
	app := Config{}

	log.Println("Starting broker service on port", webPort)

	// define http server
	serve := &http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: app.routes(),
	}

	// start the server
	err := serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
