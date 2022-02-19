package main

import (
	"log"
	"net/http"
)

// Handler for root URL '/'
func root(w http.ResponseWriter, r *http.Request) {
	// If URL request != "/" -> response HTTP 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Response for correct request
	w.Write([]byte("Tenis CCE - Belén"))
}

// Use POST request to create a new tennis session
func createSession(w http.ResponseWriter, r *http.Request) {
	// check field 'Method' from http request
	// if 'Method' is not POST, then send 405 error in the header 'Method not
	// allowed
	if r.Method != "POST" {
		// The error code must be sent in the header before the body, otherwise
		// is the default header code http.StatusOK
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	w.Write([]byte("Create a new tennis session..."))
}
func main() {
	// Define HOST and PORT
	SERVICE := "localhost:4000"
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the root function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", root)
	mux.HandleFunc("/createSession", createSession)
	// Start a TCP web server listening on PORT
	// If http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	log.Printf("Starting server on %s", SERVICE)
	err := http.ListenAndServe(SERVICE, mux)
	log.Fatal(err)
}

// Eduardo Rodriguez @erodrigufer (c) 2022
