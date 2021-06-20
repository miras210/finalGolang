package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comics", app.createComicsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comics", app.listComicsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comics/:id", app.showComicsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comics/:id", app.updateComicsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comics/:id", app.deleteComicsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	// Return the httprouter instance.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
