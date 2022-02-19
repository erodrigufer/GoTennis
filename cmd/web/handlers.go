package main

import (
	"fmt"
	"html/template"
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
	// Slice containing the paths to the two files. Note that the
	// root.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/root.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	// template.ParseFiles() func. to read the template file into template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		// http.Error() function to send a generic 500 Internal Server Error
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Execute() method on the template set to write the template content
	// as the response body. The last parameter to Execute() represents
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

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
		log.Println(err.Error())
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
