package main

import (
	"fmt"
	"log"

	// for handling http req/res
	"net/http"
)

func main() {
	// static files and assets
	fileServer := http.FileServer(http.Dir("./static"))

	// routes
	http.Handle("/", fileServer)

	fmt.Println("Starting front end service on port 8080")
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Panic(err)
	}
}
