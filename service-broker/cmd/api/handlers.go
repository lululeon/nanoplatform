package main

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // allow this to be omitted when it's empty
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := JsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	// _ is error response that we're not doing anything with and don't care about
	res, _ := json.MarshalIndent(response, "", "\t")

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(res)
}
