package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *application) http.Handler {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Post("/identify", registerUser(app))

	return router
}
