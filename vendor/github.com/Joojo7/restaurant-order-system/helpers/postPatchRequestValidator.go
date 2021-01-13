package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

//PostPatchRequestValidator for creating menus
func PostPatchRequestValidator(response http.ResponseWriter, request *http.Request, err error) bool {
	response.Header().Add("Content-Type", "application/json")

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(response, msg, http.StatusBadRequest)
			return false

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(response, msg, http.StatusBadRequest)
			return false

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(response, msg, http.StatusBadRequest)
			return false

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(response, msg, http.StatusBadRequest)
			return false

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(response, msg, http.StatusBadRequest)
			return false

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(response, msg, http.StatusRequestEntityTooLarge)
			return false

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return false

		}

	}

	// Call the single body validator
	return SingleJsonValidator(response, request, err)

	// fmt.Fprintf(response, "menu: %+v", result)
}

//MaxRequestValidator uses http.MaxBytesReader to enforce a maximum read of 1MB .
func MaxRequestValidator(response http.ResponseWriter, request *http.Request) *json.Decoder {

	// Use http.MaxBytesReader to enforce a maximum read of 1MB .
	request.Body = http.MaxBytesReader(response, request.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it to check for unexpected feilds
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	return dec

	// fmt.Fprintf(response, "menu: %+v", result)
}

//ContentTypeValidator checks for content type existence and check for json validity
func ContentTypeValidator(response http.ResponseWriter, request *http.Request) {

	// check for content type existence and check for json validity
	if request.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(request.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(response, msg, http.StatusUnsupportedMediaType)
			return
		}
	}
	return
}

// SingleJsonValidator calls decode again, using a pointer to an empty anonymous struct as
// the destination. If the request body only contained a single JSON
// object this will return an io.EOF error. So if we get anything else,
// we know that there is additional data in the request body.
func SingleJsonValidator(response http.ResponseWriter, request *http.Request, err error) bool {

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = MaxRequestValidator(response, request).Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(response, msg, http.StatusBadRequest)
		return false
	}
	return true
}

func IsNilFixed(i interface{}) bool {
	fmt.Print(i)
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
