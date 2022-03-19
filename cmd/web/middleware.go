package main

import (
	"fmt"
	"net/http"
)

// This middleware sets two header values for all incoming server requests
// These two header values should protect against XSS and clickjacking
func secureHeaders(next http.Handler) http.Handler {
	// Explanation:
	// http.HandlerFunc is a type that works as an adapter to allow a function f
	// to be returned as an http.Handler (which is an interface, that requires
	// the ServeHTTP() method), since the method ServeHTTP of the type
	// HandlerFunc simply calls the HandlerFunc which has a signature:
	// func(ResponseWriter, *Request)
	// In the above code we are type casting an anonymous function into a
	// HandlerFunc, so that when this function is called with ServeHTTP, it will
	// execute, and in the last step call the ServeHTTP method of the next http
	// handler in the chain of handlers
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// Log every client's request.
// Log IP address of client, protocol used, http method and requested URL
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.RemoteAddr is the IP address of the client doing the request
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

// Send an Internal Server Error message code to a client, when the server has
// to close an http connection with a client due to a panic inside the goroutine
// handling the client
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not.
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// check if user is authenticated, if so serve next middleware in chain, if not
// redirect user to '/user/login' and return from middleware
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect the user to the login page
		// and return from the middleware so that no subsequent handlers in the
		// middleare chain are executed.
		if app.authenticatedUser(r) == 0 {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			// http.StatusFound (302) URI of requested resource has been changed
			// temporarily
			return // return so that other middlewares are not executed
		}

		// call the next handler in the chain if the user has been authenticated
		next.ServeHTTP(w, r)
	})
}
