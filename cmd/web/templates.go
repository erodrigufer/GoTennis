package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/erodrigufer/GoTennis/pkg/forms"
	"github.com/erodrigufer/GoTennis/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that is passeed to the HTML templates
type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Session           *models.Session
	Sessions          []*models.Session // a slice of sessions, useful to store the latest sessions
	Form              *forms.Form
}

// Return a human readable representation of a time.Time object (at UTC)
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// Convert time to UTC before formatting it
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object in a global variable.
// This is a string-keyed map which acts as a lookup between the names of of
// custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages one-by-one
	// pages is a string slice with all the *.page.tmpl files found in a
	// specific directory
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl')
		// and assign it to the name variable
		name := filepath.Base(page)
		// Register the template.FuncMap template set before parsing the files
		// First create an empty template set with template.New(), then register
		// the custom template functions and finally parse the files
		// Finally parse the page template file in to a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in this case, it's just the 'base' layout at the
		// moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in this case, it's just the 'footer' partial at the
		// moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}
	// Return the template set's map
	return cache, nil
}
