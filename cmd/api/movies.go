package main

import (
	"time"
	"net/http"
	"fmt"
	
	"github.com/mohamidsaiid/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// read the id value from the id parameter in the URL
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return 
	}

	movie := data.Movie {
		ID : id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Runtime: 102,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}
	
	err = app.writeJSON(w, http.StatusOK, envlope{"movie":movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
