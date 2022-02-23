package main

import (
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
