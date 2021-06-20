package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//
	router.HandlerFunc(http.MethodPost, "/v1/comics", app.requirePermission("comics:write", app.createComicsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/comics", app.requirePermission("comics:read", app.listComicsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/comics/:id", app.requirePermission("comics:read", app.showComicsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/comics/:id", app.requirePermission("comics:write", app.updateComicsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/comics/:id", app.requirePermission("comics:write", app.deleteComicsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Return the httprouter instance.
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
