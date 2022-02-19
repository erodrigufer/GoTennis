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
	// Start a TCP web server listening on PORT
	// If http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	log.Printf("Starting server on %s", SERVICE)
	err := http.ListenAndServe(SERVICE, mux)
	log.Fatal(err)
}

// Eduardo Rodriguez @erodrigufer (c) 2022
