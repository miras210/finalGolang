package main

import (
	"encoding/json"
	"fmt"
	"github.com/miras210/finalGolang/internal/data"
	"net/http"
	"time"
)

// Add a createComicsHandler for the "POST /v1/comics" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createComicsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title string `json:"title"`
		Year  int32  `json:"year"`
		Pages int32  `json:"runtime"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showComicsHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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
		app.serverErrorResponse(w, r, err)
	}
}
