package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// Show session's information
func showSession(w http.ResponseWriter, r *http.Request) {
	// Get from the URL the string value for id and convert it to an integer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// The string to int convertion failed or the int is less than 1
	// id could be less than 1 if the query was wrong, since it will return nil
	if err != nil || id < 1 {
		// Reply to the request with a 404 Not Found error
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	_, err = fmt.Fprintf(w, "Display tennis session with ID %d...", id)
	if err != nil {
		log.Println(err)
	}

}

// Use POST request to create a new tennis session
func createSession(w http.ResponseWriter, r *http.Request) {
	// check field 'Method' from http request
	// if 'Method' is not POST, then send 405 error in the header 'Method not
	// allowed
	if r.Method != "POST" {
		// Set 'Allow' field in response header to POST, to tell the client
		// which method(s) are valid
		w.Header().Set("Allow", "POST")
		// The error code must be sent in the header before the body, otherwise
		// is the default header code http.StatusOK
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method Not Allowed"))

		// Use the http.Error() function to send a 405 status code
		// and "Method Not Allowed" string as the response body.
		// This helper function replaces both calls to WriteHeader and Write
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

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
