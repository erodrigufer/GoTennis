package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/erodrigufer/GoTennis/pkg/forms"
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

// After GET request, respond with a form to do a POST request for a
// new tennis session
func (app *application) createSessionForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// Create a new tennis session after receiving a POST request: validate the data
// from the POST request, if valid Insert() data into the database. If the data
// is not valid, render the forms page with error messages in the invalid fields
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
	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm) // the parameter are the url.Values POSTed
	// into the form
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	// If the form is not valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	// Because the form data (with type url.Values) has been anonymously embeded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field, to then add Insert the
	// data into a new row in the SQL database
	id, err := app.session.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Add a string value to the corresponding key ("flash") to the session data
	// Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware
	app.sessionManager.Put(r, "flash", "Tennis session was successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/session/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// parse POST request data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the contents  of the form
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// Try to create a new user record in the database. If the email already
	// exists, add an error message to the form and re-display it
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	// email address is already stored in the users table
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		// another kind of error happened, on the server-side
		app.serverError(w, err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
