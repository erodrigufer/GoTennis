package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/erodrigufer/GoTennis/pkg/models"

	"github.com/justinas/nosurf" // anti-CSRF mechanisms
)

// This middleware sets two header values for all incoming server requests
// These two header values should instruct the client's web browser to implement
// some additional security measures to protect against XSS and clickjacking.
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
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			// http.StatusFound (302) URI of requested resource has been changed
			// temporarily
			return // return so that other middlewares are not executed
		}

		// call the next handler in the chain if the user has been authenticated
		next.ServeHTTP(w, r)
	})
}

// Anti-CSRF middleware with a customized cookie with Secure, Path and HttpOnly
// flags set
func noSurf(next http.Handler) http.Handler {
	// If the CSRF check succeeds, the "next" handler will be called
	csrfHandler := nosurf.New(next)

	// sets the cookie token used for the CSRF token
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// Add user info stored in the db as a struct to the context of a request, if
// the user is authenticated (logged in/userID can be found in the session) and
// the db has info for a user with that userID (the user has not been deleted)
// if not then pass the unchanged request with the original context to the other
// handlers.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a userID value Exists in the session. If the userID *is not
		// present* then call the next handler in the chain as normal
		exists := app.sessionManager.Exists(r, "userID")
		// userID does not exist, call the next http.Handler
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		// Fetch the details of the current user from the database. If
		// no matching record is found, remove the (invalid) userID from
		// their session and call the next handler in the chain as normal
		user, err := app.users.Get(app.sessionManager.GetInt(r, "userID"))
		// the user was eliminated from the db in the meantime, remove the
		// userID from the session as well, and serve next http.Handler as usual
		if err == models.ErrNoRecord {
			app.sessionManager.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		// Otherwise, the request is coming from a valid, authenticated
		// (logged in) user. A new copy of the request is created with the user
		// information added to the request context, and the next handler in
		// the chain is called *using this new copy of the request*
		// - r.Context() retrieves the existing context for request r
		// - the method context.WithValue creates a new copy of the context from
		// r.Context() and appends the struct data 'user' with the key
		// 'contextKeyUser'
		// - r.WithContext creates a copy of the request r with the new context
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
