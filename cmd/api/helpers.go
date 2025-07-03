package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mohamidsaiid/greenlight/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envlope map[string]interface{}

// retrive the id param from the current request context
// then convert it into an integer and return it
// if the operation isn't succecful return 0 and an error
func (app *application) readIdParam(r *http.Request) (int64, error) {
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// To limit the size of the recived request 
	// We would use the MaxBytesReader() method
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decode the req body into to targeted destaniation
	dec := json.NewDecoder(r.Body)
	// To disallow and unknown fileds to be in the response which cannot be mapped to the target destination
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
	
		switch {
			// Using errors.As() to specify the type of the error
		case errors.As(err, &syntaxError) :
			return fmt.Errorf("body contains badly-formed JSON (at charcter %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at charcter %d)", unmarshalTypeError.Offset)
		
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")


		case strings.HasPrefix(err.Error(), "json:unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
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

func (app *application) readString(qs url.Values, key, defaultValue string) string {
	// extract the value of a given key. it it was provided the get method would return 
	// otherwise would return ""
	s := qs.Get(key)
	
	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
	}()

	fn()
}
