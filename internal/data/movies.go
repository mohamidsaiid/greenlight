package data

import (
	"time"

	"github.com/mohamidsaiid/greenlight/internal/validator"
)

type Movie struct {
	// Unique integer id for each movie
	ID int64 `json:"id"`
	// Timestamp to indecate when the movie added to our database
	CreatedAt time.Time `json:"-"`
	// Movie title
	Title string `json:"title"`
	// Movie release year
	Year int32 `json:"year,omitempty"`
	// Movie runtime (in minutes)
	Runtime Runtime `json:"runtime,omitempty"`
	// Slice for genres of the movie (action, romance, ...etc)
	Genres []string `json:"genres,omitempty"`
	// The version no. starts with one and is goning to be increamented each time the movie information is updated
	Version int32 `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {

	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be greater than 500 bytes long") 


	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be after year 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a postive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")


}