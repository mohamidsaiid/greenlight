package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envlope map[string]interface{}

// retrive the id param from the current request context
// then convert it into an integer and return it
// if the operation isn't succecful return 0 and an error
func (app *application) readIdParam(r *http.Request) (int64, error){
	// to get the prameters from the url
	params := httprouter.ParamsFromContext(r.Context())

	// parse the parameter of the id into an int to be used later 
	// using the strconv to convert it from string to int
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// to write the data into JSON and write it to the request using respWriter 
// adding any specific http header to the response
// setting the Content-type header to application/json
func (app *application) writeJSON(w http.ResponseWriter, status int, data envlope, headers http.Header) error {
	// here to parse the data from its own structure to a JSON output the write
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	// adding any needed http header to the response
	for key, value := range headers {
		w.Header()[key] = value
	}
		
	// setting the content-type header to JSON insted of application/text
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}