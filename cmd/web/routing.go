package main

import (
	"net/http"
)

// Method to create mux, routing paths and initialize a secure header middleware
// before routing requests with the mux
func (app *application) routes() http.Handler {

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the root app method as the handler for the "/" URL pattern.
	// instead of using functions, use methods to inject the packet-wide
	// dependencies (like the loggers) in every method without using global
	// variables
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.root)
	mux.HandleFunc("/session/create", app.createSession)
	mux.HandleFunc("/session", app.showSession)

	// Create a handler/fileServer for all files in the static directory
	// Type Dir implements the interface required by FileServer and makes the
	// code portable by using the native file system (which could be different
	// for Windows and other Unix systems)
	// ./ui/static/ will be the root (like a root jail) of the fileServer
	// it will serve files relative to this path. Nonetheless, a security
	// concern is that symlink that points outside the 'jail' can also be
	// followed (check documentation of type Dir)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// the fileServer is now the handler for all URL paths starting with
	// '/static/'
	// http.StripPrefix, will create a new http.Handler that first strips the
	// prefix "/static" from the URL request, and passes the new request to
	// the fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// secureHeaders executes its instructions and then returns the next http
	// Handler in the chain of events, in this case the mux
	return secureHeaders(mux)
}
