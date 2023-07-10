package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"auth/pkg"
)

const webPort = "7567"

type Config struct {
	SupertokensServerUrl string
}

func main() {
	stServerUrl := os.Getenv("SUPERTOKENS_SERVER_URL")
	if len(stServerUrl) == 0 {
		log.Fatal("No value found for supertokens server url, which cannot be blank. Exiting.")
	}

	app := Config{}

	stErr := pkg.InitSupertokensAuth(stServerUrl)
	if stErr != nil {
		log.Println("Cannot init supertokens!!")

		// stop
		log.Fatal(stErr.Error())
	}
	log.Println("supertokens initialized...")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("auth service up and listening at :%s", webPort)
}
