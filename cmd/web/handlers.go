package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Handler for root URL '/'
func (app *application) root(w http.ResponseWriter, r *http.Request) {
	// If URL request != "/" -> response HTTP 404
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// Slice containing the paths to the two files. Note that the
	// root.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/root.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	// template.ParseFiles() func. to read the template file into template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// 500 Internal Server Error
		app.serverError(w, err)
		return
	}
	// Execute() method on the template set to write the template content
	// as the response body. The last parameter to Execute() represents
	// dynamic data that we want to pass in
	err = ts.Execute(w, nil)
	if err != nil {
		// 500 Internal Server Error
		app.serverError(w, err)
		return
	}

}

// Show session's information
func (app *application) showSession(w http.ResponseWriter, r *http.Request) {
	// Get from the URL the string value for id and convert it to an integer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// The string to int convertion failed or the int is less than 1
	// id could be less than 1 if the query was wrong, since it will return nil
	if err != nil || id < 1 {
		// Reply to the request with a 404 Not Found error
		app.notFound(w)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	_, err = fmt.Fprintf(w, "Display tennis session with ID %d...", id)
	if err != nil {
		// Internal Server Error
		app.serverError(w, err)
	}

}

// Use POST request to create a new tennis session
func (app *application) createSession(w http.ResponseWriter, r *http.Request) {
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

		// This helper function replaces both calls to WriteHeader and Write
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new tennis session..."))
}
