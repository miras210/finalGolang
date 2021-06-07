package main

import (
	"fmt"
	"net/http"
)

// Add a createComicsHandler for the "POST /v1/comics" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createComicsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new comic book")
}

// Add a showComicsHandler for the "GET /v1/comics/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showComicsHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of comics %d\n", id)
}
