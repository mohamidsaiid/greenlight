package data

import "time"

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
