package main

import (
	"fmt"
	"log"
	"net/http"

	"auth/pkg"
)

const webPort = "7567"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("auth service up and listening at :%s", webPort)

	stErr := pkg.InitSupertokensAuth()
	if stErr != nil {
		fmt.Println("Cannot reach supertokens...")

		// stop
		log.Panic(stErr)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
