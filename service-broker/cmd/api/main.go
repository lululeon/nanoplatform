package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "3000"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Service broker up and listening at :%s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
