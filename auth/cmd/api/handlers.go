package main

import (
	"auth/pkg"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // allow this to be omitted when it's empty
}

func respond(w http.ResponseWriter, payload JsonResponse) {
	// _ is error response that we're not doing anything with and don't care about
	res, _ := json.MarshalIndent(payload, "", "\t")

	// send response
	w.Header().Set("Content-Type", "application/json")
	if !payload.Error {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(res)
}

func (app *Config) AuthService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := JsonResponse{
		Error:   false,
		Message: "auth service",
	}

	respond(w, response)
}

func (app *Config) AddRolePerm(w http.ResponseWriter, r *http.Request) {
	var rp pkg.RolePerm
	response := JsonResponse{
		Error:   false,
		Message: "auth service",
	}

	err := json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		response.Error = true
		response.Message = "Invalid JSON body"
		respond(w, response)
		return
	}

	// TODO: send to supertokens
	// ...

	response.Message = fmt.Sprintf("got: %+v", rp)
	respond(w, response)
}
