package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/erodrigufer/GoTennis/pkg/models"
)

// Handler for root URL '/'
func (app *application) root(w http.ResponseWriter, r *http.Request) {
	// The following section of code was only necessary with the default
	// servemux of the standard library, but pat (library) only serves exactly
	// '/' as request to this method

	// If URL request != "/" -> response HTTP 404
	//	if r.URL.Path != "/" {
	//		app.notFound(w)
	//		return
	//	}
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
	// Get from the URL the string value for :id and convert it to an integer
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

// Use GET request to get a form to do a POST request for a tennis session
func (app *application) createSessionForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

// Use POST request to create a new tennis session
func (app *application) createSession(w http.ResponseWriter, r *http.Request) {

	// First call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to
	// send a 400 Bad Request response to the client
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	// Create a new session record in the database using the form data.
	id, err := app.session.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect to the resource of the newly created tennis session
	// The client's browser should then automatically redirect to the new
	// resource
	http.Redirect(w, r, fmt.Sprintf("/session/%d", id), http.StatusSeeOther)
}
