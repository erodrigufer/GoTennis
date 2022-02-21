package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Send error message and stack trace to error logger
// then send a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// the first parameter of Output equals the calldepth, which is the count
	// of the number of frames to skip when computing the file name
	// and line number. So basically, just go back on the stack trace to display
	// the name of function (file) which called the error logging helper
	// function
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// send specific status upon certain errors inquired by the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// convenience wrapper for a 404 not found resource
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
