package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comics", app.createComicsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comics", app.showAllComicsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comics/:id", app.showComicsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comics/:id", app.updateComicsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comics/:id", app.deleteComicsHandler)
	// Return the httprouter instance.
	return router
}
