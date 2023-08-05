package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...),
		ExposedHeaders: []string{"Link"},

		//warning! this is incompatible with allowedOrigins "*"
		AllowCredentials: true,

		MaxAge: 300,
	}))

	// SuperTokens Middleware - adds these apis: https://app.swaggerhub.com/apis/supertokens/FDI/1.16.0
	mux.Use(supertokens.Middleware)

	// healthcheck
	mux.Use(middleware.Heartbeat("/ping"))

	// authZ
	mux.Post("/health", app.AuthService)
	mux.Put("/add-role-perm", app.AddRolePerm)
	mux.Post("/remove-role-perm", app.RemoveRolePerm)

	return mux
}
