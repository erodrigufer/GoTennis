package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Retrieve the appropriate template set from the cache based on the page name
// (like 'root.page.tmpl'). If no entry exists in the cache with the
// provided name, call the serverError helper method
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, dynamicData *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		// fmt.Errorf returns an object that fulfills the error interface
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	// Execute the template set, passing in any dynamic data
	err := ts.Execute(w, dynamicData)
	if err != nil {
		app.serverError(w, err)
	}
}

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
