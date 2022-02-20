package main

import (
	"log"
	"net/http"
)

func main() {
	// Define HOST and PORT
	SERVICE := "localhost:4000"
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the root function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", root)
	mux.HandleFunc("/session/create", createSession)
	mux.HandleFunc("/session", showSession)

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

	// Start a TCP web server listening on PORT
	// If http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	log.Printf("Starting server at %s", SERVICE)
	err := http.ListenAndServe(SERVICE, mux)
	log.Fatal(err)
}

// Eduardo Rodriguez @erodrigufer (c) 2022
