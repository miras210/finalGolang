package main

import (
	"fmt"
	"github.com/miras210/finalGolang/internal/data"
	"net/http"
	"time"
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

	comics := data.Comics{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Batman",
		Pages:     165,
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comics": comics}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
