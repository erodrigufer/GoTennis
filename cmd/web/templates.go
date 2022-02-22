package main

import (
	"html/template"
	"path/filepath"

	"github.com/erodrigufer/GoTennis/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
type templateData struct {
	Session *models.Session
	// a slice of sessions, useful to store the latest sessions
	Sessions []*models.Session
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
		// Parse the page template file in to a template set
		ts, err := template.ParseFiles(page)
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
