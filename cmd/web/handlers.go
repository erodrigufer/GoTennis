package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/erodrigufer/GoTennis/pkg/models"
)

// Handler for root URL '/'
func (app *application) root(w http.ResponseWriter, r *http.Request) {
	// If URL request != "/" -> response HTTP 404
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// fetch the last sessions from database
	s, err := app.session.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Create an instance of the templateData struct holding the slice of
	// the latest sessions
	dynamicData := &templateData{Sessions: s}

	// render page
	app.render(w, r, "root.page.tmpl", dynamicData)

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

	// Use the SessionModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.session.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
		// another kind of error
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// structure holding dynamic data passed on to the template for page
	// generation
	dynamicData := &templateData{Session: s}

	// render page
	app.render(w, r, "show.page.tmpl", dynamicData)
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

	// Create some variables holding dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"
	// Pass the data to the Session.Insert() method, receiving the
	// ID of the new record back
	id, err := app.session.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)

		return
	}
	// Redirect the user to the relevant page for the session
	// With the HTTP Status See Other, the request after the redirect on the
	// client-side will be a GET method
	http.Redirect(w, r, fmt.Sprintf("/session?id=%d", id), http.StatusSeeOther)
}
