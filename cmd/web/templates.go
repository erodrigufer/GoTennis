package main

import "github.com/erodrigufer/GoTennis/pkg/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
type templateData struct {
	Session *models.Session
	// a slice of sessions, useful to store the latest sessions
	Sessions []*models.Session
}
