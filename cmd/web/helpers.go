package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/erodrigufer/GoTennis/pkg/models"
	"github.com/justinas/nosurf"
)

// Default data is automatically added every time a template is rendered to the
// dynamic data being passed to the app.render method
// Default data added: CurrentYear, Flash messages and AuthenticatedUser
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		// if pointer is nil, create a new instance of templateData
		td = &templateData{}
	}
	// CSRF token is available in all templates per default
	// The token for a particular request is embedded into the html, so that
	// only the legitimate agent could get that token and actually responds
	// with this token value as a hidden input field in all POST and 'non-safe'
	// requests
	td.CSRFToken = nosurf.Token(r)
	// check if user has already been authenticated
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	// Add any flash message (if it exists) to the template data
	td.Flash = app.sessionManager.PopString(r, "flash")
	return td
}

// Retrieve the appropriate template set from the cache based on the page name
// (like 'root.page.tmpl'). If no entry exists in the cache with the
// provided name, call the serverError helper method
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, dynamicData *templateData) {
	ts, ok := app.templateCache[name]
	// the object did not exist in the cache map
	if !ok {
		// fmt.Errorf returns an object that fulfills the error interface
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	// Initialize a buffer to first execute template into buffer, if there is an
	// error, then the data will not be half-written to the client, but instead
	// will remain in the buffer, and an Internal Server Error will be sent to
	// the client. Something like this can happen, when there is an error in a
	// template, then the Execute() method will return an error
	buf := new(bytes.Buffer)
	// Execute the template set, passing in any dynamic data
	err := ts.Execute(buf, app.addDefaultData(dynamicData, r))
	if err != nil {
		app.serverError(w, err)
		return // Do not send the template back to the client
	}
	// There was no error while executing/rendering the template, so send the
	// whole template back to the client
	buf.WriteTo(w)
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

// check if the current session's user has been successfully authenticated,
// if so, return its user struct from the db data, if not return nil
func (app *application) authenticatedUser(r *http.Request) int {
	// return the value inside the context associated with the key
	// 'contextKeyUser' and type-cast it to the models.User type/struct
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	// type-cast failed
	if !ok {
		return nil
	}
	return user
}
