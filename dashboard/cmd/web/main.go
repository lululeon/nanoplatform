package main

import (
	"fmt"
	"log"

	// for handling http req/res
	"net/http"
)

const webPort = "3001"

func main() {
	// static files and assets - note abs path from where we (tend to) invoke go run!! Not portable.
	// Eventually / for prod, you'd look @ embedding: https://pkg.go.dev/embed#hdr-File_Systems
	fileServer := http.FileServer(http.Dir("./dashboard/cmd/web/static"))
	listenPort := fmt.Sprintf(":%s", webPort)

	// routes
	http.Handle("/", fileServer)

	fmt.Printf("Starting front end service on port %s", listenPort)
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		log.Panic(err)
	}
}
