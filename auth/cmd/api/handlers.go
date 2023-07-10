package main

import (
	"auth/pkg"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // allow this to be omitted when it's empty
}

func respond(w http.ResponseWriter, payload JsonResponse) {
	// _ is error response that we're not doing anything with and don't care about
	res, _ := json.MarshalIndent(payload, "", "\t")

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if payload.Ok {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(res)
}

func (app *Config) AuthService(w http.ResponseWriter, r *http.Request) {
	response := JsonResponse{
		Ok:      true,
		Message: "auth service ok",
	}

	respond(w, response)
}

func (app *Config) AddRolePerm(w http.ResponseWriter, r *http.Request) {
	var rp pkg.RolePerm
	var err error
	response := JsonResponse{
		Ok:      true,
		Message: "",
	}

	err = json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		e := fmt.Sprintf("AddRolePerm: Could not parse request body %s", err.Error())
		response.Ok = false
		response.Message = e
		respond(w, response)
		return
	}

	// send to supertokens
	err = pkg.STAddRolePerm(rp.Role, rp.Permissions)
	if err != nil {
		e := fmt.Sprintf("AddRolePerm: failure @ st endpoint: %s", err.Error())
		response.Ok = false
		response.Message = e
		respond(w, response)
		return
	}

	response.Message = fmt.Sprintf("Processed: %+v", rp)
	respond(w, response)
}

func (app *Config) RemoveRolePerm(w http.ResponseWriter, r *http.Request) {
	var rp pkg.RolePerm
	var err error
	response := JsonResponse{
		Ok:      true,
		Message: "",
	}

	err = json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		e := fmt.Sprintf("RemoveRolePerm: Could not parse request body %s", err.Error())
		response.Ok = false
		response.Message = e
		respond(w, response)
		return
	}

	// send to supertokens
	err = pkg.STDelRolePerm(rp.Role, rp.Permissions)
	if err != nil {
		e := fmt.Sprintf("RemoveRolePerm: failure @ st endpoint: %s", err.Error())
		response.Ok = false
		response.Message = e
		respond(w, response)
		return
	}

	response.Message = fmt.Sprintf("Processed: %+v", rp)
	respond(w, response)
}
