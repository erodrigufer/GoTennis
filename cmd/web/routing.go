package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// Method to create mux, routing paths and initialize multiple middlewares,
// before returning a servemux
func (app *application) routes() http.Handler {

	// Use the pat.New() function to initialize a new servemux, then
	// register the root app method as the handler for the "/" URL pattern.
	// instead of using functions, use methods to inject the packet-wide
	// dependencies (like the loggers) in every method without using global
	// variables
	// the routing uses clean URLs and is method-based
	// Use a sessionManager to store and load session data in all routes that
	// are dynamic, and might need access to the cookie-based session data
	// The static paths do not need access to the session data, since they are
	// stateless
	mux := pat.New()
	mux.Get("/", app.sessionManager.Enable(http.HandlerFunc(app.root)))
	mux.Get("/session/create", app.sessionManager.Enable(http.HandlerFunc(app.createSessionForm)))
	mux.Post("/session/create", app.sessionManager.Enable(http.HandlerFunc(app.createSession)))
	mux.Get("/session/:id", app.sessionManager.Enable(http.HandlerFunc(app.showSession)))
	mux.Get("/user/signup", app.sessionManager.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", app.sessionManager.Enable(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", app.sessionManager.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", app.sessionManager.Enable(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", app.sessionManager.Enable(http.HandlerFunc(app.logoutUser)))

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
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// chain of middlewares being executed before the mux, e.g.
	// a defer function to recover from a panic from within a client's connec.
	// (the go routine for the client), a logger for all requests and then
	// secureHeaders executes its instructions and then returns the next http
	// Handler in the chain of events, in this case the mux
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
